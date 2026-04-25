<script setup>
import HelloWorld from "@/components/HelloWorld.vue";
import { ref } from "vue";

const result = ref("");
const loading = ref(false);
const error = ref("");

async function handleRequest() {
  loading.value = true;
  error.value = "";
  result.value = "";

  try {
    const res = await fetch("http://localhost:9009/api/example/test", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({}),
    });
    const data = await res.json();
    result.value = JSON.stringify(data, null, 2);
  } catch (e) {
    error.value = e.message;
  } finally {
    loading.value = false;
  }
}

async function handleToContent() {
  const [tab] = await browser.tabs.query({ active: true, currentWindow: true })
  await browser.tabs.sendMessage(tab.id, {type: "form_sidepanel"});
  console.log("成功发送消息给页面脚本")
}
</script>

<template>
  <p>sssssssss</p>
  <HelloWorld></HelloWorld>

  <div style="margin-top: 16px;">
    <button @click="handleRequest" :disabled="loading">
      {{ loading ? "请求中..." : "发送 POST 请求" }}
    </button>

    <p v-if="error" style="color: red;">请求失败：{{ error }}</p>

    <pre v-if="result" style="background: #f5f5f5; padding: 12px; border-radius: 4px; white-space: pre-wrap;">{{ result }}</pre>

  </div>

  <button @click="handleToContent" >
    call content
  </button>
</template>

<style scoped></style>
