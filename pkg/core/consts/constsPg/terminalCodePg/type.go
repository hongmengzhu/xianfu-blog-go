package terminalCodePg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// TerminalCode 域类型;普通;系统;商户
type TerminalCode string

const (
	General  TerminalCode = "general"  //普通
	System   TerminalCode = "system"   //系统
	Merchant TerminalCode = "merchant" //商户
	Manage   TerminalCode = "manage"   //管理
	Tenant   TerminalCode = "tenant"   //租户
	Shop     TerminalCode = "shop"     //店铺
	Customer TerminalCode = "customer" //客户
)

// Name 名称
func (this TerminalCode) Name() string {
	switch this {
	case "general":
		return "普通"
	case "system":
		return "系统"
	case "merchant":
		return "商户"
	case "manage":
		return "管理"
	case "tenant":
		return "租户"
	case "shop":
		return "店铺"
	case "customer":
		return "客户"
	default:
		return "未知"
	}
}

// 值
func (this TerminalCode) String() string {
	return string(this)
}

// 值
func (this TerminalCode) Index() string {
	return string(this)
}

// IsEqual 值是否相等
func (this TerminalCode) IsEqual(id string) bool {
	return string(this) == id
}

var TerminalCodeMap = map[string]enumBasePg.EnumString{
	General.String():  enumBasePg.EnumString{General.String(), General.Name()},
	System.String():   enumBasePg.EnumString{System.String(), System.Name()},
	Merchant.String(): enumBasePg.EnumString{Merchant.String(), Merchant.Name()},
	Tenant.String():   enumBasePg.EnumString{Tenant.String(), Merchant.Name()},
	Shop.String():     enumBasePg.EnumString{Shop.String(), Shop.Name()},
	Manage.String():   enumBasePg.EnumString{Manage.String(), Manage.Name()},
	Customer.String(): enumBasePg.EnumString{Customer.String(), Customer.Name()},
}

func IsExistTerminalCode(id string) (TerminalCode, bool) {
	_, ok := TerminalCodeMap[id]
	return TerminalCode(id), ok
}
