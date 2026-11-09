package modRamIdentitySource

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type CreateUpdateCt struct {
	ID                 typePg.Uint64String `json:"id" form:"id" label:"id" `
	Name               string              `json:"name" form:"name" validate:"required,min=1,max=255" label:"名称" ` // 名称
	NameFl             string              `json:"nameFl" label:"名称外文" `                                           // 名称外文
	Code               string              `json:"code" form:"code"  label:"标志" `
	NameFull           string              `json:"nameFull" label:"全称" `    // 全称
	Description        string              `json:"description" label:"描述" ` // 描述
	Platform           string              `json:"platform" label:"平台" `
	Icon               string              `json:"icon" label:"图标" `
	Protocol           string              `json:"protocol" label:"协议" `
	Idp                string              `json:"idp" label:"身份提供商" `
	BaseConfig         string              `json:"baseConfig" label:"基础配置" `
	SyncStrategyConfig string              `json:"syncStrategyConfig" label:"同步策略" `
	EnableSlo          typePg.Int8         `json:"enableSlo" label:"单点登出" `
	AutoCreateUser     typePg.Int8         `json:"autoCreateUser" label:"自动创建用户" `
	AttributeMapping   string              `json:"attributeMapping" label:"属性映射" `
	Sort               int64               `json:"sort" label:"排序" `
}
