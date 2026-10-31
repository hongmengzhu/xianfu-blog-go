package modCacheBasicEvent

type EventCache struct {
	ID         int64  `json:"id" comment:"" `
	Sort       int64  `json:"sort" comment:"排序" `
	TenantNo   string `json:"tenant_no" comment:"租户编号" `
	OrgNo      string `json:"org_no" comment:"组织编号" `
	StoreNo    string `json:"store_no" comment:"店铺编号" `
	MerchantNo string `json:"merchant_no" comment:"商户" `
	Owner      string `json:"owner" comment:"所属/拥有者" `
	No         string `json:"no" comment:"编号" `
	Name       string `json:"name" comment:"名称" `
	//
	SharedScope   []string    `json:"shared_scope" comment:"共享范围" `
	ModelNo       string      `json:"model_no" comment:"模型编号" `
	Model         string      `json:"model" comment:"模型" `
	Field         string      `json:"field" comment:"字段名称/事件码值" `
	FieldSource   string      `json:"field_source" comment:"字段来源/原始字段名称"`
	KindUnique    string      `json:"kind_unique" comment:"模型字段种类唯一:model_no+field" `
	Module        string      `json:"module" comment:"模块" `
	ModuleSub     string      `json:"module_sub" comment:"子模块" `
	Description   string      `json:"description" comment:"描述" `
	Show          int8        `json:"show" comment:"1显示2隐藏" `
	ExtraData     interface{} `json:"extraData" label:"额外数据" `
	Value         string      `json:"value" comment:"值" `
	Category      string      `json:"category" comment:"分类" `
	CategoryTab   string      `json:"category_tab" comment:"分类选项卡" `
	Tags          []string    `json:"tags" comment:"标签" `
	SharedEventNo string      `json:"shared_event_no" comment:"共享事件编号" `
}
