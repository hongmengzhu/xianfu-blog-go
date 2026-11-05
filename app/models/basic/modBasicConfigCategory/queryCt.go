package modBasicConfigCategory

import (
	"time"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type QueryCt struct {
	model.BaseQueryCt
	ID          typePg.Int64String `json:"id" label:"id"`
	Name        string             `json:"name" label:"名称"`
	NameFl      string             `json:"nameFl" label:"名称外文"`
	NameFull    string             `json:"nameFull" label:"全称"`
	Code        string             `json:"code" label:"标志"`
	Description string             `json:"description" label:"描述"`
	State       typePg.Int8        `json:"state" label:"状态:1启用;2禁用"`
	CreateBy    string             `json:"createBy" label:"创建人"`
	CreateAt    *time.Time         `json:"createAt" label:"创建时间"`
	ParentId    string             `json:"parentId" label:"上级id"`
	ParentNo    string             `json:"parentNo" label:"上级编号"`
}
