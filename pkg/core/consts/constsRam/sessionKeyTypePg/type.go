package sessionKeyTypePg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// SessionKeyType 密钥类型
type SessionKeyType string

const (
	General            SessionKeyType = "general"            //普通
	RefreshToken       SessionKeyType = "refreshToken"       //刷新token
	AccessToken        SessionKeyType = "accessToken"        //访问令牌
	PasswordEncryption SessionKeyType = "passwordEncryption" //管理
	Login              SessionKeyType = "login"              //登陆加密
)

// Name 名称
func (this SessionKeyType) Name() string {
	switch this {
	case "general":
		return "普通"
	case "refreshToken":
		return "刷新token"
	case "accessToken":
		return "访问令牌"
	case "passwordEncryption":
		return "管理"
	case "login":
		return "登陆加密"
	default:
		return "未知"
	}
}

// 值
func (this SessionKeyType) String() string {
	return string(this)
}

// 值
func (this SessionKeyType) Code() string {
	return string(this)
}

// IsEqual 值是否相等
func (this SessionKeyType) IsEqual(id string) bool {
	return string(this) == id
}

var SessionKeyTypeMap = map[string]enumBasePg.EnumString{
	General.String():            enumBasePg.EnumString{General.String(), General.Name()},
	RefreshToken.String():       enumBasePg.EnumString{RefreshToken.String(), RefreshToken.Name()},
	AccessToken.String():        enumBasePg.EnumString{AccessToken.String(), AccessToken.Name()},
	PasswordEncryption.String(): enumBasePg.EnumString{PasswordEncryption.String(), PasswordEncryption.Name()},
	Login.String():              enumBasePg.EnumString{Login.String(), Login.Name()},
}

func IsExistSessionKeyType(id string) (SessionKeyType, bool) {
	_, ok := SessionKeyTypeMap[id]
	return SessionKeyType(id), ok
}
