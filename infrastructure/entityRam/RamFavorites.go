package entityRam

import (
	"time"
)

// RamFavoritesEntity 收藏
type RamFavoritesEntity struct {
	ID       int64      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id" comment:"" `
	Type     string     `gorm:"column:type;type:varchar(80);index;comment:类型" json:"type" comment:"类型" `
	Module   string     `gorm:"column:module;type:varchar(180);index;comment:模块" json:"module" comment:"模块" `
	CreateAt *time.Time `gorm:"column:create_at;type:datetime;index;autoCreateTime;default:current_timestamp;comment:创建时间" json:"create_at" comment:"创建时间" ` // 创建时间
	Value    string     `gorm:"column:value;type:varchar(80);comment:值id" json:"value" comment:"值id" `
	TenantNo string     `gorm:"column:tenant_no;type:varchar(80);index;default:;comment:租户编号" json:"tenant_no" comment:"租户编号" ` // 租户
	OrgNo    string     `gorm:"column:org_no;type:varchar(80);index;default:;comment:组织编号" json:"org_no" comment:"组织编号" `
	StoreNo  string     `gorm:"column:store_no;type:varchar(80);index;default:;comment:店编号" json:"store_no" comment:"店编号" `
	Name     string     `gorm:"column:name;type:varchar(255);comment:名称" json:"name" comment:"名称" ` // 名称
	AId      int64      `gorm:"column:aid;type:bigint;not null;index;default:0;comment:账号id" json:"aid" comment:"账号id" `
}

// TableName RamFavoritesEntity's table name
func (*RamFavoritesEntity) TableName() string {
	return "ram_favorites"
}

func (*RamFavoritesEntity) TableComment() string {
	return "收藏"
}
