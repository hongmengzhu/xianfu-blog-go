package entityRam

import "time"

// RamIdentitySourceCallbackEntity 认证源回调白名单
type RamIdentitySourceCallbackEntity struct {
	ID          int64      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id" comment:""`
	No          string     `gorm:"column:no;type:varchar(80);index;default:;comment:编号" json:"no" comment:"编号" `
	TenantNo    string     `gorm:"column:tenant_no;type:varchar(80);index;default:;comment:租户编号" json:"tenant_no" comment:"租户编号" ` // 租户
	OrgNo       string     `gorm:"column:org_no;type:varchar(80);index;default:;comment:组织编号" json:"org_no" comment:"组织编号" `
	StoreNo     string     `gorm:"column:store_no;type:varchar(80);index;default:;comment:店编号" json:"store_no" comment:"店编号" `
	CreateAt    *time.Time `gorm:"column:create_at;type:datetime;index;autoCreateTime;default:current_timestamp;comment:创建时间" json:"create_at" comment:"创建时间" ` // 创建时间
	UpdateAt    *time.Time `gorm:"column:update_at;type:datetime;autoUpdateTime;comment:更新时间" json:"update_at" comment:"更新时间" `                                 // 更新时间
	CreateBy    string     `gorm:"column:create_by;type:varchar(80);index;default:;comment:创建人" json:"create_by" comment:"创建人" `
	UpdateBy    string     `gorm:"column:update_by;type:varchar(80);default:;comment:更新人" json:"update_by" comment:"更新人" `
	Callback    string     `gorm:"column:callback;type:varchar(255);comment:回调地址" json:"callback" comment:"回调地址" `
	Name        string     `gorm:"column:name;type:varchar(255);comment:名称" json:"name" comment:"名称" `                          // 名称
	State       int8       `gorm:"column:state;not null;index;default:1;comment:状态|1启用|2禁用" json:"state" comment:"状态:1启用;2禁用" ` // 状态:1启用;2禁用
	Description string     `gorm:"column:description;type:varchar(255);comment:描述" json:"description" comment:"描述" `            // 描述
	Platform    string     `gorm:"column:platform;type:varchar(80);index;default:;comment:平台" json:"platform" comment:"平台" `
	Protocol    string     `gorm:"column:protocol;type:varchar(80);index;default:;comment:协议|openid|oauth|oidc,oauth2,saml2,cas,ldap" json:"protocol" comment:"协议|openid|oauth" `
	Idp         string     `gorm:"column:idp;type:varchar(80);index;default:;comment:身份 提供商" json:"idp" comment:"身份 提供商" `
}

// TableName RamIdentitySourceCallbackEntity's table name
func (*RamIdentitySourceCallbackEntity) TableName() string {
	return "ram_identity_source_callback"
}

func (*RamIdentitySourceCallbackEntity) TableComment() string {
	return "认证源回调白名单"
}
