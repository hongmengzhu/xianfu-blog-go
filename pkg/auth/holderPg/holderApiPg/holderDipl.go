package holderApiPg

import (
	_ "github.com/hongmengzhu/xianfu-blog-go/pkg/interfaces"
)

type DiplHolder struct {
	No       string `json:"no" comment:"编码"`
	Code     string `json:"code" comment:"编码" `
	Name     string `json:"name" label:"名称" `
	TenantNo string `json:"tenantNo" comment:"租户编码"`
	Ano      string `json:"ano" comment:"账号"`
}
