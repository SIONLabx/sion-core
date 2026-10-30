// 可选计数、缓存/推理明细与原始用量。

package usage

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Usage 表示一次模型调用中上游报告的用量。
// 计数字段为 nil 表示未报告，指向 0 表示明确报告为零。
type Usage struct {
	InputTokens      *int64          // 输入 token 总量，包含缓存输入。
	OutputTokens     *int64          // 输出 token 总量，包含思考输出。
	TotalTokens      *int64          // 上游报告的 token 总量，缺失时不自动计算。
	CacheReadTokens  *int64          // 缓存命中的 token 数，已包含在 InputTokens 中。
	CacheWriteTokens *int64          // 写入缓存的 token 数，已包含在 InputTokens 中。
	ReasoningTokens  *int64          // 思考消耗的 token 数，已包含在 OutputTokens 中。
	GeneratedImages  *int64          // 上游报告的生成图片数量，单位为张。
	Raw              json.RawMessage // 原生 usage JSON，保留尚未统一的字段。
}

// Validate 检查已提供的计数是否非负，并汇总所有问题。
// 校验不修改数据；上层决定如何记录异常及处理模型调用结果。
func (u *Usage) Validate() error {
	var errs []error
	if u.InputTokens != nil && *u.InputTokens < 0 {
		errs = append(errs, fmt.Errorf("usage.InputTokens must be >= 0, got %d", *u.InputTokens))
	}
	if u.OutputTokens != nil && *u.OutputTokens < 0 {
		errs = append(errs, fmt.Errorf("usage.OutputTokens must be >= 0, got %d", *u.OutputTokens))
	}
	if u.TotalTokens != nil && *u.TotalTokens < 0 {
		errs = append(errs, fmt.Errorf("usage.TotalTokens must be >= 0, got %d", *u.TotalTokens))
	}
	if u.CacheReadTokens != nil && *u.CacheReadTokens < 0 {
		errs = append(errs, fmt.Errorf("usage.CacheReadTokens must be >= 0, got %d", *u.CacheReadTokens))
	}
	if u.CacheWriteTokens != nil && *u.CacheWriteTokens < 0 {
		errs = append(errs, fmt.Errorf("usage.CacheWriteTokens must be >= 0, got %d", *u.CacheWriteTokens))
	}
	if u.ReasoningTokens != nil && *u.ReasoningTokens < 0 {
		errs = append(errs, fmt.Errorf("usage.ReasoningTokens must be >= 0, got %d", *u.ReasoningTokens))
	}
	if u.GeneratedImages != nil && *u.GeneratedImages < 0 {
		errs = append(errs, fmt.Errorf("usage.GeneratedImages must be >= 0, got %d", *u.GeneratedImages))
	}
	return errors.Join(errs...)
}
