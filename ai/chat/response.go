// 完整结果、usage、元信息与转换诊断。

package chat

import "github.com/SIONLabx/sion-core/ai/usage"

type Response struct {
	ID              string       `json:"id"`                //响应ID，可空
	RequestID       string       `json:"request_id"`        //请求ID，可空
	Model           string       `json:"model"`             //上游实际返回的模型 ID，未报告时为空
	Message         *Message     `json:"message"`           //响应的消息，可空
	FinishReason    FinishReason `json:"finish_reason"`     //兼容不同模型的完成原因，可空
	RawFinishReason string       `json:"raw_finish_reason"` //原始模型完成原因，可空
	Usage           *usage.Usage `json:"usage"`             //用量，可空
	Diagnostics     []Diagnostic `json:"diagnostics"`       //不阻断正常结果的提示，例如某个 usage 计数不可用；nil 或空列表表示没有这类提示
}

type FinishReason string

const (
	FinishReasonStop          FinishReason = "stop"           // 正常结束
	FinishReasonContentFilter FinishReason = "content_filter" // 内容过滤
	FinishReasonToolCalls     FinishReason = "tool_calls"     // 工具调用
	FinishReasonLength        FinishReason = "length"         // 输出内容超过长度限制，因此被截断
	FinishReasonUnknown       FinishReason = "unknown"        // 未知原因
)

type Diagnostic struct {
	Code   string `json:"code"`   // 错误码
	Field  string `json:"field"`  // 字段名
	Reason string `json:"reason"` // 原因
}
