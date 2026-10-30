// 定义消息角色与消息结构。

package chat

import (
	"errors"
	"fmt"
)

// Role 表示消息的角色。
type Role string

// 角色
const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleSystem    Role = "system"
)

var roleSet = map[Role]bool{
	RoleUser:      true,
	RoleAssistant: true,
	RoleSystem:    true,
}

// Message 表示一条消息，包含角色和有序的内容块。
type Message struct {
	Role    Role           // 消息角色，如用户或模型。
	Content []ContentBlock // 按顺序保存这条消息的内容块。
}

// Validate 验证消息合法性。
func (m *Message) Validate() error {
	var isEmpty = true
	var errs []error
	if _, ok := roleSet[m.Role]; !ok {
		// 角色不合法
		errs = append(errs, fmt.Errorf("invalid role: %s", m.Role))
	}
	// 验证内容块合法性
	for idx, content := range m.Content {
		if content.Text != "" {
			isEmpty = false
		}
		if err := content.Validate(); err != nil {
			// 内容块不合法
			errs = append(errs, fmt.Errorf("message(role=%s).content[%d]: %w", m.Role, idx, err))
		}
	}

	if len(m.Content) == 0 || isEmpty {
		errs = append(errs, fmt.Errorf("message(role=%s) is empty", m.Role))
	}

	return errors.Join(errs...)
}

func NewMessage(role Role, content ...ContentBlock) (Message, error) {
	msg := Message{
		Role:    role,
		Content: append([]ContentBlock{}, content...),
	}
	// 验证消息合法性
	if err := msg.Validate(); err != nil {
		return Message{}, err
	}
	return msg, nil
}
