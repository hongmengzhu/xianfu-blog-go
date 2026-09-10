package entityBasic

import (
	"time"

	"gorm.io/datatypes"
)

type BasicConfigEventEntity struct {
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
	SharedScope   datatypes.JSONType[[]string] `gorm:"column:shared_scope;type:jsonb;index;default:'[]';comment:共享范围|全局|租户|商户" json:"shared_scope" comment:"共享范围" `
	ModelNo       string                       `gorm:"column:model_no;type:varchar(80);index;default:;comment:模型编号" json:"model_no" comment:"模型编号" `
	Model         string                       `gorm:"column:model;type:varchar(80);comment:模型" json:"model" comment:"模型" `
	Field         string                       `gorm:"column:field;type:varchar(80);comment:字段名称/事件码值" json:"field" comment:"字段名称/事件码值" `
	FieldSource   string                       `gorm:"column:field_source;type:varchar(80);comment:字段来源/原始字段名称" json:"field_source" comment:"字段来源/原始字段名称"`
	KindUnique    string                       `gorm:"column:kind_unique;type:varchar(80);not null;index;default:;comment:模型字段种类唯一" json:"kind_unique" comment:"模型字段种类唯一:model_no+field" `
	Module        string                       `gorm:"column:module;type:varchar(80);index;comment:模块" json:"module" comment:"模块" `
	ModuleSub     string                       `gorm:"column:module_sub;type:varchar(80);index;comment:子模块" json:"module_sub" comment:"子模块" `
	Description   string                       `gorm:"column:description;type:varchar(255);comment:描述" json:"description" comment:"描述" `
	Show          int8                         `gorm:"column:show;type:int8;not null;index;default:1;comment:1显示2隐藏" json:"show" comment:"1显示2隐藏" `
	ExtraData     datatypes.JSON               `gorm:"column:extra_data;type:json;comment:额外数据" json:"extraData" label:"额外数据" `
	Value         string                       `gorm:"column:value;type:text;comment:值" json:"value" comment:"值" `
	Category      string                       `gorm:"column:category;type:varchar(80);index;comment:分类" json:"category" comment:"分类" `
	CategoryTab   string                       `gorm:"column:category_tab;type:varchar(80);index;comment:分类选项卡" json:"category_tab" comment:"分类选项卡" `
	Tags          datatypes.JSONType[[]string] `gorm:"column:tags;type:jsonb;index;default:'[]';comment:标签" json:"tags" comment:"标签" `
	SharedEventNo string                       `gorm:"column:shared_event_no;type:varchar(80);index;default:;comment:共享事件编号" json:"shared_event_no" comment:"共享事件编号" `
}

func (*BasicConfigEventEntity) TableName() string {
	return "basic_config_event"
}

func (*BasicConfigEventEntity) TableComment() string {
	return "事件"
}
