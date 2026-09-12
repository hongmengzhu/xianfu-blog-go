package authTokenPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/interfaces"
)

// Param 用户 会话信息 登录人信息
type Param struct {
	UniqueId      string                    `json:"uniqueId" label:"唯一身份标识，主要用来作为一次性token,从而回避重放攻击"`
	MultiTenant   interfaces.IMultiTenantPg `json:"mTenant,omitempty" label:"多租户"`
	LoginNo       string                    `json:"loginNo"  label:"登录用户编号"`
	No            string                    `json:"no" label:"系统编号"`
	LoginUserName string                    `json:"loginUserName" label:"登录用户名"`
	Name          string                    `json:"name" label:"显示名称"`
	TenantNo      string                    `json:"tNo,omitempty" label:"租户编号"`
	OrgNo         string                    `json:"oNo,omitempty" label:"组织编号"`
	Type          string                    `json:"type,omitempty"  label:"类型"`
	Other         string                    `json:"other,omitempty" label:"其他信息"`
	Version       string                    `json:"version,omitempty"  label:"版本"`
	Founder       string                    `json:"founder,omitempty" label:"创始人"`
	Extra         map[string]interface{}    `json:"extra,omitempty"  label:"额外的，扩展"`
	Result        Result                    `json:"result,omitempty" label:"密钥对"`
}

func NewJwtPg() *Param {
	return new(Param)
}

func (c Param) toNotPoint() Param {
	return c
}
