package modelBasePg

import "github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"

type CreateUpdateCt struct {
	ID          typePg.Int64String `json:"id" form:"id" label:"id" `
	Name        string             `json:"name" form:"name" binding:"required,min=1,max=255" validate:"required,min=1,max=255" label:"名称" `
	NameFl      string             `json:"nameFl" label:"名称外文" `
	Code        string             `json:"code" label:"标记" `
	NameFull    string             `json:"nameFull" label:"全称" `
	Description string             `json:"description" label:"描述" `
}
