package entityBasic

import (
	"time"

	"gorm.io/datatypes"
)

type BasicConfigModelRulesEntity struct {
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
	SharedScope  datatypes.JSONType[[]string] `gorm:"column:shared_scope;type:jsonb;index;default:'[]';comment:共享范围|全局|租户|商户" json:"shared_scope" comment:"共享范围" `
	ModelNo      string                       `gorm:"column:model_no;type:varchar(80);index;default:;comment:模型编号" json:"model_no" comment:"模型编号" `
	Description  string                       `gorm:"column:description;type:varchar(255);comment:描述" json:"description" comment:"描述" `
	ValueType    string                       `gorm:"column:value_type;type:varchar(80);comment:字段值类型" json:"value_type" comment:"字段值类型" `
	Show         int8                         `gorm:"column:show;type:int8;not null;index;default:1;comment:1显示2隐藏" json:"show" comment:"1显示2隐藏" `
	ExtraData    datatypes.JSON               `gorm:"column:extra_data;type:json;comment:额外数据" json:"extraData" label:"额外数据" `
	RuleMode     string                       `gorm:"column:rule_mode;type:varchar(80);comment:验证模式类型" json:"rule_mode" comment:"验证模式类型" `
	Coding       string                       `gorm:"column:coding;type:text;comment:代码" json:"coding" comment:"代码" `
	Condition    string                       `gorm:"column:condition;type:varchar(80);comment:条件" json:"condition" comment:"条件" `
	ErrorMessage string                       `gorm:"column:error_message;type:varchar(255);comment:错误提示" json:"error_message" comment:"错误提示" `
	Structure    string                       `gorm:"column:structure;type:varchar(80);comment:验证结构" json:"structure" comment:"验证结构" `
	RuleTarget   datatypes.JSONType[[]string] `gorm:"column:rule_target;type:jsonb;index;comment:目标前端后端" json:"rule_target" comment:"目标" `
	SharedRuleNo string                       `gorm:"column:shared_rule_no;type:varchar(80);index;comment:共享规则编号" json:"shared_rule_no" comment:"共享规则编号" `
	TypeSys      string                       `gorm:"column:type_sys;type:varchar(80);index;default:'general';comment:类型|普通|系统;" json:"type_sys" comment:"类型;普通;系统;" `
	TypeCategory string                       `gorm:"column:type_category;type:varchar(80);index;default:'general';comment:模型|事件;" json:"type_category" comment:"模型|事件;" `
}

func (*BasicConfigModelRulesEntity) TableName() string {
	return "basic_config_model_rules"
}

func (*BasicConfigModelRulesEntity) TableComment() string {
	return "模型规则"
}
