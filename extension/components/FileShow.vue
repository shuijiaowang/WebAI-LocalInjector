<script setup>
import { computed, ref } from "vue";
import { usePost } from "@/core/request.js";
import { useAppStore } from "@/pinia/app.js";

const appStore = useAppStore()

const loading = ref(false)
const error = ref("")
const expandedPaths = ref(new Set())

const vIndeterminate = (el, binding) => {
  el.indeterminate = binding.value
}

const fileTree = computed(() => {
  const selected = appStore.appState.fileSelected
  return Array.isArray(selected) ? selected : []
})

const visibleRows = computed(() => {
  const rows = []
  const walk = (nodes = [], depth = 0) => {
    nodes.forEach((node) => {
      rows.push({ node, depth })
      if (node.isDir && node.children?.length && expandedPaths.value.has(node.path)) {
        walk(node.children, depth + 1)
      }
    })
  }
  walk(fileTree.value)
  return rows
})

const createSelectionMap = (nodes = [], map = new Map()) => {
  nodes.forEach((node) => {
    map.set(node.path, node.selected === true)
    if (node.children?.length) {
      createSelectionMap(node.children, map)
    }
  })
  return map
}

const mergeFileTree = (nodes = [], selectedMap = new Map()) => {
  return nodes.map((node) => ({
    name: node.name,
    path: node.path,
    isDir: node.isDir,
    selected: selectedMap.get(node.path) ?? false,
    children: node.children?.length ? mergeFileTree(node.children, selectedMap) : [],
  }))
}

const setNodeSelected = (node, selected) => {
  node.selected = selected
  node.children?.forEach((child) => setNodeSelected(child, selected))
}

const syncDirSelected = (nodes = []) => {
  nodes.forEach((node) => {
    if (!node.children?.length) {
      return
    }
    syncDirSelected(node.children)
    node.selected = node.children.every((child) => child.selected)
  })
}

const hasSelectedChildren = (node) => {
  return node.children?.some((child) => child.selected || hasSelectedChildren(child)) ?? false
}

const isIndeterminate = (node) => {
  if (!node.children?.length) {
    return false
  }
  return !node.selected && hasSelectedChildren(node)
}

const isExpanded = (node) => expandedPaths.value.has(node.path)

const handleExpand = (node) => {
  const nextPaths = new Set(expandedPaths.value)
  if (nextPaths.has(node.path)) {
    nextPaths.delete(node.path)
  } else {
    nextPaths.add(node.path)
  }
  expandedPaths.value = nextPaths
}

const saveFileSelected = async () => {
  if (appStore.appState.saveFileSelected) {
    await appStore.appState.saveFileSelected(fileTree.value)
  }
}

const handleToggle = async (node, event) => {
  setNodeSelected(node, event.target.checked)
  syncDirSelected(fileTree.value)
  await saveFileSelected()
}

const parseResponseData = (data) => {
  if (typeof data === "string") {
    return JSON.parse(data)
  }
  return data
}

const handleUpdate = async () => {
  loading.value = true
  error.value = ""

  const { data, error: requestError } = await usePost(
    "/file/tree",
    appStore.getSelectedIgnoreConfig(),
  )

  loading.value = false

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

    const selectedMap = createSelectionMap(fileTree.value)
    appStore.appState.fileSelected = mergeFileTree(result.data, selectedMap)
    syncDirSelected(appStore.appState.fileSelected)
    expandedPaths.value = new Set(appStore.appState.fileSelected.map((node) => node.path))
    await saveFileSelected()
  } catch (e) {
    error.value = e.message
  }
}
</script>

<template>
  <section class="file-show">
    <div class="toolbar">
      <div>
        <h2>目录更新与勾选</h2>
        <p>选择要合并到本次提问里的文件或目录。</p>
      </div>
      <button type="button" :disabled="loading" @click="handleUpdate">
        {{ loading ? "更新中..." : "更新目录" }}
      </button>

      <span v-if="error" class="error">请求失败：{{ error }}</span>
    </div>

    <p v-if="!visibleRows.length" class="empty">暂无目录数据，请点击更新。</p>

    <div v-else class="tree">
      <div
        v-for="{ node, depth } in visibleRows"
        :key="node.path"
        class="tree-row"
        :style="{ paddingLeft: `${depth * 16}px` }"
      >
        <button
          v-if="node.isDir && node.children?.length"
          type="button"
          class="expand-button"
          :aria-label="isExpanded(node) ? '折叠目录' : '展开目录'"
          @click.prevent.stop="handleExpand(node)"
        >
          {{ isExpanded(node) ? "-" : "+" }}
        </button>
        <span v-else class="expand-placeholder"></span>
        <input
          type="checkbox"
          :checked="node.selected"
          v-indeterminate="isIndeterminate(node)"
          @change="handleToggle(node, $event)"
        />
        <span class="node-type">{{ node.isDir ? "目录" : "文件" }}</span>
        <span>{{ node.name }}</span>
      </div>
    </div>
  </section>
</template>

<style scoped>
.file-show {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.toolbar {
  display: flex;
  gap: 8px;
  align-items: center;
  justify-content: space-between;
}

.toolbar h2 {
  margin: 0;
  font-size: 16px;
}

.toolbar p {
  margin: 4px 0 0;
  color: #666;
  font-size: 12px;
}

.tree {
  display: flex;
  flex-direction: column;
  gap: 4px;
  max-height: 360px;
  overflow: auto;
  padding: 8px;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  background: #fff;
}

.tree-row {
  display: flex;
  gap: 8px;
  align-items: center;
  line-height: 1.6;
}

.expand-button,
.expand-placeholder {
  width: 20px;
  flex: 0 0 20px;
}

.expand-button {
  border: 0;
  background: transparent;
  cursor: pointer;
}

.node-type {
  color: #666;
  font-size: 12px;
}

.empty {
  margin: 0;
  color: #888;
}

.error {
  color: #d93025;
}
</style>