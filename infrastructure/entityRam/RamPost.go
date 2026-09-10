package entityRam

import (
	"time"
)

// RamPostEntity 岗位
type RamPostEntity struct {
	ID          int64      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id" comment:"" `
	No          string     `gorm:"column:no;type:varchar(80);index;default:;comment:编号" json:"no" comment:"编号" `
	Name        string     `gorm:"column:name;type:varchar(255);comment:名称" json:"name" comment:"名称" `                                                          // 名称
	NameFl      string     `gorm:"column:name_fl;type:varchar(255);comment:名称外文" json:"name_fl" comment:"名称外文" `                                                // 名称外文
	Code        string     `gorm:"column:code;type:varchar(80);comment:标志" json:"code" comment:"标志" `                                                           // 编号代号
	NameFull    string     `gorm:"column:name_full;type:varchar(255);comment:全称" json:"name_full" comment:"全称" `                                                // 全称
	State       int8       `gorm:"column:state;not null;index;default:1;comment:状态|1启用|2禁用" json:"state" comment:"状态:1启用;2禁用" `                                 // 状态:1启用;2禁用
	Description string     `gorm:"column:description;type:varchar(255);comment:描述" json:"description" comment:"描述" `                                            // 描述
	CreateAt    *time.Time `gorm:"column:create_at;type:datetime;index;autoCreateTime;default:current_timestamp;comment:创建时间" json:"create_at" comment:"创建时间" ` // 创建时间
	UpdateAt    *time.Time `gorm:"column:update_at;type:datetime;autoUpdateTime;comment:更新时间" json:"update_at" comment:"更新时间" `                                 // 更新时间
	CreateBy    string     `gorm:"column:create_by;type:varchar(80);index;default:;comment:创建人" json:"create_by" comment:"创建人" `
	UpdateBy    string     `gorm:"column:update_by;type:varchar(80);default:;comment:更新人" json:"update_by" comment:"更新人" `
	TenantNo    string     `gorm:"column:tenant_no;type:varchar(80);index;default:;comment:租户编号" json:"tenant_no" comment:"租户编号" ` // 租户
	OrgNo       string     `gorm:"column:org_no;type:varchar(80);index;default:;comment:组织编号" json:"org_no" comment:"组织编号" `
	StoreNo     string     `gorm:"column:store_no;type:varchar(80);index;default:;comment:店编号" json:"store_no" comment:"店编号" `
	Sort        int64      `gorm:"column:sort;type:bigint;not null;default:0;;comment:排序" json:"sort" comment:"排序" ` // 排序
}

// TableName RamPostEntity's table name
func (*RamPostEntity) TableName() string {
	return "ram_post"
}

func (*RamPostEntity) TableComment() string {
	return "岗位"
}
