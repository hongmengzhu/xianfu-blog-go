package appDomainPg

// 域模式
type AppDomain string

const (
	Tenant      AppDomain = "tenant"
	Manage      AppDomain = "manage"
	System      AppDomain = "system"
	WEB         AppDomain = "WEB"
	ManageOwner AppDomain = "manageOwner"
)
