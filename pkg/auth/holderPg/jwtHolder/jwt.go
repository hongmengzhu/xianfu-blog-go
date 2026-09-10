package jwtHolder

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/interfaces"
)

// JwtPg 用户 会话信息 登录人信息
type JwtPg struct {
	MultiTenant   interfaces.IMultiTenantPg `json:"mTenant,omitempty" multiTenant` //多租户
	LoginNo       string                    `json:"loginNo"`                       //登录用户No,随时可以修改变动
	No            string                    `json:"no"`                            //系统编号
	LoginUserName string                    `json:"loginUserName"`                 //登录用户名
	Name          string                    `json:"name"`                          //显示名称
	OrgName       string                    `json:"OrgName,omitempty"`             //组织名称
	TenantName    string                    `json:"tName,omitempty"`               //组织名称
	Type          string                    `json:"type,omitempty"`                //类型
	Other         string                    `json:"other,omitempty"`               //其他信息
	Version       string                    `json:"version,omitempty"`             //版本
	Founder       string                    `json:"founder,omitempty" commont:"创始人"`
	Extra         map[string]interface{}    `json:"extra,omitempty"` //额外的，扩展
}

func NewJwtPg() *JwtPg {
	return new(JwtPg)
}

func (c JwtPg) toNotPoint() JwtPg {
	return c
}
