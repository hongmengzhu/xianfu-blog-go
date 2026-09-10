package holderPg

import _ "github.com/hongmengzhu/xianfu-blog-go/pkg/interfaces"

// AccountHolderOs 组织架构
// @Description:
type AccountHolderOs struct {
	Departments []string `json:"departments" comment:"部门" `
	Groups      []string `json:"groups" comment:"组" `
	Levels      []string `json:"levels" comment:"级别" `
	Merchants   []string `json:"merchants" comment:"商户" `
	Orgs        []string `json:"orgs" comment:"组织" `
	Projects    []string `json:"projects" comment:"项目" `
	Roles       []string `json:"roles" comment:"角色" `
	Shops       []string `json:"shops" comment:"店铺"`
	Stores      []string `json:"Stores" comment:"店铺"`
	Teams       []string `json:"teams" comment:"团队" `
	Tenants     []string `json:"tenants" comment:"租户" `
}

func NewAccountHolderOsP() *AccountHolderOs {
	return &AccountHolderOs{
		Departments: make([]string, 0),
		Groups:      make([]string, 0),
		Levels:      make([]string, 0),
		Merchants:   make([]string, 0),
		Orgs:        make([]string, 0),
		Projects:    make([]string, 0),
		Roles:       make([]string, 0),
		Shops:       make([]string, 0),
		Teams:       make([]string, 0),
		Tenants:     make([]string, 0),
	}
}

func NewAccountHolderOs() AccountHolderOs {
	return AccountHolderOs{
		Departments: make([]string, 0),
		Groups:      make([]string, 0),
		Levels:      make([]string, 0),
		Merchants:   make([]string, 0),
		Orgs:        make([]string, 0),
		Projects:    make([]string, 0),
		Roles:       make([]string, 0),
		Shops:       make([]string, 0),
		Teams:       make([]string, 0),
		Tenants:     make([]string, 0),
	}
}
