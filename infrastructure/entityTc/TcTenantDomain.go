package entityTc

import (
	"time"
)

// TcTenantDomainEntity 域名
type TcTenantDomainEntity struct {
	ID            int64      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id" comment:"" `
	No            string     `gorm:"column:no;type:varchar(80);index;default:;comment:编号代号" json:"no" comment:"编号代号" `
	Name          string     `gorm:"column:name;type:varchar(255);comment:名称" json:"name" comment:"名称" `
	Code          string     `gorm:"column:code;type:varchar(80);index;default:;comment:标志" json:"code" comment:"标志" `
	State         int8       `gorm:"column:state;not null;index;default:1;comment:状态|1启用|2禁用" json:"state" comment:"状态:1启用;2禁用" `
	Description   string     `gorm:"column:description;type:varchar(255);comment:描述" json:"description" comment:"描述" `
	CreateAt      *time.Time `gorm:"column:create_at;index;autoCreateTime;default:current_timestamp;comment:创建时间" json:"create_at" comment:"创建时间" `
	UpdateAt      *time.Time `gorm:"column:update_at;autoUpdateTime;comment:更新时间" json:"update_at" comment:"更新时间" `
	CreateBy      string     `gorm:"column:create_by;type:varchar(80);index;default:'';comment:创建人" json:"createBy" comment:"创建人" `
	UpdateBy      string     `gorm:"column:update_by;type:varchar(80);default:'';comment:更新人" json:"updateBy" comment:"更新人" `
	Sort          int64      `gorm:"column:sort;type:bigint;not null;default:0;;comment:排序" json:"sort" comment:"排序" `
	TenantNo      string     `gorm:"column:tenant_no;type:varchar(80);index;default:;comment:租户编号" json:"tenant_no" comment:"租户编号" `
	Defaulted     int8       `gorm:"column:defaulted;type:int8;not null;index:defaulted,priority:2;default:2;comment:默认1是2否" json:"defaulted" comment:"默认1是2否" `
	DomainType    string     `gorm:"column:domain_type;type:varchar(80);index;default:;comment:域名类型|租户|组织|客户|租户模块" json:"domain_type" comment:"域名类型" `
	DomainTypeSub string     `gorm:"column:domain_type_sub;type:varchar(80);index;default:;comment:域名类型子类|普通|模块" json:"domain_type_sub" comment:"域名类型子类" `
	Module        string     `gorm:"column:module;type:varchar(80);index;default:;comment:模块" json:"module" comment:"模块" `
	Theme         string     `gorm:"column:theme;type:varchar(80);index;default:;comment:主题" json:"theme" comment:"主题" `
}

func (*TcTenantDomainEntity) TableName() string {
	return "tc_tenant_domain"
}

func (*TcTenantDomainEntity) TableComment() string {
	return "域名"
}
