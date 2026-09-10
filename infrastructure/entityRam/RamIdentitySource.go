package entityRam

import (
	"time"
)

// RamIdentitySourceEntity 身份 认证源
type RamIdentitySourceEntity struct {
	ID                 int64      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id" comment:""`
	No                 string     `gorm:"column:no;type:varchar(80);index;default:;comment:编号" json:"no" comment:"编号" `
	Name               string     `gorm:"column:name;type:varchar(255);comment:名称" json:"name" comment:"名称" `                                                          // 名称
	NameFl             string     `gorm:"column:name_fl;type:varchar(255);comment:名称外文" json:"name_fl" comment:"名称外文" `                                                // 名称外文
	Code               string     `gorm:"column:code;type:varchar(80);comment:标志" json:"code" comment:"标志" `                                                           // 编号代号
	NameFull           string     `gorm:"column:name_full;type:varchar(255);comment:全称" json:"name_full" comment:"全称" `                                                // 全称
	State              int8       `gorm:"column:state;not null;index;default:1;comment:状态|1启用|2禁用" json:"state" comment:"状态:1启用;2禁用" `                                 // 状态:1启用;2禁用
	Description        string     `gorm:"column:description;type:varchar(255);comment:描述" json:"description" comment:"描述" `                                            // 描述
	CreateAt           *time.Time `gorm:"column:create_at;type:datetime;index;autoCreateTime;default:current_timestamp;comment:创建时间" json:"create_at" comment:"创建时间" ` // 创建时间
	UpdateAt           *time.Time `gorm:"column:update_at;type:datetime;autoUpdateTime;comment:更新时间" json:"update_at" comment:"更新时间" `                                 // 更新时间
	CreateBy           string     `gorm:"column:create_by;type:varchar(80);index;default:;comment:创建人" json:"create_by" comment:"创建人" `
	UpdateBy           string     `gorm:"column:update_by;type:varchar(80);default:;comment:更新人" json:"update_by" comment:"更新人" `
	TenantNo           string     `gorm:"column:tenant_no;type:varchar(80);index;default:;comment:租户编号" json:"tenant_no" comment:"租户编号" ` // 租户
	OrgNo              string     `gorm:"column:org_no;type:varchar(80);index;default:;comment:组织编号" json:"org_no" comment:"组织编号" `
	StoreNo            string     `gorm:"column:store_no;type:varchar(80);index;default:;comment:店编号" json:"store_no" comment:"店编号" `
	Sort               int64      `gorm:"column:sort;type:bigint;not null;default:0;;comment:排序" json:"sort" comment:"排序" `
	Platform           string     `gorm:"column:platform;type:varchar(80);index;default:;comment:平台" json:"platform" comment:"平台" `
	Icon               string     `gorm:"column:icon;type:varchar(255);comment:图标" json:"icon" comment:"图标" `
	Protocol           string     `gorm:"column:protocol;type:varchar(80);index;default:;comment:协议|openid|oauth|oidc,oauth2,saml2,cas,ldap" json:"protocol" comment:"协议|openid|oauth" `
	Idp                string     `gorm:"column:idp;type:varchar(80);index;default:;comment:身份 提供商" json:"idp" comment:"身份 提供商" `
	BaseConfig         string     `gorm:"column:base_config;type:text;comment:基础配置" json:"base_config" comment:"基础配置" `
	SyncStrategyConfig string     `gorm:"column:sync_strategy_config;type:text;comment:同步策略" json:"sync_strategy_config" comment:"同步策略" `
	Configured         int8       `gorm:"column:configured;type:int8;not null;index;default:2;comment:是否配置|2未配置|1已配置" json:"configured" comment:"是否配置:2未配置;1已配置" `
	// 单点登出开关
	EnableSlo int8 `gorm:"column:enable_slo;type:int8;default:2;comment:是否启用单点登出 1启用 2关闭" json:"enable_slo"`
	// 用户自动创建策略（外部登录不存在本地账号时）
	AutoCreateUser int8 `gorm:"column:auto_create_user;type:int8;default:2;comment:自动创建本地用户1开启2关闭" json:"auto_create_user"`
	// 属性映射配置（外部身份字段映射本地user字段，json）
	AttributeMapping string `gorm:"column:attribute_mapping;type:text;comment:身份属性映射规则" json:"attribute_mapping"`
}

// TableName RamIdentitySourceEntity's table name
func (*RamIdentitySourceEntity) TableName() string {
	return "ram_identity_source"
}

func (*RamIdentitySourceEntity) TableComment() string {
	return "身份 认证源"
}
