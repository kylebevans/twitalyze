(this.webpackJsonpoptions=this.webpackJsonpoptions||[]).push([[0],{26:function(t,e,n){"use strict";n.r(e);var o=n(14),r=n(15),i=n(25),s=n(24),a=n(3),c=n(1),d=n.n(c),l=n(16),u=n.n(l),h=n(23),f=(n(47),n(48),{colors:["#1f77b4","#ff7f0e","#2ca02c","#d62728","#9467bd","#8c564b"],enableTooltip:!0,deterministic:!1,fontFamily:"impact",fontSizes:[5,60],fontStyle:"normal",fontWeight:"normal",padding:1,rotations:0,rotationAngles:[0,90],scale:"sqrt",spiral:"archimedean",transitionDuration:1e3}),j=function(t){Object(i.a)(n,t);var e=Object(s.a)(n);function n(t){var r;return Object(o.a)(this,n),(r=e.call(this,t)).state={error:null,isLoaded:!1,wordValues:[]},r}return Object(r.a)(n,[{key:"componentDidMount",value:function(){var t=this;fetch("https://kyle.evans.dev/words").then((function(t){return t.json()})).then((function(e){t.setState({isLoaded:!0,wordValues:e.wordValues})}),(function(e){t.setState({isLoaded:!0,error:e})}))}},{key:"render",value:function(){var t=this.state,e=t.error,n=t.isLoaded,o=t.wordValues;return e?Object(a.jsxs)("div",{children:["Error: ",e.message]}):n?Object(a.jsx)("div",{children:Object(a.jsx)("div",{style:{height:400,width:600},children:Object(a.jsx)(h.a,{options:f,words:o})})}):Object(a.jsx)("div",{children:"Loading..."})}}]),n}(d.a.Component),b=document.getElementById("root");u.a.render(Object(a.jsx)(j,{}),b)},45:function(t,e){}},[[26,1,2]]]);
//# sourceMappingURL=main.3de46c6c.chunk.js.map