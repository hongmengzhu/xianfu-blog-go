package modBasicTagsRelation

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type AllCt struct {
	model.BaseQueryCt
	ID          typePg.Uint64String `json:"id" label:"" `
	OrgId       string              `json:"orgId" label:"组织id" `       // 组织id
	Name        string              `json:"name" label:"名称" `          // 名称
	NameFl      string              `json:"nameFl" label:"名称外文" `      // 名称外文
	Code        string              `json:"code" label:"编号代号" `        // 编号代号
	NameFull    string              `json:"nameFull" label:"全称" `      // 全称
	Description string              `json:"description" label:"描述" `   // 描述
	State       typePg.Int8         `json:"state" label:"状态:1启用;2禁用" ` // 状态:1启用;2禁用
	CreateBy    string              `json:"createBy" label:"创建人" `     // 创建人
	CategoryNo  string              `json:"categoryNo" label:"分类" validate:"required,min=1,max=255"`
	Module      string              `json:"module" label:"模块" `
	ParentCode  string              `json:"parentCode" label:"上级"`
}
