package modBasicTagsCategory

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type CreateUpdateCt struct {
	ID          typePg.Uint64String `json:"id" form:"id" label:"id" `
	Name        string              `json:"name" form:"name" validate:"required,min=1,max=255" label:"名称" ` // 名称
	NameFl      string              `json:"nameFl" label:"名称外文" `                                           // 名称外文
	No          string              `json:"no" label:"编号代号"`
	NameFull    string              `json:"nameFull" label:"全称" `    // 全称
	Description string              `json:"description" label:"描述" ` // 描述
	ParentId    string              `json:"parentId" label:"上级" `
	ParentNo    string              `json:"parentNo" comment:"分类" `
	Code        string              `json:"code" label:"码值" `
	Module      string              `json:"module" label:"模块" `
}
