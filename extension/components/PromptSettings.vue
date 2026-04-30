<script setup>
import { computed, watch } from "vue";
import { useAppStore } from "@/pinia/app.js";

const appStore = useAppStore()
const askToAI = computed(() => appStore.appState.askToAI)

const savePromptSettings = async () => {
  if (!askToAI.value || !appStore.appState.saveAskToAI) {
    return
  }
  askToAI.value.globalPrompt = askToAI.value.globalPrompts?.[0]?.value ?? ""
  await appStore.appState.saveAskToAI(askToAI.value)
}

const addGlobalPrompt = async () => {
  askToAI.value.globalPrompts.push({ value: "", enabled: true })
  await savePromptSettings()
}

const removeGlobalPrompt = async (index) => {
  askToAI.value.globalPrompts.splice(index, 1)
  await savePromptSettings()
}

watch(
  () => askToAI.value?.globalPrompts,
  savePromptSettings,
  { deep: true },
)
</script>

<template>
  <section v-if="askToAI" class="prompt-settings">
    <div class="section-header">
      <span>全局提示词</span>
      <button type="button" @click="addGlobalPrompt">新增提示词</button>
    </div>

    <div v-if="askToAI.globalPrompts.length" class="prompt-list">
      <div
        v-for="(prompt, index) in askToAI.globalPrompts"
        :key="index"
        class="prompt-item"
      >
        <label class="prompt-enabled">
          <input v-model="prompt.enabled" type="checkbox" />
          <span>启用</span>
        </label>
        <textarea
          v-model="prompt.value"
          rows="3"
          placeholder="例如：代码风格、项目文档、输出要求"
        />
        <button type="button" class="delete-button" @click="removeGlobalPrompt(index)">
          删除
        </button>
      </div>
    </div>

    <p v-else class="empty-tip">暂无全局提示词，点击新增后配置。</p>
  </section>
</template>

<style scoped>
.prompt-settings {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.section-header,
.prompt-enabled {
  display: flex;
  align-items: center;
}

.section-header {
  justify-content: space-between;
  gap: 8px;
}

.section-header span {
  font-weight: 600;
}

.prompt-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.prompt-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 8px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  background: #f8fafc;
}

.prompt-enabled {
  gap: 6px;
  color: #444;
}

textarea {
  width: 100%;
  box-sizing: border-box;
  resize: vertical;
}

.delete-button {
  align-self: flex-end;
}

.empty-tip {
  margin: 0;
  color: #6b7280;
  font-size: 12px;
}
</style>
