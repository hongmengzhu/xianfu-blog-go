package modBmBookmarkCategory

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type UpdateCt struct {
	ID          typePg.Uint64String `json:"id" form:"id" validate:"required" label:"id" `
	Name        string              `json:"name" form:"name" validate:"required,min=1,max=255" label:"名称" `
	NameFl      string              `json:"nameFl" label:"名称外文" `
	Code        string              `json:"code" label:"码值" `
	NameFull    string              `json:"nameFull" label:"全称" `
	Description string              `json:"description" label:"描述" `
	ParentNo    string              `json:"parentNo" label:"上级" `
}
