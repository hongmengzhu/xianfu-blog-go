package clientPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// MiniProgram 小程序
type MiniProgram string

const (
	MiniProgramWeChat MiniProgram = "weChat" //微信小程序
	MiniProgramWeCom  MiniProgram = "weCom"  //企业微信小程序
	MiniProgramAlipay MiniProgram = "alipay" //支付宝
	MiniProgramHuawei MiniProgram = "huawei" //华为
	MiniProgramDouyin MiniProgram = "douyin" //抖音
)

// Name 名称
func (this MiniProgram) Name() string {
	switch this {
	case "weChat":
		return "微信小程序"
	case "weCom":
		return "企业微信小程序"
	case "alipay":
		return "支付宝"
	case "huawei":
		return "华为"
	default:
		return "未知"
	}
}

// String 值
func (this MiniProgram) String() string {
	return string(this)
}

// Index 值
func (this MiniProgram) Index() string {
	return string(this)
}

// IsEqual 值是否相等
func (this MiniProgram) IsEqual(id string) bool {
	return string(this) == id
}

var MiniProgramMap = map[string]enumBasePg.EnumString{
	MiniProgramWeChat.String(): enumBasePg.EnumString{MiniProgramWeChat.String(), MiniProgramWeChat.Name()},
	MiniProgramWeCom.String():  enumBasePg.EnumString{MiniProgramWeCom.String(), MiniProgramWeCom.Name()},
	MiniProgramAlipay.String(): enumBasePg.EnumString{MiniProgramAlipay.String(), MiniProgramAlipay.Name()},
	MiniProgramHuawei.String(): enumBasePg.EnumString{MiniProgramHuawei.String(), MiniProgramHuawei.Name()},
	MiniProgramDouyin.String(): enumBasePg.EnumString{MiniProgramDouyin.String(), MiniProgramDouyin.Name()},
}

func IsExistMiniProgram(id string) (MiniProgram, bool) {
	_, ok := MiniProgramMap[id]
	return MiniProgram(id), ok
}
