package modRamMenu

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type QueryPublicCt struct {
	model.BaseQueryNodeCt
	ID          typePg.Uint64String `json:"id" label:"" `
	Name        string              `json:"name" label:"名称" ` // 名称
	Code        string              `json:"code" label:"码值" `
	Description string              `json:"description" label:"描述" `   // 描述
	State       typePg.Int8         `json:"state" label:"状态:1启用;2禁用" ` // 状态:1启用;2禁用
	ParentId    string              `json:"parentId" label:"上级" `
	TypeMenu    string              `json:"typeMenu" label:"菜单类型" `
	ALL         bool                `json:"all" label:"所有数据" `
}
