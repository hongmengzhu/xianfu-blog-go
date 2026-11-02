package entityRam

// RamAccountJsonOs 组织架构 organizational structure
type RamAccountJsonOs struct {
	Departments []string `json:"departments"` // 部门
	Groups      []string `json:"groups"`      // 组
	Levels      []string `json:"levels"`      // 级别
	Merchants   []string `json:"merchants"`   // 商户
	Orgs        []string `json:"orgs"`        // 组织
	Projects    []string `json:"projects"`    // 项目
	Roles       []string `json:"roles"`       // 角色
	Shops       []string `json:"shops"`       // 店铺
	Stores      []string `json:"stores"`      // 店铺
	Teams       []string `json:"teams"`       // 团队
	Tenants     []string `json:"tenants"`     // 租户
}

// {"orgs": [], "roles": [], "shops": [], "teams": [], "groups": [], "tenants": [], "projects": [], "departments": [],"levels":[],"merchants":[],"stores":[]}
