package modCacheBasicRules

type RulesCache struct {
	ID         int64  `json:"id" comment:"" `
	State      int8   `json:"state" comment:"1有效2停用" `
	Sort       int64  `json:"sort" comment:"排序" `
	TenantNo   string `json:"tenant_no" comment:"租户编号" `
	OrgNo      string `json:"org_no" comment:"组织编号" `
	StoreNo    string `json:"store_no" comment:"店铺编号" `
	MerchantNo string `json:"merchant_no" comment:"商户" `
	No         string `json:"no" comment:"编号" `
	Name       string `json:"name" comment:"名称" `
	//
	Description  string      `json:"description" comment:"描述" `
	ValueType    string      `json:"value_type" comment:"字段值类型" `
	Show         int8        `json:"show" comment:"1显示2隐藏" `
	ExtraData    interface{} `json:"extraData" label:"额外数据" `
	RuleMode     string      `json:"rule_mode" comment:"验证模式类型" `
	Coding       string      `json:"coding" comment:"代码" `
	Condition    string      `json:"condition" comment:"条件" `
	ErrorMessage string      `json:"error_message" comment:"错误提示" `
	Structure    string      `json:"structure" comment:"验证结构" `
	RuleTarget   []string    `json:"rule_target" comment:"目标" `
	SharedRuleNo string      `json:"shared_rule_no" comment:"共享规则编号" `
	TypeSys      string      `json:"type_sys" comment:"类型;普通;系统;" `
	TypeModel    string      `json:"type_model" comment:"模型类型|模型|事件;" `
	ValueNo      string      `json:"value_no" comment:"值编号/模块编号" `
	FieldNo      string      `json:"field_no" comment:"字段编号" `
}
