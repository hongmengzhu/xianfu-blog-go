package modBasicConfigCategory

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type CreateUpdateCt struct {
	ID          typePg.Uint64String `json:"id" form:"id" label:"id"`
	Name        string              `json:"name" form:"name" validate:"required,min=1,max=255" label:"名称"`
	NameFl      string              `json:"nameFl" label:"名称外文"`
	NameFull    string              `json:"nameFull" label:"全称"`
	Code        string              `json:"code" label:"标志"`
	Description string              `json:"description" label:"描述"`
	ParentId    string              `json:"parentId" label:"上级"`
	ParentNo    string              `json:"parentNo" label:"上级编号"`
	Sort        typePg.Int64String  `json:"sort" label:"排序"`
}
