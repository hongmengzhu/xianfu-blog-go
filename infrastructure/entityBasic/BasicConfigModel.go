package entityBasic

import (
	"time"

	"gorm.io/datatypes"
)

type BasicConfigModelEntity struct {
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
	TypeSys         string                       `gorm:"column:type_sys;type:varchar(80);index;default:'general';comment:类型|普通|系统;" json:"type_sys" comment:"类型;普通;系统;" `
	ModelCategory   string                       `gorm:"column:model_category;type:varchar(80);index;comment:模型种类|表模型|配置模型" json:"model_category" comment:"模型种类" `
	Model           string                       `gorm:"column:model;type:varchar(80);comment:模型" json:"model" comment:"模型" `
	Module          string                       `gorm:"column:module;type:varchar(80);index;comment:模块" json:"module" comment:"模块" `
	ModuleSub       string                       `gorm:"column:module_sub;type:varchar(80);index;comment:子模块" json:"module_sub" comment:"子模块" `
	Table           string                       `gorm:"column:table;type:varchar(80);index;comment:表名" json:"table" comment:"表名" `
	Description     string                       `gorm:"column:description;type:varchar(255);comment:描述" json:"description" comment:"描述" `
	Show            int8                         `gorm:"column:show;type:int8;not null;index;default:1;comment:1显示2隐藏" json:"show" comment:"1显示2隐藏" `
	ExtraData       datatypes.JSON               `gorm:"column:extra_data;type:json;comment:额外数据" json:"extraData" label:"额外数据" `
	Value           string                       `gorm:"column:value;type:text;comment:值" json:"value" comment:"值" `
	Client          string                       `gorm:"column:client;type:varchar(80);comment:端" json:"client" comment:"端" `
	ClientProgram   string                       `gorm:"column:client_program;type:varchar(80);comment:端内程序" json:"client_program" comment:"端内程序|隔开" `
	ClientSync      string                       `gorm:"column:client_sync;type:text;comment:端同步" json:"client_sync" comment:"端同步" `
	LoadingLocation string                       `gorm:"column:loading_location;type:varchar(80);index;default:;comment:加载位置" json:"loading_location" comment:"加载位置" `
	Cache           string                       `gorm:"column:cache;type:varchar(80);comment:缓存" json:"cache" comment:"缓存" `
	CacheKey        string                       `gorm:"column:cache_key;type:varchar(80);index;default:;comment:缓存key" json:"cache_key" comment:"缓存key" `
	DesignerSave    string                       `gorm:"column:designer_save;type:varchar(80);index;default:noDeletion;comment:设计保存条件" json:"designer_save" comment:"设计保存条件" `
	KindUnique      string                       `gorm:"column:kind_unique;type:varchar(80);not null;index;default:;comment:模型字段种类唯一" json:"kind_unique" comment:"模型字段种类唯一:model_no+field" `
	Category        string                       `gorm:"column:category;type:varchar(80);index;comment:分类" json:"category" comment:"分类" `
	CategoryTab     string                       `gorm:"column:category_tab;type:varchar(80);index;comment:分类选项卡" json:"category_tab" comment:"分类选项卡" `
	Tags            datatypes.JSONType[[]string] `gorm:"column:tags;type:jsonb;index;default:'[]';comment:标签" json:"tags" comment:"标签" `
	From            string                       `gorm:"column:from;type:varchar(80);index;default:;comment:来自" json:"from" comment:"来自" `
	SharedModelNo   string                       `gorm:"column:shared_model_no;type:varchar(80);index;default:;comment:共享模型" json:"shared_model_no" comment:"共享模型" `
}

func (*BasicConfigModelEntity) TableName() string {
	return "basic_config_model"
}

func (*BasicConfigModelEntity) TableComment() string {
	return "模型"
}
