package entityRam

import (
	"time"
)

// RamResourceRelationEntity 资源权限关系
type RamResourceRelationEntity struct {
	ID           int64      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id" comment:"" `
	No           string     `gorm:"column:no;type:varchar(80);index;default:;comment:编号" json:"no" comment:"编号" `
	Code         string     `gorm:"column:code;type:varchar(80);comment:标志" json:"code" comment:"标志" `                                                           // 编号代号
	CreateAt     *time.Time `gorm:"column:create_at;type:datetime;index;autoCreateTime;default:current_timestamp;comment:创建时间" json:"create_at" comment:"创建时间" ` // 创建时间
	CreateBy     string     `gorm:"column:create_by;type:varchar(80);index;default:;comment:创建人" json:"create_by" comment:"创建人" `
	Sort         int64      `gorm:"column:sort;type:bigint;not null;default:0;;comment:排序" json:"sort" comment:"排序" `               // 排序
	TenantNo     string     `gorm:"column:tenant_no;type:varchar(80);index;default:;comment:租户编号" json:"tenant_no" comment:"租户编号" ` // 租户
	OrgNo        string     `gorm:"column:org_no;type:varchar(80);index;default:;comment:组织编号" json:"org_no" comment:"组织编号" `
	StoreNo      string     `gorm:"column:store_no;type:varchar(80);index;default:;comment:店编号" json:"store_no" comment:"店编号" `
	TypeSys      string     `gorm:"column:type_sys;type:varchar(100);index;default:;comment:类型|普通|系统|api" json:"type_sys" comment:"类型;普通;系统;api" `
	TypeCategory string     `gorm:"column:type_category;type:varchar(100);index;default:;comment:类型\\;角色\\;资源组\\;部门" json:"type_category" comment:"类型;角色;资源组;部门" `
	TypeValue    string     `gorm:"column:type_value;type:varchar(80);index;default:;comment:对应类型id" json:"type_value" comment:"对应类型id" `
	TypeAttr     string     `gorm:"column:type_attr;type:varchar(100);index;default:;comment:属性|菜单分类|资源" json:"type_attr" comment:"属性;菜单menu;按钮button;资源:resource;其他other" `
	TypeDomain   string     `gorm:"column:type_domain;type:varchar(80);index;default:'default';comment:域|平台platform|店铺shop|其他other" json:"type_domain" comment:"域;平台platform;店铺shop;其他other" `
	AuthorityId  int64      `gorm:"column:authority_id;type:bigint;not null;index;default:0;comment:权限id" json:"authority_id" comment:"权限id" `
	SourceValue  string     `gorm:"column:source_value;type:varchar(80);index;default:;comment:原始权限id" json:"source_value" comment:"原始权限id" `
	ResourceId   int64      `gorm:"column:resource_id;type:bigint;not null;index;default:0;comment:资源id" json:"resource_id" comment:"资源id" `
	Mark         string     `gorm:"column:mark;type:varchar(80);index;default:;comment:标记" json:"mark" comment:"标记" `
	Path         string     `gorm:"column:path;type:varchar(5000);comment:路径" json:"path" comment:"路径" `
	Method       string     `gorm:"column:method;type:varchar(255);comment:方法" json:"method" comment:"方法" `
}

func (*RamResourceRelationEntity) TableName() string {
	return "ram_resource_relation"
}

func (*RamResourceRelationEntity) TableComment() string {
	return "资源权限关系"
}
