const BASE_URL = "http://localhost:9009/api";

/**
 * 通用 POST 请求工具（带基础路径封装）
 * @param {string} url - 接口路径（如 /example/test），自动拼接 BASE_URL
 * @param {object} data - 请求体数据
 * @returns {Promise<{ data: any, error: string | null }>}
 */
export async function usePost(url, data = {}) {
    let result = null;
    let error = null;
    try {
        // 自动拼接基础地址，无需手动写
        const fullUrl = BASE_URL + url;
        const res = await fetch(fullUrl, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(data),
        });
        // 处理 HTTP 错误状态码
        if (!res.ok) throw new Error(`请求错误：${res.status}`);
        result = await res.json();
        result = JSON.stringify(result, null, 2);
    } catch (e) {
        error = e.message;
        console.error("请求异常：", e);
    }
    return { data: result, error };
}