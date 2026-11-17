package modRamResourceGroupRelation

type QueryByTypeValueCt struct {
	TypeValue    string `json:"typeValue" form:"typeValue" binding:"required" validate:"required" label:"资源组id" `
	TypeCategory string `json:"typeCategory" label:"类型" `
}
