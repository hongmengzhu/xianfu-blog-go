package entityPg

import (
	"time"
)

type BaseTenantEntity struct {
	ID          int64      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id" comment:"" `
	No          string     `gorm:"column:no;type:varchar(80);index;default:;comment:编号代号" json:"no" comment:"编号代号" `
	Code        string     `gorm:"column:code;type:varchar(80);index;default:;comment:标志" json:"code" comment:"标志" `
	Name        string     `gorm:"column:name;type:varchar(255);default:;comment:名称" json:"name" comment:"名称" `
	NameFl      string     `gorm:"column:name_fl;type:varchar(255);default:;comment:名称外文" json:"name_fl" comment:"名称外文" `
	NameFull    string     `gorm:"column:name_full;type:varchar(255);default:;comment:全称" json:"name_full" comment:"全称" `
	State       int8       `gorm:"column:state;not null;index;default:1;comment:1有效2停用11取消(对应有效)12弃置(对应停用)13批量删除(无状态)" json:"state" comment:"1有效2停用11取消(对应有效)12弃置(对应停用)13批量删除(无状态)" ` // 1有效2停用11取消(对应有效)12弃置(对应停用)13批量删除(无状态)
	Description string     `gorm:"column:description;type:varchar(255);default:;comment:描述" json:"description" comment:"描述" `
	CreateAt    *time.Time `gorm:"column:create_at;type:datetime;index;autoCreateTime;default:current_timestamp;comment:创建时间" json:"create_at" comment:"创建时间" `
	UpdateAt    *time.Time `gorm:"column:update_at;type:datetime;autoUpdateTime;comment:更新时间" json:"update_at" comment:"更新时间" `
	CreateBy    string     `gorm:"column:create_by;type:varchar(80);index;default:'';comment:创建人" json:"createBy" comment:"创建人" `
	UpdateBy    string     `gorm:"column:update_by;type:varchar(80);default:'';comment:更新人" json:"updateBy" comment:"更新人" `
	Sort        int64      `gorm:"column:sort;type:bigint;not null;default:0;;comment:排序" json:"sort" comment:"排序" `
	TenantNo    string     `gorm:"column:tenant_no;type:varchar(80);index;default:;comment:租户编号" json:"tenant_no" comment:"租户编号" `
	Ano         string     `gorm:"column:ano;type:varchar(80);index;default:'';comment:操作员" json:"ano" comment:"操作员" `
}
