package modBasicConfigList

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type CreateUpdateCt struct {
	ID          typePg.Uint64String `json:"id" form:"id" binding:"required" validate:"required" label:"id" `
	Name        string              `json:"name" form:"name" binding:"required,min=1,max=255" validate:"required,min=1,max=255" label:"名称" ` // 名称
	Description string              `json:"description" label:"描述" `                                                                         // 描述
	EventNo     string              `son:"eventNo" comment:"事件编号" `
	Field       string              `json:"field" comment:"字段名称" `
	FieldPath   string              `json:"fieldPath" comment:"路径字段名称" `
	Show        int8                `json:"show" comment:"1显示2隐藏" `
	Content     string              `json:"content" comment:"内容" `
	TypeDomain  string              `json:"typeDomain" comment:"域类型|系统|租户|商户|模块|" `
}
