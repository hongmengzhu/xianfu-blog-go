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
	General.String():            enumBasePg.EnumString{Id: General.String(), Name: General.Name()},
	RefreshToken.String():       enumBasePg.EnumString{Id: RefreshToken.String(), Name: RefreshToken.Name()},
	AccessToken.String():        enumBasePg.EnumString{Id: AccessToken.String(), Name: AccessToken.Name()},
	PasswordEncryption.String(): enumBasePg.EnumString{Id: PasswordEncryption.String(), Name: PasswordEncryption.Name()},
	Login.String():              enumBasePg.EnumString{Id: Login.String(), Name: Login.Name()},
}

func IsExistSessionKeyType(id string) (SessionKeyType, bool) {
	_, ok := SessionKeyTypeMap[id]
	return SessionKeyType(id), ok
}
