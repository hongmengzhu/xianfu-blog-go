package entityRam

import (
	"time"
)

// RamAccountDeviceEntity 设备
type RamAccountDeviceEntity struct {
	ID          int64      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id" comment:"" `
	No          string     `gorm:"column:no;type:varchar(80);index;default:;comment:编号" json:"no" comment:"编号" `
	Name        string     `gorm:"column:name;type:varchar(255);comment:名称" json:"name" comment:"名称" `
	NameFl      string     `gorm:"column:name_fl;type:varchar(255);comment:名称外文" json:"name_fl" comment:"名称外文" `
	Code        string     `gorm:"column:code;type:varchar(80);index;default:;comment:标志" json:"code" comment:"标志" `
	NameFull    string     `gorm:"column:name_full;type:varchar(255);comment:全称" json:"name_full" comment:"全称" `
	State       int8       `gorm:"column:state;not null;index;default:1;comment:1有效2停用11取消(对应有效)12弃置(对应停用)13批量删除(无状态)" json:"state" comment:"1有效2停用11取消(对应有效)12弃置(对应停用)13批量删除(无状态)" `
	Description string     `gorm:"column:description;type:varchar(255);comment:描述" json:"description" comment:"描述" `
	CreateAt    *time.Time `gorm:"column:create_at;type:datetime;index;autoCreateTime;default:current_timestamp;comment:创建时间" json:"create_at" comment:"创建时间" `
	UpdateAt    *time.Time `gorm:"column:update_at;type:datetime;autoUpdateTime;comment:更新时间" json:"update_at" comment:"更新时间" `
	Sort        int64      `gorm:"column:sort;type:bigint;not null;default:0;;comment:排序" json:"sort" comment:"排序" `
	TenantNo    string     `gorm:"column:tenant_no;type:varchar(80);index;default:;comment:租户编号" json:"tenant_no" comment:"租户编号" `
	Ano         string     `gorm:"column:ano;type:varchar(80);index;default:;comment:账号" json:"ano" comment:"账号"`
	Version     string     `gorm:"column:version;type:varchar(80);default:'';comment:版本" json:"version" comment:"版本" `
	UserAgent   string     `gorm:"column:user_agent;type:text;comment:用户代理" json:"user_agent" label:"用户代理" `
	System      string     `gorm:"column:system;type:varchar(80);index;default:;comment:系统" json:"system" comment:"系统"`
}

func (*RamAccountDeviceEntity) TableName() string {
	return "ram_account_device"
}

func (*RamAccountDeviceEntity) TableComment() string {
	return "设备"
}
