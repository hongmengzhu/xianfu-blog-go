package entityApi

import (
	"time"
)

// ApiDiplAccessKeyEntity key
type ApiDiplAccessKeyEntity struct {
	ID          int64      `gorm:"column:id;primaryKey;autoIncrement:true;comment:" json:"id" comment:"" `
	Name        string     `gorm:"column:name;type:varchar(255);default:;comment:名称" json:"name" comment:"名称" `
	Description string     `gorm:"column:description;type:varchar(255);comment:描述" json:"description" comment:"描述" `
	State       int8       `gorm:"column:state;not null;index;default:1;comment:1有效2停用11取消(对应有效)12弃置(对应停用)13批量删除(无状态)" json:"state" comment:"1有效2停用11取消(对应有效)12弃置(对应停用)13批量删除(无状态)" `
	CreateAt    *time.Time `gorm:"column:create_at;type:datetime;index;autoCreateTime;default:current_timestamp;comment:创建时间" json:"create_at" comment:"创建时间" `
	UpdateAt    *time.Time `gorm:"column:update_at;type:datetime;autoUpdateTime;comment:更新时间;comment:更新时间" json:"update_at" comment:"更新时间" `
	CreateBy    string     `gorm:"column:create_by;type:varchar(80);index;default:;comment:创建人" json:"create_by" comment:"创建人" `
	UpdateBy    string     `gorm:"column:update_by;type:varchar(80);default:;comment:更新人" json:"update_by" comment:"更新人" `
	TenantNo    string     `gorm:"column:tenant_no;type:varchar(80);index;default:;comment:租户编号" json:"tenant_no" comment:"租户编号" `
	Ano         string     `gorm:"column:ano;type:varchar(80);index;default:;comment:账号" json:"ano" comment:"账号" `
	Key         string     `gorm:"column:key;type:varchar(200);comment:键" json:"key" comment:"键" `
	Secret      string     `gorm:"column:secret;type:varchar(200);comment:密钥" json:"secret" comment:"密钥" `
	Type        string     `gorm:"column:type;type:varchar(80);comment:类型" json:"type" comment:"类型" `
	KindUnique  string     `gorm:"column:kind_unique;type:varchar(80);index;default:;comment:种类唯一" json:"kind_unique" comment:"种类唯一" `
	DiplNo      string     `gorm:"column:dipl_no;type:varchar(80);index;default:;comment:应用编号" json:"dipl_no" comment:"应用编号" `
	ExpiryDate  *time.Time `gorm:"column:expiry_date;type:datetime;comment:有效期" json:"expiry_date" comment:"有效期" `
}

func (*ApiDiplAccessKeyEntity) TableName() string {
	return "api_dipl_access_key"
}

func (*ApiDiplAccessKeyEntity) TableComment() string {
	return "key"
}
