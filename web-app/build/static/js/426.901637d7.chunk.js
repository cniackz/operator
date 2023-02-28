"use strict";(self.webpackChunkweb_app=self.webpackChunkweb_app||[]).push([[426],{426:function(e,n,t){t.r(n);var r=t(29439),o=t(1413),a=t(72791),s=t(78687),i=t(57689),c=t(11135),l=t(25787),u=t(23814),d=t(61889),Z=t(45248),f=t(81207),v=t(1261),h=t(5450),p=t(87995),x=t(80184),j=(0,s.$j)((function(e){return{loadingTenant:e.tenants.loadingTenant}}),null);n.default=(0,l.Z)((function(e){return(0,c.Z)((0,o.Z)((0,o.Z)((0,o.Z)((0,o.Z)({},u.OR),u.qg),u.VX),u.Bz))}))(j((function(e){var n=e.classes,t=(0,v.TL)(),o=(0,i.UO)(),c=(0,s.v9)((function(e){return e.tenants.loadingTenant})),l=(0,a.useState)([]),u=(0,r.Z)(l,2),j=u[0],m=u[1],g=(0,a.useState)(!0),b=(0,r.Z)(g,2),w=b[0],T=b[1],k=o.tenantName||"",C=o.tenantNamespace||"";return(0,a.useEffect)((function(){c&&T(!0)}),[c]),(0,a.useEffect)((function(){w&&f.Z.invoke("GET","/api/v1/namespaces/".concat(C,"/tenants/").concat(k,"/events")).then((function(e){for(var n=0;n<e.length;n++){var t=Date.now()/1e3|0;e[n].seen=(0,Z.v1)((t-e[n].last_seen).toString())}m(e),T(!1)})).catch((function(e){t((0,p.Ih)(e)),T(!1)}))}),[w,C,k,t]),(0,x.jsxs)(a.Fragment,{children:[(0,x.jsx)("h1",{className:n.sectionTitle,children:"Events"}),(0,x.jsx)(d.ZP,{item:!0,xs:12,children:(0,x.jsx)(h.Z,{events:j,loading:w})})]})})))},5450:function(e,n,t){var r=t(29439),o=t(72791),a=t(57482),s=t(79836),i=t(53382),c=t(53994),l=t(56890),u=t(35855),d=t(64554),Z=t(56125),f=t(20890),v=t(85172),h=t(95678),p=t(39281),x=t(35527),j=t(80184),m=function(e){var n=e.event,t=o.useState(!1),a=(0,r.Z)(t,2),s=a[0],i=a[1];return(0,j.jsxs)(o.Fragment,{children:[(0,j.jsxs)(u.Z,{sx:{"& > *":{borderBottom:"unset"},cursor:"pointer"},children:[(0,j.jsx)(c.Z,{component:"th",scope:"row",onClick:function(){return i(!s)},children:n.event_type}),(0,j.jsx)(c.Z,{onClick:function(){return i(!s)},children:n.reason}),(0,j.jsx)(c.Z,{onClick:function(){return i(!s)},children:n.seen}),(0,j.jsx)(c.Z,{onClick:function(){return i(!s)},children:n.message.length>=30?"".concat(n.message.slice(0,30),"..."):n.message}),(0,j.jsx)(c.Z,{onClick:function(){return i(!s)},children:s?(0,j.jsx)(h.Z,{}):(0,j.jsx)(v.Z,{})})]}),(0,j.jsx)(u.Z,{children:(0,j.jsx)(c.Z,{style:{paddingBottom:0,paddingTop:0},colSpan:5,children:(0,j.jsx)(Z.Z,{in:s,timeout:"auto",unmountOnExit:!0,children:(0,j.jsx)(d.Z,{sx:{margin:1},children:(0,j.jsx)(f.Z,{style:{background:"#efefef",border:"1px solid #dedede",padding:4,fontSize:14,color:"#666666"},children:n.message})})})})})]})};n.Z=function(e){var n=e.events;return e.loading?(0,j.jsx)(a.Z,{}):(0,j.jsx)(p.Z,{component:x.Z,children:(0,j.jsxs)(s.Z,{"aria-label":"collapsible table",children:[(0,j.jsx)(l.Z,{children:(0,j.jsxs)(u.Z,{children:[(0,j.jsx)(c.Z,{children:"Type"}),(0,j.jsx)(c.Z,{children:"Reason"}),(0,j.jsx)(c.Z,{children:"Age"}),(0,j.jsx)(c.Z,{children:"Message"}),(0,j.jsx)(c.Z,{})]})}),(0,j.jsx)(i.Z,{children:n.map((function(e){return(0,j.jsx)(m,{event:e},"".concat(e.event_type,"-").concat(e.seen))}))})]})})}},85172:function(e,n,t){var r=t(64836);n.Z=void 0;var o=r(t(45649)),a=t(80184),s=(0,o.default)((0,a.jsx)("path",{d:"M7.41 8.59 12 13.17l4.59-4.58L18 10l-6 6-6-6 1.41-1.41z"}),"KeyboardArrowDown");n.Z=s},95678:function(e,n,t){var r=t(64836);n.Z=void 0;var o=r(t(45649)),a=t(80184),s=(0,o.default)((0,a.jsx)("path",{d:"M7.41 15.41 12 10.83l4.59 4.58L18 14l-6-6-6 6z"}),"KeyboardArrowUp");n.Z=s},39281:function(e,n,t){t.d(n,{Z:function(){return p}});var r=t(87462),o=t(63366),a=t(72791),s=t(28182),i=t(94419),c=t(31402),l=t(66934),u=t(75878),d=t(21217);function Z(e){return(0,d.Z)("MuiTableContainer",e)}(0,u.Z)("MuiTableContainer",["root"]);var f=t(80184),v=["className","component"],h=(0,l.ZP)("div",{name:"MuiTableContainer",slot:"Root",overridesResolver:function(e,n){return n.root}})({width:"100%",overflowX:"auto"}),p=a.forwardRef((function(e,n){var t=(0,c.Z)({props:e,name:"MuiTableContainer"}),a=t.className,l=t.component,u=void 0===l?"div":l,d=(0,o.Z)(t,v),p=(0,r.Z)({},t,{component:u}),x=function(e){var n=e.classes;return(0,i.Z)({root:["root"]},Z,n)}(p);return(0,f.jsx)(h,(0,r.Z)({ref:n,as:u,className:(0,s.Z)(x.root,a),ownerState:p},d))}))},56890:function(e,n,t){t.d(n,{Z:function(){return m}});var r=t(87462),o=t(63366),a=t(72791),s=t(28182),i=t(94419),c=t(829),l=t(31402),u=t(66934),d=t(75878),Z=t(21217);function f(e){return(0,Z.Z)("MuiTableHead",e)}(0,d.Z)("MuiTableHead",["root"]);var v=t(80184),h=["className","component"],p=(0,u.ZP)("thead",{name:"MuiTableHead",slot:"Root",overridesResolver:function(e,n){return n.root}})({display:"table-header-group"}),x={variant:"head"},j="thead",m=a.forwardRef((function(e,n){var t=(0,l.Z)({props:e,name:"MuiTableHead"}),a=t.className,u=t.component,d=void 0===u?j:u,Z=(0,o.Z)(t,h),m=(0,r.Z)({},t,{component:d}),g=function(e){var n=e.classes;return(0,i.Z)({root:["root"]},f,n)}(m);return(0,v.jsx)(c.Z.Provider,{value:x,children:(0,v.jsx)(p,(0,r.Z)({as:d,className:(0,s.Z)(g.root,a),ref:n,role:d===j?null:"rowgroup",ownerState:m},Z))})}))}}]);
//# sourceMappingURL=426.901637d7.chunk.js.map