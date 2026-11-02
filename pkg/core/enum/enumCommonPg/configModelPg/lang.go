package configModelPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// Lang 语言
type Lang string

const (
	LangJS   Lang = "js"   //js
	LangTs   Lang = "ts"   //ts
	LangGo   Lang = "go"   //go
	LangJava Lang = "java" //java
)

// Name 名称
func (this Lang) Name() string {
	switch this {
	case "js":
		return "js"
	case "ts":
		return "ts"
	case "go":
		return "go"
	case "java":
		return "java"
	default:
		return "未知"
	}
}

// 值
func (this Lang) String() string {
	return string(this)
}

// Index 值
func (this Lang) Index() string {
	return string(this)
}

// IsEqual 值是否相等
func (this Lang) IsEqual(id string) bool {
	return string(this) == id
}

var LangMap = map[string]enumBasePg.EnumString{
	LangJS.String(): enumBasePg.EnumString{LangJS.String(), LangJS.Name()},
	LangTs.String(): enumBasePg.EnumString{LangTs.String(), LangTs.Name()},
	LangGo.String(): enumBasePg.EnumString{LangGo.String(), LangGo.Name()},
}

func IsExistLang(id string) (Lang, bool) {
	_, ok := LangMap[id]
	return Lang(id), ok
}
