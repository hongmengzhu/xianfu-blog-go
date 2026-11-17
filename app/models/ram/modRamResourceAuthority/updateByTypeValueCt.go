package modRamResourceAuthority

type UpdateByTypeValueCt struct {
	TypeValue string   `json:"typeValue" form:"typeValue" binding:"required" validate:"required" label:"资源组id" `
	Ids       []string `json:"ids" form:"ids" binding:"required" validate:"required" label:"数据" `
}
