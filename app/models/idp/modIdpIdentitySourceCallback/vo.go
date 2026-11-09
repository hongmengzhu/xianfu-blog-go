package modRamIdentitySourceCallback

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
	"time"
)

type Vo struct {
	ID          typePg.Uint64String `json:"id" label:"id" `
	Name        string              `json:"name" label:"名称" `     // 名称
	Callback    string              `json:"callback" label:"回调地址" `
	State       typePg.Int8         `json:"state" label:"状态:1启用;2禁用" ` // 状态:1启用;2禁用
	Description string              `json:"description" label:"描述" ` // 描述
	Platform    string              `json:"platform" label:"平台" `
	Protocol    string              `json:"protocol" label:"协议" `
	Idp         string              `json:"idp" label:"身份提供商" `
	CreateAt    *time.Time          `json:"createAt" label:"创建时间" ` // 创建时间
	UpdateAt    *time.Time          `json:"updateAt" label:"更新时间" ` // 更新时间
	CreateBy    string              `json:"createBy" label:"创建人" `  // 创建人
	UpdateBy    string              `json:"updateBy" label:"更新人" `  // 更新人
}
