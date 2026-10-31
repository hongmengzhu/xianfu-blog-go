package modBmBookmarkCategory

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type Vo struct {
	ID          typePg.Uint64String `json:"id" label:"id" `
	Name        string              `json:"name" label:"名称" `
	NameFl      string              `json:"nameFl" label:"名称外文" `
	No          string              `json:"no" label:"编号" `
	Code        string              `json:"code" label:"码值" `
	NameFull    string              `json:"nameFull" label:"全称" `
	State       typePg.Int8         `json:"state" label:"状态:1启用;2禁用" `
	Description string              `json:"description" label:"描述" `
	CreateAt    *typePg.Time        `json:"createAt" label:"创建时间" `
	UpdateAt    *typePg.Time        `json:"updateAt" label:"更新时间" `
	CreateBy    string              `json:"createBy" label:"创建人" `
	Ano         string              `json:"ano" label:"操作员" `
	ParentNo    string              `json:"parentNo" label:"上级" `
	ParentId    string              `json:"parentId" label:"上级" `
	TypeOwner   string              `json:"typeOwner" label:"类型" `
	TeamNo      string              `json:"teamNo" label:"团队编号" `
}
