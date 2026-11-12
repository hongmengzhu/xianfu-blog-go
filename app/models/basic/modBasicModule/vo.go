package modBasicModule

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type Vo struct {
	ID          typePg.Uint64String `json:"id" label:"id" `
	No          string              `json:"no" label:"编号" `
	Code        string              `json:"code" label:"码值" `
	Name        string              `json:"name" label:"名称" `
	NameFl      string              `json:"nameFl" label:"名称外文" `
	NameFull    string              `json:"nameFull" label:"全称" `
	State       typePg.Int8         `json:"state" label:"状态:1启用;2禁用" `
	Description string              `json:"description" label:"描述" `
	CreateAt    *typePg.Time        `json:"createAt" label:"创建时间" `
	UpdateAt    *typePg.Time        `json:"updateAt" label:"更新时间" `
	CreateBy    string              `json:"createBy" label:"创建人" `
	Ano         string              `json:"ano" label:"操作员" `
	ParentNo    string              `json:"parentNo" label:"上级" `
	ParentName  string              `json:"parentNane" label:"上级" `
	Vhost       string              `json:"vhost" label:"虚拟域名" `
}
