<script setup>
import { computed, ref, watch } from "vue";
import { usePost } from "@/core/request.js";
import { useAppStore } from "@/pinia/app.js";

const appStore = useAppStore()

const contentLoading = ref(false)
const error = ref("")
const copyStatus = ref("")
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

const buildFullPrompt = (content) => {
  const parts = []
  const globalPrompt = askToAI.value.globalPrompt.trim()
  const currentQuestion = askToAI.value.currentQuestion.trim()

  if (globalPrompt) {
    parts.push(`【全局提示词】\n${globalPrompt}`)
  }
  if (currentQuestion) {
    parts.push(`【本次提问需求】\n${currentQuestion}`)
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

const handlePrintContent = async () => {
  contentLoading.value = true
  error.value = ""
  copyStatus.value = ""

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
  } catch (e) {
    error.value = e.message
  }
}

watch(
  () => appStore.appState.askToAI,
  async (value) => {
    if (!value || !appStore.appState.saveAskToAI) {
      return
    }
    await appStore.appState.saveAskToAI(value)
  },
  { deep: true },
)
</script>

<template>
  <section class="ask-to-ai">
    <label class="field">
      <span>全局提示词</span>
      <textarea
        v-model="askToAI.globalPrompt"
        rows="4"
        placeholder="例如：代码风格、项目文档、输出要求"
      />
    </label>

    <label class="field">
      <span>本次提问需求</span>
      <textarea
        v-model="askToAI.currentQuestion"
        rows="3"
        placeholder="请输入本次要问 AI 的具体需求"
      />
    </label>

    <div v-if="askToAI.questionHistory.length" class="history">
      <div class="history-title">近三次提问需求</div>
      <button
        v-for="question in askToAI.questionHistory"
        :key="question"
        type="button"
        class="history-item"
        @click="useHistoryQuestion(question)"
      >
        {{ question }}
      </button>
    </div>

    <div class="toolbar">
      <button type="button" :disabled="contentLoading" @click="handlePrintContent">
        {{ contentLoading ? "输出中..." : "输出内容" }}
      </button>
      <span v-if="error" class="error">请求失败：{{ error }}</span>
    </div>

    <label class="field">
      <span>完整提问文本</span>
      <textarea v-model="fullPrompt" rows="5" readonly placeholder="点击输出内容后生成完整提问文本" />
    </label>

    <div class="toolbar">
      <button type="button" :disabled="!fullPrompt" @click="copyFullPrompt">复制完整提问</button>
      <span v-if="copyStatus" class="copy-status">{{ copyStatus }}</span>
    </div>
  </section>
</template>

<style scoped>
.ask-to-ai {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.field span,
.history-title {
  font-weight: 600;
}

textarea {
  width: 100%;
  box-sizing: border-box;
  resize: vertical;
}

.history {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.history-item {
  overflow: hidden;
  text-align: left;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.toolbar {
  display: flex;
  gap: 8px;
  align-items: center;
}

.error {
  color: #d93025;
}

.copy-status {
  color: #188038;
}
</style>