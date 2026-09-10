package entityBasic

import (
	"time"
)

// BasicAccountApplyDenyListEntity 账号申请被拒绝列表
type BasicAccountApplyDenyListEntity struct {
	ID          int64      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id" comment:"" `
	No          string     `gorm:"column:no;type:varchar(80);index;default:;comment:编号" json:"no" comment:"编号" `
	Name        string     `gorm:"column:name;type:varchar(255);comment:名称" json:"name" comment:"名称" `
	State       int8       `gorm:"column:state;not null;index;default:1;comment:1有效2停用11取消(对应有效)12弃置(对应停用)13批量删除(无状态)" json:"state" comment:"1有效2停用11取消(对应有效)12弃置(对应停用)13批量删除(无状态)" `
	Description string     `gorm:"column:description;type:varchar(255);comment:描述" json:"description" comment:"描述" `
	CreateAt    *time.Time `gorm:"column:create_at;type:datetime;index;autoCreateTime;default:current_timestamp;comment:创建时间" json:"create_at" comment:"创建时间" `
	UpdateAt    *time.Time `gorm:"column:update_at;type:datetime;autoUpdateTime;comment:更新时间" json:"update_at" comment:"更新时间" `
	CreateBy    string     `gorm:"column:create_by;type:varchar(80);index;default:;comment:创建人" json:"create_by" comment:"创建人" `
	UpdateBy    string     `gorm:"column:update_by;type:varchar(80);default:;comment:更新人" json:"update_by" comment:"更新人" `
	Sort        int64      `gorm:"column:sort;type:bigint;not null;default:0;;comment:排序" json:"sort" comment:"排序" `
	TenantNo    string     `gorm:"column:tenant_no;type:varchar(80);index;default:;comment:租户编号" json:"tenant_no" comment:"租户编号" `
	Module      string     `gorm:"column:module;type:varchar(80);index;comment:模块" json:"module" comment:"模块" `
	TypeSys     string     `gorm:"column:type_sys;type:varchar(80);index;default:'general';comment:类型|普通|系统;" json:"type_sys" comment:"类型;普通;系统;" `
	TypeDomain  string     `gorm:"column:type_domain;type:varchar(80);index;default:'general';comment:域类型" json:"type_domain" comment:"域类型系统-商户" `
	TypeField   string     `gorm:"column:type_field;type:varchar(80);index;default:'username';comment:字段类型|字段;" json:"type_field" comment:"类型;字段;" `
	TypeExpr    string     `gorm:"column:type_expr;type:varchar(80);index;default:'general';comment:类型|普通|正则;" json:"type_expr" comment:"类型;普通;正则;" `
	Expr        string     `gorm:"column:expr;type:text;comment:表达式" json:"expr" comment:"表达式" `
	Message     string     `gorm:"column:message;type:text;comment:错误时消息" json:"message" comment:"消息" `
}

func (*BasicAccountApplyDenyListEntity) TableName() string {
	return "basic_account_apply_deny_list"
}

func (*BasicAccountApplyDenyListEntity) TableComment() string {
	return "账号值申请被拒绝列表"
}
