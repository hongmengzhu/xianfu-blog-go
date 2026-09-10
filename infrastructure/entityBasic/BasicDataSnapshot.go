package entityBasic

import (
	"time"
)

// BasicDataSnapshotEntity 数据快照
type BasicDataSnapshotEntity struct {
	ID       int64      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id" comment:"" `
	CreateAt *time.Time `gorm:"column:create_at;type:datetime;index;autoCreateTime;default:current_timestamp;comment:创建时间" json:"create_at" comment:"创建时间" `
	CreateBy string     `gorm:"column:create_by;type:varchar(80);index;default:;comment:创建人" json:"create_by" comment:"创建人" `
	Name     string     `gorm:"column:name;type:varchar(255);comment:名称" json:"name" comment:"名称" `
	TenantNo string     `gorm:"column:tenant_no;type:varchar(80);index;default:;comment:租户编号" json:"tenant_no" comment:"租户编号" `
	Module   string     `gorm:"column:module;type:varchar(180);index;comment:模块" json:"module" comment:"模块" `
	Version  string     `gorm:"column:version;type:varchar(80);comment:版本" json:"version" comment:"版本" `
	Value    string     `gorm:"column:value;type:varchar(80);index;default:;comment:值" json:"value" comment:"值" `
	Mark     string     `gorm:"column:mark;type:varchar(80);index;default:;comment:标记" json:"mark" comment:"标记" `
	Data     string     `gorm:"column:data;type:text;comment:内容" json:"data" comment:"内容" `
	Extend   string     `gorm:"column:extend;type:text;comment:扩展参数" json:"extend" comment:"扩展参数" `
}

// TableName IamGroupEntity's table name
func (*BasicDataSnapshotEntity) TableName() string {
	return "basic_data_snapshot"
}

func (*BasicDataSnapshotEntity) TableComment() string {
	return "数据快照"
}
