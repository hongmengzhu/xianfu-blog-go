package modRamMenu

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type QueryCt struct {
	model.BaseQueryCt
	ID          typePg.Uint64String `json:"id" label:"" `
	Name        string              `json:"name" label:"名称" `     // 名称
	NameFl      string              `json:"nameFl" label:"名称外文" ` // 名称外文
	Code        string              `json:"code" label:"码值" `
	NameFull    string              `json:"nameFull" label:"全称" `      // 全称
	Description string              `json:"description" label:"描述" `   // 描述
	State       typePg.Int8         `json:"state" label:"状态:1启用;2禁用" ` // 状态:1启用;2禁用
	ParentId    string              `json:"parentId" label:"上级" `
	TypeMenu    string              `json:"typeMenu" label:"菜单类型" `
	AuthCode    string              `json:"authCode" label:"权限标识" `
	ALL         bool                `json:"all" label:"所有数据" `
}
