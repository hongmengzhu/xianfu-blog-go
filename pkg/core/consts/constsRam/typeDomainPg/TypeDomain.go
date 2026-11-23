package typeDomainPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// TypeDomain 域类型;普通;系统;商户
type TypeDomain string

const (
	General  TypeDomain = "general"  //普通
	System   TypeDomain = "system"   //系统
	Merchant TypeDomain = "merchant" //商户
	Manage   TypeDomain = "manage"   //管理
	Tenant   TypeDomain = "tenant"   //租户
	Shop     TypeDomain = "shop"     //店铺
	Customer TypeDomain = "customer" //客户
)

// Name 名称
func (this TypeDomain) Name() string {
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
func (this TypeDomain) String() string {
	return string(this)
}

// 值
func (this TypeDomain) Code() string {
	return string(this)
}

// IsEqual 值是否相等
func (this TypeDomain) IsEqual(id string) bool {
	return string(this) == id
}

var TypeDomainMap = map[string]enumBasePg.EnumString{
	General.String():  enumBasePg.EnumString{Id: General.String(), Name: General.Name()},
	System.String():   enumBasePg.EnumString{Id: System.String(), Name: System.Name()},
	Merchant.String(): enumBasePg.EnumString{Id: Merchant.String(), Name: Merchant.Name()},
	Tenant.String():   enumBasePg.EnumString{Id: Tenant.String(), Name: Merchant.Name()},
	Shop.String():     enumBasePg.EnumString{Id: Shop.String(), Name: Shop.Name()},
	Manage.String():   enumBasePg.EnumString{Id: Manage.String(), Name: Manage.Name()},
	Customer.String(): enumBasePg.EnumString{Id: Customer.String(), Name: Customer.Name()},
}

func IsExistTypeDomain(id string) (TypeDomain, bool) {
	_, ok := TypeDomainMap[id]
	return TypeDomain(id), ok
}
