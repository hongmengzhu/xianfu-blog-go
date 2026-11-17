package modBasicTagsRelation

type CreateCt struct {
	Name         string                 `json:"name" form:"name" binding:"required,min=1,max=255" validate:"required,min=1,max=255" label:"名称" `
	NameFl       string                 `json:"nameFl" label:"名称外文" `
	No           string                 `json:"no" label:"标签" binding:"required,min=1,max=255" validate:"required,min=1,max=255"`
	Category     string                 `json:"category" label:"分类" binding:"required,min=1,max=255" validate:"required,min=1,max=255"`
	NameFull     string                 `json:"nameFull" label:"全称" `
	Description  string                 `json:"description" label:"描述" `
	OrgId        string                 `json:"orgId" label:"组织id" `
	AttributeMap map[string]interface{} `json:"attributeMap" label:"属性" `
	TypeSys      string                 `json:"typeSys" label:"类型" binding:"required,min=1,max=255" validate:"required,min=1,max=255"`
	NameShort    string                 `json:"nameShort" label:"名称简称" `
}
