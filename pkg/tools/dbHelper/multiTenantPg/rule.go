package multiTenantPg

import "strings"

// parseRuleParam
//
//	@Description: 解析规则
//	@param ruleParam
//	@return rule
func parseRuleParam(ruleParam string) (rule MultiRule) {
	//租户
	rule.MultipleTenant = strings.Contains(ruleParam, RuleParamTenant)
	//商户
	rule.Merchant = strings.Contains(ruleParam, RuleParamMerchant)
	//店铺
	rule.MultipleStore = strings.Contains(ruleParam, RuleParamStore)
	return rule
}
