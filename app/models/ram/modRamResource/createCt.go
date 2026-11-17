package modRamResource

type CreateCt struct {
	Name        string `json:"name" form:"name" binding:"required,min=1,max=255" validate:"required,min=1,max=255" label:"名称" ` // 名称
	NameFl      string `json:"nameFl" label:"名称外文" `                                                                            // 名称外文
	Code        string `json:"code" label:"码值" `
	NameFull    string `json:"nameFull" label:"全称" `    // 全称
	Description string `json:"description" label:"描述" ` // 描述
	ParentId    string `json:"parentId" label:"上级" `
	TypeSys     string `json:"typeSys" binding:"required" validate:"required" label:"类型" `  //类型;普通;系统;api
	TypeAttr    string `json:"typeAttr" binding:"required" validate:"required" label:"属性" ` //属性;菜单分类;资源
	Path        string `json:"path" label:"路径" `
	Method      string `json:"method" label:"方法" `
	ParentNo    string `json:"parentNo" label:"上级" `
}
