<script setup>
import { computed, ref } from "vue";
import { usePost } from "@/core/request.js";
import { useAppStore } from "@/pinia/app.js";

const appStore = useAppStore()

const contentLoading = ref(false)
const error = ref("")

const fileTree = computed(() => {
  const selected = appStore.appState.fileSelected
  return Array.isArray(selected) ? selected : []
})

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

const handlePrintContent = async () => {
  contentLoading.value = true
  error.value = ""

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

    console.log(result.data)
  } catch (e) {
    error.value = e.message
  }
}
</script>

<template>
  <button type="button" :disabled="contentLoading" @click="handlePrintContent">
    {{ contentLoading ? "输出中..." : "输出内容" }}
  </button>
  <span v-if="error" class="error">请求失败：{{ error }}</span>
</template>

<style scoped>
.error {
  color: #d93025;
}
</style>