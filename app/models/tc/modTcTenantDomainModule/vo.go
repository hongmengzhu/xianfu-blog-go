package modTcTenantDomainModule

import (
	"time"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type Vo struct {
	ID          typePg.Int64String `json:"id" label:"id" `
	No          string             `json:"no" label:"编号代号"`
	Name        string             `json:"name" label:"名称" `          // 名称
	NameFl      string             `json:"nameFl" label:"名称外文" `      // 名称外文
	Code        string             `json:"code" label:"编号代号" `        // 编号代号
	NameFull    string             `json:"nameFull" label:"全称" `      // 全称
	State       typePg.Int8        `json:"state" label:"状态:1启用;2禁用" ` // 状态:1启用;2禁用
	Description string             `json:"description" label:"描述" `   // 描述
	CreateAt    *time.Time         `json:"createAt" label:"创建时间" `    // 创建时间
	CreateBy    string             `json:"createBy" label:"创建人" `     // 创建人
}
