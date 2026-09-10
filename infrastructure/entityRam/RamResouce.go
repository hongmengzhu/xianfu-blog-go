package entityRam

import (
	"time"
)

// RamResourceEntity 资源
type RamResourceEntity struct {
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
	Sort         int64      `gorm:"column:sort;type:bigint;not null;default:0;;comment:排序" json:"sort" comment:"排序" `                                                        // 排序
	TypeSys      string     `gorm:"column:type_sys;type:varchar(100);index;default:;comment:类型|普通|系统|api" json:"type_sys" comment:"类型;普通;系统;api" `                           // 类型
	TypeAttr     string     `gorm:"column:type_attr;type:varchar(100);index;default:;comment:属性|菜单分类|资源" json:"type_attr" comment:"属性;菜单menu;按钮button;资源:resource;其他other" ` // 类型
	TypeDomain   string     `gorm:"column:type_domain;type:varchar(80);index;default:'default';comment:域|平台platform|店铺shop|其他other" json:"type_domain" comment:"域;平台platform;店铺shop;其他other" `
	Path         string     `gorm:"column:path;type:varchar(5000);comment:路径" json:"path" comment:"路径" `
	Method       string     `gorm:"column:method;type:varchar(255);comment:方法" json:"method" comment:"方法" `
	MenuNo       string     `gorm:"column:menu_no;type:varchar(80);index;default:;comment:菜单编号" json:"menu_no" comment:"菜单编号" `
	ParentId     string     `gorm:"column:parent_id;type:varchar(80);index;default:;comment:上级" json:"parent_id" comment:"上级" `
	IdLink       string     `gorm:"column:id_link;type:text;comment:上级" json:"parent_id_link" comment:"上级" `
	ParentNo     string     `gorm:"column:parent_no;type:varchar(80);index;default:;comment:上级编号" json:"parent_no" comment:"上级编号" `
	NoLink       string     `gorm:"column:no_link;type:text;comment:上级编号" json:"no_link" comment:"上级编号" `
	TerminalCode string     `gorm:"column:terminal_code;type:varchar(80);index;default:;comment:终端类型" json:"terminal_code" comment:"终端类型" `
	SystemNo     string     `gorm:"column:system_no;type:varchar(80);index;default:'default';comment:系统编号" json:"systemNo" comment:"系统编号"`
}

// TableName RamRoleEntity's table name
func (*RamResourceEntity) TableName() string {
	return "ram_resource"
}

func (*RamResourceEntity) TableComment() string {
	return "资源"
}
