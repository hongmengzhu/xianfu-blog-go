package interfaces

// 实体总
type IHolderPg interface {
}

// 多租户
type IMultiTenantPg interface {
}

// 用户数据
type IHolderDataPg interface {
}

// jwt数据
type IHolderJwtPg interface {
}

// 租户数据
type ITenantPg interface {
}

// 其他
type IHolderOtherPg interface {
}

// 其他
type IHolderRule interface {
}

type StandardHolder struct {
	Jwt         IHolderJwtPg           `json:"jwt,omitempty" commont:"jwt"`
	MultiTenant IMultiTenantPg         `json:"mult,omitempty" commont:"多租户"`
	HolderData  IHolderDataPg          `json:"hdata,omitempty" commont:"用户数据"`
	Tenant      ITenantPg              `json:"tenant,omitempty" commont:"租户数据"`
	Other       IHolderOtherPg         `json:"other,omitempty" commont:"其他数据"`
	Rule        IHolderRule            `json:"rule,omitempty" commont:"规则"`
	Founder     string                 `json:"founder,omitempty" commont:"创始人"`
	TypeDomain  string                 `json:"typeDomain,omitempty" commont:"类型"`
	Map         map[string]interface{} `json:"map,omitempty" commont:"map"`
}

func NewStandardHolder() *StandardHolder {
	return new(StandardHolder)
}

func (c *StandardHolder) IsFounder() bool {
	return c.Founder == "1"
}
