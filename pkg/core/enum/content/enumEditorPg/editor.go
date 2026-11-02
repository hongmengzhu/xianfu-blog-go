package enumEditorPg

import (
	"strings"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// Editor 编辑器
type Editor string

const (
	Html     Editor = "html"     //html
	Text     Editor = "text"     //文本
	RichText Editor = "richText" //富文本
	Markdown Editor = "markdown" //markdown
)

// Name 名称
func (this Editor) Name() string {
	switch this {
	case "html":
		return "html"
	case "text":
		return "文本"
	case "richText":
		return "富文本"
	case "markdown":
		return "markdown"
	default:
		return "未知"
	}
}

// 值
func (this Editor) String() string {
	return string(this)
}

// 值
func (this Editor) Index() string {
	return string(this)
}

// IsEqual 值是否相等
func (this Editor) IsEqual(id string) bool {
	return string(this) == id
}

// Find 字符串是否存在
//
//	@Description:
//	@receiver t
//	@param str
//	@return bool
func (t Editor) Find(str string) bool {
	if len(str) <= 0 {
		return false
	}
	return strings.Contains(str, t.String())
}

var EditorMap = map[string]enumBasePg.EnumString{
	Html.String():     enumBasePg.EnumString{Html.String(), Html.Name()},
	Text.String():     enumBasePg.EnumString{Text.String(), Text.Name()},
	RichText.String(): enumBasePg.EnumString{RichText.String(), RichText.Name()},
	Markdown.String(): enumBasePg.EnumString{Markdown.String(), Markdown.Name()},
}

func IsExistEditor(id string) bool {
	_, ok := EditorMap[id]
	return ok
}
