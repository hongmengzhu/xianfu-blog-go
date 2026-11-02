package typeFieldPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// TypeField 表达式类型;普通;正则
type TypeField string

const (
	Username TypeField = "username" //用户名
	Phone    TypeField = "phone"    //手机号
	Mail     TypeField = "mail"     //邮箱
	Password TypeField = "password" //密码
	Name     TypeField = "name"     //名称
)

// Name 名称
func (this TypeField) Name() string {
	switch this {
	case "username":
		return "用户名"
	case "phone":
		return "手机号"
	case "mail":
		return "邮箱"
	case "password":
		return "密码"
	case "name":
		return "名称"
	default:
		return "未知"
	}
}

// 值
func (this TypeField) String() string {
	return string(this)
}

// 值
func (this TypeField) Index() string {
	return string(this)
}

// IsEqual 值是否相等
func (this TypeField) IsEqual(id string) bool {
	return string(this) == id
}

var TypeFieldMap = map[string]enumBasePg.EnumString{
	Username.String(): enumBasePg.EnumString{Username.String(), Username.Name()},
	Phone.String():    enumBasePg.EnumString{Phone.String(), Phone.Name()},
	Mail.String():     enumBasePg.EnumString{Mail.String(), Mail.Name()},
	Password.String(): enumBasePg.EnumString{Password.String(), Password.Name()},
	Name.String():     enumBasePg.EnumString{Name.String(), Name.Name()},
}

func IsExistTypeField(id string) (TypeField, bool) {
	_, ok := TypeFieldMap[id]
	return TypeField(id), ok
}
