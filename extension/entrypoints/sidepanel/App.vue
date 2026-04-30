<script setup>
import FileRequest from "@/components/FileRequest.vue";
import FileShow from "@/components/FileShow.vue";
import AskToAI from "@/components/AskToAI.vue";
import PromptSettings from "@/components/PromptSettings.vue";
</script>

<template>
  <div class="sidepanel-page">
    <header class="page-header">
      <h1>WebAI Local Injector</h1>
      <p>选择文件内容，生成并发送到当前网页。</p>
    </header>

    <main class="panel-stack">
      <section class="card">
        <AskToAI />
      </section>

      <section class="card">
        <FileShow />
      </section>

      <details class="card collapsible">
        <summary>忽略配置</summary>
        <FileRequest />
      </details>

      <details class="card collapsible">
        <summary>全局提示词</summary>
        <PromptSettings />
      </details>
    </main>
  </div>
</template>

<style scoped>
.sidepanel-page {
  min-height: 100vh;
  box-sizing: border-box;
  overflow-x: hidden;
  padding: 14px;
  background: #f3f4f6;
  color: #1f2937;
  font-size: 14px;
}

.page-header {
  margin-bottom: 12px;
}

.page-header h1 {
  margin: 0;
  font-size: 18px;
}

.page-header p {
  margin: 6px 0 0;
  color: #6b7280;
  font-size: 12px;
}

.panel-stack {
  display: flex;
  min-width: 0;
  flex-direction: column;
  gap: 12px;
}

.card {
  min-width: 0;
  box-sizing: border-box;
  padding: 12px;
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  background: #fff;
  box-shadow: 0 1px 2px rgb(15 23 42 / 6%);
}

.collapsible summary {
  cursor: pointer;
  font-weight: 600;
}

.collapsible > :not(summary) {
  margin-top: 12px;
}
</style>
