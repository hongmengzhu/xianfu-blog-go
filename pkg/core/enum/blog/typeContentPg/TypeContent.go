package typeContentPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// TypeContent 内容类型
type TypeContent string

const (
	ORIGINAL    TypeContent = "original"
	TRANSLATION TypeContent = "translation"
	REPOST      TypeContent = "repost"
)

func (this TypeContent) Name() string {
	switch this {
	case "original":
		return "原创"
	case "translation":
		return "翻译"
	case "repost":
		return "转载"
	default:
		return "未知"
	}
}
func (this TypeContent) String() string {
	return string(this)
}

func (this TypeContent) Index() string {
	return string(this)
}

var MapTypeContent = map[string]enumBasePg.EnumString{
	ORIGINAL.String():    enumBasePg.EnumString{ORIGINAL.String(), ORIGINAL.Name()},
	TRANSLATION.String(): enumBasePg.EnumString{TRANSLATION.String(), TRANSLATION.Name()},
	REPOST.String():      enumBasePg.EnumString{REPOST.String(), REPOST.Name()},
}

func IsExistTypeContent(id string) bool {
	_, ok := MapTypeContent[id]
	return ok
}
