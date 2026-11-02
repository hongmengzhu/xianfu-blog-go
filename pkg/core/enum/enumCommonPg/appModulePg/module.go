package appModulePg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/consts/constsRam/typeDomainPg"
)

type AppModule string

const (
	Tenant      AppModule = "tenant"
	Manage      AppModule = "manage"
	Customer    AppModule = "customer"
	App         AppModule = "app"
	System      AppModule = "system"
	Merchant    AppModule = "merchant"
	WEB         AppModule = "WEB"
	ManageOwner AppModule = "manageOwner"
)

func (c AppModule) ToTypeDomain() typeDomainPg.TypeDomain {
	switch c {
	case System:
		return typeDomainPg.System
	case Merchant:
		return typeDomainPg.Merchant
	case Manage:
		return typeDomainPg.Manage
	case Tenant:
		return typeDomainPg.Tenant
	case Customer:
		return typeDomainPg.Customer
	case App:
		return typeDomainPg.Customer
	case WEB:
		return typeDomainPg.General
	default:
		return typeDomainPg.General
	}
}
func (this AppModule) String() string {
	return string(this)
}

// IsEqual 值是否相等
func (this AppModule) IsEqual(id string) bool {
	return string(this) == id
}
