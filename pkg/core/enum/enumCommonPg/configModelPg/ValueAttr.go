package configModelPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// ValueAttr 字段值属性|单值|对象
type ValueAttr string

const (
	ValueAttrObject ValueAttr = "object" //对象
	ValueAttrSingle ValueAttr = "single" //单个
)

// Name 名称
func (this ValueAttr) Name() string {
	switch this {
	case "single":
		return "单个"
	case "object":
		return "对象"
	default:
		return "未知"
	}
}

// 值
func (this ValueAttr) String() string {
	return string(this)
}

// 值
func (this ValueAttr) Index() string {
	return string(this)
}

// IsEqual 值是否相等
func (this ValueAttr) IsEqual(id string) bool {
	return string(this) == id
}

var ValueAttrMap = map[string]enumBasePg.EnumString{
	ValueAttrSingle.String(): enumBasePg.EnumString{ValueAttrSingle.String(), ValueAttrSingle.Name()},
	ValueAttrObject.String(): enumBasePg.EnumString{ValueAttrObject.String(), ValueAttrObject.Name()},
}

func IsExistValueAttr(id string) (ValueAttr, bool) {
	_, ok := ValueAttrMap[id]
	return ValueAttr(id), ok
}
