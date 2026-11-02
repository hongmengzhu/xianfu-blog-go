package typeAccessMethodPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// TypeAccessMethod 类型;控制台;api
type TypeAccessMethod string

const (
	Console TypeAccessMethod = "console" //控制台
	Api     TypeAccessMethod = "api"     //api
)

// Name 名称
func (this TypeAccessMethod) Name() string {
	switch this {
	case "console":
		return "控制台"
	case "api":
		return "api"
	default:
		return "未知"
	}
}

// 值
func (this TypeAccessMethod) String() string {
	return string(this)
}

// 值
func (this TypeAccessMethod) Index() string {
	return string(this)
}

var TypeAccessMethodMap = map[string]enumBasePg.EnumString{
	Console.String(): enumBasePg.EnumString{Console.String(), Console.Name()},
	Api.String():     enumBasePg.EnumString{Api.String(), Api.Name()},
}

func IsExistTypeAccessMethod(id string) (TypeAccessMethod, bool) {
	_, ok := TypeAccessMethodMap[id]
	return TypeAccessMethod(id), ok
}
