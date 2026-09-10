package holderPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/holderPg/jwtHolder"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/interfaces"
	"github.com/pangu-2/go-tools/tools/strPg"
)

type HolderPg struct {
	interfaces.StandardHolder
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
	jwt := c.Jwt.(jwtHolder.JwtPg)
	return jwt.Founder == "1"
}

func (c HolderPg) GetAccount() AccountHolder {
	return c.HolderData.(AccountHolder)
}

func (c HolderPg) GetAccountNo() string {
	return c.HolderData.(AccountHolder).No
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
func (c *HolderPg) GetRule() (interfaces.IHolderRule, bool) {
	if nil == c.Rule {
		return nil, false
	}
	return c.Rule, true
}

func (c *HolderPg) GetTypeDomain() string {
	return c.TypeDomain
}
