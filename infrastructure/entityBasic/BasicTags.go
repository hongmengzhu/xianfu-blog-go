package entityBasic

import (
	"time"
)

// BasicTagsEntity 标签
type BasicTagsEntity struct {
	ID          int64      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id" comment:"" `
	No          string     `gorm:"column:no;type:varchar(80);index;default:;comment:编号" json:"no" comment:"编号" `
	Name        string     `gorm:"column:name;type:varchar(255);comment:名称" json:"name" comment:"名称" `
	NameFl      string     `gorm:"column:name_fl;type:varchar(255);comment:名称外文" json:"name_fl" comment:"名称外文" `
	NameShort   string     `gorm:"column:name_short;type:varchar(255);comment:名称简称" json:"name_short" comment:"名称简称" `
	Code        string     `gorm:"column:code;type:varchar(80);index;default:;comment:标志" json:"code" comment:"标志" `
	NameFull    string     `gorm:"column:name_full;type:varchar(255);comment:全称" json:"name_full" comment:"全称" `
	State       int8       `gorm:"column:state;not null;index;default:1;comment:1有效2停用11取消(对应有效)12弃置(对应停用)13批量删除(无状态)" json:"state" comment:"1有效2停用11取消(对应有效)12弃置(对应停用)13批量删除(无状态)" `
	Description string     `gorm:"column:description;type:varchar(255);comment:描述" json:"description" comment:"描述" `
	CreateAt    *time.Time `gorm:"column:create_at;type:datetime;index;autoCreateTime;default:current_timestamp;comment:创建时间" json:"create_at" comment:"创建时间" `
	UpdateAt    *time.Time `gorm:"column:update_at;type:datetime;autoUpdateTime;comment:更新时间" json:"update_at" comment:"更新时间" `
	CreateBy    string     `gorm:"column:create_by;type:varchar(80);index;default:;comment:创建人" json:"create_by" comment:"创建人" `
	UpdateBy    string     `gorm:"column:update_by;type:varchar(80);default:;comment:更新人" json:"update_by" comment:"更新人" `
	Sort        int64      `gorm:"column:sort;type:bigint;not null;default:0;;comment:排序" json:"sort" comment:"排序" `
	TenantNo    string     `gorm:"column:tenant_no;type:varchar(80);index;default:;comment:租户编号" json:"tenant_no" comment:"租户编号" `
	ParentId    string     `gorm:"column:parent_id;type:varchar(80);index;comment:上级" json:"parent_id" comment:"上级" `
	IdLink      string     `gorm:"column:id_link;type:varchar(5000);comment:上级" json:"id_link" comment:"上级链" `
	ParentNo    string     `gorm:"column:parent_no;type:varchar(80);index;default:;comment:上级编号" json:"parent_no" comment:"上级编号" `
	NoLink      string     `gorm:"column:no_link;type:text;comment:上级编号" json:"no_link" comment:"上级编号" `
	Type        string     `gorm:"column:type;type:varchar(80);index;default:;comment:类型|" json:"type" comment:"类型" `
	TypeSys     string     `gorm:"column:type_sys;type:varchar(80);index;default:'general';comment:类型|普通|系统;" json:"type_sys" comment:"类型;普通;系统;" `
	Module      string     `gorm:"column:module;type:varchar(80);index;comment:模块" json:"module" comment:"模块" `
	Source      string     `gorm:"column:source;type:varchar(80);index;default:;comment:源" json:"source" comment:"源" `
	CategoryNo  string     `gorm:"column:category_no;type:varchar(80);index;comment:分类编号" json:"category_no" comment:"分类编号" `
}

func (*BasicTagsEntity) TableName() string {
	return "basic_tags"
}

func (*BasicTagsEntity) TableComment() string {
	return "标签"
}
