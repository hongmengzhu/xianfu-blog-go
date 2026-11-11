package modBlogArticleCategory

type CreateCt struct {
	Name        string `json:"name" form:"name" validate:"required,min=1,max=255" label:"名称" ` // 名称
	NameFl      string `json:"nameFl" label:"名称外文" `
	Code        string `json:"code" label:"码值" `
	NameFull    string `json:"nameFull" label:"全称" `
	Description string `json:"description" label:"描述" `
	ParentNo    string `json:"parentNo" label:"上级" `
	ParentId    string `json:"parentId" label:"上级" `
	Module      string `json:"module" label:"模块" `
}
