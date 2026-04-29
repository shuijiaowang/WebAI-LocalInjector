// import {defineContentScript} from "#imports"; //可省略
export default defineContentScript({
    matches: ['<all_urls>'], //全匹配
    runAt: 'document_idle',//页面完全加载完成
    // 脚本注入后执行的核心逻辑
    async main() {
        console.log("插件初始化")
    },
});
