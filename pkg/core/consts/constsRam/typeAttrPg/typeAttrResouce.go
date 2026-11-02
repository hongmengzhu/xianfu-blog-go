package typeAttrPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// TypeAttr 资源属性
type TypeAttr string

const (
	Category     TypeAttr = "category"     //分类
	CategoryLast TypeAttr = "categoryLast" //最后一级分类
	Resource     TypeAttr = "resource"     //资源
)

// Name 名称
func (this TypeAttr) Name() string {
	switch this {
	case "category":
		return "分类"
	case "categoryLast":
		return "最后一级分类"
	case "resource":
		return "资源"
	default:
		return "未知"
	}
}

// 值
func (this TypeAttr) String() string {
	return string(this)
}

// 值
func (this TypeAttr) Index() string {
	return string(this)
}

// IsEqual 值是否相等
func (this TypeAttr) IsEqual(id string) bool {
	return string(this) == id
}

var TypeAttrMap = map[string]enumBasePg.EnumString{
	Resource.String():     enumBasePg.EnumString{Resource.String(), Resource.Name()},
	Category.String():     enumBasePg.EnumString{Category.String(), Category.Name()},
	CategoryLast.String(): enumBasePg.EnumString{CategoryLast.String(), CategoryLast.Name()},
}

func IsExistTypeAttr(id string) (TypeAttr, bool) {
	_, ok := TypeAttrMap[id]
	return TypeAttr(id), ok
}
