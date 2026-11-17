package modRamResourceMenu

type ResourceAndGroupVo struct {
	Type string `json:"type" form:"type" binding:"required" validate:"required" label:"类型" `
	Id   string `json:"id" form:"id" binding:"required" validate:"required" label:"资源组或资源" `
}
