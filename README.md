<!-- 项目介绍、模块职责与目录入口。 -->

# sion-core

面向 SION 和其他 Go 应用的模型调用与 Agent Runtime SDK。

sion-core 分为两个模块：**ai 负责调用模型，agent 负责使用模型和工具完成多轮运行。** 调用方用 SDK 提供的方法组织 System Prompt、消息、图片、工具和参数，获得统一的结果、事件与 usage；需要自动执行循环时，再组合 Agent Runtime。

当前已按下面的目录创建源码骨架，并初始化 `github.com/SIONLabx/sion-core` 模块，使用 Go 1.25。每个 Go 文件仅包含首行职责注释与 `package` 声明；具体接口签名尚未冻结，模型调用和 Agent Runtime 功能尚未实现。

## 模块概览

| 模块 | 职责 | 适合什么场景 |
|---|---|---|
| **ai** | 多入口、多协议的模型调用；文本与图片、工具声明、推理、流、错误与 usage | 直接请求模型，或给其他运行时提供模型能力 |
| **agent** | 工具执行与模型调用循环、运行状态、事件、限制及用量汇总 | 注册工具后自动完成多轮交互 |
| **examples** | 最小可运行使用示例，随实现补充 | 学习单次调用、流、图片和 Agent 运行 |

一个 Go module，两个顶层功能目录。模型调用入口为 `ai.Client`，运行入口为 `agent.Agent`；项目根目录保存模块配置与文档。普通模型调用可以只使用 ai。

这种介绍与分层参考 Pi 的 [Packages 概览](https://github.com/earendil-works/pi/blob/d5b3de22d0e24a3e9c38a9410d3329b099717283/README.md#packages)、[pi-ai](https://github.com/earendil-works/pi/blob/d5b3de22d0e24a3e9c38a9410d3329b099717283/packages/ai/README.md) 和 [pi-agent-core](https://github.com/earendil-works/pi/blob/d5b3de22d0e24a3e9c38a9410d3329b099717283/packages/agent/README.md)。这里采用 Go package 组织，当前不拆成两个独立发布的 Go module。

## 目录结构

~~~text
sion-core/
├── AGENTS.md                         # 协作方式与架构约定
├── README.md                         # 项目介绍、模块职责与目录入口
├── go.mod                            # 模块路径、Go 版本与依赖声明
│
├── ai/                               # 统一模型调用 SDK，package ai
│   ├── client.go                     # Client 构造、目标解析与调用分派
│   ├── config.go                     # 默认入口、HTTP 客户端与扩展注入
│   ├── presets.go                    # DeepSeek、火山等入口便利配置
│   ├── builtins.go                   # 默认装配内置协议
│   ├── chat.go                       # 单次对话与流式对话调用
│   ├── images.go                     # 图片生成与受支持的图片流调用
│   │
│   ├── chat/                         # 对话公共契约
│   │   ├── message.go                # 消息构造、历史、工具结果与重放信息
│   │   ├── content.go                # 有序的文本、图片、推理与工具内容块
│   │   ├── tool.go                   # 工具名称、描述与参数 Schema
│   │   ├── reasoning.go              # 推理开关、等级、预算及相互约束
│   │   ├── options.go                # 采样、输出限制、工具选择与兼容策略
│   │   ├── request.go                # 模型 ID、入口选择、消息、工具与参数
│   │   ├── response.go               # 完整结果、usage、元信息与转换诊断
│   │   ├── event.go                  # 文本、推理、工具、用量与终止事件
│   │   └── stream.go                 # 消费事件、最终结果、取消与关闭
│   │
│   ├── images/                       # 图片生成公共契约
│   │   ├── request.go                # 模型 ID、入口选择、提示词与参考图
│   │   ├── options.go                # 尺寸、数量、输出格式等选项
│   │   ├── response.go               # 图片集合、逐项状态与 usage
│   │   ├── event.go                  # 图片进展、完成、失败与整体结束
│   │   └── stream.go                 # 受支持的图片流消费与关闭
│   │
│   ├── content/
│   │   └── image.go                  # 共享图片表示：URL、内存数据、MIME
│   ├── usage/
│   │   └── usage.go                  # 可选计数、缓存/推理明细与原始用量
│   │
│   ├── provider/                     # 调用入口与目标配置
│   │   ├── provider.go               # 地址、认证、请求头、按操作配置的路由
│   │   ├── model.go                  # 可选模型描述与模型级覆盖
│   │   ├── capabilities.go           # 支持、不支持、未知及已知限制
│   │   └── target.go                 # 合并配置，得到一次调用的目标
│   │
│   ├── api/                          # 公开协议接口及实现
│   │   ├── chat.go                   # 对话接口与流式扩展接口
│   │   ├── images.go                 # 图片接口与流式扩展接口
│   │   ├── errors.go                 # 调用错误、上游状态与原因
│   │   ├── transport.go              # 可注入的 HTTP 执行契约
│   │   ├── openaichat/               # OpenAI Chat Completions
│   │   ├── openairesponses/          # OpenAI Responses
│   │   ├── anthropic/                # Anthropic Messages
│   │   ├── openaiimages/             # OpenAI Images
│   │   └── arkimages/                # 方舟图片生成
│   │
│   └── internal/                     # 仅供 ai 内部实现使用
│       ├── sse/                      # SSE 帧解码，不解释模型内容
│       │   └── decoder.go            # SSE 帧解码入口
│       └── httpx/                    # 共用 HTTP 执行与资源清理
│           └── client.go             # 共用 HTTP 执行入口
│
├── agent/                            # Agent Runtime，package agent
│   ├── agent.go                      # 构造、运行配置与 Run 入口
│   ├── client.go                     # Runtime 需要的最小模型调用接口
│   ├── loop.go                       # 模型 → 工具 → 模型循环及终止判断
│   ├── tool.go                       # chat.Tool 声明与执行函数绑定
│   ├── executor.go                   # 参数校验、执行调度与结果收集
│   ├── state.go                      # 每次运行独立的消息、回合与工具状态
│   ├── options.go                    # 调用次数、超时、执行及上下文策略
│   ├── event.go                      # Run/Turn/Tool 事件，关联 chat 流事件
│   └── result.go                     # 停止原因、完整/部分结果与逐次/累计用量
│
└── examples/                         # 后续最小可运行使用示例
    ├── chat/                         # 文本调用与 usage
    ├── stream/                       # 消费事件、取消与最终结果
    ├── vision/                       # 文字与图片混合输入
    ├── tools/                        # 单次调用的工具声明与两轮历史
    ├── images/                       # 图片生成与逐项结果
    └── agent/                        # 注册工具后自动完成多轮调用
~~~

上面的源码文件均已创建。测试跟随对应 package，协议样本放在各 adapter 的 testdata；当前测试文件只有职责注释和 package 声明，尚无测试用例。examples 各子目录及 testdata 的 requests / responses 目录用 README.md 记录职责，后续补充可运行示例和测试样本。文件可以随实际实现复杂度合并。

## ai 模型调用

ai 处理一次模型操作：解析目标、校验输入、转换协议、发起请求、解析完整响应或流，并返回统一结果。

| 位置 | 主要职责 |
|---|---|
| `ai` | Client、配置、入口预设、内置协议装配和调用分派 |
| `ai/chat` | 对话请求、消息与工具构造、推理意图、完整响应和流 |
| `ai/images` | 图片生成请求、生成参数、图片集合与逐项结果 |
| `ai/content` | 对话与图片生成共享的图片表示 |
| `ai/usage` | 用量及计数包含关系，保留缺失状态 |
| `ai/provider` | 入口、模型级覆盖、能力与调用目标 |
| `ai/api` | 可供自定义实现使用的协议契约及具体 adapter |
| `ai/internal` | 私有的传输与帧解码实现 |

首期接入目标为 DeepSeek、火山方舟和第三方 Coding Plan。对话协议覆盖 Chat Completions、Responses 与 Anthropic Messages；图片生成覆盖 OpenAI Images 和方舟图片生成。范围为文本和图片。

**图片理解使用 chat，图片生成使用 images。** 同一条 User 消息可以保留文字与图片的交错顺序，adapter 将它编码为目标协议的内容块。生成图片有独立的输入、结果与流语义。

模型 ID 由调用方提供，不要求预先注册，也不根据 `deepseek/...` 等前缀推断入口或协议。一个入口的对话和图片生成可配置不同路由。各项能力仍取决于实际入口与模型。

### 协议实现的职责

各 adapter 在自己的 package 内完成请求、历史和响应的转换。以 Chat Completions 为例：

~~~text
ai/api/openaichat/
├── adapter.go                        # 实现公开协议接口，组织调用
├── compat.go                         # 类型化兼容配置、默认值与校验
├── options.go                        # 推理、预算、采样、缓存和输出限制转换
├── messages.go                       # 角色、图文、工具历史与状态重放
├── tools.go                          # 工具 Schema、选择与约束支持
├── request.go                        # 原生请求、端点、请求头
├── response.go                       # 完整结果、usage、结束原因与原生错误
├── stream.go                         # 原生增量事件与最终结果汇总
├── adapter_test.go                   # 从公开协议接口验证行为
└── testdata/
    ├── requests/                     # 完整请求 JSON 快照
    └── responses/                    # 完整响应、SSE、错误和中断回放
~~~

`compat.go` 定义规则，其他文件应用规则。例如推理开关与强度由 options 转换，工具历史和签名由 messages 处理，推理增量及 usage 由 stream 解析。自定义 adapter 使用公共契约接入 Client。

三个对话 adapter 已创建上述文件骨架。两个图片 adapter 保留 adapter、compat、options、request、response、stream 和 adapter_test 文件；图片参数、逐项结果与受支持的图片流由它们各自处理。

## agent 运行时

agent 将模型调用组织为一次运行：

1. 接收输入、模型目标、工具实现及运行限制。
2. 调用模型，消费响应或流，记录本次结果。
3. 遇到完整、可执行的工具调用，校验参数并执行已注册函数。
4. 追加 Assistant 消息与工具结果，再次调用模型。
5. 得到无需继续执行工具的正常响应，或触发限制、错误、取消时结束。

| 概念 | 说明 |
|---|---|
| Agent | 组合模型调用能力、工具实现与运行配置 |
| Run | 一次完整运行，可以包含多次模型请求和工具调用 |
| 工具声明 | `ai/chat` 中的 `Tool`：模型能看到的名称、描述和参数 Schema |
| 可执行工具 | `agent.Tool`：声明加上应用提供的执行函数 |
| 运行状态 | 本次 Run 的历史、回合与工具状态；不写进共享 Client |
| 运行结果 | 最终或部分结果、停止原因、逐次 usage 与累计已知用量 |

agent 通过自己声明的最小模型调用接口接收 `ai.Client`，复用 `ai/chat`、`ai/usage` 等契约。循环不需要识别模型的原生协议，模型调用层也不依赖 agent。

工具的业务逻辑由应用提供，Runtime 负责实际调度与执行。长期存储、跨进程恢复和上下文压缩策略可通过应用或后续扩展接入。图片生成也可以被应用包装为工具，复用 ai 的 images 操作。

## 调用流程

~~~mermaid
flowchart LR
    App[调用方] -->|单次调用| AI[ai.Client]
    App -->|自动循环| Agent[agent.Agent]
    Agent -->|注入的模型调用接口| AI
    Agent -->|调度执行| Tools[应用注册的工具函数]
    Tools -->|工具结果| Agent
    AI --> Protocol[ai/api 协议实现]
    Protocol --> Endpoint[配置的模型入口]
~~~

这是运行时调用关系。Go 的 import 关系由公共契约和注入接口组织：`ai` 不导入 `agent`，具体 adapter 不反向导入 `ai` 入口 package，`agent` 不导入具体协议实现。

## 设计约定

- **单次调用和自动循环分别使用。** Client 返回一次模型操作的结果，Agent 负责多次调用与工具执行。
- **统一语义，明确反馈差异。** 推理未指定、关闭、开启、等级与预算分别表达；已知不能满足的显式要求报错，允许的调整返回诊断。
- **保留完整结果。** 工具调用、必要的重放信息、usage、结束原因和请求标识随结果返回。
- **缺失用量不补零。** 推理与缓存明细的包含关系明确，流式累计快照不重复求和；Run 保留逐次结果及累计值的完整性。
- **停止和取消有明确含义。** 流中断不当作成功；Runtime 配置调用次数和超时等限制，返回正常完成、限制、取消或错误的区别。
- **配置由调用方提供。** SDK 接收地址、凭证、模型与策略，不绑定应用配置文件或模型清单。
- **模块独立扩展。** 新协议集中在 ai，新工具或运行策略集中在 agent 及调用方。

## 文档

- [目录设计、依赖规则与验收](../sion-docs/sion-core/01-directory-structure.md)
- [Agent Runtime：循环、状态、执行与用量](../sion-docs/sion-core/agent-runtime.md)
- [协议请求格式：完整 JSON 与差异](../sion-docs/sion-core/protocol-request-examples.md)
- [Pi 多协议调用链](../sion-docs/sion-core/pi-protocol-flow.md)
- [Pi 推理参数与兼容清单](../sion-docs/sion-core/pi-compatibility.md)
- [Pi 源码对照表](../sion-docs/sion-core/pi-source-map.md)
- [手写 Go 学习路线](../sion-docs/sion-core/README.md)
