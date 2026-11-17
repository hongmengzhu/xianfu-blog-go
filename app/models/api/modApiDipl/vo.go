package modApiDipl

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
	"time"
)

type Vo struct {
	ID          typePg.Uint64String `json:"id" label:"id" `
	Name        string              `json:"name" label:"名称" `     // 名称
	NameFl      string              `json:"nameFl" label:"名称外文" ` // 名称外文
	No          string              `json:"no" label:"编号代号"`
	Code        string              `json:"code" label:"码值" `
	NameFull    string              `json:"nameFull" label:"全称" `      // 全称
	State       typePg.Int8         `json:"state" label:"状态:1启用;2禁用" ` // 状态:1启用;2禁用
	Description string              `json:"description" label:"描述" `   // 描述
	CreateAt    *time.Time          `json:"createAt" label:"创建时间" `    // 创建时间
	UpdateAt    *time.Time          `json:"updateAt" label:"更新时间" `    // 更新时间
	CreateBy    string              `json:"createBy" label:"创建人" `     // 创建人
	UpdateBy    string              `json:"updateBy" label:"更新人" `     // 更新人
	OrgId       string              `json:"orgId" label:"组织id" `       // 组织id
	Type        string              `json:"type" label:"类型" binding:"required,min=1,max=255" validate:"required,min=1,max=255"`
	Url         string              `json:"url" label:"建议url" binding:"required,min=1,max=255" validate:"required,min=1,max=255"`
	Category    string              `json:"category" label:"分类" `
	CategoryNo  string              `json:"categoryNo" label:"分类编号" `
}
