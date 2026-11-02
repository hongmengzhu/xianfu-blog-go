package modRamDepartment

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type QueryPublicCt struct {
	model.BaseQueryNodeCt
	ID       typePg.Uint64String `json:"id" label:"" `
	OrgId    string              `json:"orgId" label:"组织id" ` // 组织id
	Name     string              `json:"name" label:"名称" `    // 名称
	Code     string              `json:"code" label:"码值" `
	State    typePg.Int8         `json:"state" label:"状态:1启用;2禁用" ` // 状态:1启用;2禁用
	CreateBy string              `json:"createBy" label:"创建人" `     // 创建人
	ParentId string              `json:"parentId" label:"上级" `
	No       string              `json:"no" label:"编号" `
}
