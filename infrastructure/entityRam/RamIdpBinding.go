package entityRam

import "time"

// RamIdpBindingEntity 身份 身份绑定
type RamIdpBindingEntity struct {
	ID            int64      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id" comment:""`
	No            string     `gorm:"column:no;type:varchar(80);index;default:;comment:编号" json:"no" comment:"编号" `
	Description   string     `gorm:"column:description;type:varchar(255);comment:描述" json:"description" comment:"描述" `
	CreateAt      *time.Time `gorm:"column:create_at;type:datetime;index;autoCreateTime;default:current_timestamp;comment:创建时间" json:"create_at" comment:"创建时间" `
	UpdateAt      *time.Time `gorm:"column:update_at;type:datetime;autoUpdateTime;comment:更新时间" json:"update_at" comment:"更新时间" `
	CreateBy      string     `gorm:"column:create_by;type:varchar(80);index;default:;comment:创建人" json:"create_by" comment:"创建人" `
	UpdateBy      string     `gorm:"column:update_by;type:varchar(80);default:;comment:更新人" json:"update_by" comment:"更新人" `
	TenantNo      string     `gorm:"column:tenant_no;type:varchar(80);index;default:;comment:租户编号" json:"tenant_no" comment:"租户编号" ` // 租户
	OrgNo         string     `gorm:"column:org_no;type:varchar(80);index;default:;comment:组织编号" json:"org_no" comment:"组织编号" `
	StoreNo       string     `gorm:"column:store_no;type:varchar(80);index;default:;comment:店编号" json:"store_no" comment:"店编号" `
	TypeDomain    string     `gorm:"column:type_domain;type:varchar(80);index;default:'general';comment:域类型" json:"type_domain" comment:"域类型系统-商户" `
	Idp           string     `gorm:"column:idp;type:varchar(80);index;default:;comment:身份 提供商" json:"idp" comment:"身份 提供商" `
	OpenId        string     `gorm:"column:open_id;type:varchar(255);comment:应用唯一标识openid" json:"open_id" comment:"openid" `
	UnionId       string     `gorm:"column:union_id;type:varchar(255);comment:应用唯一标识unionid" json:"union_id" comment:"unionid" `
	AppMark       string     `gorm:"column:app_mark;type:varchar(80);index;default:;comment:同一IDP下多个应用（如多个小程序）的隔离标识" json:"app_mark" comment:"同一IDP下多个应用（如多个小程序）的隔离标识" `
	AccessToken   string     `gorm:"column:access_token;type:varchar(255);comment:访问令牌" json:"access_token" comment:"访问令牌" `
	RefreshToken  string     `gorm:"column:refresh_token;type:varchar(255);comment:刷新令牌" json:"refresh_token" comment:"刷新令牌" `
	State         int8       `gorm:"column:state;not null;index;default:1;comment:状态|1启用|2禁用" json:"state" comment:"状态:1启用;2禁用" ` // 状态:1启用;2禁用
	StateBind     int8       `gorm:"column:state_bind;type:int8;not null;index;default:1;comment:绑定状态|1未绑定|2已绑定" json:"state_bind" comment:"绑定状态:1未绑定;2已绑定" `
	BindTime      *time.Time `gorm:"column:bind_time;type:datetime;comment:绑定时间" json:"bind_time" comment:"绑定时间" `
	UnBindTime    *time.Time `gorm:"column:un_bind_time;type:datetime;comment:解绑时间" json:"un_bind_time" comment:"解绑时间" `
	LastLoginTime *time.Time `gorm:"column:last_login_time;type:datetime;comment:最后登录时间" json:"last_login_time" comment:"最后登录时间" `
	Platform      string     `gorm:"column:platform;type:varchar(80);index;default:;comment:平台" json:"platform" comment:"平台" `
	Protocol      string     `gorm:"column:protocol;type:varchar(80);index;default:;comment:协议|openid|oauth" json:"protocol" comment:"协议|openid|oauth" `
	Mail          string     `gorm:"column:mail;type:varchar(80);index;default:;comment:邮箱" json:"mail" comment:"邮箱" `
	Phone         string     `gorm:"column:phone;type:varchar(80);index;default:;comment:手机" json:"phone" comment:"手机" `
	NickName      string     `gorm:"column:nick_name;type:varchar(80);index;default:;comment:昵称" json:"nick_name" comment:"昵称" `
	Avatar        string     `gorm:"column:avatar;type:varchar(255);comment:头像" json:"avatar" comment:"头像" `
	ExtraData     string     `gorm:"column:extra_data;type:text;comment:扩展数据" json:"extra_data" comment:"扩展数据" `
	BindAno       string     `gorm:"column:bind_ano;type:varchar(80);index;default:;comment:绑定账号" json:"bind_ano" comment:"绑定账号" `
	ExternalSub   string     `gorm:"column:external_sub;type:varchar(255);index;comment:外部身份唯一标识(sub/NameID/uid)" json:"external_sub" comment:"外部唯一身份标识"`
}

// TableName RamIdpBindingEntity's table name
func (*RamIdpBindingEntity) TableName() string {
	return "ram_idp_binding"
}

func (*RamIdpBindingEntity) TableComment() string {
	return "身份 身份绑定"
}
