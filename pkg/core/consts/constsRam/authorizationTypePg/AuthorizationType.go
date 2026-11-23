package authorizationTypePg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// AuthorizationType 密码/pin码/双因子码/openid
type AuthorizationType string

const (
	PASSWORD AuthorizationType = "password"
	PIN      AuthorizationType = "pin"
	MFA2     AuthorizationType = "MFA2"
	OPENID   AuthorizationType = "openid"
)

func (this AuthorizationType) Name() string {
	switch this {
	case "password":
		return "密码"
	case "pin":
		return "pin码"
	case "MFA2":
		return "双因子码"
	case "openid":
		return "openid"
	default:
		return ""
	}
}
func (this AuthorizationType) String() string {
	return string(this)
}

func (this AuthorizationType) Index() string {
	return string(this)
}

var Map = map[string]enumBasePg.EnumString{
	PASSWORD.String(): enumBasePg.EnumString{Id: PASSWORD.String(), Name: PASSWORD.Name()},
	PIN.String():      enumBasePg.EnumString{Id: PIN.String(), Name: PIN.Name()},
	MFA2.String():     enumBasePg.EnumString{Id: MFA2.String(), Name: MFA2.Name()},
	OPENID.String():   enumBasePg.EnumString{Id: OPENID.String(), Name: OPENID.Name()},
}

func IsExist(id string) bool {
	_, ok := Map[id]
	return ok
}
