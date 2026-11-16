package holderPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/sdk/auth/holderPg/holderIfPg"
	"github.com/pangu-2/go-tools/tools/strPg"
)

type HolderPg struct {
	holderIfPg.StandardHolder
}

func (c *HolderPg) GetTenantNo() string {
	if nil == c.HolderData {
		return "-1"
	}
	holder := c.HolderData.(AccountHolder)
	if strPg.IsBlank(holder.TenantNo) {
		return "-1"
	}
	return holder.TenantNo
}
func (c *HolderPg) GetOrgNo() string {
	if nil == c.HolderData {
		return "-1"
	}
	holder := c.HolderData.(AccountHolder)
	if strPg.IsBlank(holder.OrgNo) {
		return "-1"
	}
	return holder.OrgNo
}

// IsFounder
//
//	@Description: 是否创始人
//	@receiver c
//	@return bool
func (c *HolderPg) IsFounder() bool {
	if nil == c.HolderData {
		return false
	}
	if holder, ok := c.HolderData.(AccountHolder); ok {
		return holder.Founder == 1
	}
	return false
}

func (c HolderPg) GetAccount() AccountHolder {
	if nil == c.HolderData {
		return AccountHolder{}
	}
	if holder, ok := c.HolderData.(AccountHolder); ok {
		return holder
	}
	return AccountHolder{}
}

func (c HolderPg) GetAccountNo() string {
	if nil == c.HolderData {
		return ""
	}
	if holder, ok := c.HolderData.(AccountHolder); ok {
		return holder.No
	}
	return ""
}

// GetMerchantNo
//
//	@Description: 商户编号
//	@receiver c
//	@return string
func (c *HolderPg) GetMerchantNo() string {
	return ""
}

func (c *HolderPg) GetStoreNo() string {
	return ""
}

func (c *HolderPg) GetOwnerNo() string {
	return c.GetAccount().OwnerNo
}

// GetRule
//
//	@Description: 规则
//	@receiver c
//	@return interfaces.IHolderRule
//	@return bool
func (c *HolderPg) GetRule() (holderIfPg.IHolderRule, bool) {
	if nil == c.Rule {
		return nil, false
	}
	return c.Rule, true
}

func (c *HolderPg) GetTypeDomain() string {
	return c.TypeDomain
}
