package modBlogCategory

type CreateCt struct {
	Name        string `json:"name" form:"name" binding:"required,min=1,max=255" validate:"required,min=1,max=255" label:"名称" ` // 名称
	NameFl      string `json:"nameFl" label:"名称外文" `
	No          string `json:"no" label:"编号代号" `
	NameFull    string `json:"nameFull" label:"全称" `
	Description string `json:"description" label:"描述" `
	ParentNo    string `json:"parentNo" label:"上级" `
	ParentId    string `json:"parentId" label:"上级" `
}
