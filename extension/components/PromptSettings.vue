<script setup>
import { computed, watch } from "vue";
import { useAppStore } from "@/pinia/app.js";

const appStore = useAppStore()
const askToAI = computed(() => appStore.appState.askToAI)

watch(
  () => askToAI.value?.globalPrompt,
  async () => {
    if (!askToAI.value || !appStore.appState.saveAskToAI) {
      return
    }
    await appStore.appState.saveAskToAI(askToAI.value)
  },
)
</script>

<template>
  <section v-if="askToAI" class="prompt-settings">
    <label class="field">
      <span>全局提示词</span>
      <textarea
        v-model="askToAI.globalPrompt"
        rows="4"
        placeholder="例如：代码风格、项目文档、输出要求"
      />
    </label>
  </section>
</template>

<style scoped>
.prompt-settings,
.field {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.field span {
  font-weight: 600;
}

textarea {
  width: 100%;
  box-sizing: border-box;
  resize: vertical;
}
</style>
