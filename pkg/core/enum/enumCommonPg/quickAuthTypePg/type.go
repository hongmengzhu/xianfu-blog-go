package quickAuthTypePg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// Type 快捷登录类型
type Type string

const (
	App    Type = "app"
	Qrcode Type = "qrcode"
	Sim    Type = "sim"
)

// Name 名称
func (this Type) Name() string {
	switch this {
	case "app":
		return "app"
	case "qrcode":
		return "扫码"
	case "sim":
		return "sim"
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

var TypeMap = map[string]enumBasePg.EnumString{
	App.String():    enumBasePg.EnumString{Id: App.String(), Name: App.Name()},
	Qrcode.String(): enumBasePg.EnumString{Id: Qrcode.String(), Name: Qrcode.Name()},
	Sim.String():    enumBasePg.EnumString{Id: Sim.String(), Name: Sim.Name()},
}

func IsExistType(id string) (Type, bool) {
	_, ok := TypeMap[id]
	return Type(id), ok
}
