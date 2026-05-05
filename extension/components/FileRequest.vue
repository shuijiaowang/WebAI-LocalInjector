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
const contentFilters = computed(() => fileTreeRequest.value?.contentFilters ?? [])

const hasChildren = (item) => Boolean(item.children?.length)

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

    <div class="group">
      <div class="group-header">
        <span>内容过滤</span>
      </div>

      <div v-for="item in contentFilters" :key="item.value" class="filter-group">
        <div class="filter-row">
          <label v-if="item.request !== false" class="enable">
            <input v-model="item.enabled" type="checkbox" />
            {{ item.label || item.value }}
          </label>
          <span v-else class="filter-label">{{ item.label || item.value }}</span>
        </div>

        <div v-if="hasChildren(item)" class="filter-children">
          <div v-for="child in item.children" :key="child.value" class="filter-child">
            <label class="enable">
              <input v-model="child.enabled" type="checkbox" />
              {{ child.label || child.value }}
            </label>

            <div v-if="hasChildren(child)" class="filter-children">
              <label v-for="grandChild in child.children" :key="grandChild.value" class="enable filter-child">
                <input v-model="grandChild.enabled" type="checkbox" />
                {{ grandChild.label || grandChild.value }}
              </label>
            </div>
          </div>
        </div>
      </div>

      <p v-if="!contentFilters.length" class="empty">暂无配置</p>
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
.group-header,
.filter-label {
  font-weight: 600;
}

.add-row,
.item-row,
.filter-row {
  display: flex;
  gap: 8px;
  align-items: center;
}

.filter-group,
.filter-child {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.filter-children {
  display: flex;
  flex-direction: column;
  gap: 6px;
  margin-left: 24px;
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