package entityBlog

import (
	"time"
)

// BlogArticleCategoryEntity -
type BlogArticleCategoryEntity struct {
	ID          int64      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id" comment:"" `
	No          string     `gorm:"column:no;type:varchar(80);default:;comment:编号代号" json:"no" comment:"编号代号" `
	Code        string     `gorm:"column:code;type:varchar(80);comment:标志" json:"code" comment:"标志" `
	Name        string     `gorm:"column:name;type:varchar(255);comment:名称" json:"name" comment:"名称" `
	NameFl      string     `gorm:"column:name_fl;type:varchar(255);comment:名称外文" json:"name_fl" comment:"名称外文" `
	NameFull    string     `gorm:"column:name_full;type:varchar(255);comment:全称" json:"name_full" comment:"全称" `
	State       int8       `gorm:"column:state;not null;index;default:1;comment:1有效2停用11取消(对应有效)12弃置(对应停用)13批量删除(无状态)" json:"state" comment:"1有效2停用11取消(对应有效)12弃置(对应停用)13批量删除(无状态)" `
	Description string     `gorm:"column:description;type:varchar(255);comment:描述" json:"description" comment:"描述" `
	CreateAt    *time.Time `gorm:"column:create_at;type:datetime;index;autoCreateTime;default:current_timestamp;comment:创建时间" json:"create_at" comment:"创建时间" `
	UpdateAt    *time.Time `gorm:"column:update_at;type:datetime;autoUpdateTime;comment:更新时间" json:"update_at" comment:"更新时间" `
	CreateBy    string     `gorm:"column:create_by;type:varchar(80);index;default:'';comment:创建人" json:"createBy" comment:"创建人" `
	UpdateBy    string     `gorm:"column:update_by;type:varchar(80);default:'';comment:更新人" json:"updateBy" comment:"更新人" `
	Sort        int64      `gorm:"column:sort;type:bigint;not null;default:0;;comment:排序" json:"sort" comment:"排序" `
	TenantNo    string     `gorm:"column:tenant_no;type:varchar(80);index;default:;comment:租户编号" json:"tenant_no" comment:"租户编号" `
	ParentId    string     `gorm:"column:parent_id;type:varchar(80);index;comment:上级" json:"parent_id" comment:"上级" `
	ParentNo    string     `gorm:"column:parent_no;type:varchar(80);index;comment:上级" json:"parent_no" comment:"上级" `
	IdLink      string     `gorm:"column:id_link;type:text;comment:上级" json:"id_link" comment:"上级链" `
	NoLink      string     `gorm:"column:no_link;type:text;comment:上级" json:"no_link" comment:"上级链" `
	TypeSys     string     `gorm:"column:type_sys;type:varchar(80);index;default:'general';comment:类型|普通|系统;" json:"type_sys" comment:"类型;普通;系统;" `
	Ano         string     `gorm:"column:ano;type:varchar(80);index;default:'';comment:操作员" json:"ano" comment:"操作员" `
}

func (*BlogArticleCategoryEntity) TableName() string {
	return "blog_article_category"
}

func (*BlogArticleCategoryEntity) TableComment() string {
	return "分类"
}
