package multiTenantPg

import (
	"strings"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/holderPg"
	"github.com/pangu-2/go-tools/tools/strPg"
)

// 多租户
type MultiTenantPg struct {
	TenantNo  []string `json:"tenantNo,omitempty"`  //租户
	OrgNo     []string `json:"orgNo,omitempty"`     //组织
	RuleParam string   `json:"ruleParam,omitempty"` //参数规则
}

func (c *MultiTenantPg) TenantNoToArr() ([]string, bool) {
	ids := make([]string, 0)
	if nil == c.TenantNo {
		return ids, false
	}
	if len(c.TenantNo) > 0 {
		for _, item := range c.TenantNo {
			if strPg.IsBlank(item) {
				continue
			}
			ids = append(ids, strings.TrimSpace(item))
		}
		return ids, true
	}
	return ids, false
}

func (c *MultiTenantPg) OrgNoToArr() ([]string, bool) {
	ids := make([]string, 0)
	if nil == c.OrgNo {
		return ids, false
	}
	if len(c.OrgNo) > 0 {
		for _, item := range c.OrgNo {
			if strPg.IsBlank(item) {
				continue
			}
			ids = append(ids, strings.TrimSpace(item))
		}
		return ids, true
	}
	return ids, false
}
func GetMultiTenant(source holderPg.HolderPg) (MultiTenantPg, bool) {
	if nil == source.MultiTenant {
		return MultiTenantPg{}, false
	}
	return source.MultiTenant.(MultiTenantPg), true
}

func GetMultiOrg(source holderPg.HolderPg) (MultiTenantPg, bool) {
	if nil == source.MultiTenant {
		return MultiTenantPg{}, false
	}
	return source.MultiTenant.(MultiTenantPg), true
}
