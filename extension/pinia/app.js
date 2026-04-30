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

    const normalizeProjectPaths = (paths = [], currentRootPath = '') => {
        const normalized = Array.isArray(paths)
            ? paths.map((path) => String(path).trim()).filter(Boolean)
            : []
        if (currentRootPath) {
            normalized.unshift(currentRootPath)
        }
        return Array.from(new Set(normalized))
    }

    const normalizeFileSelectedByRootPath = (selectedMap = {}) => {
        if (!selectedMap || typeof selectedMap !== 'object' || Array.isArray(selectedMap)) {
            return {}
        }
        return Object.fromEntries(
            Object.entries(selectedMap).filter(([, value]) => Array.isArray(value)),
        )
    }

    const createAskToAIFallback = () => ({
        globalPrompt: '',
        globalPrompts: [],
        currentQuestion: '',
        questionHistory: [],
    })

    const normalizeGlobalPrompts = (items, legacyPrompt = '') => {
        const normalizedItems = Array.isArray(items)
            ? items
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
            : []

        if (normalizedItems.length) {
            return normalizedItems
        }

        return legacyPrompt
            ? [{ value: legacyPrompt, enabled: true }]
            : []
    }

    const normalizeAskToAI = (askToAI = {}) => {
        const globalPrompts = normalizeGlobalPrompts(askToAI.globalPrompts, askToAI.globalPrompt ?? '')
        return {
            globalPrompt: globalPrompts[0]?.value ?? askToAI.globalPrompt ?? '',
            globalPrompts,
            currentQuestion: askToAI.currentQuestion ?? '',
            questionHistory: Array.isArray(askToAI.questionHistory) ? askToAI.questionHistory.slice(0, 3) : [],
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
        projectPaths: [],
        fileSelected: [],
        fileSelectedByRootPath: {},
        askToAI: createAskToAIFallback(),
    })
    const initAppState=async ()=>{
        appState.fileTreeRequestStorage=storage.defineItem(`local:fileTreeRequest`, {
            fallback: createFileTreeRequestFallback()
        })
        appState.saveFileTreeRequest=async (request = appState.fileTreeRequest) => {
            await appState.fileTreeRequestStorage.setValue(toPlain(request))
        }
        appState.fileTreeRequest=normalizeFileTreeRequest(await appState.fileTreeRequestStorage.getValue())

        appState.projectPathsStorage=storage.defineItem(`local:projectPaths`, {
            fallback: []
        })
        appState.projectPaths=normalizeProjectPaths(
            await appState.projectPathsStorage.getValue(),
            appState.fileTreeRequest.rootPath,
        )
        appState.saveProjectPaths=async (paths = appState.projectPaths) => {
            appState.projectPaths = normalizeProjectPaths(paths, appState.fileTreeRequest.rootPath)
            await appState.projectPathsStorage.setValue(toPlain(appState.projectPaths))
        }

        appState.fileSelectedStorage=storage.defineItem(`local:fileSelected`, {
            fallback: []
        })
        const storedFileSelected=await appState.fileSelectedStorage.getValue()
        appState.fileSelectedByRootPathStorage=storage.defineItem(`local:fileSelectedByRootPath`, {
            fallback: {}
        })
        appState.fileSelectedByRootPath=normalizeFileSelectedByRootPath(
            await appState.fileSelectedByRootPathStorage.getValue(),
        )
        if (
            appState.fileTreeRequest.rootPath
            && !appState.fileSelectedByRootPath[appState.fileTreeRequest.rootPath]
            && Array.isArray(storedFileSelected)
        ) {
            appState.fileSelectedByRootPath[appState.fileTreeRequest.rootPath] = storedFileSelected
        }
        appState.fileSelected=Array.isArray(appState.fileSelectedByRootPath[appState.fileTreeRequest.rootPath])
            ? appState.fileSelectedByRootPath[appState.fileTreeRequest.rootPath]
            : []
        appState.saveFileSelected=async (fileSelected = appState.fileSelected) => {
            appState.fileSelected = fileSelected
            appState.fileSelectedByRootPath[appState.fileTreeRequest.rootPath] = toPlain(fileSelected)
            await appState.fileSelectedByRootPathStorage.setValue(toPlain(appState.fileSelectedByRootPath))
            await appState.fileSelectedStorage.setValue(toPlain(fileSelected))
        }
        appState.switchProject=async (rootPath) => {
            const nextRootPath = String(rootPath ?? '').trim()
            if (!nextRootPath) {
                return
            }
            await appState.saveFileSelected()
            appState.fileTreeRequest.rootPath = nextRootPath
            if (!appState.projectPaths.includes(nextRootPath)) {
                appState.projectPaths.push(nextRootPath)
                await appState.saveProjectPaths()
            }
            appState.fileSelected = Array.isArray(appState.fileSelectedByRootPath[nextRootPath])
                ? appState.fileSelectedByRootPath[nextRootPath]
                : []
            await appState.saveFileTreeRequest()
            await appState.saveFileSelected()
        }

        appState.askToAIStorage=storage.defineItem(`local:askToAI`, {
            fallback: createAskToAIFallback()
        })
        appState.askToAI=normalizeAskToAI(await appState.askToAIStorage.getValue())
        appState.saveAskToAI=async (askToAI = appState.askToAI) => {
            await appState.askToAIStorage.setValue(toPlain(normalizeAskToAI(askToAI)))
        }
    }
    initAppState()
    // 统一导出
    return {
        appState,
        getSelectedIgnoreConfig,
    }
})