
import {defineContentScript} from "#imports";
import {init} from "../../core/init.js";

// 定义内容脚本：仅注入到B站视频播放页
export default defineContentScript({
    // 精准匹配B站视频页（支持通配符）
    matches: ['https://chatglm.cn/main/alltoolsdetail/*'],
    // 可选：B站是SPA，监听路由变化确保页面切换后仍生效
    runAt: 'document_idle',
    allFrames: false,

    // 脚本注入后执行的核心逻辑
    main() {

        //监听发来的消息，
        browser.runtime.onMessage.addListener(async (message) => {
            //监听插件是否启用
            if (message.type === 'form_sidepanel') {
                console.log("来自sidepanel的消息",message)
                console.log("插件初始化")
                const textarea = document.querySelector(".input-box-inner textarea")
                console.log("找到输入框",textarea)
            }
        });
    },
});
