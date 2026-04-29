<script setup>
import { ref } from "vue";
import {usePost} from "@/core/request.js";
import FileRequest from "@/components/FileRequest.vue";
import FileShow from "@/components/FileShow.vue";

const result = ref("");
const loading = ref(false);
const error = ref("");

async function handleRequest() {
  loading.value = true;
  // 只需要传接口后半段即可！
  const res = await usePost("/example/test", {
    example:"测试请求"
  });
  if (res.data) {
    result.value=res.data
  }
  if (res.error) {
    error.value = res.error;
  }
  loading.value = false;
}
async function handleToContent() {
  const [tab] = await browser.tabs.query({ active: true, currentWindow: true })
  await browser.tabs.sendMessage(tab.id, {type: "form_sidepanel"});
  console.log("成功发送消息给页面脚本")
}
</script>

<template>
  <FileRequest></FileRequest>
  <FileShow></FileShow>
  <p>sssssssss</p>
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
