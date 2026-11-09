package modRamIdentityProvider

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type CreateUpdateCt struct {
	ID          typePg.Uint64String `json:"id" form:"id" label:"id" `
	Name        string              `json:"name" form:"name" validate:"required,min=1,max=255" label:"名称" ` // 名称
	NameFl      string              `json:"nameFl" label:"名称外文" `                                           // 名称外文
	Code        string              `json:"code" form:"code"  label:"标志" `
	NameFull    string              `json:"nameFull" label:"全称" `    // 全称
	Description string              `json:"description" label:"描述" ` // 描述
	Icon        string              `json:"icon" label:"图标" `
	Protocol    string              `json:"protocol" label:"协议" `
	Platform    string              `json:"platform" label:"平台" `
	Sort        int64               `json:"sort" label:"排序" `
}
