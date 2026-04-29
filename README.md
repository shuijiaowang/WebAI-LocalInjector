# WebAI-LocalInjector
本地go服务提供接口返回代码文本，浏览器插件运行勾选目录然后接收文本直接注入AI网页输入框然后提问

任务拆解，
1.运行一个后端go服务，允许跨域，暴露两个接口。
2.初始化WXT侧边栏实现请求后端并获得数据。
3.初始化WXT注入逻辑。

【任务类型】修改小需求
【核心目标/问题】
初始化go后端，非常简单，
1.无用户，无jwt，无数据库
2.默认监听9006端口，配置文件可以不存在。
3.主要暴露两个接口，用于返回请求。需要传参，具体逻辑暂时不屑，测试阶段返回文本即可
4.一切从简
【明确前提&假设】列清前置判断，不确定项直接提问
【只允许查看】backend
【不要做】不要全项目扫描；不要改无关代码；不要跑全量构建，不要自动执行全量测试；不要循环修复；不要测试我会手动测试的。
【已知信息】报错/现象/复现步骤/已有能力（最多5行，极简表述）
【验收标准】可量化、可验证的成功条件（拒绝“能用就行”）
【输出要求】先给核心结论，再给修改点；代码改动尽量小


【任务类型】新增后端接口和逻辑处理
【核心目标/问题】
当前已完成后端的一个接口，传入参数返回文件目录。
现在需要新建一个vue组件，用于展示和勾选对应的文件，并存储到storage中
1.添加一个配置输入框，可以设置这几个参数配置（全局的,通用的）
```
{
  "rootPath": "D:\\project\\test\\WebAI-LocalInjector",
  "ignoreDirs": [".idea", "node_modules", ".git", ".wxt","extension"],
  "ignoreFiles": ["go.sum", "package-lock.json"],
  "ignoreExts": [".png", ".svg", ".ico"]
}
```
2.添加一个按钮，更新目录的按钮。点击会请求后端数据
3.将后端返回的目录结构进行展示(上述参数过滤后的)
4.可勾选选择需要的Dirs，Files，并同步存储到storage中，之后会错误参数传给后端继续别的请求。
【明确前提&假设】列清前置判断，不确定项直接提问
【只允许查看】backend
【不要做】不要全项目扫描；不要改无关代码；不要跑全量构建，不要自动执行全量测试；不要循环修复；不要测试我会手动测试的。
【已知信息】报错/现象/复现步骤/已有能力（最多5行，极简表述）
【验收标准】可量化、可验证的成功条件（拒绝“能用就行”）
【输出要求】先给核心结论，再给修改点；代码改动尽量小

http://127.0.0.1:9009/api/file/tree
{
  "rootPath": "D:\\project\\test\\WebAI-LocalInjector",
  "ignoreDirs": [".idea", "node_modules", ".git", ".wxt","extension"],
  "ignoreFiles": ["go.sum", "package-lock.json"],
  "ignoreExts": [".png", ".svg", ".ico"]
}


【任务类型】新增后端接口和逻辑处理
【核心目标/问题】
当前已完成后端的一个接口，传入参数返回文件目录。
现在需要新建一个vue组件，用于前端存储参数，不用每次请求都要重写
extension/components/FileRequest.vue
1.修改配置输入框，测试可以存储了，但是现在页面不适合用户编辑，需要修改，例如输入框，标签化，勾选，right/false？
```
{
  "rootPath": "D:\\project\\test\\WebAI-LocalInjector",
  "ignoreDirs": [".idea", "node_modules", ".git", ".wxt","extension"],
  "ignoreFiles": ["go.sum", "package-lock.json"],
  "ignoreExts": [".png", ".svg", ".ico"]
}
```
2.每次修改都要同步存储到storage
3.好像确实，新增的要存储，勾选true才作为参数进行网络请求，没勾选的就是false不用管？//这里可能得修改一下数据结构
【明确前提&假设】列清前置判断，不确定项直接提问
【只允许查看】extension/pinia/app.js，extension/components/FileRequest.vue
【不要做】不要全项目扫描；不要改无关代码；不要跑全量构建，不要自动执行全量测试；不要循环修复；不要测试我会手动测试的。
【已知信息】报错/现象/复现步骤/已有能力（最多5行，极简表述）
【验收标准】可量化、可验证的成功条件（拒绝“能用就行”）
【输出要求】先给核心结论，再给修改点；代码改动尽量小




【任务类型】修改前端vue组件和storage存储
【核心目标/问题】
现在需要新建一个vue组件，用于前端存储参数，不用每次请求都要重写，已经简单初始化。extension/components/FileRequest.vue

1.新增配置输入框，适合用户编辑，需要修改，例如输入框，标签化，勾选，right/false？
如：ignoreDirs：新增
如：复选框，输入框，删除
基本数据如下，但需要给ignoreDirs/ignoreFiles/ignoreExts的每个项对应true/false，如{.idea,true},不要照搬下面的，数据结构你自己考量，仅需要修改
appState.fileTreeRequestStorage=storage.defineItem(`local:fileTreeRequest`, {
            fallback: {
                //示例：
                rootPath: "D:\\project\\test\\WebAI-LocalInjector",
                ignoreDirs: [".idea", "node_modules", ".git", ".wxt","extension"],
                ignoreFiles: ["go.sum", "package-lock.json"],
                ignoreExts: [".png", ".svg", ".ico"]
            }
        })
```
```

2.每次修改都要同步存储到storage，监听change，然后修改。全量查全量写，方法已经给出。

【明确前提&假设】列清前置判断，不确定项直接提问
【只允许查看】extension/pinia/app.js，extension/components/FileRequest.vue
【不要做】不要全项目扫描；不要改无关代码；不要跑全量构建，不要自动执行全量测试；不要循环修复；不要测试我会手动测试的。
【已知信息】报错/现象/复现步骤/已有能力（最多5行，极简表述）
【验收标准】可量化、可验证的成功条件（拒绝“能用就行”）
【输出要求】先给核心结论，再给修改点；代码改动尽量小