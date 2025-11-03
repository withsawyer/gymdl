package video

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core"
	"github.com/nichuanfang/gymdl/processor"
	"github.com/nichuanfang/gymdl/utils"
	"github.com/playwright-community/playwright-go"
	"github.com/withsawyer/gopher-tools/datetime"
)

const (
	// videoURLPattern 视频URL匹配模式
	videoURLPattern = `/video/(\d+)`
	// videoIDPattern 视频ID匹配模式
	videoIDPattern = `video_id=([a-zA-Z0-9]+)`
	// jsonPattern JSON数据匹配模式
	jsonPattern = `({.*?"errors":\s*null\s*})`
	// scriptPattern 脚本标签匹配模式
	scriptPattern = `<script[^>]*>(.*?)</script>`
	// titlePattern 标题匹配模式
	titlePattern = `<title>(.*?)</title>`
)

// Platform 表示平台类型
type Platform string

// DouYinProcessor 抖音视频处理器，实现视频下载功能
type DouYinProcessor struct {
	cfg       *config.Config
	tempDir   string
	videos    []*VideoInfo
	videoInfo *VideoInfo
}

// Init 初始化抖音处理器
func (p *DouYinProcessor) Init(cfg *config.Config) {
	p.cfg = cfg
	p.videos = make([]*VideoInfo, 0)
	p.tempDir = processor.BuildOutputDir(DouyinTempDir)
	p.videoInfo = &VideoInfo{}
}

// Name 返回处理器名称
func (p *DouYinProcessor) Name() processor.LinkType {
	return processor.LinkDouyin
}

// Videos 返回已下载的视频信息列表
func (p *DouYinProcessor) Videos() []*VideoInfo {
	return p.videos
}

// Download 下载抖音视频
func (p *DouYinProcessor) Download(link string) error {
	err := p.method1(link)
	if err != nil {
		return fmt.Errorf("下载抖音视频失败: %v", err)
	}

	return nil
}

func (p *DouYinProcessor) method1(link string) error {
	//// 初始化 Playwright 和浏览器
	//ctx, page, pw, err := p.initPlaywrightAndBrowser()
	//if err != nil {
	//	return err
	//}
	//defer func() {
	//	page.Close()
	//	ctx.Close()
	//	pw.Stop()
	//}()
	//
	//// 加载 cookies
	//if err = p.loadCookies(ctx); err != nil {
	//	utils.InfoWithFormat("加载 cookies 失败: %v", err)
	//}
	//
	//// 提取视频ID
	//videoID, err := p._extractVideoID(page, link)
	//if err != nil {
	//	return err
	//}
	//utils.InfoWithFormat("提取视频ID成功: %s", videoID)
	//// 提取视频内容和URL
	//html, err := page.Content()
	//if err != nil {
	//	return fmt.Errorf("获取页面内容失败: %v", err)
	//}
	html := `<!DOCTYPE html><html style="font-size: 50px; --vh: 961px;"><head>
    <script data-sdk-glue-in="pre-handler" nonce="">    !function(){var e="1243",t="16720";function n(){ try{ var n="gfkadpd",o=e+","+t,i=function(e){for(var t=document.cookie.split(";"),n=0;n<t.length;n++){var o=t[n].trim();if(o.startsWith(e+"="))return o.substring(e.length+1)}return null}(n);if(i){if(-1!=i.indexOf(o))return; o+="|"+i}document.cookie =n+"="+o+";expires="+new Date((new Date).getTime()+3*24*60*60*1e3).toUTCString()+"; path=/; SameSite=None; Secure;"}catch(e){}}window.__t_version='AOSdCc';var o=function(n,o,i,a){if(Math.ceil(100*Math.random())<=100*o){var r={ev_type:"batch",list:[{ev_type:"custom",payload:{name:"sdk_glue_load",type:"event",metrics:{},categories:{sdk_glue_load_status:n,sdk_glue_load_err_src:i,payload_bdms_aid:e,payload_bdms_page_id:t, duration:a}},common:{context:{ctx_bdms_aid:e,ctx_bdms_page_id:t},bid:"web_bdms_cn",pid:window.location.pathname,view_id:"/_1",user_id:"",session_id:"0-a-1-2-c",release:"",env:"production", url:window.location.href,timestamp: +new Date,sdk_version:"1.6.1",sdk_name:"SDK_SLARDAR_WEB"}}]},d=new XMLHttpRequest;d.open("POST","https://mon.zijieapi.com/monitor_browser/collect/batch/?biz_id=web_bdms_cn",!0),d.setRequestHeader("Content-type","application/json"),d.send(JSON.stringify(r))}};!function(){try{n(),document.cookie="wdglgl=; expires=Mon, 20 Sep 2010 00:00:00 UTC; path=/;",o("before_load",.001,"","")}catch(e){}var e=performance.now();window.addEventListener("error",(function(t){try{var n=t.target||t.srcElement;if(n instanceof HTMLElement&&"SCRIPT"==n.nodeName)if(-1!=(n.src||"").indexOf("sdk-glue")){var i=(performance.now()-e).toFixed(2);a=i,document.cookie="wdglgl=; expires=Mon, 20 Sep 2010 00:00:00 UTC; path=/;",document.cookie="wdglgl="+a+"; expires="+new Date((new Date).getTime()+3*24*60*60*1e3).toUTCString()+"; path=/; SameSite=None; Secure;",o("load_error",1,n.src,i)}}catch(e){}var a}),!0),window.__glue_t=+new Date}()}();</script><script nonce="" data-dt="RfDZYq" src="https://lf-c-flwb.bytetos.com/obj/rc-client-security/web/glue/1.0.0.65/sdk-glue.js"></script><script nonce="">;(function (){    var sdkInfo = {      csrf: {        init: function (options) {window.secsdk.csrf.setOptions(options)},        isLoaded: function () { return !!window.secsdk },        srcList: ["https://lf1-cdn-tos.bytegoofy.com/obj/goofy/secsdk/secsdk-lastest.umd.js","https://lf3-cdn-tos.bytegoofy.com/obj/goofy/secsdk/secsdk-lastest.umd.js","https://lf6-cdn-tos.bytegoofy.com/obj/goofy/secsdk/secsdk-lastest.umd.js"],                      },      bdms: {        init: function (options) {window.bdms.init(options)},        isLoaded: function () { return !!window.bdms },        srcList: ["https://lf-c-flwb.bytetos.com/obj/rc-client-security/web/dn/1.0.2.0-alpha.5/bdms.NVI6Q9.js"],                      },      verifyCenter: {        init: function (options) {window.TTGCaptcha.init(options)},        isLoaded: function () { return !!window.TTGCaptcha },        srcList: ["https://lf-rc1.yhgfb-cn-static.com/obj/rc-verifycenter/sec_sdk_build/4.0.10/captcha/index.js","https://lf-rc2.yhgfb-cn-static.com/obj/rc-verifycenter/sec_sdk_build/4.0.10/captcha/index.js"],                      }    };    if (window._SdkGlueInit) {      window._SdkGlueInit({        bdms: {"_0xffe":"uee4cI","aid":1243,"ddrt":3,"pageId":16720,"paths":{"include":["/aweme/v1","/aweme/v2","/web/api"]}},self: {aid:1243,pageId:16720,}      }, sdkInfo);    }    })()</script><script no-entry="true" src="https://lf-c-flwb.bytetos.com/obj/rc-client-security/web/dn/1.0.2.0-alpha.5/bdms.NVI6Q9.js"></script>
    <meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1,shrink-to-fit=no,viewport-fit=cover,minimum-scale=1,maximum-scale=1,user-scalable=no"><meta http-equiv="x-ua-compatible" content="ie=edge"><meta name="renderer" content="webkit"><meta name="layoutmode" content="standard"><meta name="imagemode" content="force"><meta name="wap-font-scale" content="no"><meta name="format-detection" content="telephone=no"><link href="//lf-douyin-mobile.bytecdn.com/obj/growth-douyin-share/growth/douyin_ug/static/css/video.09525075.css" rel="stylesheet" crossorigin="anonymous"><script nonce="" defer="defer">!function(){"use strict";var e,t,n,a,r,i,o,d,c={},u={};function f(e){var t=u[e];if(void 0!==t)return t.exports;var n=u[e]={id:e,loaded:!1,exports:{}};return c[e].call(n.exports,n,n.exports,f),n.loaded=!0,n.exports}if(f.m=c,f.n=function(e){var t=e&&e.__esModule?function(){return e.default}:function(){return e};return f.d(t,{a:t}),t},t=Object.getPrototypeOf?function(e){return Object.getPrototypeOf(e)}:function(e){return e.__proto__},f.t=function(n,a){if(1&a&&(n=this(n)),8&a||"object"==typeof n&&n&&(4&a&&n.__esModule||16&a&&"function"==typeof n.then))return n;var r=Object.create(null);f.r(r);var i={};e=e||[null,t({}),t([]),t(t)];for(var o=2&a&&n;"object"==typeof o&&!~e.indexOf(o);o=t(o))Object.getOwnPropertyNames(o).forEach((function(e){i[e]=function(){return n[e]}}));return i.default=function(){return n},f.d(r,i),r},f.d=function(e,t){for(var n in t)f.o(t,n)&&!f.o(e,n)&&Object.defineProperty(e,n,{enumerable:!0,get:t[n]})},f.f={},f.e=function(e){return Promise.all(Object.keys(f.f).reduce((function(t,n){return f.f[n](e,t),t}),[]))},f.u=function(e){return"9737"===e?"static/js/9737.5471907a.js":"6430"===e?"static/js/6430.c4525e5d.js":"static/js/async/"+({1016:"mix_(id)/page",1519:"sticker_page",1743:"mv_template_(id)/page",2005:"VideoMEssage",2191:"challenge_page",2513:"user_(id)/page",2616:"music_page",276:"delightdesign_h5",3981:"note_page",4112:"playlet_page",4628:"video_page",4647:"mix_page",4748:"user_page",5102:"sticker_(id)/page",54:"challenge_(id)/page",5404:"playlet_(id)/page",575:"jx-video_(id)/page",6576:"slides_(id)/page",6684:"note_(id)/page",6764:"slides_page",6818:"group_(id)/page",7263:"musicplaylist_page",7590:"mv_template_page",7989:"billboard_page",8317:"music_(id)/page",8318:"byted-lina-douyin-search-comment-board",8365:"jx-video_page",9166:"group_page",944:"trends_page",9706:"musicplaylist_(id)/page",9826:"AdaptLoginHeader",9973:"video_(id)/page"}[e]||e)+"."+{1006:"fc12eee8",1016:"f18aac5c",1519:"3b68d8ce",1743:"e7a34a10",1865:"8eb063f5",2005:"a6d76d21",2191:"987fbaa8",2513:"f512f80e",2616:"0e83d259",276:"74cdb36c",3384:"f0fc9bd9",3633:"295c106d",3683:"532a8402",3981:"ce9eae02",4002:"cf735d2e",4097:"16d21948",4112:"74993d69",4628:"0512652a",4647:"a04ace0e",4748:"fa6fe594",5102:"30585b8b",5118:"1442fa2d",54:"ad2a6e66",5404:"8eee2110",575:"b2db73dc",6576:"7f379cf9",6684:"390fc151",6764:"c80a3354",6818:"84f02879",7263:"9d802fb2",7431:"689d5bb9",7590:"f5642752",7989:"c4f7ee21",8317:"bdcbe90e",8318:"e19b25a7",8365:"1ab03031",9166:"d9e8e47a",944:"abcf8670",9543:"83ada03c",9706:"55ea72fb",9826:"d113e854",9973:"15670c72"}[e]+".js"},f.miniCssF=function(e){return"static/css/async/"+{1016:"mix_(id)/page",1519:"sticker_page",1743:"mv_template_(id)/page",2005:"VideoMEssage",2191:"challenge_page",2513:"user_(id)/page",2616:"music_page",276:"delightdesign_h5",3981:"note_page",4112:"playlet_page",4628:"video_page",4647:"mix_page",4748:"user_page",5102:"sticker_(id)/page",54:"challenge_(id)/page",5404:"playlet_(id)/page",575:"jx-video_(id)/page",6576:"slides_(id)/page",6684:"note_(id)/page",6764:"slides_page",6818:"group_(id)/page",7263:"musicplaylist_page",7590:"mv_template_page",7989:"billboard_page",8317:"music_(id)/page",8318:"byted-lina-douyin-search-comment-board",8365:"jx-video_page",9166:"group_page",944:"trends_page",9706:"musicplaylist_(id)/page",9826:"AdaptLoginHeader",9973:"video_(id)/page"}[e]+"."+{1016:"db291758",1519:"2ec23390",1743:"2d880233",2005:"e1cf9bc9",2191:"2ec23390",2513:"5afcb9c4",2616:"2ec23390",276:"81f9cfdb",3981:"2ec23390",4112:"2ec23390",4628:"2ec23390",4647:"2ec23390",4748:"2ec23390",5102:"3d32d68d",54:"7f7f670f",5404:"249a1684",575:"ca54cc38",6576:"9845db1a",6684:"255ce86f",6764:"2ec23390",6818:"32c20c94",7263:"2ec23390",7590:"2ec23390",7989:"860592c8",8317:"60636dff",8318:"0dc804dc",8365:"2ec23390",9166:"2ec23390",944:"a492d002",9706:"8a56833b",9826:"09bdaf0a",9973:"7f707a39"}[e]+".css"},f.h=function(){return"4ee028efd66a0421"},f.g=function(){if("object"==typeof globalThis)return globalThis;try{return this||Function("return this")()}catch(e){if("object"==typeof window)return window}}(),f.o=function(e,t){return Object.prototype.hasOwnProperty.call(e,t)},n={},a="douyin_ug_edenx:",f.l=function(e,t,r,i){if(n[e])n[e].push(t);else{if(void 0!==r)for(var o,d,c=document.getElementsByTagName("script"),u=0;u<c.length;u++){var l=c[u];if(l.getAttribute("src")==e||l.getAttribute("data-webpack")==a+r){o=l;break}}o||(d=!0,(o=document.createElement("script")).charset="utf-8",o.timeout=120,f.nc&&o.setAttribute("nonce",f.nc),o.setAttribute("data-webpack",a+r),o.src=e,0!==o.src.indexOf(window.location.origin+"/")&&(o.crossOrigin="anonymous")),n[e]=[t];var s=function(t,a){o.onerror=o.onload=null,clearTimeout(p);var r=n[e];if(delete n[e],o.parentNode&&o.parentNode.removeChild(o),r&&r.forEach((function(e){return e(a)})),t)return t(a)},p=setTimeout(s.bind(null,void 0,{type:"timeout",target:o}),12e4);o.onerror=s.bind(null,o.onerror),o.onload=s.bind(null,o.onload),d&&document.head.appendChild(o)}},f.r=function(e){"undefined"!=typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(e,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(e,"__esModule",{value:!0})},f.nmd=function(e){return e.paths=[],e.children||(e.children=[]),e},r=[],f.O=function(e,t,n,a){if(!t){var i=1/0;for(u=0;u<r.length;u++){t=r[u][0],n=r[u][1],a=r[u][2];for(var o=!0,d=0;d<t.length;d++)(!1&a||i>=a)&&Object.keys(f.O).every((function(e){return f.O[e](t[d])}))?t.splice(d--,1):(o=!1,a<i&&(i=a));if(o){r.splice(u--,1);var c=n();void 0!==c&&(e=c)}}return e}a=a||0;for(var u=r.length;u>0&&r[u-1][2]>a;u--)r[u]=r[u-1];r[u]=[t,n,a]},f.p="//lf-douyin-mobile.bytecdn.com/obj/growth-douyin-share/growth/douyin_ug/",f.rv=function(){return"1.3.12"},"undefined"!=typeof document){var l={8242:0};f.f.miniCss=function(e,t){l[e]?t.push(l[e]):0!==l[e]&&{1016:1,2616:1,6576:1,4628:1,944:1,575:1,1519:1,9826:1,4112:1,9973:1,7590:1,9706:1,6818:1,7989:1,8318:1,7263:1,1743:1,2513:1,3981:1,4748:1,8365:1,8317:1,2191:1,9166:1,2005:1,276:1,6764:1,5102:1,5404:1,6684:1,4647:1,54:1}[e]&&t.push(l[e]=new Promise((function(t,n){var a=f.miniCssF(e),r=f.p+a;if(function(e,t){for(var n=document.getElementsByTagName("link"),a=0;a<n.length;a++){var r=(o=n[a]).getAttribute("data-href")||o.getAttribute("href");if("stylesheet"===o.rel&&(r===e||r===t))return o}var i=document.getElementsByTagName("style");for(a=0;a<i.length;a++){var o;if((r=(o=i[a]).getAttribute("data-href"))===e||r===t)return o}}(a,r))return t();!function(e,t,n,a,r){var i=document.createElement("link");i.rel="stylesheet",i.type="text/css",f.nc&&(i.nonce=f.nc),i.onerror=i.onload=function(n){if(i.onerror=i.onload=null,"load"===n.type)a();else{var o=n&&("load"===n.type?"missing":n.type),d=n&&n.target&&n.target.href||t,c=Error("Loading CSS chunk "+e+" failed.\\n("+d+")");c.code="CSS_CHUNK_LOAD_FAILED",c.type=o,c.request=d,i.parentNode&&i.parentNode.removeChild(i),r(c)}},i.href=t,0!==i.href.indexOf(window.location.origin+"/")&&(i.crossOrigin="anonymous"),n?n.parentNode.insertBefore(i,n.nextSibling):document.head.appendChild(i)}(e,r,null,t,n)})).then((function(){l[e]=0}),(function(t){throw delete l[e],t})))}}i={8242:0},f.f.j=function(e,t){var n=f.o(i,e)?i[e]:void 0;if(0!==n)if(n)t.push(n[2]);else if(8242!=e){var a=new Promise((function(t,a){n=i[e]=[t,a]}));t.push(n[2]=a);var r=f.p+f.u(e),o=Error();f.l(r,(function(t){if(f.o(i,e)&&(0!==(n=i[e])&&(i[e]=void 0),n)){var a=t&&("load"===t.type?"missing":t.type),r=t&&t.target&&t.target.src;o.message="Loading chunk "+e+" failed.\n("+a+": "+r+")",o.name="ChunkLoadError",o.type=a,o.request=r,n[1](o)}}),"chunk-"+e,e)}else i[e]=0},f.O.j=function(e){return 0===i[e]},o=function(e,t){var n,a,r=t[0],o=t[1],d=t[2],c=0;if(r.some((function(e){return 0!==i[e]}))){for(n in o)f.o(o,n)&&(f.m[n]=o[n]);if(d)var u=d(f)}for(e&&e(t);c<r.length;c++)a=r[c],f.o(i,a)&&i[a]&&i[a][0](),i[a]=0;return f.O(u)},(d=self.__LOADABLE_LOADED_CHUNKS__=self.__LOADABLE_LOADED_CHUNKS__||[]).forEach(o.bind(null,0)),d.push=o.bind(null,d.push.bind(d)),f.ruid="bundler=rspack@1.3.12"}()</script><script nonce="">
    ;(function(){
        window._MODERNJS_ROUTE_MANIFEST = {"routeAssets":{"video":{"chunkIds":["8242","3361","2072","2126","2118","5817","5283","9737","4101","8734"],"assets":["static/js/lib-react.8c02d96d.js","static/js/lib-axios.251165c1.js","static/js/lib-polyfill.5f7f1705.js","static/js/lib-router.1e0a06eb.js","static/js/argus-builder-strategy.45d542af.js","static/js/5283.f72d3f68.js","static/js/9737.5471907a.js","static/js/4101.d275a5ca.js","static/css/video.09525075.css","static/js/video.4ee51a9d.js"],"referenceCssAssets":["static/css/video.09525075.css"]},"video_(id)/page":{"chunkIds":["7431","4097","5118","1006","3384","9973"],"assets":["static/js/async/7431.689d5bb9.js","static/js/async/4097.16d21948.js","static/js/async/5118.1442fa2d.js","static/js/async/1006.fc12eee8.js","static/js/async/3384.f0fc9bd9.js","static/css/async/video_(id)/page.7f707a39.css","static/js/async/video_(id)/page.15670c72.js"],"referenceCssAssets":["static/css/async/video_(id)/page.7f707a39.css"]},"video_page":{"chunkIds":["4628"],"assets":["static/css/async/video_page.2ec23390.css","static/js/async/video_page.0512652a.js"],"referenceCssAssets":["static/css/async/video_page.2ec23390.css"]}}};
    })();
</script><script nonce="" defer="defer" src="//lf-douyin-mobile.bytecdn.com/obj/growth-douyin-share/growth/douyin_ug/static/js/lib-react.8c02d96d.js" crossorigin="anonymous"></script><script nonce="" defer="defer" src="//lf-douyin-mobile.bytecdn.com/obj/growth-douyin-share/growth/douyin_ug/static/js/lib-axios.251165c1.js" crossorigin="anonymous"></script><script nonce="" defer="defer" src="//lf-douyin-mobile.bytecdn.com/obj/growth-douyin-share/growth/douyin_ug/static/js/lib-polyfill.5f7f1705.js" crossorigin="anonymous"></script><script nonce="" defer="defer" src="//lf-douyin-mobile.bytecdn.com/obj/growth-douyin-share/growth/douyin_ug/static/js/lib-router.1e0a06eb.js" crossorigin="anonymous"></script><script nonce="" defer="defer" src="//lf-douyin-mobile.bytecdn.com/obj/growth-douyin-share/growth/douyin_ug/static/js/argus-builder-strategy.45d542af.js" crossorigin="anonymous"></script><script nonce="" defer="defer" src="//lf-douyin-mobile.bytecdn.com/obj/growth-douyin-share/growth/douyin_ug/static/js/5283.f72d3f68.js" crossorigin="anonymous"></script><script nonce="" defer="defer" src="//lf-douyin-mobile.bytecdn.com/obj/growth-douyin-share/growth/douyin_ug/static/js/9737.5471907a.js" crossorigin="anonymous"></script><script nonce="" defer="defer" src="//lf-douyin-mobile.bytecdn.com/obj/growth-douyin-share/growth/douyin_ug/static/js/4101.d275a5ca.js" crossorigin="anonymous"></script><script nonce="" defer="defer" src="//lf-douyin-mobile.bytecdn.com/obj/growth-douyin-share/growth/douyin_ug/static/js/video.4ee51a9d.js" crossorigin="anonymous"></script><title data-react-helmet="true">恐怖小狗《The Puppy》一命速通 #恐怖游戏 #steam - 抖音</title><link rel="dns-prefetch" href="//lf-douyin-mobile.bytecdn.com"><link rel="dns-prefetch" href="//lf-cdn-tos.bytescm.com"><link rel="dns-prefetch" href="//lf3-cdn-tos.bytescm.com"><link rel="dns-prefetch" href="//lf3-short.ibytedapm.com"><link rel="dns-prefetch" href="//lf3-short.bytegoofy.com"><link rel="dns-prefetch" href="//lf-c-flwb.bytetos.com"><link rel="dns-prefetch" href="//mon.snssdk.com"><link rel="dns-prefetch" href="//mcs.snssdk.com"><link rel="dns-prefetch" href="//sf1-cdn-tos.douyinstatic.com"><meta charset="utf-8"><meta http-equiv="Cache-Control" content="no-transform"><meta http-equiv="Cache-Control" content="no-siteapp"><meta name="applicable-device" content="mobile"><meta name="viewport" content="width=device-width,initial-scale=1"><meta name="baidu-site-verification" content="szjdG38sKy"><link rel="shortcut icon" href="https://sf1-cdn-tos.douyinstatic.com/obj/eden-cn/kpchkeh7upepld/fe_app_new/favicon_v2.ico" type="image/x-icon"><link rel="apple-touch-icon-precomposed" href="https://sf1-cdn-tos.douyinstatic.com/obj/eden-cn/kpchkeh7upepld/fe_app_new/logo_launcher_v2.png"><link rel="preload" crossorigin="anonymous" as="script" nonce="" href="https://lf3-cdn-tos.bytescm.com/obj/static/log-sdk/collect/5.0/collect.js"><link rel="preload" crossorigin="anonymous" as="script" nonce="" href="https://lf3-short.ibytedapm.com/slardar/fe/sdk-web/browser.cn.js"><link rel="preload" crossorigin="anonymous" as="script" nonce="" href="https://lf-security.bytegoofy.com/obj/security-secsdk/runtime.js"><meta name="apple-mobile-web-app-capable" content="yes"><meta name="apple-mobile-web-app-status-bar-style" content="default"><meta name="screen-orientation" content="portrait"><meta name="format-detection" content="telephone=no"><meta name="x5-orientation" content="portrait"><style>a,abbr,acronym,address,applet,article,aside,audio,b,big,blockquote,body,canvas,caption,center,cite,code,dd,del,details,dfn,div,dl,dt,em,embed,fieldset,figcaption,figure,footer,form,h1,h2,h3,h4,h5,h6,header,hgroup,html,i,iframe,img,ins,kbd,label,legend,li,mark,menu,nav,object,ol,output,p,pre,q,ruby,s,samp,section,small,span,strike,strong,sub,summary,sup,table,tbody,td,tfoot,th,thead,time,tr,tt,u,ul,var,video{margin:0;padding:0;border:0;font-size:100%;font:inherit;vertical-align:baseline}article,aside,details,figcaption,figure,footer,header,hgroup,menu,nav,section{display:block}body{line-height:1}ol,ul{list-style:none}blockquote,q{quotes:none}blockquote:after,blockquote:before,q:after,q:before{content:&quot}table{border-collapse:collapse;border-spacing:0}*{box-sizing:border-box}html{font-family:PingFang SC,system ui,-apple-system,BlinkMacSystemFont,Helvetica Neue,Helvetica,STHeiTi,sans-serif;-webkit-font-smoothing:antialiased;-webkit-tap-highlight-color:transparent;text-rendering:optimizeLegibility}</style><script nonce="" src="https://lf-security.bytegoofy.com/obj/security-secsdk/runtime.js" project-id="207" x-nonce="argus-csp-token"></script><script src="https://lf-security.bytegoofy.com/obj/security-secsdk/runtime-stable.js" nonce="" project-id="207"></script><script nonce="" src="https://lf-security.bytegoofy.com/obj/security-secsdk/config_207.js"></script>
    <script nonce="" src="https://lf-security.bytegoofy.com/obj/security-secsdk/project_207.js"></script>
    <script nonce="" src="https://lf-security.bytegoofy.com/obj/security-secsdk/strategy_207.js?v=1"></script>

    <script nonce="">!function(){function e(e){return this.config=e,this}"undefined"!=typeof window&&(e.prototype={reset:function(){var e=this.config&&this.config.baseline&&!isNaN(this.config.baseline)?this.config.baseline:750,t=document.documentElement.clientWidth>750?750:document.documentElement.clientWidth,n=t/e*100;t>=750?(document.documentElement.style.fontSize="40px",document.documentElement.style.setProperty("width","375px"),document.documentElement.style.setProperty("margin","0 auto"),document.documentElement.style.setProperty("background-color","#f4f4f4")):(document.documentElement.style.fontSize="".concat(n,"px"),document.documentElement.style.removeProperty("width"),document.documentElement.style.removeProperty("margin"),document.documentElement.style.removeProperty("background-color"))}},window.Adapter=new e(window.ADAPTER_CONF||{}),window.Adapter.reset(),window.addEventListener("load",(function(){window.Adapter.reset()}),!0),window.addEventListener("pageshow",(function(){var e=window.navigator.userAgent;/android/i.test(e)&&window.Adapter.reset()}),!0),window.addEventListener("resize",(function(){window.Adapter.reset()}),!0))}()</script><script nonce="">!function(){function n(n){document.documentElement.style.setProperty("--vh",window.innerHeight+"px")}"undefined"!=typeof window&&(n(),window.addEventListener("load",(function(){n()}),!0),window.addEventListener("resize",(function(){n()}),!0))}()</script><link crossorigin="anonymous" href="//lf-douyin-mobile.bytecdn.com/obj/growth-douyin-share/growth/douyin_ug/static/css/async/video_(id)/page.7f707a39.css" rel="stylesheet">  <link data-react-helmet="true" rel="canonical" href="https://www.douyin.com/video/7535852404919586100">
    <meta data-react-helmet="true" name="description" content="恐怖小狗《The Puppy》一命速通 #恐怖游戏 #steam - 式子大王于20250807发布在抖音，已经收获了7007个喜欢，来抖音，记录美好生活！"><meta data-react-helmet="true" name="keywords" content="抖音,抖音短视频,抖音官网"><meta data-react-helmet="true" http-equiv="mobile-agent" content="format=html5;url=https://m.douyin.com/share/video/7535852404919586100"><meta data-react-helmet="true" name="mobile-agent" content="format=html5;url=https://m.douyin.com/share/video/7535852404919586100">
    <script src="https://lf3-short.ibytedapm.com/slardar/fe/sdk-web/browser.cn.js?bid=aweme_share&amp;globalName=slardar" crossorigin="anonymous"></script></head><body><div id="root"><div class="container"><div class="video-container horizontal-video"><img src="https://p3-sign.douyinpic.com/tos-cn-i-dy/dd954010955847ea8f5493d60d984341~tplv-dy-resize-walign-adapt-aq:540:q75.webp?lk3s=138a59ce&amp;x-expires=1763388000&amp;x-signature=x6p6PYRGDbgA5TziiC6TBWT%2FJDs%3D&amp;from=327834062&amp;s=PackSourceEnum_DOUYIN_REFLOW&amp;se=false&amp;sc=cover&amp;biz_tag=aweme_video&amp;l=202511032212173D058217B6CB452C34B1" class="poster" loading="eager"><div class="video-msg-container"></div></div></div></div>
<div id="douyin_reflow_abtest" ssrabtest="{&quot;reflow_video_page_optimize&quot;:&quot;2&quot;,&quot;search_video_mark_abtest&quot;:&quot;0&quot;,&quot;share_page_dark_mode_adaptation&quot;:&quot;1&quot;,&quot;reflow_page_open_app_optimize&quot;:&quot;4&quot;,&quot;reflow_page_auto_open_optimize&quot;:&quot;1&quot;,&quot;reflow_to_featured_app&quot;:&quot;1&quot;,&quot;select_pool_data&quot;:{&quot;use_new_select_scope&quot;:0},&quot;jx_label_display_position_list&quot;:[&quot;main&quot;,&quot;related_card&quot;],&quot;reflow_page_615883_optimaze&quot;:false}"></div>
<div id="douyin_reflow_location" location="{&quot;isMainland&quot;:true,&quot;country&quot;:&quot;cn&quot;}"></div>
<div id="douyin_reflow_token" xsstoken="ed17286642009416be9e1edd3344be9c59b9a613d90bf0553909341dafa1145151e8fcbf428b84c1af9431a28376abc6"></div>
<div id="douyin_reflow_tcc" tccconfig="{&quot;block_path_list&quot;:[],&quot;delayed_open_download_millisecond&quot;:1000,&quot;forbid_copy_board_list&quot;:{&quot;ios_os_version&quot;:&quot;16&quot;,&quot;scene_from&quot;:[&quot;click_schema_hwqs_cloud&quot;,&quot;qs_os_list&quot;,&quot;douyin_h5_flow&quot;,&quot;douyin_h5_search&quot;,&quot;ug_page_delivery&quot;],&quot;from_aid&quot;:[&quot;615883&quot;,&quot;2329&quot;],&quot;app&quot;:[&quot;douyin_search&quot;],&quot;app_name_android&quot;:[],&quot;app_name_ios&quot;:[],&quot;forbid_launch_method&quot;:&quot;1&quot;},&quot;forbid_copy_board_operation&quot;:0,&quot;forbid_schema_call_browser_app_list&quot;:[&quot;firefox&quot;],&quot;forbid_schema_call_scene_from_list&quot;:[&quot;click_schema_hwqs_cloud&quot;,&quot;qs_os_list&quot;,&quot;douyin_h5_flow&quot;,&quot;douyin_h5_search&quot;],&quot;forbid_ui_open_change_list&quot;:{&quot;scene_from&quot;:[&quot;click_schema_hwqs_cloud&quot;,&quot;qs_os_list&quot;,&quot;douyin_h5_flow&quot;,&quot;douyin_h5_search&quot;,&quot;ug_page_delivery&quot;],&quot;from_aid&quot;:[&quot;615883&quot;],&quot;app&quot;:[&quot;douyin_search&quot;],&quot;app_name_android&quot;:[],&quot;app_name_ios&quot;:[],&quot;show_scroll_toast&quot;:[]},&quot;forbid_ui_optimize_browser_app_list&quot;:[&quot;baidu&quot;],&quot;forbid_ui_optimize_scene_from_list&quot;:[&quot;click_schema_hwqs_cloud&quot;,&quot;qs_os_list&quot;,&quot;douyin_h5_flow&quot;],&quot;forbid_zlink_change_list&quot;:{&quot;scene_from&quot;:[&quot;click_schema_hwqs_cloud&quot;,&quot;qs_os_list&quot;,&quot;douyin_h5_flow&quot;,&quot;douyin_h5_search&quot;,&quot;ug_page_delivery&quot;],&quot;from_aid&quot;:[&quot;615883&quot;,&quot;1349&quot;],&quot;app&quot;:[&quot;douyin_search&quot;,&quot;maya&quot;,&quot;new_maya&quot;],&quot;app_name_android&quot;:[&quot;baidu&quot;,&quot;weixin&quot;,&quot;qq&quot;,&quot;qqbrowser&quot;,&quot;safari&quot;],&quot;app_name_ios&quot;:[&quot;baidu&quot;,&quot;weixin&quot;,&quot;qqbrowser&quot;,&quot;uc&quot;,&quot;quark&quot;,&quot;qq&quot;]},&quot;gov_domain_white_list_path&quot;:[&quot;/share/video/&quot;],&quot;overseas_config&quot;:{&quot;is_oversea_download_switch&quot;:1,&quot;effective_ratio&quot;:100,&quot;is_popup_en_switch&quot;:0,&quot;english_word_line1&quot;:&quot;Watch¥Video.¥Have¥fun.&quot;,&quot;english_word_line2&quot;:&quot;Be¥creative!&quot;,&quot;english_btn_word:&quot;:&quot;Open¥Now&quot;,&quot;cn_word_line1&quot;:&quot;记录美好生活&quot;,&quot;cn_word_line2&quot;:&quot;&quot;,&quot;cn_btn_word&quot;:&quot;现在打开&quot;,&quot;en_btn_text&quot;:&quot;Get¥Douyin¥App&quot;,&quot;en_banner_open_text&quot;:&quot;Get¥App&quot;,&quot;android_download_url&quot;:&quot;https://z.douyin.com/dBkb1&quot;,&quot;ios_download_url&quot;:&quot;https://z.douyin.com/dBkb1&quot;,&quot;btn_text_en_switch&quot;:1,&quot;pop_up_enable_switch&quot;:0,&quot;no_redirect_m_station_arr&quot;:[&quot;hk&quot;,&quot;mo&quot;]},&quot;qps_limit&quot;:20,&quot;ssr_down_grade_csr&quot;:0,&quot;token_encry_cooperation&quot;:{&quot;switch&quot;:true,&quot;block_ratio&quot;:1000,&quot;post_list_relation_block_ratio&quot;:0,&quot;fe_key&quot;:&quot;douyin_reflow_token&quot;,&quot;new_fe_key&quot;:&quot;e2rere3sb2rsvc&quot;,&quot;new_fe_key_ratio&quot;:0,&quot;algorithm&quot;:[{&quot;encry_algorithm&quot;:&quot;aes-128-cbc&quot;,&quot;sk&quot;:&quot;sykKwe59_q11peDD&quot;,&quot;block_size&quot;:16,&quot;token_expire_time_millisecond&quot;:100000},{&quot;encry_algorithm&quot;:&quot;aes-128-cbc&quot;,&quot;sk&quot;:&quot;sykKwe59_q11peDz&quot;,&quot;block_size&quot;:16,&quot;token_expire_time_millisecond&quot;:100000,&quot;algorithm_expire_time_millisecond&quot;:1717751102312}]}}"></div>
<div id="douyin_reflow_webId" webid="7568501739219011082" usercip="171.213.179.198"></div>
<div id="douyin_reflow_page"></div>
<script nonce="" src="https://privacy.zijieapi.com/api/web-cmp/sdk/?project_key=61566ccba19eeb6d"></script><script nonce="">!function(n,t){if(n.LogAnalyticsObject=t,!n[t]){function c(){c.q.push(arguments)}c.q=c.q||[],n[t]=c}n[t].l=+new Date}(window,"collectEvent")</script><script nonce="" async="" src="https://lf3-cdn-tos.bytescm.com/obj/static/log-sdk/collect/5.0/collect.js" crossorigin="anonymous"></script><script nonce="">!function(e,r,n,t,s,a,i,o,c,l,d,m,p,f){a="precollect",i="getAttribute",o="addEventListener",l=function(e){(d=[].slice.call(arguments)).push(Date.now(),location.href),(e==a?l.p.a:l.q).push(d)},l.q=[],l.p={a:[]},e[s]=l,(m=document.createElement("script")).src=n+"?bid=aweme_share&globalName="+s,m.crossOrigin=n.indexOf("sdk-web")>0?"anonymous":"use-credentials",r.getElementsByTagName("head")[0].appendChild(m),o in e&&(l.pcErr=function(r){r=r||e.event,(p=r.target||r.srcElement)instanceof Element||p instanceof HTMLElement?p[i]("integrity")?e[s](a,"sri",p[i]("href")||p[i]("src")):e[s](a,"st",{tagName:p.tagName,url:p[i]("href")||p[i]("src")}):e[s](a,"err",r.error||r.message)},l.pcRej=function(r){r=r||e.event,e[s](a,"err",r.reason||r.detail&&r.detail.reason)},e[o]("error",l.pcErr,!0),e[o]("unhandledrejection",l.pcRej,!0)),"PerformanceLongTaskTiming"in e&&((f=l.pp={entries:[]}).observer=new PerformanceObserver((function(e){f.entries=f.entries.concat(e.getEntries())})),f.observer.observe({entryTypes:["longtask","largest-contentful-paint","layout-shift"]}))}(window,document,"https://lf3-short.ibytedapm.com/slardar/fe/sdk-web/browser.cn.js",0,"slardar")</script><script nonce="">var regExpRes=location.pathname.match(/\/([^/?#]*)\/?$/);window.slardar("context.merge",{item_id:regExpRes&&regExpRes[1]}),window.slardar("init",{bid:"aweme_share",pid:"video_ssr",release:"1.0.0.2579",env:"online"}),window.slardar("on","beforeSend",(function(e){e.ev_type&&"performance"===e.ev_type&&("lcp"===e.payload.name&&(console.log("[dev]\u6267\u884c=>slardarLCP",(new Date).getTime()),window.slardarLcp=(new Date).getTime()));return e}))</script><script nonce="" id="__LOADABLE_REQUIRED_CHUNKS__" type="application/json">["7431","4097","5118","1006","3384","9973"]</script><script nonce="" id="__LOADABLE_REQUIRED_CHUNKS___ext" type="application/json">{"namedChunks":["video_(id)/page"]}</script><script nonce="" crossorigin="anonymous" defer="true" src="//lf-douyin-mobile.bytecdn.com/obj/growth-douyin-share/growth/douyin_ug/static/js/async/7431.689d5bb9.js"></script><script nonce="" crossorigin="anonymous" defer="true" src="//lf-douyin-mobile.bytecdn.com/obj/growth-douyin-share/growth/douyin_ug/static/js/async/4097.16d21948.js"></script><script nonce="" crossorigin="anonymous" defer="true" src="//lf-douyin-mobile.bytecdn.com/obj/growth-douyin-share/growth/douyin_ug/static/js/async/5118.1442fa2d.js"></script><script nonce="" crossorigin="anonymous" defer="true" src="//lf-douyin-mobile.bytecdn.com/obj/growth-douyin-share/growth/douyin_ug/static/js/async/1006.fc12eee8.js"></script><script nonce="" crossorigin="anonymous" defer="true" src="//lf-douyin-mobile.bytecdn.com/obj/growth-douyin-share/growth/douyin_ug/static/js/async/3384.f0fc9bd9.js"></script><script nonce="" crossorigin="anonymous" defer="true" src="//lf-douyin-mobile.bytecdn.com/obj/growth-douyin-share/growth/douyin_ug/static/js/async/video_(id)/page.15670c72.js"></script><script nonce="">window._SSR_DATA = {"data":{},"context":{"request":{"params":{},"query":{"region":"CN","mid":"7535852694951480091","u_code":"ej0dk7l2g004jm","did":"MS4wLjABAAAA506pteC-jo5menwXlQdVeSPo9ZDGF9bw7il9PPlJa1SBQnBGUFwmjq9zZ3n8O0CP","iid":"MS4wLjABAAAANwkJuWIRFOzg5uCpDRpMj4OX-QryoDgn-yYlXQnRwQQ","with_sec_did":"1","video_share_track_ver":"","titleType":"title","share_sign":"b3bf1vN0jxpuYyeTfmAtHCcIdS8VfzSGoBFpwnce2bM-","share_version":"170400","ts":"1762177424","from_aid":"6383","from_ssr":"1","share_track_info":"{\"link_description_type\":\"\"}","from":"web_code_link"},"pathname":"\u002Fshare\u002Fvideo\u002F7535852404919586100\u002F","host":"www.iesdouyin.com","url":"https:\u002F\u002Fwww.iesdouyin.com\u002Fshare\u002Fvideo\u002F7535852404919586100\u002F?region=CN&mid=7535852694951480091&u_code=ej0dk7l2g004jm&did=MS4wLjABAAAA506pteC-jo5menwXlQdVeSPo9ZDGF9bw7il9PPlJa1SBQnBGUFwmjq9zZ3n8O0CP&iid=MS4wLjABAAAANwkJuWIRFOzg5uCpDRpMj4OX-QryoDgn-yYlXQnRwQQ&with_sec_did=1&video_share_track_ver=&titleType=title&share_sign=b3bf1vN0jxpuYyeTfmAtHCcIdS8VfzSGoBFpwnce2bM-&share_version=170400&ts=1762177424&from_aid=6383&from_ssr=1&share_track_info=%7B%22link_description_type%22%3A%22%22%7D&from=web_code_link"},"reporter":{}},"mode":"string","renderLevel":2}</script>
<script nonce="">window._ROUTER_DATA = {"loaderData":{"video_layout":null,"video_(id)\u002Fpage":{"ua":"Mozilla\u002F5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit\u002F605.1.15 (KHTML, like Gecko) Version\u002F16.0 Mobile\u002F15E148 Safari\u002F604.1","isSpider":false,"webId":"7568501739219011082","query":{"region":"CN","mid":"7535852694951480091","u_code":"ej0dk7l2g004jm","did":"MS4wLjABAAAA506pteC-jo5menwXlQdVeSPo9ZDGF9bw7il9PPlJa1SBQnBGUFwmjq9zZ3n8O0CP","iid":"MS4wLjABAAAANwkJuWIRFOzg5uCpDRpMj4OX-QryoDgn-yYlXQnRwQQ","with_sec_did":"1","video_share_track_ver":"","titleType":"title","share_sign":"b3bf1vN0jxpuYyeTfmAtHCcIdS8VfzSGoBFpwnce2bM-","share_version":"170400","ts":"1762177424","from_aid":"6383","from_ssr":"1","share_track_info":"{\"link_description_type\":\"\"}","from":"web_code_link"},"renderInSSR":1,"lastPath":"7535852404919586100","appName":"safari","host":"www.iesdouyin.com","isNotSupportWebp":false,"commonContext":{"ua":"Mozilla\u002F5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit\u002F605.1.15 (KHTML, like Gecko) Version\u002F16.0 Mobile\u002F15E148 Safari\u002F604.1","isSpider":false,"webId":"7568501739219011082","query":{"region":"CN","mid":"7535852694951480091","u_code":"ej0dk7l2g004jm","did":"MS4wLjABAAAA506pteC-jo5menwXlQdVeSPo9ZDGF9bw7il9PPlJa1SBQnBGUFwmjq9zZ3n8O0CP","iid":"MS4wLjABAAAANwkJuWIRFOzg5uCpDRpMj4OX-QryoDgn-yYlXQnRwQQ","with_sec_did":"1","video_share_track_ver":"","titleType":"title","share_sign":"b3bf1vN0jxpuYyeTfmAtHCcIdS8VfzSGoBFpwnce2bM-","share_version":"170400","ts":"1762177424","from_aid":"6383","from_ssr":"1","share_track_info":"{\"link_description_type\":\"\"}","from":"web_code_link"},"renderInSSR":1,"lastPath":"7535852404919586100","appName":"safari","host":"www.iesdouyin.com","isNotSupportWebp":false},"isAiDouyinOpt":"","videoInfoRes":{"extra":{"logid":"202511032212173D058217B6CB452C34B1","now":1762179137769},"filter_list":[],"is_oversea":0,"item_list":[{"aweme_id":"7535852404919586100","desc":"恐怖小狗《The Puppy》一命速通 #恐怖游戏 #steam","create_time":1754577401,"author":{"short_id":"45599235625","nickname":"式子大王","signature":"恐游哥布林主播\n直男","avatar_thumb":{"uri":"100x100\u002Faweme-avatar\u002Ftos-cn-avt-0015_00ac6c591671011509c48dd7b4c27b5d","url_list":["https:\u002F\u002Fp3.douyinpic.com\u002Faweme\u002F100x100\u002Faweme-avatar\u002Ftos-cn-avt-0015_00ac6c591671011509c48dd7b4c27b5d.jpeg?from=327834062","https:\u002F\u002Fp6.douyinpic.com\u002Faweme\u002F100x100\u002Faweme-avatar\u002Ftos-cn-avt-0015_00ac6c591671011509c48dd7b4c27b5d.jpeg?from=327834062","https:\u002F\u002Fp9.douyinpic.com\u002Faweme\u002F100x100\u002Faweme-avatar\u002Ftos-cn-avt-0015_00ac6c591671011509c48dd7b4c27b5d.jpeg?from=327834062"],"width":720,"height":720},"avatar_medium":{"uri":"100x100\u002Faweme-avatar\u002Ftos-cn-avt-0015_00ac6c591671011509c48dd7b4c27b5d","url_list":["https:\u002F\u002Fp3.douyinpic.com\u002Faweme\u002F100x100\u002Faweme-avatar\u002Ftos-cn-avt-0015_00ac6c591671011509c48dd7b4c27b5d.jpeg?from=327834062","https:\u002F\u002Fp6.douyinpic.com\u002Faweme\u002F100x100\u002Faweme-avatar\u002Ftos-cn-avt-0015_00ac6c591671011509c48dd7b4c27b5d.jpeg?from=327834062","https:\u002F\u002Fp9.douyinpic.com\u002Faweme\u002F100x100\u002Faweme-avatar\u002Ftos-cn-avt-0015_00ac6c591671011509c48dd7b4c27b5d.jpeg?from=327834062"],"width":720,"height":720},"follow_status":0,"aweme_count":28,"following_count":0,"favoriting_count":0,"unique_id":"45599235625","mplatform_followers_count":0,"followers_detail":null,"platform_sync_info":null,"geofencing":null,"policy_version":null,"sec_uid":"MS4wLjABAAAAwJYnMkU3OD9tFsyxKQDJmIGfgYZ2lgi10ta9BsejXOA-m2pDSN0v9XVg2dJaeVvP","type_label":null,"card_entries":null,"mix_info":null},"music":{"mid":"7535852694951480091","title":"@式子大王创作的原声一式子大王（原声中的歌曲：Apex-Jamvana）","author":"式子大王","cover_hd":{"uri":"1080x1080\u002Faweme-avatar\u002Ftos-cn-avt-0015_00ac6c591671011509c48dd7b4c27b5d","url_list":["https:\u002F\u002Fp3.douyinpic.com\u002Faweme\u002F1080x1080\u002Faweme-avatar\u002Ftos-cn-avt-0015_00ac6c591671011509c48dd7b4c27b5d.jpeg?from=327834062","https:\u002F\u002Fp6.douyinpic.com\u002Faweme\u002F1080x1080\u002Faweme-avatar\u002Ftos-cn-avt-0015_00ac6c591671011509c48dd7b4c27b5d.jpeg?from=327834062","https:\u002F\u002Fp9.douyinpic.com\u002Faweme\u002F1080x1080\u002Faweme-avatar\u002Ftos-cn-avt-0015_00ac6c591671011509c48dd7b4c27b5d.jpeg?from=327834062"],"width":720,"height":720},"cover_large":{"uri":"1080x1080\u002Faweme-avatar\u002Ftos-cn-avt-0015_00ac6c591671011509c48dd7b4c27b5d","url_list":["https:\u002F\u002Fp3.douyinpic.com\u002Faweme\u002F1080x1080\u002Faweme-avatar\u002Ftos-cn-avt-0015_00ac6c591671011509c48dd7b4c27b5d.jpeg?from=327834062","https:\u002F\u002Fp6.douyinpic.com\u002Faweme\u002F1080x1080\u002Faweme-avatar\u002Ftos-cn-avt-0015_00ac6c591671011509c48dd7b4c27b5d.jpeg?from=327834062","https:\u002F\u002Fp9.douyinpic.com\u002Faweme\u002F1080x1080\u002Faweme-avatar\u002Ftos-cn-avt-0015_00ac6c591671011509c48dd7b4c27b5d.jpeg?from=327834062"],"width":720,"height":720},"cover_medium":{"uri":"720x720\u002Faweme-avatar\u002Ftos-cn-avt-0015_00ac6c591671011509c48dd7b4c27b5d","url_list":["https:\u002F\u002Fp3.douyinpic.com\u002Faweme\u002F720x720\u002Faweme-avatar\u002Ftos-cn-avt-0015_00ac6c591671011509c48dd7b4c27b5d.jpeg?from=327834062","https:\u002F\u002Fp9.douyinpic.com\u002Faweme\u002F720x720\u002Faweme-avatar\u002Ftos-cn-avt-0015_00ac6c591671011509c48dd7b4c27b5d.jpeg?from=327834062","https:\u002F\u002Fp6.douyinpic.com\u002Faweme\u002F720x720\u002Faweme-avatar\u002Ftos-cn-avt-0015_00ac6c591671011509c48dd7b4c27b5d.jpeg?from=327834062"],"width":720,"height":720},"cover_thumb":{"uri":"168x168\u002Faweme-avatar\u002Ftos-cn-avt-0015_00ac6c591671011509c48dd7b4c27b5d","url_list":["https:\u002F\u002Fp3.douyinpic.com\u002Fimg\u002Faweme-avatar\u002Ftos-cn-avt-0015_00ac6c591671011509c48dd7b4c27b5d~c5_168x168.jpeg?from=327834062","https:\u002F\u002Fp6.douyinpic.com\u002Fimg\u002Faweme-avatar\u002Ftos-cn-avt-0015_00ac6c591671011509c48dd7b4c27b5d~c5_168x168.jpeg?from=327834062","https:\u002F\u002Fp9.douyinpic.com\u002Fimg\u002Faweme-avatar\u002Ftos-cn-avt-0015_00ac6c591671011509c48dd7b4c27b5d~c5_168x168.jpeg?from=327834062"],"width":720,"height":720},"duration":352,"position":null,"status":1},"cha_list":null,"video":{"play_addr":{"uri":"v0200fg10000d2abjffog65h5b4rsmng","url_list":["https:\u002F\u002Faweme.snssdk.com\u002Faweme\u002Fv1\u002Fplaywm\u002F?video_id=v0200fg10000d2abjffog65h5b4rsmng&ratio=720p&line=0"]},"cover":{"uri":"tos-cn-i-dy\u002Fdd954010955847ea8f5493d60d984341","url_list":["https:\u002F\u002Fp3-sign.douyinpic.com\u002Ftos-cn-i-dy\u002Fdd954010955847ea8f5493d60d984341~tplv-dy-resize-walign-adapt-aq:540:q75.webp?lk3s=138a59ce&x-expires=1763388000&x-signature=x6p6PYRGDbgA5TziiC6TBWT%2FJDs%3D&from=327834062&s=PackSourceEnum_DOUYIN_REFLOW&se=false&sc=cover&biz_tag=aweme_video&l=202511032212173D058217B6CB452C34B1","https:\u002F\u002Fp11-sign.douyinpic.com\u002Ftos-cn-i-dy\u002Fdd954010955847ea8f5493d60d984341~tplv-dy-resize-walign-adapt-aq:540:q75.webp?lk3s=138a59ce&x-expires=1763388000&x-signature=5qNIvugr%2BpyMDF3dMtZ03uNUgwc%3D&from=327834062&s=PackSourceEnum_DOUYIN_REFLOW&se=false&sc=cover&biz_tag=aweme_video&l=202511032212173D058217B6CB452C34B1","https:\u002F\u002Fp26-sign.douyinpic.com\u002Ftos-cn-i-dy\u002Fdd954010955847ea8f5493d60d984341~tplv-dy-resize-walign-adapt-aq:540:q75.webp?lk3s=138a59ce&x-expires=1763388000&x-signature=jHELHOtD5LvOXEJ7f9wEWiBwV2g%3D&from=327834062&s=PackSourceEnum_DOUYIN_REFLOW&se=false&sc=cover&biz_tag=aweme_video&l=202511032212173D058217B6CB452C34B1","https:\u002F\u002Fp3-sign.douyinpic.com\u002Ftos-cn-i-dy\u002Fdd954010955847ea8f5493d60d984341~tplv-dy-resize-walign-adapt-aq:540:q75.jpeg?lk3s=138a59ce&x-expires=1763388000&x-signature=VnwZGbPOCkDLnLzE%2BcjmRABPonk%3D&from=327834062&s=PackSourceEnum_DOUYIN_REFLOW&se=false&sc=cover&biz_tag=aweme_video&l=202511032212173D058217B6CB452C34B1"],"width":540,"height":720},"height":1080,"width":1920,"bit_rate":null,"duration":352596,"big_thumbs":null},"statistics":{"aweme_id":"7535852404919586100","comment_count":178,"digg_count":7007,"play_count":0,"share_count":474,"collect_count":2804},"text_extra":[{"start":20,"end":25,"type":1,"hashtag_name":"恐怖游戏","hashtag_id":1767586622634062},{"start":26,"end":32,"type":1,"hashtag_name":"steam","hashtag_id":1583312161022990}],"video_labels":null,"aweme_type":4,"image_infos":null,"risk_infos":{"warn":true,"type":53,"content":"作者声明：可能引人不适，请谨慎观看","reflow_unplayable":0},"comment_list":null,"geofencing":null,"video_text":null,"label_top_text":null,"promotions":null,"long_video":null,"images":null,"group_id_str":"7535849207987588392","chapter_list":null,"interaction_stickers":null,"img_bitrate":null,"chapter_bar_color":null,"common_labels":null}],"status_code":0},"itemId":"7535852404919586100","isVideoOptimize":false,"openDirectGroup":"4","isVideoStyleOptimize":true,"isButtonOptimize":true,"isAutoOpenApp":false,"darkModeAdaptation":false,"serverToken":"","abParams":{"reflow_video_page_optimize":"2","search_video_mark_abtest":"0","share_page_dark_mode_adaptation":"1","reflow_page_open_app_optimize":"4","reflow_page_auto_open_optimize":"1","reflow_to_featured_app":"1","select_pool_data":{"use_new_select_scope":0},"jx_label_display_position_list":["main","related_card"],"reflow_page_615883_optimaze":false}}},"errors":null}</script></body></html>`
	err := p.extractDataFromHTML(html)
	// 保存视频信息，当前只获取一个视频，所以直接保存
	p.videos = append(p.videos, p.videoInfo)
	if err != nil {
		return fmt.Errorf("提取视频 URL 失败: %v", err)
	}
	// 下载视频
	if err = p.downloadVideo(); err != nil {
		return fmt.Errorf("下载视频失败: %v", err)
	}
	return nil
}

// initPlaywrightAndBrowser 初始化 Playwright 和浏览器
func (p *DouYinProcessor) initPlaywrightAndBrowser() (playwright.BrowserContext, playwright.Page, *playwright.Playwright, error) {
	pw, err := playwright.Run()
	if err != nil {
		err = playwright.Install()
		if err != nil {
			return nil, nil, nil, fmt.Errorf("启动 Playwright 失败: %v", err)
		}
	}

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		pw.Stop()
		return nil, nil, nil, fmt.Errorf("启动浏览器失败: %v", err)
	}

	// 随机选择用户代理
	selectedUserAgent := p.getRandomUserAgent()

	// 创建浏览器上下文
	contextOptions := playwright.BrowserNewContextOptions{
		UserAgent:         playwright.String(selectedUserAgent),
		Viewport:          &playwright.Size{Width: 375, Height: 667},
		DeviceScaleFactor: playwright.Float(2),
		Locale:            playwright.String("zh-CN"),
		TimezoneId:        playwright.String("Asia/Shanghai"),
		IsMobile:          playwright.Bool(true),
		HasTouch:          playwright.Bool(true),
		ColorScheme:       (*playwright.ColorScheme)(playwright.String("light")),
		ExtraHttpHeaders: map[string]string{
			"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,video/mp4,*/*;q=0.8",
			"Accept-Language":           "zh-CN,zh;q=0.9,en;q=0.8",
			"Connection":                "keep-alive",
			"Upgrade-Insecure-Requests": "1",
		},
	}

	ctx, err := browser.NewContext(contextOptions)
	if err != nil {
		browser.Close()
		pw.Stop()
		return nil, nil, nil, fmt.Errorf("创建上下文失败: %v", err)
	}

	page, err := ctx.NewPage()
	if err != nil {
		ctx.Close()
		browser.Close()
		pw.Stop()
		return nil, nil, nil, fmt.Errorf("创建页面失败: %v", err)
	}

	return ctx, page, pw, nil
}

// getRandomUserAgent 获取随机用户代理
func (p *DouYinProcessor) getRandomUserAgent() string {
	userAgents := []string{
		"Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36 QuarkPC/4.6.0.558",
	}
	// 随机选择一条 userAgent
	rand.New(rand.NewSource(time.Now().Unix()))
	return userAgents[rand.Intn(len(userAgents))]
}

// loadCookies 加载 cookies
func (p *DouYinProcessor) loadCookies(ctx playwright.BrowserContext) error {
	cookies := p.parseDouYinCookiesFile()
	if cookies != nil && len(cookies) > 0 {
		if err := ctx.AddCookies(cookies); err != nil {
			return err
		}
		utils.InfoWithFormat("成功加载 %d 个 cookies", len(cookies))
	}
	return nil
}

// _extractVideoID 提取视频ID
func (p *DouYinProcessor) _extractVideoID(page playwright.Page, link string) (string, error) {
	videoID := ""

	// 监听网络请求中的 video_id
	page.On("request", func(request playwright.Request) {
		requestURL := request.URL()
		if strings.Contains(requestURL, "video_id=") {
			m := regexp.MustCompile(videoIDPattern).FindStringSubmatch(requestURL)
			if len(m) > 1 {
				videoID = m[1]
				utils.DebugWithFormat("网络请求中捕获到 video_id: %s", videoID)
			}
		}
	})

	// 访问 URL - 等待网络空闲状态以确保页面完全加载
	if _, err := page.Goto(link, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle, // 等待网络空闲
		Timeout:   playwright.Float(60000),
	}); err != nil {
		return "", fmt.Errorf("访问页面失败: %v", err)
	}

	// 确保页面完全加载 - 等待所有资源加载完成
	if err := page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State: playwright.LoadStateNetworkidle, // 再次确认网络空闲
	}); err != nil {
		utils.InfoWithFormat("等待页面网络空闲失败: %v", err)
	}

	// 等待特定元素出现，确保关键内容已加载
	waitForElements(page)

	// 尝试滚动页面，确保动态加载的内容也被加载
	if _, err := page.Evaluate(`window.scrollTo(0, document.body.scrollHeight);`); err != nil {
		utils.InfoWithFormat("页面滚动失败: %v", err)
	}

	// 滚动后再次等待网络空闲
	if err := page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State: playwright.LoadStateNetworkidle,
	}); err != nil {
		utils.InfoWithFormat("滚动后等待网络空闲失败: %v", err)
	}

	// 使用智能等待替代硬编码延时
	if err := p.waitForVideoContent(page); err != nil {
		utils.InfoWithFormat("等待视频内容超时: %v", err)
	}

	// 从当前URL提取视频ID
	currentURL := page.URL()
	if m := regexp.MustCompile(videoURLPattern).FindStringSubmatch(currentURL); len(m) > 1 {
		videoID = m[1]
		utils.DebugWithFormat("从当前 URL 提取到 video_id: %s", videoID)
	}

	// 从原始URL提取视频ID作为备选
	if videoID == "" {
		log.Println("未捕获到 video_id，尝试从 URL 直接提取")
		if m := regexp.MustCompile(videoURLPattern).FindStringSubmatch(link); len(m) > 1 {
			videoID = m[1]
			utils.DebugWithFormat("从原始 URL 提取到 video_id: %s", videoID)
		}
	}

	// 如果仍然没有获取到videoID，尝试从页面内容中搜索aweme_id
	if videoID == "" {
		utils.DebugWithFormat("尝试从页面内容中搜索aweme_id")
		html, err := page.Content()
		if err == nil {
			// 尝试直接匹配aweme_id
			awemeIDRegex := regexp.MustCompile(`"aweme_id"\s*:\s*"([^"]+)"`)
			if m := awemeIDRegex.FindStringSubmatch(html); len(m) > 1 {
				videoID = m[1]
				utils.DebugWithFormat("从页面内容中提取到 aweme_id: %s", videoID)
			}
		}
	}

	if videoID == "" {
		return "", errors.New("未能捕获到视频数据")
	}

	return videoID, nil
}

// parseDouYinCookiesFile 解析抖音 cookies 文件
func (p *DouYinProcessor) parseDouYinCookiesFile() []playwright.OptionalCookie {
	playwrightCookies := make([]playwright.OptionalCookie, 0)
	domains := []string{
		".douyin.com",
		"douyin.com",
		"www.douyin.com",
		"v.douyin.com",
		"www.iesdouyin.com",
		"iesdouyin.com",
	}
	for _, domain := range domains {
		cookies := utils.GetCookiesByDomain(p.cfg.CookieCloud.CookieFilePath, domain)
		if len(cookies) > 0 {
			// 直接在这里处理cookies，避免值传递问题
			for name, value := range cookies {
				playwrightCookies = append(playwrightCookies, playwright.OptionalCookie{
					Name:     name,
					Value:    value,
					Domain:   playwright.String(".douyin.com"),
					Path:     playwright.String("/"),
					HttpOnly: playwright.Bool(true),
					Secure:   playwright.Bool(true),
					SameSite: (*playwright.SameSiteAttribute)(playwright.String("Lax")),
				})
			}
		}
	}
	return playwrightCookies
}

// extractDataFromHTML 从HTML中提取视频URL
func (p *DouYinProcessor) extractDataFromHTML(html string) error {
	utils.DebugWithFormat("[extract] HTML长度: %d 字符", len(html))

	// 提取视频标题
	titleMatches := regexp.MustCompile(titlePattern).FindAllStringSubmatch(html, -1)
	var title string
	for _, titleMatch := range titleMatches {
		title = titleMatch[1]
		if !strings.Contains(title, "-") {
			continue
		}
		if title != "" {
			sTitle := strings.Split(title, "-")
			p.videoInfo.Title = sTitle[0]
		}
	}
	utils.InfoWithFormat("[extract] 提取到视频标题: %s", title)

	// 查找包含视频数据的script标签
	scriptRegex := regexp.MustCompile(scriptPattern)
	scriptMatches := scriptRegex.FindAllStringSubmatch(html, -1)

	for _, scriptMatch := range scriptMatches {
		scriptContent := scriptMatch[1]
		// 检查是否包含关键数据标记
		if !strings.Contains(scriptContent, "aweme_id") || !strings.Contains(scriptContent, "status_code") {
			continue
		}
		// 尝试提取JSON部分
		jsonMatches := regexp.MustCompile(jsonPattern).FindAllStringSubmatch(scriptContent, -1)
		for _, jsonMatch := range jsonMatches {
			jsonStr := jsonMatch[1]
			// 清理JSON，确保匹配完整的JSON结构
			cleanJSON, err := p.cleanJSONString(jsonStr)
			if err != nil || cleanJSON == "" {
				continue
			}
			var data map[string]interface{}
			if err := json.Unmarshal([]byte(cleanJSON), &data); err != nil {
				continue
			}
			// 递归查找视频URL
			hjd := p.findDataInJson(data)
			if hjd.VideoUrl != "" {
				p.videoInfo.CoverUrl = hjd.CoverUrl
				p.videoInfo.Time = hjd.Time
				p.videoInfo.Desc = hjd.Desc
				p.videoInfo.Author = hjd.Author
				p.videoInfo.Ratio = hjd.Ratio
				p.videoInfo.DownloadUrl = hjd.VideoUrl
				utils.InfoWithFormat("[extract] 提取到视频信息: %v", hjd)
				return nil
			}
		}
	}

	return errors.New("未能提取到视频URL")
}

// cleanJSONString 清理JSON字符串，确保其完整性
func (p *DouYinProcessor) cleanJSONString(jsonStr string) (string, error) {
	braceCount := 0
	jsonEnd := -1
	for i, char := range jsonStr {
		if char == '{' {
			braceCount++
		} else if char == '}' {
			braceCount--
			if braceCount == 0 {
				jsonEnd = i + 1
				break
			}
		}
	}

	if jsonEnd > 0 {
		return jsonStr[:jsonEnd], nil
	}
	return "", errors.New("无法找到完整的JSON结构")
}

type htmlJsonData struct {
	VideoUrl string
	CoverUrl string
	Time     string
	Desc     string
	Author   string
	Ratio    string
}

// waitForVideoContent 智能等待视频内容加载完成
func (p *DouYinProcessor) waitForVideoContent(page playwright.Page) error {
	// 使用轮询检查页面是否包含视频关键数据
	deadline := time.Now().Add(30 * time.Second) // 最多等待30秒
	for time.Now().Before(deadline) {
		html, err := page.Content()
		if err == nil && (strings.Contains(html, "aweme_id") || strings.Contains(html, "video")) {
			utils.DebugWithFormat("检测到视频内容已加载")
			return nil
		}
		time.Sleep(500 * time.Millisecond) // 每500ms检查一次
	}
	return errors.New("等待视频内容超时")
}

// waitForElements 等待关键元素出现
func waitForElements(page playwright.Page) {
	// 尝试等待几个关键元素出现，但不阻塞主流程
	go func() {
		// 等待视频容器元素
		if _, err := page.WaitForSelector("title", playwright.PageWaitForSelectorOptions{
			Timeout: playwright.Float(10000),
		}); err != nil {
			utils.DebugWithFormat("等待视频元素超时: %v", err)
		}
	}()
}

// findDataInJson 在数据结构中查找视频URL
func (p *DouYinProcessor) findDataInJson(data map[string]interface{}) *htmlJsonData {
	var hjd = &htmlJsonData{}

	// 安全提取字段，避免 panic
	getString := func(m map[string]interface{}, key string) string {
		if v, ok := m[key].(string); ok {
			return v
		}
		return ""
	}

	getMap := func(m map[string]interface{}, key string) map[string]interface{} {
		if key == "" {
			return m
		}
		if v, ok := m[key].(map[string]interface{}); ok {
			return v
		}
		return nil
	}

	getList := func(m map[string]interface{}, key string) []interface{} {
		if v, ok := m[key].([]interface{}); ok {
			return v
		}
		return nil
	}

	// 从 data 开始逐层解析
	if loaderData := getMap(data, "loaderData"); loaderData != nil {
		if videoPage := getMap(loaderData, "video_(id)/page"); videoPage != nil {
			if videoInfoRes := getMap(videoPage, "videoInfoRes"); videoInfoRes != nil {
				if itemList := getList(videoInfoRes, "item_list"); len(itemList) > 0 {
					if item := getMap(itemList[0].(map[string]interface{}), ""); item != nil {
						hjd.Desc = getString(item, "desc")

						if author := getMap(item, "author"); author != nil {
							hjd.Author = getString(author, "nickname")
						}

						if createTime, ok := item["create_time"].(int64); ok {
							hjd.Time = datetime.FormatTimeToStr(time.Unix(createTime, 0), "yyyy-mm-dd hh:mm:ss")
						}

						if video := getMap(item, "video"); video != nil {
							if playAddr := getMap(video, "play_addr"); playAddr != nil {
								if urls := getList(playAddr, "url_list"); len(urls) > 0 {
									hjd.VideoUrl = urls[0].(string)
									if hjd.VideoUrl != "" {
										hjd.VideoUrl = strings.Replace(hjd.VideoUrl, "playwm", "play", 1)
									}
								}
							}

							if cover := getMap(video, "cover"); cover != nil {
								if urls := getList(cover, "url_list"); len(urls) > 0 {
									hjd.CoverUrl = urls[0].(string)
								}
							}
						}
					}
				}
			}
		}
	}

	return hjd
}

// extractURLFromField 从指定字段提取URL
func (p *DouYinProcessor) extractURLFromField(data map[string]interface{}, fieldName string, isVideoWatermarked bool) string {
	field, ok := data[fieldName]
	if !ok {
		return ""
	}

	fieldMap, ok := field.(map[string]interface{})
	if !ok {
		return ""
	}

	urlList, ok := fieldMap["url_list"].([]interface{})
	if !ok || len(urlList) == 0 {
		return ""
	}

	// 遍历所有URL，跳过无效地址
	for _, item := range urlList {
		if rowUrl, ok := item.(string); ok && strings.HasPrefix(rowUrl, "http") {
			utils.DebugWithFormat("[extract] 从%s.url_list找到: %s", fieldName, rowUrl)

			// 检查URL是否能正常访问
			if !p.isURLAccessible(rowUrl) {
				utils.DebugWithFormat("[extract] URL不可访问，跳过: %s", rowUrl)
				continue
			}

			// 如果是有水印的视频，尝试去除水印
			if isVideoWatermarked && strings.Contains(rowUrl, "playwm") {
				// 替换playwm为play，转换为无水印视频
				rowUrl = strings.Replace(rowUrl, "playwm", "play", 1)
				// 检查转换后的URL是否可访问
				if !p.isURLAccessible(rowUrl) {
					utils.DebugWithFormat("[extract] 转换后的URL不可访问，使用原URL: %s", rowUrl)
					// 恢复原URL
					rowUrl = strings.Replace(rowUrl, "play", "playwm", 1)
				}
			}
			return rowUrl
		}
		utils.DebugWithFormat("[extract] 跳过无效URL: %v", item)
	}

	return ""
}

// isURLAccessible 检查URL是否可正常访问
func (p *DouYinProcessor) isURLAccessible(url string) bool {
	// 使用HEAD请求快速检查URL可用性，避免下载整个文件
	client := &http.Client{
		Timeout: 5 * time.Second, // 设置超时时间
	}

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return false
	}

	// 设置User-Agent，避免被拦截
	req.Header.Set("User-Agent", p.getRandomUserAgent())

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// 检查响应状态码，200-299表示成功
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

// isVideoURL 判断是否为有效的视频URL
func (p *DouYinProcessor) isVideoURL(url string) bool {
	videoExtensions := []string{".mp4", ".m3u8", ".ts", "douyinvod.com", "snssdk.com"}
	for _, ext := range videoExtensions {
		if strings.Contains(strings.ToLower(url), ext) {
			return true
		}
	}
	return false
}

// getMapKeys 获取map的所有键
func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// downloadResource 下载资源到指定路径
func (p *DouYinProcessor) downloadVideo() error {
	// 下载视频逻辑
	for _, videoInfo := range p.videos {
		// 首先下载视频文件
		fp := filepath.Join(p.tempDir, videoInfo.Author)
		fn := videoInfo.Desc + ".mp4"
		downloadSize, err := p._downloadResource(videoInfo.DownloadUrl, fp, fn)
		if err != nil {
			return err
		}
		videoInfo.Size = downloadSize
		utils.InfoWithFormat("[download] 下载完成: %s", filepath.Join(fp, fn))

		if videoInfo.CoverUrl != "" {
			// 下载封面图片
			fn = videoInfo.Title + ".png"
			_, err = p._downloadResource(videoInfo.DownloadUrl, fp, fn)
			if err != nil {
				return err
			}
			utils.InfoWithFormat("[download] 下载完成: %s", filepath.Join(fp, fn))
		}
	}
	return nil
}

func (p *DouYinProcessor) _downloadResource(url, filepath, filename string) (string, error) {
	downloader, err := utils.DownloadFile(url, &utils.DownloadOptions{
		SavePath:  filepath,
		FileName:  filename,
		Timeout:   1200, //  下载超时时间，单位秒
		IgnoreSSL: true,
		ProgressFunc: func(progress *utils.DownloadProgress) {
			utils.DebugWithFormat("[download] 下载进度: %d/%d - %.2f%%", progress.Downloaded, progress.TotalBytes, progress.FormattedSpeed)
		},
		MaxRetries: 2,
		ChunkSize:  10,
	})
	// 启动下载
	if err = downloader.Start(); err != nil {
		utils.ErrorWithFormat("[downloader] 下载失败: %v", err)
		return "", err
	}
	for {
		progress := downloader.GetProgress()
		if progress.Status == utils.StatusCompleted || progress.Status == utils.StatusFailed {
			if downloader.GetProgress().Status == utils.StatusCompleted {
				return "", nil
			} else {
				return "", errors.New(fmt.Sprintf("[downloader] 下载失败: %s", url))
			}
		}
		time.Sleep(2 * time.Second)
	}
}

// _extractURLParams 从 URL 中提取查询参数，并返回一个键值对映射。
func (p *DouYinProcessor) _extractURLParams(rawURL string) (map[string]string, error) {
	// 解析 URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	// 提取查询参数
	queryParams := parsedURL.Query()
	params := make(map[string]string)
	// 将查询参数转换为 map
	for key, values := range queryParams {
		// 如果参数有多个值，只取第一个值
		if len(values) > 0 {
			params[key] = values[0]
		}
	}
	return params, nil
}

func (p *DouYinProcessor) Tidy() error {
	files, err := os.ReadDir(p.tempDir)
	if err != nil {
		return fmt.Errorf("读取临时目录失败: %w", err)
	}
	if len(files) == 0 {
		utils.WarnWithFormat("[DouYinVideo] ⚠️ 未找到待整理的资源文件")
		return errors.New("未找到待整理的资源文件")
	}

	switch p.cfg.Tidy.Mode {
	case 1:
		return p.tidyToLocal(files)
	case 2:
		return p.tidyToWebDAV(files, core.GlobalWebDAV)
	default:
		return fmt.Errorf("未知整理模式: %d", p.cfg.Tidy.Mode)
	}
}

// 整理到本地
func (p *DouYinProcessor) tidyToLocal(files []os.DirEntry) error {
	dstDir := p.cfg.Tidy.DistDir
	if dstDir == "" {
		_ = processor.RemoveTempDir(p.tempDir)
		return errors.New("未配置输出目录")
	}
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		_ = processor.RemoveTempDir(p.tempDir)
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	for _, f := range files {
		src := filepath.Join(p.tempDir, f.Name())
		dst := filepath.Join(dstDir, "douyin", utils.SanitizeFileName(f.Name()))
		if err := os.Rename(src, dst); err != nil {
			utils.WarnWithFormat("[DouYinVideo] ⚠️ 移动失败 %s → %s: %v", src, dst, err)
			continue
		}
		utils.InfoWithFormat("[DouYinVideo] 📦 已整理: %s", dst)
	}
	//清除临时目录
	err := processor.RemoveTempDir(p.tempDir)
	if err != nil {
		utils.WarnWithFormat("[DouYinVideo] ⚠️ 删除临时目录失败: %s (%v)", p.tempDir, err)
		return err
	}
	utils.DebugWithFormat("[DouYinVideo] 🧹 已删除临时目录: %s", p.tempDir)
	return nil
}

// 整理到webdav
func (p *DouYinProcessor) tidyToWebDAV(files []os.DirEntry, webdav *core.WebDAV) error {
	if webdav == nil {
		_ = processor.RemoveTempDir(p.tempDir)
		return errors.New("WebDAV 未初始化")
	}

	for _, f := range files {
		filePath := filepath.Join(p.tempDir, "douyin", f.Name())
		if err := webdav.Upload(filePath); err != nil {
			utils.WarnWithFormat("[DouYinVideo] ☁️ 上传失败 %s: %v", f.Name(), err)
			continue
		}
		utils.InfoWithFormat("[DouYinVideo] ☁️ 已上传: %s", f.Name())
	}
	//清除临时目录
	err := processor.RemoveTempDir(p.tempDir)
	if err != nil {
		utils.WarnWithFormat("[DouYinVideo] ⚠️ 删除临时目录失败: %s (%v)", p.tempDir, err)
		return err
	}
	utils.DebugWithFormat("[DouYinVideo] 🧹 已删除临时目录: %s", p.tempDir)
	return nil
}
