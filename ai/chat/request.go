// 模型 ID、入口选择、消息、工具与参数。

package chat

import (
	"errors"
	"fmt"
	"strings"
)

type Request struct {
	Model    string    `json:"model"`    // 模型ID，非空
	Provider string    `json:"provider"` // 本地配置的入口标识
	Messages []Message `json:"messages"`
}

// 验证请求体
func (r *Request) Validate() error {
	var errs []error
	// 纯空白也不可
	if strings.TrimSpace(r.Model) == "" {
		errs = append(errs, fmt.Errorf("model is required"))
	}
	if len(r.Messages) == 0 {
		errs = append(errs, fmt.Errorf("messages are required"))
	}
	for i, msg := range r.Messages {
		if err := msg.Validate(); err != nil {
			errs = append(errs, fmt.Errorf("messages[%d]: %w", i, err))
		}
	}
	return errors.Join(errs...)
}
