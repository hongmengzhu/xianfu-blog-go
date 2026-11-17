package modRamResourceAuthority

type CreateCt struct {
	Name        string `json:"name" form:"name" binding:"required,min=1,max=255" validate:"required,min=1,max=255" label:"名称" ` // 名称
	NameFl      string `json:"nameFl" label:"名称外文" `                                                                            // 名称外文
	Code        string `json:"code" label:"码值" `
	NameFull    string `json:"nameFull" label:"全称" `    // 全称
	Description string `json:"description" label:"描述" ` // 描述
	OrgId       string `json:"orgId" label:"组织id" `     // 组织id
	TypeSys     string `json:"typeSys" label:"类型" `
	Type        string `json:"type" label:"类型" `
	TypeValue   string `json:"typeValue" label:"对应类型id" `
	Mark        string `json:"mark" label:"标记" `
}
