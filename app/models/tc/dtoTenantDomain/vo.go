package dtoTenantDomain

// Dto
// @Description: 租户域名数据模型
type Dto struct {
	Domain   string `json:"domain" label:"域名" `
	TenantNo string `json:"tenantNo" label:"租户编号" `
	Ano      string `json:"ano" label:"账号编号" `
	Theme    string `json:"theme" label:"主题"`
}
