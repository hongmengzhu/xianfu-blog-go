package entityBasic

import (
	"time"

	"gorm.io/datatypes"
)

type BasicConfigEntity struct {
	ID         int64      `gorm:"column:id;primaryKey;autoIncrement:true;comment:" json:"id" comment:"" `
	CreateAt   *time.Time `gorm:"column:create_at;type:datetime;index;autoCreateTime;default:current_timestamp;comment:创建时间" json:"create_at" comment:"创建时间" `
	UpdateAt   *time.Time `gorm:"column:update_at;type:datetime;autoUpdateTime;comment:更新时间;comment:更新时间" json:"update_at" comment:"更新时间" `
	CreateBy   string     `gorm:"column:create_by;type:varchar(80);index;default:;comment:创建人" json:"create_by" comment:"创建人" `
	UpdateBy   string     `gorm:"column:update_by;type:varchar(80);default:;comment:更新人" json:"update_by" comment:"更新人" `
	State      int8       `gorm:"column:state;not null;index;default:1;comment:1有效2停用" json:"state" comment:"1有效2停用" `
	Sort       int64      `gorm:"column:sort;type:bigint;not null;index;default:0;;comment:排序" json:"sort" comment:"排序" `
	TenantNo   string     `gorm:"column:tenant_no;type:varchar(80);index;default:;comment:租户编号" json:"tenant_no" comment:"租户编号" `
	OrgNo      string     `gorm:"column:org_no;type:varchar(80);index;default:;comment:组织编号" json:"org_no" comment:"组织编号" `
	StoreNo    string     `gorm:"column:store_no;type:varchar(80);index;default:;comment:店铺编号" json:"store_no" comment:"店铺编号" `
	MerchantNo string     `gorm:"column:merchant_no;type:varchar(80);index;default:;comment:商户" json:"merchant_no" comment:"商户" `
	Owner      string     `gorm:"column:owner;type:varchar(80);index;comment:所属/拥有者" json:"owner" comment:"所属/拥有者" `
	No         string     `gorm:"column:no;type:varchar(80);index;default:;comment:编号" json:"no" comment:"编号" `
	Name       string     `gorm:"column:name;type:varchar(255);comment:名称" json:"name" comment:"名称" `
	//
	TypeSys      string         `gorm:"column:type_sys;type:varchar(80);index;default:'general';comment:类型|普通|系统;" json:"type_sys" comment:"类型;普通;系统;" `
	EventNo      string         `gorm:"column:event_no;type:varchar(80);index;default:;comment:事件编号" json:"event_no" comment:"事件编号" `
	ModelNo      string         `gorm:"column:model_no;type:varchar(80);index;default:;comment:模型编号" json:"model_no" comment:"模型编号" `
	Model        string         `gorm:"column:model;type:varchar(80);comment:模型" json:"model" comment:"模型" `
	Module       string         `gorm:"column:module;type:varchar(80);index;comment:模块" json:"module" comment:"模块" `
	ModuleSub    string         `gorm:"column:module_sub;type:varchar(80);index;comment:子模块" json:"module_sub" comment:"子模块" `
	Field        string         `gorm:"column:field;type:varchar(80);comment:字段名称" json:"field" comment:"字段名称" `
	FieldPath    string         `gorm:"column:field_path;type:varchar(80);comment:路径字段名称" json:"field_path" comment:"路径字段名称" `
	Description  string         `gorm:"column:description;type:varchar(255);comment:描述" json:"description" comment:"描述" `
	Show         int8           `gorm:"column:show;type:int8;not null;index;default:1;comment:1显示2隐藏" json:"show" comment:"1显示2隐藏" `
	Binary       int8           `gorm:"column:binary;type:int8;not null;index;default:2;comment:值二进制1是2否" json:"binary" comment:"值二进制1是2否" `
	ExtraData    datatypes.JSON `gorm:"column:extra_data;type:json;comment:额外数据" json:"extraData" label:"额外数据" `
	Value        string         `gorm:"column:value;type:text;comment:值" json:"value" comment:"值" `
	ValueBinary  string         `gorm:"column:value_binary;type:text;comment:值二进制" json:"value_binary" comment:"值二进制" `
	DefaultValue string         `gorm:"column:default_value;type:varchar(80);comment:默认值" json:"default_value" comment:"默认值" `
	FormCode     string         `gorm:"column:form_code;type:varchar(80);index;comment:表单" json:"form_code" comment:"表单" `
	ValueType    string         `gorm:"column:value_type;type:varchar(80);comment:字段值类型" json:"value_type" comment:"字段值类型" `
	KindUnique   string         `gorm:"column:kind_unique;type:varchar(80);not null;index;default:;comment:字段种类唯一" json:"kind_unique" comment:"字段种类唯一:model_no+field" `
}

func (*BasicConfigEntity) TableName() string {
	return "basic_config"
}

func (*BasicConfigEntity) TableComment() string {
	return "配置"
}
