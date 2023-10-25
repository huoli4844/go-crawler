import{g as L,m as O,I as R,J as j,_ as P,r as U,d as z,i as F,c as S,M as _,k as d,l as b,N as X,O as V,Q as J,S as Q,F as q,W as G,P as B,U as K,V as Y}from"./index-4b771267.js";const v=(e,t,l)=>{const n=R(l);return{[`${e.componentCls}-${t}`]:{color:e[`color${l}`],background:e[`color${n}Bg`],borderColor:e[`color${n}Border`],[`&${e.componentCls}-borderless`]:{borderColor:"transparent"}}}},Z=e=>j(e,(t,l)=>{let{textColor:n,lightBorderColor:a,lightColor:o,darkColor:c}=l;return{[`${e.componentCls}-${t}`]:{color:n,background:o,borderColor:a,"&-inverse":{color:e.colorTextLightSolid,background:c,borderColor:c},[`&${e.componentCls}-borderless`]:{borderColor:"transparent"}}}}),ee=e=>{const{paddingXXS:t,lineWidth:l,tagPaddingHorizontal:n,componentCls:a}=e,o=n-l,c=t-l;return{[a]:P(P({},U(e)),{display:"inline-block",height:"auto",marginInlineEnd:e.marginXS,paddingInline:o,fontSize:e.tagFontSize,lineHeight:`${e.tagLineHeight}px`,whiteSpace:"nowrap",background:e.tagDefaultBg,border:`${e.lineWidth}px ${e.lineType} ${e.colorBorder}`,borderRadius:e.borderRadiusSM,opacity:1,transition:`all ${e.motionDurationMid}`,textAlign:"start",[`&${a}-rtl`]:{direction:"rtl"},"&, a, a:hover":{color:e.tagDefaultColor},[`${a}-close-icon`]:{marginInlineStart:c,color:e.colorTextDescription,fontSize:e.tagIconSize,cursor:"pointer",transition:`all ${e.motionDurationMid}`,"&:hover":{color:e.colorTextHeading}},[`&${a}-has-color`]:{borderColor:"transparent",[`&, a, a:hover, ${e.iconCls}-close, ${e.iconCls}-close:hover`]:{color:e.colorTextLightSolid}},"&-checkable":{backgroundColor:"transparent",borderColor:"transparent",cursor:"pointer",[`&:not(${a}-checkable-checked):hover`]:{color:e.colorPrimary,backgroundColor:e.colorFillSecondary},"&:active, &-checked":{color:e.colorTextLightSolid},"&-checked":{backgroundColor:e.colorPrimary,"&:hover":{backgroundColor:e.colorPrimaryHover}},"&:active":{backgroundColor:e.colorPrimaryActive}},"&-hidden":{display:"none"},[`> ${e.iconCls} + span, > span + ${e.iconCls}`]:{marginInlineStart:o}}),[`${a}-borderless`]:{borderColor:"transparent",background:e.tagBorderlessBg}}},H=L("Tag",e=>{const{fontSize:t,lineHeight:l,lineWidth:n,fontSizeIcon:a}=e,o=Math.round(t*l),c=e.fontSizeSM,g=o-n*2,C=e.colorFillAlter,i=e.colorText,r=O(e,{tagFontSize:c,tagLineHeight:g,tagDefaultBg:C,tagDefaultColor:i,tagIconSize:a-2*n,tagPaddingHorizontal:8,tagBorderlessBg:e.colorFillTertiary});return[ee(r),Z(r),v(r,"success","Success"),v(r,"processing","Info"),v(r,"error","Error"),v(r,"warning","Warning")]}),oe=()=>({prefixCls:String,checked:{type:Boolean,default:void 0},onChange:{type:Function},onClick:{type:Function},"onUpdate:checked":Function}),le=z({compatConfig:{MODE:3},name:"ACheckableTag",inheritAttrs:!1,props:oe(),setup(e,t){let{slots:l,emit:n,attrs:a}=t;const{prefixCls:o}=F("tag",e),[c,g]=H(o),C=r=>{const{checked:u}=e;n("update:checked",!u),n("change",!u),n("click",r)},i=S(()=>_(o.value,g.value,{[`${o.value}-checkable`]:!0,[`${o.value}-checkable-checked`]:e.checked}));return()=>{var r;return c(d("span",b(b({},a),{},{class:[i.value,a.class],onClick:C}),[(r=l.default)===null||r===void 0?void 0:r.call(l)]))}}}),m=le,ne=()=>({prefixCls:String,color:{type:String},closable:{type:Boolean,default:!1},closeIcon:B.any,visible:{type:Boolean,default:void 0},onClose:{type:Function},onClick:K(),"onUpdate:visible":Function,icon:B.any,bordered:{type:Boolean,default:!0}}),h=z({compatConfig:{MODE:3},name:"ATag",inheritAttrs:!1,props:ne(),slots:Object,setup(e,t){let{slots:l,emit:n,attrs:a}=t;const{prefixCls:o,direction:c}=F("tag",e),[g,C]=H(o),i=X(!0);V(()=>{e.visible!==void 0&&(i.value=e.visible)});const r=s=>{s.stopPropagation(),n("update:visible",!1),n("close",s),!s.defaultPrevented&&e.visible===void 0&&(i.value=!1)},u=S(()=>J(e.color)||Q(e.color)),D=S(()=>_(o.value,C.value,{[`${o.value}-${e.color}`]:u.value,[`${o.value}-has-color`]:e.color&&!u.value,[`${o.value}-hidden`]:!i.value,[`${o.value}-rtl`]:c.value==="rtl",[`${o.value}-borderless`]:!e.bordered})),M=s=>{n("click",s)};return()=>{var s,p,f;const{icon:k=(s=l.icon)===null||s===void 0?void 0:s.call(l),color:$,closeIcon:y=(p=l.closeIcon)===null||p===void 0?void 0:p.call(l),closable:w=!1}=e,A=()=>w?y?d("span",{class:`${o.value}-close-icon`,onClick:r},[y]):d(Y,{class:`${o.value}-close-icon`,onClick:r},null):null,N={backgroundColor:$&&!u.value?$:void 0},T=k||null,I=(f=l.default)===null||f===void 0?void 0:f.call(l),W=T?d(q,null,[T,d("span",null,[I])]):I,E=e.onClick!==void 0,x=d("span",b(b({},a),{},{onClick:M,class:[D.value,a.class],style:[N,a.style]}),[W,A()]);return g(E?d(G,null,{default:()=>[x]}):x)}}});h.CheckableTag=m;h.install=function(e){return e.component(h.name,h),e.component(m.name,m),e};const re=h;export{re as _};
