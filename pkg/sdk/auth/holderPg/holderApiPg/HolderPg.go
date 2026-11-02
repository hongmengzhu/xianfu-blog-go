package holderApiPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/sdk/auth/holderPg/holderIfPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/sdk/auth/holderPg/jwtHolder"
)

type HolderPg struct {
	holderIfPg.StandardHolder
}

func (c *HolderPg) GetTenantNo() string {
	return c.HolderData.(DiplHolder).TenantNo
}

func (c *HolderPg) GetJwt() jwtHolder.JwtPg {
	return c.Jwt.(jwtHolder.JwtPg)
}

func (c HolderPg) GetDipl() DiplHolder {
	return c.HolderData.(DiplHolder)
}

func (c HolderPg) GetDiplNo() string {
	return c.HolderData.(DiplHolder).No
}

func (c HolderPg) GetAno() string {
	return c.HolderData.(DiplHolder).Ano
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
