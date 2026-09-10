package entityRam

import (
	"time"
)

// RamAccountDenyListEntity 拒绝名单
type RamAccountDenyListEntity struct {
	ID          int64      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id" comment:"" `
	Name        string     `gorm:"column:name;type:varchar(255);comment:名称" json:"realName" comment:"名称" `
	Description string     `gorm:"column:description;type:varchar(255);comment:描述" json:"description" comment:"描述" `
	CreateAt    time.Time  `gorm:"column:create_at;type:datetime;index;autoCreateTime;default:current_timestamp;comment:创建时间" json:"create_at" comment:"创建时间" `
	UpdateAt    time.Time  `gorm:"column:update_at;type:datetime;autoUpdateTime;comment:更新时间" json:"update_at" comment:"更新时间" `
	CreateBy    string     `gorm:"column:create_by;type:varchar(80);index;default:;comment:创建人" json:"create_by" comment:"创建人" `
	UpdateBy    string     `gorm:"column:update_by;type:varchar(80);default:;comment:更新人" json:"update_by" comment:"更新人" `
	State       int8       `gorm:"column:state;index;default:1;comment:1有效2停用11取消(对应有效)12弃置(对应停用)13批量删除(无状态)" json:"state" comment:"1有效2停用11取消(对应有效)12弃置(对应停用)13批量删除(无状态)" `
	TenantNo    string     `gorm:"column:tenant_no;type:varchar(80);index;default:;comment:租户编号" json:"tenant_no" comment:"租户编号" `
	LoginSource string     `gorm:"column:login_source;type:varchar(80);index;default:;comment:登陆来源" json:"login_source" comment:"登陆来源" `
	Ano         string     `gorm:"column:ano;type:varchar(80);index;default:;comment:账号" json:"ano" comment:"账号"`
	Client      string     `gorm:"column:client;type:varchar(80);index;comment:客户端" json:"client" comment:"客户端" `
	System      string     `gorm:"column:system;type:varchar(80);index;default:;comment:系统" json:"system" comment:"系统"`
	AppNo       string     `gorm:"column:app_no;type:varchar(80);index;default:;comment:应用编号" json:"app_no" comment:"应用编号" `
	Device      string     `gorm:"column:device;type:varchar(80);index;default:'';comment:设备" json:"device" comment:"设备" `
	DeviceNo    string     `gorm:"column:device_no;type:varchar(80);index;default:'';comment:设备编号" json:"device_no" comment:"设备编号" `
	StartAt     *time.Time `gorm:"column:start_at;type:datetime;comment:活动开始时间" json:"start_at" comment:"开始时间" `
	EndAt       *time.Time `gorm:"column:end_at;type:datetime;comment:活动结束时间" json:"end_at" comment:"结束时间" `
	Version     string     `gorm:"column:version;type:varchar(80);default:'';comment:版本" json:"version" comment:"版本" `
}

func (*RamAccountDenyListEntity) TableName() string {
	return "ram_account_deny_list"
}

func (*RamAccountDenyListEntity) TableComment() string {
	return "拒绝名单"
}
