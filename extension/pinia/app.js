import { defineStore } from 'pinia'
import { reactive } from 'vue'
import {toPlain} from "@/core/toPlain.js";
export const useAppStore = defineStore('app', () => {

    const createFileTreeRequestFallback = () => ({
        rootPath: "D:\\project\\test\\WebAI-LocalInjector",
        ignoreDirs: [
            { value: ".idea", enabled: true },
            { value: "node_modules", enabled: true },
            { value: ".git", enabled: true },
            { value: ".wxt", enabled: true },
            { value: "extension", enabled: true },
        ],
        ignoreFiles: [
            { value: "go.sum", enabled: true },
            { value: "package-lock.json", enabled: true },
        ],
        ignoreExts: [
            { value: ".png", enabled: true },
            { value: ".svg", enabled: true },
            { value: ".ico", enabled: true },
        ],
    })

    const normalizeIgnoreItems = (items = []) => items
        .map((item) => {
            if (typeof item === 'string') {
                return { value: item, enabled: true }
            }
            return {
                value: item?.value ?? '',
                enabled: item?.enabled !== false,
            }
        })
        .filter((item) => item.value)

    const normalizeFileTreeRequest = (request = {}) => {
        const fallback = createFileTreeRequestFallback()
        return {
            rootPath: request.rootPath ?? fallback.rootPath,
            ignoreDirs: normalizeIgnoreItems(request.ignoreDirs ?? fallback.ignoreDirs),
            ignoreFiles: normalizeIgnoreItems(request.ignoreFiles ?? fallback.ignoreFiles),
            ignoreExts: normalizeIgnoreItems(request.ignoreExts ?? fallback.ignoreExts),
        }
    }
    //用于提取fileTreeRequest中选中true的项作为请求参数请求后端
    const getSelectedIgnoreConfig = () => {
        // 解构当前文件树配置
        const { rootPath, ignoreDirs, ignoreFiles, ignoreExts } = appState.fileTreeRequest;

        // 通用过滤函数：只保留 enabled=true 的项，提取 value
        const filterEnabledItems = (itemList) => {
            return itemList
                .filter(item => item.enabled) // 过滤启用的项
                .map(item => item.value);     // 只保留值
        };

        // 返回你需要的格式
        return {
            rootPath: rootPath,
            ignoreDirs: filterEnabledItems(ignoreDirs),
            ignoreFiles: filterEnabledItems(ignoreFiles),
            ignoreExts: filterEnabledItems(ignoreExts)
        };
    };

    //全局状态
    const appState =reactive({
        fileTreeRequest: createFileTreeRequestFallback(),
        fileSelected: [],
    })
    const initAppState=async ()=>{
        appState.fileTreeRequestStorage=storage.defineItem(`local:fileTreeRequest`, {
            fallback: createFileTreeRequestFallback()
        })
        appState.saveFileTreeRequest=async (request = appState.fileTreeRequest) => {
            await appState.fileTreeRequestStorage.setValue(toPlain(request))
        }
        appState.fileTreeRequest=normalizeFileTreeRequest(await appState.fileTreeRequestStorage.getValue())

        appState.fileSelectedStorage=storage.defineItem(`local:fileSelected`, {
            fallback: []
        })
        const storedFileSelected=await appState.fileSelectedStorage.getValue()
        appState.fileSelected=Array.isArray(storedFileSelected) ? storedFileSelected : []
        appState.saveFileSelected=async (fileSelected = appState.fileSelected) => {
            await appState.fileSelectedStorage.setValue(toPlain(fileSelected))
        }
    }
    initAppState()
    // 统一导出
    return {
        appState,
        getSelectedIgnoreConfig,
    }
})