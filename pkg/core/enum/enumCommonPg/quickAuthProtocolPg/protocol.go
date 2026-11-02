package quickAuthProtocolPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// Protocol 快捷登录协议
type Protocol string

const (
	Oauth2 Protocol = "oauth2"
	Qrcode Protocol = "qrcode"
	Sim    Protocol = "sim"
)

// Name 名称
func (this Protocol) Name() string {
	switch this {
	case "oauth2":
		return "oauth2"
	case "qrcode":
		return "扫码"
	case "sim":
		return "sim"
	default:
		return "未知"
	}
}

// 值
func (this Protocol) String() string {
	return string(this)
}

// 值
func (this Protocol) Index() string {
	return string(this)
}

var ProtocolMap = map[string]enumBasePg.EnumString{
	Oauth2.String(): enumBasePg.EnumString{Oauth2.String(), Oauth2.Name()},
	Qrcode.String(): enumBasePg.EnumString{Qrcode.String(), Qrcode.Name()},
	Sim.String():    enumBasePg.EnumString{Sim.String(), Sim.Name()},
}

func IsExistProtocol(id string) (Protocol, bool) {
	_, ok := ProtocolMap[id]
	return Protocol(id), ok
}
