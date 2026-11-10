package modBasicTags

import (
	"github.com/hongmengzhu/xianfu-blog-go/app/models/basic/modBasicTagsRelation"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type CreateUpdateCt struct {
	ID           typePg.Uint64String              `json:"id" form:"id" label:"id" `
	Name         string                           `json:"name" form:"name" validate:"required,min=1,max=255" label:"名称" `
	NameFl       string                           `json:"nameFl" label:"名称外文" `
	Code         string                           `json:"code" label:"标记" `
	CategoryNo   string                           `json:"categoryNo" label:"分类"`
	Type         string                           `json:"type" label:"类型" validate:"required,min=1,max=255"`
	NameFull     string                           `json:"nameFull" label:"全称" `    // 全称
	Description  string                           `json:"description" label:"描述" ` // 描述
	AttributeMap modBasicTagsRelation.AttributeVo `json:"attributeMap" label:"属性" `
}
