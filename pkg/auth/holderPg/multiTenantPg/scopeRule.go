package multiTenantPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/interfaces"
)

// 解析关键字
const RuleParamTenant = "Tenant;"
const RuleParamMerchant = "Merchant;"
const RuleParamStore = "Store;"

var _ interfaces.IHolderRule = (*MultiRule)(nil)

// 多租户规则
type MultiRule struct {
	MultipleTenant bool //租户多个
	Tenant         bool //租户
	Org            bool //组织
	MultipleOrg    bool //组织多个
	Merchant       bool //商户
	MultipleStore  bool //店铺多个
	Store          bool //店铺
	Owner          bool //所有者
	MultiOwner     bool //所有者 多个
}

//用 MultiRule2 实现接口 interfaces.IHolderRule

// NewMultiRuleDefault
//
//	@Description: 默认规则
//	@return MultiRule
func NewMultiRuleDefault() MultiRule {
	return MultiRule{
		MultipleTenant: true,
		Tenant:         true,
		Org:            true,
		MultipleOrg:    true,
		Merchant:       true,
		MultipleStore:  true,
		Store:          true,
		Owner:          true,
		MultiOwner:     true,
	}
}

// NewMultiRuleDefaultBySystem
//
//	@Description: 默认规则
//	@return MultiRule
func NewMultiRuleDefaultBySystem() MultiRule {
	return MultiRule{
		MultipleTenant: false,
		Tenant:         false,
		Org:            false,
		MultipleOrg:    false,
		Merchant:       false,
		MultipleStore:  false,
		Store:          false,
		Owner:          false,
		MultiOwner:     false,
	}
}

// NewMultiRuleDefaultByConsumer
//
//	@Description: 默认规则
//	@return MultiRule
func NewMultiRuleDefaultByConsumer() MultiRule {
	return MultiRule{
		MultipleTenant: false,
		Tenant:         false,
		Org:            false,
		MultipleOrg:    false,
		Merchant:       false,
		MultipleStore:  false,
		Store:          false,
		Owner:          false,
		MultiOwner:     false,
	}
}
