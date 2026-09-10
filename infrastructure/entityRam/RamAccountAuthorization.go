package entityRam

import (
	"time"
)

// RamAccountAuthorizationEntity 账户授权
type RamAccountAuthorizationEntity struct {
	ID          int64      `gorm:"column:id;primaryKey;autoIncrement:true;comment:" json:"id" comment:"" `
	CreateAt    *time.Time `gorm:"column:create_at;type:datetime;index;autoCreateTime;default:current_timestamp;comment:创建时间" json:"create_at" comment:"创建时间" ` // 创建时间
	UpdateAt    *time.Time `gorm:"column:update_at;type:datetime;autoUpdateTime;comment:更新时间;comment:更新时间" json:"update_at" comment:"更新时间" `                    // 更新时间
	CreateBy    string     `gorm:"column:create_by;type:varchar(80);index;default:;comment:创建人" json:"create_by" comment:"创建人" `
	UpdateBy    string     `gorm:"column:update_by;type:varchar(80);default:;comment:更新人" json:"update_by" comment:"更新人" `
	TenantNo    string     `gorm:"column:tenant_no;type:varchar(80);index;default:;comment:租户编号" json:"tenant_no" comment:"租户编号" ` // 租户
	OrgNo       string     `gorm:"column:org_no;type:varchar(80);index;default:;comment:组织编号" json:"org_no" comment:"组织编号" `
	Ano         string     `gorm:"column:ano;type:varchar(80);index;default:;comment:账号" json:"ano" comment:"账号"`
	Value       string     `gorm:"column:value;type:text;comment:值" json:"value" comment:"值" `                                // 值
	Key         string     `gorm:"column:key;type:varchar(80);comment:键" json:"key" comment:"键" `                             // 键
	ExtraData   string     `gorm:"column:extra_data;type:text;comment:额外参数/密钥/干扰码" json:"extraData" label:"额外参数/密钥/干扰码" `     // 额外参数密钥/干扰码
	Type        string     `gorm:"column:type;type:varchar(80);comment:类型" json:"type" comment:"类型" `                         // 类型 密码/pin码/双因子码/openid
	AppNo       string     `gorm:"column:app_no;type:bigint;index;index;default:;comment:应用id" json:"app_no" comment:"应用id" ` // 应用id 微信/钉钉/支付宝/qq
	Description string     `gorm:"column:description;type:varchar(255);comment:描述" json:"description" comment:"描述" `
	KindUnique  string     `gorm:"column:kind_unique;type:varchar(80);index;default:;comment:种类唯一" json:"kind_unique" comment:"种类唯一" `
	App         string     `gorm:"column:app;type:varchar(80);index;default:;comment:应用" json:"app" comment:"应用" `
}

func (*RamAccountAuthorizationEntity) TableName() string {
	return "ram_account_authorization"
}

func (*RamAccountAuthorizationEntity) TableComment() string {
	return "账户授权"
}
