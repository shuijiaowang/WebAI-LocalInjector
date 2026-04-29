export async function toTextarea(SELECTOR, INPUT_TEXT) {
    // ================= 配置区 =================
    // const SELECTOR = '.semi-input-textarea.semi-input-textarea-autosize';
    // const INPUT_TEXT = '测试';
    // =========================================

    // 辅助函数：随机延迟
    const sleep = (ms) => new Promise(r => setTimeout(r, ms + Math.random() * 150));

    // 1. 获取元素
    const el = document.querySelector(SELECTOR);
    if (!el) {
        console.error('❌ 未找到元素，请检查选择器');
        return;
    }
    console.log('✅ 找到元素:', el);

    // 2. 模拟真实点击（聚焦）
    const rect = el.getBoundingClientRect();
    const x = rect.left + rect.width / 2;
    const y = rect.top + rect.height / 2;

    el.dispatchEvent(new PointerEvent('pointerdown', { bubbles: true, clientX: x, clientY: y }));
    await sleep(100);
    el.dispatchEvent(new PointerEvent('pointerup', { bubbles: true, clientX: x, clientY: y }));
    await sleep(100);
    el.dispatchEvent(new MouseEvent('click', { bubbles: true, clientX: x, clientY: y }));
    el.focus();
    await sleep(200);

    // 3. 获取原生 setter（绕过框架拦截）
    const nativeInputValueSetter = Object.getOwnPropertyDescriptor(
        window.HTMLTextAreaElement.prototype, 'value'
    )?.set || Object.getOwnPropertyDescriptor(
        window.HTMLInputElement.prototype, 'value'
    )?.set;

    // ==============================================
    // ✅ 关键修复：分段输入 + 逐段触发事件（解决长文本清空）
    // ==============================================
    el.dispatchEvent(new CompositionEvent('compositionstart', { bubbles: true }));
    await sleep(150);

    // 分段写入，每段 200 字符（网站最长能接受的安全长度）
    const chunkSize = 2000;
    const totalChunks = Math.ceil(INPUT_TEXT.length / chunkSize);

    for (let i = 0; i < totalChunks; i++) {
        const chunk = INPUT_TEXT.slice(i * chunkSize, (i + 1) * chunkSize);

        // 分段赋值
        nativeInputValueSetter.call(el, el.value + chunk);

        // 每段都触发更新事件，让框架同步状态
        el.dispatchEvent(new CompositionEvent('compositionupdate', { bubbles: true, data: chunk }));
        el.dispatchEvent(new InputEvent('input', {
            bubbles: true,
            cancelable: true,
            inputType: 'insertText',
            data: chunk
        }));

        // 长文本必须加延迟！否则直接被拦截清空
        await sleep(30 + Math.random() * 70);
    }

    // 结束输入
    el.dispatchEvent(new CompositionEvent('compositionend', { bubbles: true, data: INPUT_TEXT }));
    await sleep(400);

    // ==============================================
    // ✅ 关键修复：强制同步一次最终值（防止框架回滚）
    // ==============================================
    nativeInputValueSetter.call(el, INPUT_TEXT);
    el.dispatchEvent(new InputEvent('input', {
        bubbles: true,
        cancelable: true,
        inputType: 'insertText',
        data: INPUT_TEXT
    }));
    await sleep(500);

    // 4. 模拟按下回车发送
    const enterKeyOpts = {
        bubbles: true,
        cancelable: true,
        key: 'Enter',
        code: 'Enter',
        keyCode: 13,
        which: 13
    };

    el.dispatchEvent(new KeyboardEvent('keydown', enterKeyOpts));
    await sleep(100);
    el.dispatchEvent(new KeyboardEvent('keypress', enterKeyOpts));
    await sleep(100);
    el.dispatchEvent(new KeyboardEvent('keyup', enterKeyOpts));

    // 表单提交
    if (el.form) {
        el.form.dispatchEvent(new Event('submit', { bubbles: true, cancelable: true }));
    }

    console.log('✅ 超长文本输入完成，不会清空！');
}