setInterval(async () => {
    const sleep = (ms) => new Promise(r => setTimeout(r, ms + Math.random() * 150));
    async function toClick(SELECTOR) {
        // ================= 配置区 =================
        // const SELECTOR = '.semi-input-textarea.semi-input-textarea-autosize';
        // const INPUT_TEXT = '测试';
        // =========================================
        // SELECTOR='[data-empty-conversation="false"] div .relative button div'
        // 辅助函数：随机延迟

        // 1. 获取元素
        const el = document.querySelectorAll(SELECTOR)[2];
        if (!el) {
            console.error('❌ 未找到元素，请检查选择器');
            return;
        }
        console.log('✅ 找到元素:', el);

        // 2. 模拟真实点击（聚焦）
        const rect = el.getBoundingClientRect();
        const x = rect.left + rect.width / 2 + 20;
        const y = rect.top + rect.height / 2 + 20;

        el.dispatchEvent(new PointerEvent('pointerdown', {bubbles: true, clientX: x, clientY: y}));
        await sleep(100);
        el.dispatchEvent(new PointerEvent('pointerup', {bubbles: true, clientX: x, clientY: y}));
        await sleep(100);
        // el.dispatchEvent(new MouseEvent('click', { bubbles: true, clientX: x, clientY: y }));
        // await sleep(200);
    }

    await toClick('[data-empty-conversation="false"] .relative button div')
    await sleep(0);
    document.querySelector('[role="menu"] .text-dbx-function-danger').click()
    await sleep(0);
    document.querySelector('[data-slot="dialog-footer"] button:last-child').click()
    await sleep(0);

}, 1000)