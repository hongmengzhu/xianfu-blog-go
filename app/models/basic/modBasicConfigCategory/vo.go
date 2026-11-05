package modBasicConfigCategory

import (
	"time"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type Vo struct {
	ID          typePg.Int64String `json:"id" label:"id"`
	No          string             `json:"no" label:"编号"`
	Code        string             `json:"code" label:"标志"`
	Name        string             `json:"name" label:"名称"`
	NameFl      string             `json:"nameFl" label:"名称外文"`
	NameFull    string             `json:"nameFull" label:"全称"`
	State       typePg.Int8        `json:"state" label:"状态"`
	Description string             `json:"description" label:"描述"`
	CreateAt    *time.Time         `json:"createAt" label:"创建时间"`
	UpdateAt    *time.Time         `json:"updateAt" label:"更新时间"`
	CreateBy    string             `json:"createBy" label:"创建人"`
	UpdateBy    string             `json:"updateBy" label:"更新人"`
	TenantNo    string             `json:"tenantNo" label:"租户编号"`
	ParentId    string             `json:"parentId" label:"上级id"`
	ParentNo    string             `json:"parentNo" label:"上级编号"`
	NoLink      string             `json:"noLink" label:"编号链"`
	Sort        typePg.Int64String `json:"sort" label:"排序"`
	TypeSys     string             `json:"typeSys" label:"类型"`
}
