package typeSourcePg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// Module 模块
type Module string

const (
	Article  Module = "article"
	NEW      Module = "new"
	TUTORIAL Module = "tutorial"
)

func (this Module) Name() string {
	switch this {
	case "article":
		return "文章"
	case "new":
		return "新闻"
	case "tutorial":
		return "教程"
	default:
		return "未知"
	}
}
func (this Module) String() string {
	return string(this)
}

func (this Module) Index() string {
	return string(this)
}

var MapModule = map[string]enumBasePg.EnumString{
	Article.String():  enumBasePg.EnumString{Article.String(), Article.Name()},
	NEW.String():      enumBasePg.EnumString{NEW.String(), NEW.Name()},
	TUTORIAL.String(): enumBasePg.EnumString{TUTORIAL.String(), TUTORIAL.Name()},
}

func IsExistModule(id string) bool {
	_, ok := MapModule[id]
	return ok
}
