package entityRam

type RamAccountJsonExtraCond struct {
	OwnerCodes  []string `json:"ownerCodes"`  // 拥有者
	Departments []string `json:"departments"` // 部门
	Roles       []string `json:"roles"`       // 角色
	Teams       []string `json:"teams"`       // 团队
	Groups      []string `json:"groups"`      // 组
	Levels      []string `json:"levels"`      // 级别
}
