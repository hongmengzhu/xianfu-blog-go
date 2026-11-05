package modBasicConfigCategory

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type QueryPublicCt struct {
	model.BaseQueryNodeCt
	ID          typePg.Int64String `json:"id" label:"id"`
	Name        string             `json:"name" label:"名称"`
	Code        string             `json:"code" label:"标志"`
	Description string             `json:"description" label:"描述"`
	State       typePg.Int8        `json:"state" label:"状态:1启用;2禁用"`
}
