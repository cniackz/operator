// This file is part of MinIO Operator
// Copyright (c) 2021 MinIO, Inc.
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

package api

// list of all console environment constants
const (
	OperatorSAToken = "OPERATOR_SA_TOKEN"
	Marketplace     = "OPERATOR_MARKETPLACE"
	globalAppName   = "MinIO Operator"

	// Constants for prometheus annotations
	prometheusPath   = "prometheus.io/path"
	prometheusPort   = "prometheus.io/port"
	prometheusScrape = "prometheus.io/scrape"

	// Constants for DirectPV
	DirectPVMode = "DIRECTPV_MODE"

	// Image versions

	KESImageVersion = "minio/kes:v0.17.6"

	// Constants for common configuration
	MinIOServer        = "OPERATOR_MINIO_SERVER"
	OperatorUIHostname = "OPERATOR_HOSTNAME"
	OperatorUIPort     = "OPERATOR_PORT"
	OperatorUITLSPort  = "OPERATOR_TLS_PORT"

	// Constants for Secure middleware
	SecureAllowedHosts                    = "OPERATOR_SECURE_ALLOWED_HOSTS"
	SecureAllowedHostsAreRegex            = "OPERATOR_SECURE_ALLOWED_HOSTS_ARE_REGEX"
	SecureFrameDeny                       = "OPERATOR_SECURE_FRAME_DENY"
	SecureContentTypeNoSniff              = "OPERATOR_SECURE_CONTENT_TYPE_NO_SNIFF"
	SecureBrowserXSSFilter                = "OPERATOR_SECURE_BROWSER_XSS_FILTER"
	SecureContentSecurityPolicy           = "OPERATOR_SECURE_CONTENT_SECURITY_POLICY"
	SecureContentSecurityPolicyReportOnly = "OPERATOR_SECURE_CONTENT_SECURITY_POLICY_REPORT_ONLY"
	SecureHostsProxyHeaders               = "OPERATOR_SECURE_HOSTS_PROXY_HEADERS"
	SecureSTSSeconds                      = "OPERATOR_SECURE_STS_SECONDS"
	SecureSTSIncludeSubdomains            = "OPERATOR_SECURE_STS_INCLUDE_SUB_DOMAINS"
	SecureSTSPreload                      = "OPERATOR_SECURE_STS_PRELOAD"
	SecureTLSRedirect                     = "OPERATOR_SECURE_TLS_REDIRECT"
	SecureTLSHost                         = "OPERATOR_SECURE_TLS_HOST"
	SecureTLSTemporaryRedirect            = "OPERATOR_SECURE_TLS_TEMPORARY_REDIRECT"
	SecureForceSTSHeader                  = "OPERATOR_SECURE_FORCE_STS_HEADER"
	SecurePublicKey                       = "OPERATOR_SECURE_PUBLIC_KEY"
	SecureReferrerPolicy                  = "OPERATOR_SECURE_REFERRER_POLICY"
	SecureFeaturePolicy                   = "OPERATOR_SECURE_FEATURE_POLICY"
	SecureExpectCTHeader                  = "OPERATOR_SECURE_EXPECT_CT_HEADER"
	SlashSeparator                        = "/"
)
