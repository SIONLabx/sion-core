<!-- 协作方式与架构约定。 -->

# 协作方式

本仓库用于手写 Go 模型调用与 Agent Runtime SDK，面向 SION 与其他调用方；参考 Pi 理解统一调用、协议适配、流式结果与工具循环。已按用户明确要求创建目录、首行职责注释、package 声明与 go.mod，当前仅有源码骨架，模型调用和 Agent Runtime 功能均待实现。

- 用户负责手写 Go 代码。默认协作内容是解释设计、阅读 Pi、维护知识库和审阅用户的实现。
- 除非用户明确要求代写或修改具体实现，不主动生成 Go 源码、测试代码、代码骨架、go.mod 或 go.sum，也不运行会生成这些文件的初始化命令。
- 用户明确授权某一项实现后，在该项范围内执行，不重复请求确认。
- 知识库位于当前工作区的 ../sion-docs/sion-core/；其 README.md 记录目标、学习顺序和进度。
- 设计和讲解应说明：要解决的问题、对应的 Pi 源码、Go 中的组织方式、用户的手写任务和验收标准。
- 区分 Pi 的现有行为与 sion-core 的设计选择；引用 Pi 时记录源码版本。
- 记录协议格式时，使用完整、严格可解析的 JSON 请求体和具体数据，不用伪代码或占位内容代替；注明请求地址、鉴权方式、适用版本、官方来源与核对日期，并比较同一场景的协议差异。
- 协议请求示例保存在 ../sion-docs/sion-core/protocol-request-examples.md 及 examples/requests/；内嵌 JSON 与独立文件保持一致。区分文档核对、JSON 解析和真实接口调用，未经调用验证不宣称请求已成功。
- 不把尚未完成的练习或尚未验证的理解标为“已掌握”。

# 结构约定

- 当前推荐目录与职责见 ../sion-docs/sion-core/01-directory-structure.md；用户在讨论中调整设计时同步维护文档。
- 一个 Go module，顶层功能目录为 ai/ 与 agent/。ai package 提供 Client、公共模型调用方法、入口预设和内置协议装配；根目录只保留模块配置和文档，不再另设 sioncore 入口 package。
- 首版覆盖文本对话、文字加图片输入的图片理解，以及图片生成；音频和视频暂不接入。
- ai/chat/ 定义对话消息、工具声明、参数、请求、响应与流契约；ai/images/ 定义图片生成的请求、结果与流契约。图片理解复用 chat。
- agent/ 是 SDK 的运行时模块，调用方可选择是否使用；它负责模型调用循环、注册工具的执行调度、运行内状态、事件、停止条件与逐次/累计 usage。Client 的单次调用语义保持独立。
- agent/ 通过自己定义的最小模型调用接口接收 Client，依赖 ai/chat、ai/usage 等公共契约，不导入 ai 入口包或具体协议实现；ai 及其所有子 package 都不依赖 agent。
- chat.Tool 表达模型可见的声明，agent.Tool 在其上绑定调用方提供的执行函数；Runtime 负责调用函数与回传结果，业务实现由调用方提供。
- User 消息规范化为有序且带类型的 Text/Image 内容块；纯文本构造生成单个 Text 块。保留消息与块的顺序，由 adapter 按角色、协议和目标能力编码，不静默丢弃图片。
- ai/content/ 保存共享图片表示，ai/usage/ 保存共享用量契约；chat 与 images 可依赖它们，彼此独立，不依赖 provider、api 或 ai 入口包。
- ai/provider/ 取代旧 catalog/ 的连接与模型描述职责，管理入口、按操作选择的路由、可选模型配置和目标解析；不依赖 chat、images、api 或 ai 入口包。
- 模型 ID 原样传递，不根据前缀猜测厂商、协议或能力，不要求模型先注册。模型能力按入口、模型与操作声明，区分支持、不支持和未知。
- ai/api/ 定义按操作划分的公开协议接口、流式扩展接口、调用错误与 HTTP 执行契约；可依赖 chat、images、provider，不导入具体实现。
- 首版协议目标包括 OpenAI Chat Completions、OpenAI Responses、Anthropic Messages、OpenAI Images 和方舟图片生成。兼容范围以本期 SDK 操作为准，不等于实现厂商全部端点。
- 具体实现放在 ai/api/openaichat、ai/api/openairesponses、ai/api/anthropic、ai/api/openaiimages、ai/api/arkimages 等子 package；不导入 ai 入口包或其他具体 adapter。
- 协议与供应商入口分别配置；第三方入口可复用协议实现。不同协议的兼容选项由各自 compat.go 解释，不根据模型名自动选择或失败后自动切换协议。
- ai 入口包默认装配内置 adapter，也允许调用方显式注入自定义实现。第三方需要的契约必须可以公开导入。
- 协议专有兼容类型和转换规则归对应 adapter；provider 只承载配置材料，不导入具体协议 package。
- usage 保留上游缺失状态与计数包含关系；累计事件不重复求和。图片结果保留逐项成功、失败和整体状态。
- ai/internal/ 保存 AI 协议之间确实共享的私有实现，agent 不导入它；流的取消、错误终止与资源关闭从第一条流式路径开始定义。
- Agent Runtime 负责通用循环与运行内状态；工具业务逻辑、长期会话存储、上下文裁剪/压缩策略、配置文件读取、凭证查找和提示词内容由应用提供或注入。具体运行时规划见 ../sion-docs/sion-core/agent-runtime.md。
