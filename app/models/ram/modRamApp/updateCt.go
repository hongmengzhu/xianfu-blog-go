package modRamApp

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type UpdateCt struct {
	ID          typePg.Uint64String `json:"id" form:"id" binding:"required" validate:"required" label:"id" `
	Name        string              `json:"name" form:"name" binding:"required,min=1,max=255" validate:"required,min=1,max=255" label:"名称" `
	NameFl      string              `json:"nameFl" label:"名称外文" `
	No          string              `json:"no" form:"no" label:"编号代号"`
	Code        string              `json:"code" label:"码值" `
	CategoryNo  string              `json:"categoryNo" label:"分类编号" `
	Type        string              `json:"type" label:"类型" binding:"required,min=1,max=255" validate:"required,min=1,max=255"`
	Url         string              `json:"url" label:"建议url" binding:"required,min=1,max=255" validate:"required,min=1,max=255"`
	NameFull    string              `json:"nameFull" label:"全称" `    // 全称
	Description string              `json:"description" label:"描述" ` // 描述
}
