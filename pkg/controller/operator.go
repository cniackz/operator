// This file is part of MinIO Operator
// Copyright (c) 2020 MinIO, Inc.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package controller

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	runetimetime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"net"
	"net/http"
	"strings"
	"time"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/minio/operator/pkg/common"
	"github.com/minio/operator/pkg/utils"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	miniov2 "github.com/minio/operator/pkg/apis/minio.min.io/v2"
	xcerts "github.com/minio/pkg/certs"
	"github.com/minio/pkg/env"
	"k8s.io/klog/v2"
	k8sscheme "k8s.io/kubectl/pkg/scheme"
)

const (
	// CertPasswordEnv Env variable is used to decrypt the private key in the TLS certificate for operator if need it
	CertPasswordEnv = "OPERATOR_CERT_PASSWD"
	// OperatorDeploymentNameEnv Env variable to specify a custom deployment name for Operator
	OperatorDeploymentNameEnv = "MINIO_OPERATOR_DEPLOYMENT_NAME"
	// OperatorCATLSSecretName is the name of the secret for the operator CA
	OperatorCATLSSecretName = "operator-ca-tls"
	// OperatorCSRSignerCASecretName is the name of the secret for the signer-ca certificate
	// this is a copy of the secret signer-ca in namespace
	OperatorCSRSignerCASecretName = "openshift-csr-signer-ca"
	// OpenshiftKubeControllerNamespace is the namespace of kube controller manager operator in Openshift
	OpenshiftKubeControllerNamespace = "openshift-kube-controller-manager-operator"
	// OpenshiftCATLSSecretName is the secret name of the CRD's signer in kubernetes under  OpenshiftKubeControllerNamespace namespace
	OpenshiftCATLSSecretName = "csr-signer"
	// DefaultDeploymentName is the default name of the operator deployment
	DefaultDeploymentName = "minio-operator"
	// DefaultOperatorImage is the version fo the operator being used
	DefaultOperatorImage = "minio/operator:v5.0.10"
)

// ErrBucketsAlreadyExist - all buckets already exist.
var ErrBucketsAlreadyExist = errors.New("all buckets already exist")

var serverCertsManager *xcerts.Manager

// rolloutRestartDeployment - executes the equivalent to kubectl rollout restart deployment
func (c *Controller) rolloutRestartDeployment(deployName string) error {
	ctx := context.Background()
	namespace := miniov2.GetNSFromFile()
	deployment, err := c.kubeClientSet.AppsV1().Deployments(namespace).Get(ctx, deployName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if deployment.Spec.Template.ObjectMeta.Annotations == nil {
		deployment.Spec.Template.ObjectMeta.Annotations = make(map[string]string)
	}
	deployment.Spec.Template.ObjectMeta.Annotations["kubectl.kubernetes.io/restartedAt"] = time.Now().Format(time.RFC3339)
	data, err := runtime.Encode(k8sscheme.Codecs.LegacyCodec(appsv1.SchemeGroupVersion), deployment)
	if err != nil {
		return err
	}
	_, err2 := c.kubeClientSet.AppsV1().Deployments(namespace).Patch(ctx, deployName, types.StrategicMergePatchType, data, metav1.PatchOptions{FieldManager: "kubectl-rollout"})
	if err2 != nil {
		return err2
	}
	return nil
}

func (c *Controller) fetchUserCredentials(ctx context.Context, tenant *miniov2.Tenant) []*v1.Secret {
	var userCredentials []*v1.Secret
	for _, credential := range tenant.Spec.Users {
		credentialSecret, err := c.kubeClientSet.CoreV1().Secrets(tenant.Namespace).Get(ctx, credential.Name, metav1.GetOptions{})
		if err == nil && credentialSecret != nil {
			userCredentials = append(userCredentials, credentialSecret)
		}
	}
	return userCredentials
}

// getTransport returns a *http.Transport with the collection of the trusted CA certificates
// returns a cached transport if already available
func (c *Controller) getTransport() *http.Transport {
	if c.transport != nil {
		return c.transport
	}
	c.transport = c.createTransport()
	return c.transport
}

// createTransport returns a *http.Transport with the collection of the trusted CA certificates
func (c *Controller) createTransport() *http.Transport {
	rootCAs := c.fetchTransportCACertificates()
	dialer := &net.Dialer{
		Timeout:   15 * time.Second,
		KeepAlive: 15 * time.Second,
	}
	c.transport = &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		MaxIdleConnsPerHost:   1024,
		IdleConnTimeout:       15 * time.Second,
		ResponseHeaderTimeout: 15 * time.Minute,
		TLSHandshakeTimeout:   15 * time.Second,
		ExpectContinueTimeout: 15 * time.Second,
		// Go net/http automatically unzip if content-type is
		// gzip disable this feature, as we are always interested
		// in raw stream.
		DisableCompression: true,
		TLSClientConfig: &tls.Config{
			// Can't use SSLv3 because of POODLE and BEAST
			// Can't use TLSv1.0 because of POODLE and BEAST using CBC cipher
			// Can't use TLSv1.1 because of RC4 cipher usage
			MinVersion: tls.VersionTLS12,
			RootCAs:    rootCAs,
		},
	}

	return c.transport
}

// fetchTransportCACertificates retrieves a *x509.CertPool with all CA that operator will trust
func (c *Controller) fetchTransportCACertificates() (pool *x509.CertPool) {
	rootCAs := miniov2.MustGetSystemCertPool()
	// Default kubernetes CA certificate
	rootCAs.AppendCertsFromPEM(miniov2.GetPodCAFromFile())

	if utils.GetOperatorRuntime() == common.OperatorRuntimeOpenshift {
		// Openshift Service CA certificate
		if serviceCA := miniov2.GetOpenshiftServiceCAFromFile(); serviceCA != nil {
			rootCAs.AppendCertsFromPEM(serviceCA)
		}
		// Openshift csr-signer CA certificate
		if cert := miniov2.GetOpenshiftCSRSignerCAFromFile(); cert != nil {
			rootCAs.AppendCertsFromPEM(cert)
		}
	}

	// Custom ca certificate to be used by operator
	operatorCATLSCert, err := c.kubeClientSet.CoreV1().Secrets(miniov2.GetNSFromFile()).Get(context.Background(), OperatorCATLSSecretName, metav1.GetOptions{})
	if err == nil && operatorCATLSCert != nil {
		if val, ok := operatorCATLSCert.Data["tls.crt"]; ok {
			rootCAs.AppendCertsFromPEM(val)
		}
		if val, ok := operatorCATLSCert.Data["ca.crt"]; ok {
			rootCAs.AppendCertsFromPEM(val)
		}
		if val, ok := operatorCATLSCert.Data["public.crt"]; ok {
			rootCAs.AppendCertsFromPEM(val)
		}
	}
	return rootCAs
}

// GetSignerCAFromSecret Retrieves the CA certificate for Openshift CSR signed certificates from
// openshift-kube-controller-manager-operator namespace
func (c *Controller) GetSignerCAFromSecret() ([]byte, error) {
	var caCertificate []byte
	openShiftCATLSCertSecret, err := c.kubeClientSet.CoreV1().Secrets(OpenshiftKubeControllerNamespace).Get(
		context.Background(), OpenshiftCATLSSecretName, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, fmt.Errorf("%s secret in %s wasn't found", OpenshiftCATLSSecretName, OpenshiftKubeControllerNamespace)
		}
		return nil, fmt.Errorf("%s secret was found but we failed to load the secret: %#v", OpenshiftCATLSSecretName, err)
	} else if openShiftCATLSCertSecret != nil {
		if val, ok := openShiftCATLSCertSecret.Data[common.TLSCRT]; ok {
			caCertificate = val
		}
	}
	return caCertificate, nil
}

// GetOpenshiftCSRSignerCAFromSecret loads the tls certificate in openshift-csr-signer-ca secret in operator namespace
func (c *Controller) GetOpenshiftCSRSignerCAFromSecret() ([]byte, error) {
	var caCertificate []byte
	operatorNamespace := miniov2.GetNSFromFile()
	openShifCSRSignerCATLSCertSecret, err := c.kubeClientSet.CoreV1().Secrets(operatorNamespace).Get(
		context.Background(), OperatorCSRSignerCASecretName, metav1.GetOptions{})
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, fmt.Errorf("%s secret wasn't found, Skip load the certificate", OperatorCSRSignerCASecretName)
		}
		// Lack of permissions to read the secret
		return nil, fmt.Errorf("%s secret was found but we failed to get certificate: %#v", OperatorCSRSignerCASecretName, err)
	} else if openShifCSRSignerCATLSCertSecret != nil {
		// When secret was obtained with no errors
		if val, ok := openShifCSRSignerCATLSCertSecret.Data[common.TLSCRT]; ok {
			// OpenShift csr-signer secret has tls.crt certificates that we need to append in order
			// to trust the signer. If we append the val, Operator will be able to provisioning the
			// initial users and get Tenant Health, so tenant can be properly initialized and in
			// green status, otherwise if we don't append it, it will get stuck and expose this
			// issue in the log:
			// Failed to get cluster health: Get "https://minio.tenant-lite.svc.cluster.local/minio/health/cluster":
			// x509: certificate signed by unknown authority
			caCertificate = val
		}
	}
	return caCertificate, nil
}

// checkOpenshiftSignerCACertInOperatorNamespace checks if csr-signer secret in openshift changed and updates or create
// a copy of the secret in operator namespace
func (c *Controller) checkOpenshiftSignerCACertInOperatorNamespace(ctx context.Context) error {
	// get the current certificate from openshift
	csrSignerCertificate, err := c.GetSignerCAFromSecret()
	if err != nil {
		return err
	}
	namespace := miniov2.GetNSFromFile()
	// get openshift-csr-signer-ca secret in minio-operator namespace
	csrSignerSecret, err := c.kubeClientSet.CoreV1().Secrets(namespace).Get(ctx, OperatorCSRSignerCASecretName, metav1.GetOptions{})
	if err != nil {
		// if csrSignerCa doesnt exists create it
		if k8serrors.IsNotFound(err) {
			klog.Infof("'%s/%s' secret is missing, creating", namespace, OperatorCSRSignerCASecretName)
			operatorDeployment, err := c.kubeClientSet.AppsV1().Deployments(namespace).Get(ctx, getOperatorDeploymentName(), metav1.GetOptions{})
			if err != nil {
				return err
			}

			ownerReference := metav1.OwnerReference{
				APIVersion: appsv1.SchemeGroupVersion.Version,
				Kind:       "Deployment",
				Name:       operatorDeployment.Name,
				UID:        operatorDeployment.UID,
			}

			csrSignerSecret := &v1.Secret{
				Type: "Opaque",
				ObjectMeta: metav1.ObjectMeta{
					Name:            OperatorCSRSignerCASecretName,
					Namespace:       miniov2.GetNSFromFile(),
					OwnerReferences: []metav1.OwnerReference{ownerReference},
				},
				Data: map[string][]byte{
					common.TLSCRT: csrSignerCertificate,
				},
			}
			_, err = c.kubeClientSet.CoreV1().Secrets(namespace).Create(ctx, csrSignerSecret, metav1.CreateOptions{})
			// Reload CA certificates
			c.createTransport()
			return err
		}
		return err
	}

	if caCert, ok := csrSignerSecret.Data[common.TLSCRT]; ok && !bytes.Equal(caCert, csrSignerCertificate) {
		csrSignerSecret.Data[common.TLSCRT] = csrSignerCertificate
		_, err = c.kubeClientSet.CoreV1().Secrets(namespace).Update(ctx, csrSignerSecret, metav1.UpdateOptions{})
		if err != nil {
			return err
		}
		klog.Infof("'%s/%s' secret changed, updating '%s/%s' secret", OpenshiftKubeControllerNamespace, OpenshiftCATLSSecretName, namespace, OperatorCSRSignerCASecretName)
		c.fetchTransportCACertificates()
		// Reload CA certificates
		c.createTransport()
	}
	return nil
}

func (c *Controller) createUsers(ctx context.Context, tenant *miniov2.Tenant, tenantConfiguration map[string][]byte) (err error) {
	defer func() {
		if err == nil {
			if _, err = c.updateProvisionedUsersStatus(ctx, tenant, true); err != nil {
				klog.V(2).Infof(err.Error())
			}
		}
	}()

	userCredentials := c.fetchUserCredentials(ctx, tenant)
	if len(userCredentials) == 0 {
		return nil
	}

	if _, err = c.updateTenantStatus(ctx, tenant, StatusProvisioningInitialUsers, 0); err != nil {
		return err
	}

	// get a new admin client
	adminClient, err := tenant.NewMinIOAdmin(tenantConfiguration, c.getTransport())
	if err != nil {
		klog.Errorf("Error instantiating adminClnt: %v", err)
		return err
	}

	// configuration that means MinIO is running with LDAP enabled
	// and we need to skip the console user creation
	if err = tenant.CreateUsers(adminClient, userCredentials, tenantConfiguration); err != nil {
		klog.V(2).Infof("Unable to create MinIO users: %v", err)
	}

	return err
}

func (c *Controller) verifyLDAPParams(ctx context.Context, tenant *miniov2.Tenant) (err error) {
	// Verify there is a user
	user := tenant.Spec.LDAPPolicyAttachToSingleUser.User
	policy := tenant.Spec.LDAPPolicyAttachToSingleUser.Policy
	if user != "" && policy != "" {
		klog.Infof("Attach policy: %s to user: %user", policy, user)
	} else {
		klog.Errorf("User and Policy have to be provided")
		return errors.New("User and Policy have to be provided")
	}
	return nil
}

func (c *Controller) ldapPolicyAttachToSingleUser(ctx context.Context, tenant *miniov2.Tenant) (err error) {
	result := c.verifyLDAPParams(ctx, tenant)
	if result != nil {
		klog.Errorf("LDAP Params did not pass validation")
		return err
	}
	klog.Info("Configuring LDAP in the Tenant")
	tenantName := tenant.Name
	tenantNamespace := tenant.Namespace
	pool0Name := tenant.Spec.Pools[0].Name
	podName := tenantName + "-" + pool0Name + "-0"
	request := c.kubeClientSet.CoreV1().RESTClient().Post().Resource("pods").Name(podName).Namespace(
		tenantNamespace).SubResource("exec")
	scheme := runetimetime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		klog.Error("There was an error trying to configure LDAP User")
		return err
	}
	parameterCodec := runetimetime.NewParameterCodec(scheme)
	LDAPPolicyAttachToSingleUser := tenant.Spec.LDAPPolicyAttachToSingleUser
	User := LDAPPolicyAttachToSingleUser.User
	Policy := LDAPPolicyAttachToSingleUser.Policy
	MinIOServerEndpoint := tenant.MinIOServerEndpoint()
	tenant.ConfigurationSecretName()
	tenant.HasConfigurationSecret()
	tenantSpecConfigName := tenant.Spec.Configuration.Name
	secret, _ := c.kubeClientSet.CoreV1().Secrets(tenant.Namespace).Get(ctx, tenantSpecConfigName, metav1.GetOptions{})
	secretData := secret.Data
	content := string(secretData["config.env"])
	contents := strings.Split(content, "\n")
	rootAccessKey := ""
	rootSecretKey := ""
	for i := range contents {
		content := contents[i]
		if strings.Contains(content, "MINIO_ROOT_USER") {
			valor := strings.Split(content, "=")
			user := valor[1]
			plainUser := user[1 : len(user)-1]
			rootAccessKey = plainUser
		}
		if strings.Contains(content, "MINIO_ROOT_PASSWORD") {
			valor := strings.Split(content, "=")
			password := valor[1]
			plainPassword := password[1 : len(password)-1]
			rootSecretKey = plainPassword
		}
	}
	// export set HOME=/tmp is to allow mc to work when not using root user in security context to avoid
	// mc: <ERROR> Unable to save new mc config. mkdir /.mc: permission denied.
	alias := "export set HOME=/tmp; mc alias set myminio " + MinIOServerEndpoint + " " + rootAccessKey + " " + rootSecretKey
	ldapCommand := "mc idp ldap policy attach myminio " + Policy + " --user=" + "'" + User + "'"
	klog.Info(ldapCommand)
	finalCommand := alias + "; " + ldapCommand
	cmds := []string{
		"sh",
		"-c",
		finalCommand,
	}
	request.VersionedParams(&corev1.PodExecOptions{
		Stdin:     false,
		Stdout:    true,
		Stderr:    true,
		TTY:       false,
		Container: "minio",
		Command:   cmds,
	}, parameterCodec)
	// Get a rest.Config from the kubeconfig file.  This will be passed into all
	// the client objects we create.
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	restConfig, err := kubeconfig.ClientConfig()
	if err != nil {
		klog.Errorf("Error attaching LDAP Policy to user %s", err)
		return err
	}
	exec, err := remotecommand.NewSPDYExecutor(restConfig, "POST", request.URL())
	fmt.Println(exec)
	if err != nil {
		klog.Errorf("Error attaching LDAP Policy to user %s", err)
		return err
	}
	var stdout, stderr bytes.Buffer
	fmt.Println(stdout)
	fmt.Println(stderr)
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  nil,
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	})
	if err != nil {
		klog.Errorf("Error attaching LDAP Policy to user %s", err)
		return err
	}
	actualOutput := stdout.String()
	actualError := stderr.String()
	klog.Infof("Output from pod: %v", actualOutput)
	klog.Infof("Error from pod: %v", actualError)
	// If expected output is contained by the actual output, it will be
	// considered as a success.
	if strings.Contains(actualError, "The specified policy change is already in effect") {
		klog.Info("Command(s) succeeded, will not be executed again")
		// By saving the command in the status, we guarantee that will not be re-executed.
		if _, err = c.updateLDAPPolicyAttachToSingleUserStatus(ctx, tenant, true); err != nil {
			klog.V(2).Infof(err.Error())
		}
	}
	return nil
}

func (c *Controller) createBuckets(ctx context.Context, tenant *miniov2.Tenant, tenantConfiguration map[string][]byte) (err error) {
	defer func() {
		if err == nil {
			if _, err = c.updateProvisionedBucketStatus(ctx, tenant, true); err != nil {
				klog.V(2).Infof(err.Error())
			}
		}
	}()

	tenantBuckets := tenant.Spec.Buckets
	if len(tenantBuckets) == 0 {
		return nil
	}

	// get a new admin client
	minioClient, err := tenant.NewMinIOUser(tenantConfiguration, c.getTransport())
	if err != nil {
		// show the error and continue
		klog.Errorf("Error instantiating minio Client: %v ", err)
		return err
	}
	err = tenant.CheckBucketsExist(minioClient, tenantBuckets...)
	if err != nil {
		if _, err = c.updateTenantStatus(ctx, tenant, StatusProvisioningDefaultBuckets, 0); err != nil {
			return err
		}
		if err = tenant.CreateBuckets(minioClient, tenantBuckets...); err != nil {
			klog.V(2).Infof("Unable to create MinIO buckets: %v", err)
			return err
		}
	}
	return ErrBucketsAlreadyExist
}

// getOperatorDeploymentName Internal func returns the Operator deployment name from MINIO_OPERATOR_DEPLOYMENT_NAME ENV variable or the default name
func getOperatorDeploymentName() string {
	return env.Get(OperatorDeploymentNameEnv, DefaultDeploymentName)
}
