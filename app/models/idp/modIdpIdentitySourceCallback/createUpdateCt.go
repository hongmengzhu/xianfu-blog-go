package modRamIdentitySourceCallback

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type CreateUpdateCt struct {
	ID          typePg.Uint64String `json:"id" form:"id" label:"id" `
	Name        string              `json:"name" form:"name" validate:"required,min=1,max=255" label:"名称" ` // 名称
	Callback    string              `json:"callback" form:"callback" validate:"required,min=1,max=255" label:"回调地址" `
	Description string              `json:"description" label:"描述" ` // 描述
	State       typePg.Int8         `json:"state" label:"状态" `
	Platform    string              `json:"platform" label:"平台" `
	Protocol    string              `json:"protocol" label:"协议" `
	Idp         string              `json:"idp" label:"身份提供商" `
}
