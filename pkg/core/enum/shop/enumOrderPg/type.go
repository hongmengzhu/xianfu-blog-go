package enumOrderPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// Type 订单状态
type Type string

const (
	GENERAL Type = "general" //普通
)

// Name 名称
func (this Type) Name() string {
	switch this {
	case "general":
		return "普通"
	default:
		return "未知"
	}
}

// 值
func (this Type) String() string {
	return string(this)
}

// 值
func (this Type) Index() string {
	return string(this)
}

// IsEqual 值是否相等
func (this Type) IsEqual(id string) bool {
	return string(this) == id
}

var TypeMap = map[string]enumBasePg.EnumString{
	GENERAL.String(): enumBasePg.EnumString{GENERAL.String(), GENERAL.Name()},
}

func IsExistType(id string) (Type, bool) {
	_, ok := TypeMap[id]
	return Type(id), ok
}
