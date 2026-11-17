package modelBasePg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type Vo struct {
	ID          typePg.Int64String `json:"id" form:"id" label:"id" `
	Name        string             `json:"name" form:"name" binding:"required,min=1,max=255" validate:"required,min=1,max=255" label:"名称" `
	NameFl      string             `json:"nameFl" label:"名称外文" `
	Code        string             `json:"code" label:"标记" `
	NameFull    string             `json:"nameFull" label:"全称" `
	Description string             `json:"description" label:"描述" `
	State       typePg.Int8        `json:"state" label:"1有效2停用11取消(对应有效)12弃置(对应停用)13批量删除(无状态)" `
	CreateAt    *typePg.Time       `json:"createAt" label:"创建时间" `
	UpdateAt    *typePg.Time       `json:"updateAt" label:"更新时间" `
	CreateBy    string             `json:"createBy" label:"创建人" `
	UpdateBy    string             `json:"updateBy" label:"更新人" `
}
