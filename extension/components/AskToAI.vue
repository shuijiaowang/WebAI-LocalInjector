<script setup>
import { computed, ref, watch } from "vue";
import { usePost } from "@/core/request.js";
import { useAppStore } from "@/pinia/app.js";

const appStore = useAppStore()

const contentLoading = ref(false)
const error = ref("")
const copyStatus = ref("")
const insertStatus = ref("")
const insertError = ref("")
const fullPrompt = ref("")

const fileTree = computed(() => {
  const selected = appStore.appState.fileSelected
  return Array.isArray(selected) ? selected : []
})

const askToAI = computed(() => appStore.appState.askToAI)

const getSelectedPaths = (nodes = [], paths = []) => {
  nodes.forEach((node) => {
    if (node.selected) {
      paths.push(node.path)
    }
    if (node.children?.length) {
      getSelectedPaths(node.children, paths)
    }
  })
  return paths
}

const parseResponseData = (data) => {
  if (typeof data === "string") {
    return JSON.parse(data)
  }
  return data
}

const formatContent = (content) => {
  if (typeof content === "string") {
    return content
  }
  return JSON.stringify(content, null, 2)
}

const getEnabledGlobalPrompts = () => {
  const hasPromptList = Array.isArray(askToAI.value.globalPrompts)
  const globalPrompts = Array.isArray(askToAI.value.globalPrompts)
    ? askToAI.value.globalPrompts
    : []
  const enabledPrompts = globalPrompts
    .filter((prompt) => prompt?.enabled !== false)
    .map((prompt) => (typeof prompt === "string" ? prompt : prompt?.value ?? "").trim())
    .filter(Boolean)

  if (enabledPrompts.length) {
    return enabledPrompts
  }
  if (hasPromptList) {
    return []
  }

  const legacyPrompt = askToAI.value.globalPrompt?.trim()
  return legacyPrompt ? [legacyPrompt] : []
}

const buildFullPrompt = (content) => {
  const parts = []
  const globalPrompts = getEnabledGlobalPrompts()
  const currentQuestion = askToAI.value.currentQuestion.trim()
  if (currentQuestion) {
    parts.push(`【本次提问需求】\n${currentQuestion}`)
  }
  if (globalPrompts.length) {
    parts.push(`【全局提示词】\n${globalPrompts.join("\n\n")}`)
  }

  parts.push(`【项目文件内容】\n${formatContent(content)}`)

  return parts.join("\n\n")
}

const saveAskToAI = async () => {
  if (appStore.appState.saveAskToAI) {
    await appStore.appState.saveAskToAI(askToAI.value)
  }
}

const recordQuestionHistory = async () => {
  const currentQuestion = askToAI.value.currentQuestion.trim()
  if (!currentQuestion) {
    return
  }

  askToAI.value.questionHistory = [
    currentQuestion,
    ...askToAI.value.questionHistory.filter((item) => item !== currentQuestion),
  ].slice(0, 3)
  await saveAskToAI()
}

const useHistoryQuestion = async (question) => {
  askToAI.value.currentQuestion = question
  await saveAskToAI()
}

const copyFullPrompt = async () => {
  if (!fullPrompt.value) {
    return
  }

  try {
    await navigator.clipboard.writeText(fullPrompt.value)
    copyStatus.value = "已复制"
    setTimeout(() => {
      copyStatus.value = ""
    }, 1500)
  } catch (e) {
    error.value = e.message
  }
}

const handleToContent = async (content) => {
  insertStatus.value = ""
  insertError.value = ""

  try {
    const [tab] = await browser.tabs.query({ active: true, currentWindow: true })
    if (!tab?.id) {
      throw new Error("未找到当前标签页")
    }
    await browser.tabs.sendMessage(tab.id, { type: "form_sidepanel", content })
    insertStatus.value = "已发送到页面"
    console.log("成功发送消息给页面脚本")
  } catch (e) {
    insertError.value = e.message
  }
}

const handlePrintContent = async () => {
  contentLoading.value = true
  error.value = ""
  copyStatus.value = ""
  insertStatus.value = ""
  insertError.value = ""

  const { data, error: requestError } = await usePost(
    "/file/content",
    {
      ...appStore.getSelectedIgnoreConfig(),
      selectedPaths: getSelectedPaths(fileTree.value),
    },
  )

  contentLoading.value = false

  if (requestError) {
    error.value = requestError
    return
  }

  try {
    const result = parseResponseData(data)
    if (result?.code !== 0) {
      error.value = result?.msg || "请求失败"
      return
    }

    fullPrompt.value = buildFullPrompt(result.data)
    await recordQuestionHistory()
    await handleToContent(fullPrompt.value)
  } catch (e) {
    error.value = e.message
  }
}

watch(
  () => askToAI.value.currentQuestion,
  async () => {
    if (!askToAI.value || !appStore.appState.saveAskToAI) {
      return
    }
    await appStore.appState.saveAskToAI(askToAI.value)
  },
)
</script>

<template>
  <section class="ask-to-ai">
    <label class="field">
      <span>本次提问需求</span>
      <textarea
        v-model="askToAI.currentQuestion"
        rows="3"
        placeholder="请输入本次要问 AI 的具体需求"
      />
    </label>

    <details v-if="askToAI.questionHistory.length" class="history">
      <summary class="history-title">近十次提问需求</summary>
      <button
        v-for="question in askToAI.questionHistory"
        :key="question"
        type="button"
        class="history-item"
        @click="useHistoryQuestion(question)"
      >
        {{ question }}
      </button>
    </details>

    <div class="toolbar">
      <button type="button" class="primary-action" :disabled="contentLoading" @click="handlePrintContent">
        {{ contentLoading ? "输出中..." : "输出内容" }}
      </button>
      <span v-if="insertStatus" class="copy-status">{{ insertStatus }}</span>
      <span v-if="insertError" class="warn">页面插入失败：{{ insertError }}</span>
      <span v-if="error" class="error">请求失败：{{ error }}</span>
    </div>

    <div class="prompt-result">
      <div class="prompt-meta">
        <strong>完整提问文本</strong>
        <span>{{ fullPrompt ? `已生成 ${fullPrompt.length} 字` : "点击输出内容后生成" }}</span>
      </div>
      <button type="button" :disabled="!fullPrompt" @click="copyFullPrompt">复制完整提问</button>
      <span v-if="copyStatus" class="copy-status">{{ copyStatus }}</span>
    </div>

    <details v-if="fullPrompt" class="prompt-preview">
      <summary>查看完整提问文本</summary>
      <textarea v-model="fullPrompt" rows="5" readonly />
    </details>
  </section>
</template>

<style scoped>
.ask-to-ai {
  display: flex;
  min-width: 0;
  flex-direction: column;
  gap: 12px;
}

.field {
  display: flex;
  min-width: 0;
  flex-direction: column;
  gap: 6px;
}

.field span,
.history-title {
  font-weight: 600;
}

textarea {
  display: block;
  width: 100%;
  max-width: 100%;
  box-sizing: border-box;
  overflow-x: hidden;
  overflow-wrap: anywhere;
  resize: vertical;
  white-space: pre-wrap;
  word-break: break-word;
}

.history {
  display: flex;
  min-width: 0;
  flex-direction: column;
  gap: 6px;
  padding: 8px 10px;
  border: 1px solid #e6ebf2;
  border-radius: 10px;
  background: #f8fbff;
}

.history-item {
  min-width: 0;
  overflow: hidden;
  text-align: left;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.history-title {
  cursor: pointer;
}

.toolbar {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
}

.primary-action {
  border-color: #2f6df6;
  background: linear-gradient(180deg, #3b82f6 0%, #2563eb 100%);
  color: #fff;
  box-shadow: 0 6px 14px rgb(37 99 235 / 22%);
}

.primary-action:hover {
  border-color: #1d4ed8;
  filter: brightness(0.98);
}

.prompt-result {
  display: flex;
  min-width: 0;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
  justify-content: space-between;
  padding: 10px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  background: #f8fafc;
}

.prompt-meta {
  display: flex;
  min-width: 0;
  flex-direction: column;
  gap: 2px;
}

.prompt-meta span {
  overflow: hidden;
  color: #666;
  font-size: 12px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.prompt-preview {
  min-width: 0;
  color: #444;
}

.prompt-preview summary {
  cursor: pointer;
  font-weight: 600;
}

.error {
  color: #d93025;
}

.warn {
  color: #b06000;
}

.copy-status {
  color: #188038;
}
</style>