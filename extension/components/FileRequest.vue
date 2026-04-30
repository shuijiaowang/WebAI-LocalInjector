<script setup>
import { computed, reactive, watch } from "vue";
import {useAppStore} from "@/pinia/app.js";
const appStore=useAppStore()

const itemGroups = [
  { key: "ignoreDirs", label: "忽略目录", placeholder: "例如 node_modules" },
  { key: "ignoreFiles", label: "忽略文件", placeholder: "例如 package-lock.json" },
  { key: "ignoreExts", label: "忽略扩展名", placeholder: "例如 .png" },
]

const newItems = reactive({
  ignoreDirs: "",
  ignoreFiles: "",
  ignoreExts: "",
})

const fileTreeRequest = computed(() => appStore.appState.fileTreeRequest)

const addItem = (key) => {
  const value = newItems[key].trim()
  if (!value) {
    return
  }
  fileTreeRequest.value[key].push({ value, enabled: true })
  newItems[key] = ""
}

const removeItem = (key, index) => {
  fileTreeRequest.value[key].splice(index, 1)
}

watch(
  () => appStore.appState.fileTreeRequest,
  async (request) => {
    if (!request || !appStore.appState.saveFileTreeRequest) {
      return
    }
    await appStore.appState.saveFileTreeRequest(request)
  },
  { deep: true },
)
</script>

<template>
  <section v-if="fileTreeRequest" class="file-request">
    <div v-for="group in itemGroups" :key="group.key" class="group">
      <div class="group-header">
        <span>{{ group.label }}</span>
      </div>

      <div class="add-row">
        <input
          v-model="newItems[group.key]"
          type="text"
          :placeholder="group.placeholder"
          @keyup.enter="addItem(group.key)"
        />
        <button type="button" @click="addItem(group.key)">新增</button>
      </div>

      <div v-for="(item, index) in fileTreeRequest[group.key]" :key="`${group.key}-${index}`" class="item-row">
        <label class="enable">
          <input v-model="item.enabled" type="checkbox" />
          启用
        </label>
        <input v-model="item.value" type="text" />
        <button type="button" @click="removeItem(group.key, index)">删除</button>
      </div>

      <p v-if="!fileTreeRequest[group.key]?.length" class="empty">暂无配置</p>
    </div>
  </section>
</template>

<style scoped>
.file-request {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.field,
.group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.field span,
.group-header {
  font-weight: 600;
}

.add-row,
.item-row {
  display: flex;
  gap: 8px;
  align-items: center;
}

input[type="text"] {
  flex: 1;
  min-width: 0;
}

.enable {
  display: flex;
  gap: 4px;
  align-items: center;
  white-space: nowrap;
}

.empty {
  margin: 0;
  color: #888;
}
</style>