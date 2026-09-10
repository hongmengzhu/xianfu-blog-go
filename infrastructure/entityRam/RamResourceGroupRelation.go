package entityRam

import (
	"time"
)

// RamResourceGroupRelationEntity 资源组其他组关系
type RamResourceGroupRelationEntity struct {
	ID           int64      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id" comment:"" `
	No           string     `gorm:"column:no;type:varchar(80);index;default:;comment:编号" json:"no" comment:"编号" `
	Name         string     `gorm:"column:name;type:varchar(255);comment:名称" json:"name" comment:"名称" `                                                          // 名称
	NameFl       string     `gorm:"column:name_fl;type:varchar(255);comment:名称外文" json:"name_fl" comment:"名称外文" `                                                // 名称外文
	Code         string     `gorm:"column:code;type:varchar(80);comment:标志" json:"code" comment:"标志" `                                                           // 编号代号
	NameFull     string     `gorm:"column:name_full;type:varchar(255);comment:全称" json:"name_full" comment:"全称" `                                                // 全称
	State        int8       `gorm:"column:state;not null;index;default:1;comment:状态|1启用|2禁用" json:"state" comment:"状态:1启用;2禁用" `                                 // 状态:1启用;2禁用
	Description  string     `gorm:"column:description;type:varchar(255);comment:描述" json:"description" comment:"描述" `                                            // 描述
	CreateAt     *time.Time `gorm:"column:create_at;type:datetime;index;autoCreateTime;default:current_timestamp;comment:创建时间" json:"create_at" comment:"创建时间" ` // 创建时间
	UpdateAt     *time.Time `gorm:"column:update_at;type:datetime;autoUpdateTime;comment:更新时间" json:"update_at" comment:"更新时间" `                                 // 更新时间
	CreateBy     string     `gorm:"column:create_by;type:varchar(80);index;default:;comment:创建人" json:"create_by" comment:"创建人" `
	Sort         int64      `gorm:"column:sort;type:bigint;not null;default:0;;comment:排序" json:"sort" comment:"排序" ` // 排序
	TypeSys      string     `gorm:"column:type_sys;type:varchar(100);index;default:;comment:类型|普通|系统|api" json:"type_sys" comment:"类型;普通;系统;api" `
	TypeCategory string     `gorm:"column:type_category;type:varchar(100);index;default:;comment:类型\\;角色\\;资源组\\;部门" json:"type_category" comment:"类型;角色;资源组;部门" `
	TypeValue    string     `gorm:"column:type_value;type:varchar(100);index;default:;comment:对应类型id" json:"type_value" comment:"对应类型id" `
	TypeDomain   string     `gorm:"column:type_domain;type:varchar(80);index;default:'default';comment:域|平台platform|店铺shop|其他other" json:"type_domain" comment:"域;平台platform;店铺shop;其他other" `
	Mark         string     `gorm:"column:mark;type:varchar(80);index;default:;comment:标记" json:"mark" comment:"标记" `
	GroupNo      string     `gorm:"column:group_no;type:varchar(80);index;default:;comment:资源组编号" json:"group_no" comment:"资源组编号" `
	TenantNo     string     `gorm:"column:tenant_no;type:varchar(80);index;default:;comment:租户编号" json:"tenant_no" comment:"租户编号" ` // 租户
	TerminalCode string     `gorm:"column:terminal_code;type:varchar(80);index;default:;comment:终端类型" json:"terminal_code" comment:"终端类型" `
	SystemNo     string     `gorm:"column:system_no;type:varchar(80);index;default:'default';comment:系统编号" json:"systemNo" comment:"系统编号"`
}

func (*RamResourceGroupRelationEntity) TableName() string {
	return "ram_resource_group_relation"
}

func (*RamResourceGroupRelationEntity) TableComment() string {
	return "资源组其他组关系"
}
