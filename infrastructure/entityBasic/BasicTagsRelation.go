package entityBasic

import (
	"time"
)

// BasicTagsRelationEntity 标签关系
type BasicTagsRelationEntity struct {
	ID           int64      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id" comment:"" `
	Name         string     `gorm:"column:name;type:varchar(255);comment:名称" json:"name" comment:"名称" `
	NameFl       string     `gorm:"column:name_fl;type:varchar(255);comment:名称外文" json:"name_fl" comment:"名称外文" `
	NameShort    string     `gorm:"column:name_short;type:varchar(255);comment:名称简称" json:"name_short" comment:"名称简称" `
	Code         string     `gorm:"column:code;type:varchar(80);index;default:;comment:编号代号" json:"code" comment:"编号代号" `
	NameFull     string     `gorm:"column:name_full;type:varchar(255);comment:全称" json:"name_full" comment:"全称" `
	State        int8       `gorm:"column:state;not null;index;default:1;comment:1有效2停用11取消(对应有效)12弃置(对应停用)13批量删除(无状态)" json:"state" comment:"1有效2停用11取消(对应有效)12弃置(对应停用)13批量删除(无状态)" `
	CreateAt     *time.Time `gorm:"column:create_at;type:datetime;index;autoCreateTime;default:current_timestamp;comment:创建时间" json:"create_at" comment:"创建时间" `
	CreateBy     string     `gorm:"column:create_by;type:varchar(80);index;default:;comment:创建人" json:"create_by" comment:"创建人" `
	Sort         int64      `gorm:"column:sort;type:bigint;not null;default:0;;comment:排序" json:"sort" comment:"排序" `
	TenantNo     string     `gorm:"column:tenant_no;type:varchar(80);index;default:;comment:租户编号" json:"tenant_no" comment:"租户编号" `
	OwnerNo      string     `gorm:"column:owner_no;type:varchar(80);index;default:;comment:所属编号" json:"owner_no" comment:"所属编号" `
	TagNo        string     `gorm:"column:tag_no;type:varchar(80);index;default:;comment:标签编号" json:"tag_no" comment:"标签编号" `
	ParentNo     string     `gorm:"column:parent_no;type:varchar(80);index;default:;comment:上级" json:"parent_no" comment:"上级" `
	NoLink       string     `gorm:"column:no_link;type:text;comment:上级" json:"no_link" comment:"上级链" `
	TypeSys      string     `gorm:"column:type_sys;type:varchar(80);index;default:'general';comment:类型|普通|系统;" json:"type_sys" comment:"类型;普通;系统;" `
	Module       string     `gorm:"column:module;type:varchar(80);index;default:;comment:模块" json:"module" comment:"模块" `
	ActorNo      string     `gorm:"column:actor_no;type:varchar(80);index;default:;comment:执行者编号" json:"actor_no" comment:"执行者编号" `
	CategoryNo   string     `gorm:"column:category_no;type:varchar(80);index;default:;comment:分类编号" json:"category_no" comment:"分类编号" `
	CategoryRoot string     `gorm:"column:category_root;type:varchar(80);index;default:;comment:根分类" json:"category_root" comment:"根分类" `
	Attribute    string     `gorm:"column:attribute;type:text;comment:属性" json:"attribute" comment:"属性" `
	MerchantNo   string     `gorm:"column:merchant_no;type:varchar(80);index;default:;comment:商户" json:"merchant_no" comment:"商户" `
}

func (*BasicTagsRelationEntity) TableName() string {
	return "basic_tags_relation"
}

func (*BasicTagsRelationEntity) TableComment() string {
	return "标签关系"
}
