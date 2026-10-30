// 定义内容块及其类型。

package chat

import "fmt"

// ContentType 表示内容块的类型。
type ContentType string

const (
	ContentText ContentType = "text" // 文本
)

var contentTypeSet = map[ContentType]bool{
	ContentText: true,
}

// ContentBlock 表示一条消息中的一块内容。
type ContentBlock struct {
	Type ContentType // 内容类型，如文本或图片。
	Text string      // 文本块中的文字。text可以为空字符串，但是type不能为空
}

// Validate 验证内容块合法性
func (c *ContentBlock) Validate() error {
	// 如果不在指定范围内的类型，则返回错误
	// TODO：后续可能校验内容是否符合类型要求
	if _, ok := contentTypeSet[c.Type]; !ok {
		return fmt.Errorf("invalid content type: %s", c.Type)
	}
	return nil
}

func NewContentBlock(typ ContentType, text string) (ContentBlock, error) {
	content := ContentBlock{
		Type: typ,
		Text: text,
	}
	if err := content.Validate(); err != nil {
		return ContentBlock{}, err
	}
	return content, nil
}
