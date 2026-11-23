package constDomainType

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// DomainType 域名类型
type DomainType string

const (
	Tenant       DomainType = "tenant" //租户
	Org          DomainType = "org"
	Customer     DomainType = "customer" //客户
	TenantModule DomainType = "tenantModule"
)

// Name 名称
func (this DomainType) Name() string {
	switch this {
	case "org":
		return "组织"
	case "tenant":
		return "租户"
	case "tenantModule":
		return "租户模块"
	case "customer":
		return "客户"
	default:
		return "未知"
	}
}

// 值
func (this DomainType) String() string {
	return string(this)
}

// 值
func (this DomainType) Index() string {
	return string(this)
}

// IsEqual 值是否相等
func (this DomainType) IsEqual(id string) bool {
	return string(this) == id
}

var DomainTypeMap = map[string]enumBasePg.EnumString{
	Org.String():          enumBasePg.EnumString{Id: Org.String(), Name: Org.Name()},
	TenantModule.String(): enumBasePg.EnumString{Id: TenantModule.String(), Name: TenantModule.Name()},
	Tenant.String():       enumBasePg.EnumString{Id: Tenant.String(), Name: Tenant.Name()},
	Customer.String():     enumBasePg.EnumString{Id: Customer.String(), Name: Customer.Name()},
}

func IsExistDomainType(id string) (DomainType, bool) {
	_, ok := DomainTypeMap[id]
	return DomainType(id), ok
}
