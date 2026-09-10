package entityBasic

import (
	"time"
)

// BasicAttachmentRelatedEntity 附件关联
type BasicAttachmentRelatedEntity struct {
	ID          int64      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id" comment:"" `
	Name        string     `gorm:"column:name;type:varchar(255);comment:名称" json:"name" comment:"名称" `                                                                                // 名称
	State       int8       `gorm:"column:state;not null;index;default:1;comment:1有效2停用11取消(对应有效)12弃置(对应停用)13批量删除(无状态)" json:"state" comment:"1有效2停用11取消(对应有效)12弃置(对应停用)13批量删除(无状态)" ` // 1有效2停用11取消(对应有效)12弃置(对应停用)13批量删除(无状态)
	Description string     `gorm:"column:description;type:varchar(255);comment:描述" json:"description" comment:"描述" `                                                                  // 描述
	CreateAt    *time.Time `gorm:"column:create_at;type:datetime;index;autoCreateTime;default:current_timestamp;comment:创建时间" json:"create_at" comment:"创建时间" `                       // 创建时间
	UpdateAt    *time.Time `gorm:"column:update_at;type:datetime;autoUpdateTime;comment:更新时间" json:"update_at" comment:"更新时间" `                                                       // 更新时间
	CreateBy    string     `gorm:"column:create_by;type:varchar(80);index;default:;comment:创建人" json:"create_by" comment:"创建人" `
	UpdateBy    string     `gorm:"column:update_by;type:varchar(80);default:;comment:更新人" json:"update_by" comment:"更新人" `
	Sort        int64      `gorm:"column:sort;type:bigint;not null;default:0;;comment:排序" json:"sort" comment:"排序" ` // 排序
	TenantNo    string     `gorm:"column:tenant_no;type:varchar(80);index;default:;comment:租户编号" json:"tenant_no" comment:"租户编号" `
	Size        int64      `gorm:"column:size;type:bigint;not null;default:0;;comment:大小" json:"size" comment:"大小" `
	Module      string     `gorm:"column:module;type:varchar(180);index;comment:模块" json:"module" comment:"模块" `
	Type        string     `gorm:"column:type;type:varchar(80);index;comment:类型" json:"type" comment:"类型" `
	Value       string     `gorm:"column:value;type:varchar(80);comment:值id" json:"value" comment:"值id" `
	Tag         string     `gorm:"column:tag;type:varchar(80);comment:标签" json:"tag" comment:"标签" `
	Mark        string     `gorm:"column:mark;type:varchar(80);index;comment:标记" json:"mark" comment:"标记" `
	Label       string     `gorm:"column:label;type:varchar(80);index;comment:标记" json:"label" comment:"标记" `
	File        string     `gorm:"column:file;type:varchar(980);comment:路径" json:"file" comment:"路径" `
	Http        string     `gorm:"column:http;type:varchar(380);comment:域名" json:"http" comment:"域名" `
}

func (*BasicAttachmentRelatedEntity) TableName() string {
	return "basic_attachment_related"
}

func (*BasicAttachmentRelatedEntity) TableComment() string {
	return "附件关联"
}
