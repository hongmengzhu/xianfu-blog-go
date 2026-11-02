package typeSourcePg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// TypeSource 类型源|采集|手写
type TypeSource string

const (
	HANDWRITTEN TypeSource = "handwritten"
	COLLECTION  TypeSource = "collection"
)

func (this TypeSource) Name() string {
	switch this {
	case "handwritten":
		return "手写"
	case "collection":
		return "采集"
	default:
		return "未知"
	}
}
func (this TypeSource) String() string {
	return string(this)
}

func (this TypeSource) Index() string {
	return string(this)
}

var MapTypeSource = map[string]enumBasePg.EnumString{
	HANDWRITTEN.String(): enumBasePg.EnumString{HANDWRITTEN.String(), HANDWRITTEN.Name()},
	COLLECTION.String():  enumBasePg.EnumString{COLLECTION.String(), COLLECTION.Name()},
}

func IsExistTypeSource(id string) bool {
	_, ok := MapTypeSource[id]
	return ok
}
