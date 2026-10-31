package modBmBookmarkCategory

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type QueryPublicCt struct {
	model.BaseQueryNodeCt
	ID          typePg.Uint64String `json:"id" label:"" `
	Name        string              `json:"name" label:"名称" `
	Description string              `json:"description" label:"描述" `
	State       typePg.Int8         `json:"state" label:"状态:1启用;2禁用" `
	No          string              `json:"no" label:"编号代号" `
	ParentNo    string              `json:"parentNo" label:"上级" `
}
