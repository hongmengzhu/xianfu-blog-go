package entityRam

import (
	"time"

	"gorm.io/datatypes"
)

// RamAccountLoginLogEntity 登陆日志
type RamAccountLoginLogEntity struct {
	ID          int64                                `gorm:"column:id;primaryKey;autoIncrement:true" json:"id" comment:"" `
	CreateAt    *time.Time                           `gorm:"column:create_at;type:datetime;index;autoCreateTime;default:current_timestamp;comment:创建时间" json:"create_at" comment:"创建时间" `
	Ano         string                               `gorm:"column:ano;type:varchar(80);index;default:;comment:账号" json:"ano" comment:"账号"`
	TenantNo    string                               `gorm:"column:tenant_no;type:varchar(80);index;default:;comment:租户编号" json:"tenant_no" comment:"租户编号" `
	OrgNo       string                               `gorm:"column:org_no;type:varchar(80);index;default:;comment:组织编号" json:"org_no" comment:"组织编号" `
	StoreNo     string                               `gorm:"column:store_no;type:varchar(80);index;default:;comment:店编号" json:"store_no" comment:"店编号" `
	LoginSource string                               `gorm:"column:login_source;type:varchar(80);index;default:;comment:登陆来源" json:"login_source" comment:"登陆来源" `
	AppNo       string                               `gorm:"column:app_no;type:varchar(80);index;default:;comment:应用编号" json:"app_no" comment:"应用编号" `
	Client      string                               `gorm:"column:client;type:varchar(80);index;comment:客户端" json:"client" comment:"客户端" `
	Os          datatypes.JSONType[RamAccountJsonOs] `gorm:"column:os;type:jsonb;index;default:'{}';comment:组织架构" json:"os" comment:"组织架构" `
	UserAgent   string                               `gorm:"column:user_agent;type:text;comment:用户代理" json:"user_agent" label:"用户代理" `
	System      string                               `gorm:"column:system;type:varchar(80);index;default:;comment:系统" json:"system" comment:"系统"`
	Device      string                               `gorm:"column:device;type:varchar(80);index;default:'';comment:设备" json:"device" comment:"设备" `
	DeviceNo    string                               `gorm:"column:device_no;type:varchar(80);index;default:'';comment:设备编号" json:"device_no" comment:"设备编号" `
	Version     string                               `gorm:"column:version;type:varchar(80);default:'';comment:版本" json:"version" comment:"版本" `
	Ip          string                               `gorm:"column:ip;type:varchar(80);index;default:'';comment:ip" json:"ip" comment:"ip" `
}

// TableName RamAccountLoginLog's table name
func (*RamAccountLoginLogEntity) TableName() string {
	return "ram_account_login_log"
}

func (*RamAccountLoginLogEntity) TableComment() string {
	return "登陆日志"
}
