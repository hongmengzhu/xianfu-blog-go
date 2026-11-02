package typeSysPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// TypeSys 类型;普通;系统;api
type TypeSys string

const (
	General TypeSys = "general" //普通
	System  TypeSys = "system"  //系统
	Api     TypeSys = "api"     //api
)

// Name 名称
func (this TypeSys) Name() string {
	switch this {
	case "general":
		return "普通"
	case "system":
		return "系统"
	case "api":
		return "api"
	default:
		return "未知"
	}
}

// 值
func (this TypeSys) String() string {
	return string(this)
}

// 值
func (this TypeSys) Index() string {
	return string(this)
}

var TypeSysMap = map[string]enumBasePg.EnumString{
	General.String(): enumBasePg.EnumString{General.String(), General.Name()},
	System.String():  enumBasePg.EnumString{System.String(), System.Name()},
	Api.String():     enumBasePg.EnumString{Api.String(), Api.Name()},
}

func IsExistTypeSys(id string) (TypeSys, bool) {
	_, ok := TypeSysMap[id]
	return TypeSys(id), ok
}
