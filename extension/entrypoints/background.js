import {storage} from "#imports";
export default defineBackground(async () => {
    console.log('Background script started', {id: browser.runtime.id});
    const currentValue = await counter.getValue();
    await counter.setValue(currentValue + 1);


    // 点击扩展图标时，打开侧边栏
    browser.action.onClicked.addListener(async (tab) => {
        await browser.sidePanel.open({ tabId: tab.id });
    });

    // 全局配置：所有标签页都启用侧边栏
    await browser.sidePanel.setPanelBehavior({openPanelOnActionClick: true});
});
//存储
const counter = storage.defineItem('local:counter', {
    fallback: 0,
});




