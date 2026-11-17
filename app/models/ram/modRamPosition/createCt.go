package modRamPosition

type CreateCt struct {
	Name        string `json:"name" form:"name" binding:"required,min=1,max=255" validate:"required,min=1,max=255" label:"名称" ` // 名称
	NameFl      string `json:"nameFl" label:"名称外文" `                                                                            // 名称外文
	Code        string `json:"code" label:"码值" `
	NameFull    string `json:"nameFull" label:"全称" `    // 全称
	Description string `json:"description" label:"描述" ` // 描述
	OrgId       string `json:"orgId" label:"组织id" `     // 组织id
}
