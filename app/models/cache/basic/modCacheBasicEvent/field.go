package modCacheBasicEvent

type FieldCache struct {
	ID         int64  `json:"id" comment:"" `
	State      int8   `json:"state" comment:"1有效2停用" `
	Sort       int64  `json:"sort" comment:"排序" `
	TenantNo   string `json:"tenantNo" comment:"租户编号" `
	OrgNo      string `json:"orgNo" comment:"组织编号" `
	StoreNo    string `json:"storeNo" comment:"店铺编号" `
	MerchantNo string `json:"merchantNo" comment:"商户" `
	Owner      string `json:"owner" comment:"所属/拥有者" `
	No         string `json:"no" comment:"编号" `
	Name       string `json:"name" comment:"名称" `
	//
	TypeSys          string      `json:"typeSys" comment:"类型;普通;系统;" `
	SharedScope      []string    `json:"sharedScope" comment:"共享范围" `
	EventNo          string      `json:"eventNo" comment:"事件编号" `
	ModelNo          string      `json:"modelNo" comment:"模型编号" `
	Model            string      `json:"model" comment:"模型" `
	Module           string      `json:"module" comment:"模块" `
	ModuleSub        string      `json:"moduleSub" comment:"子模块" `
	ParentField      string      `json:"parentField" comment:"上级" `
	Field            string      `json:"field" comment:"字段名称" `
	FieldPath        string      `json:"fieldPath" comment:"路径字段名称" `
	Description      string      `json:"description" comment:"描述" `
	Show             int8        `json:"show" comment:"1显示2隐藏" `
	Independent      int8        `json:"independent" comment:"独立 1 是 2 否" `
	Binary           int8        `json:"binary" comment:"值二进制 1 是 2 否" `
	ExtraData        interface{} `json:"extraData" label:"额外数据" `
	Value            string      `json:"value" comment:"值" `
	ValueBinary      string      `json:"valueBinary" comment:"值二进制" `
	DefaultValue     string      `json:"defaultValue" comment:"默认值" `
	FormCode         string      `json:"formCode" comment:"表单" `
	ValueType        string      `json:"valueType" comment:"字段值类型" `
	ValueAttr        string      `json:"valueAttr" comment:"字段值属性 | 单值 | 对象" `
	FormAttr         interface{} `json:"formAttr" label:"表单属性" `
	Rules            []string    `json:"rules" comment:"表单验证规则" `
	ParameterSource  string      `json:"parameterSource" comment:"参数源" `
	ParameterContent string      `json:"parameterContent" comment:"参数内容" `
	ParameterConfig  string      `json:"parameterConfig" comment:"参数配置" `
	KindUnique       string      `json:"kindUnique" comment:"字段种类唯一:model_no+field" `
	SharedFieldNo    string      `json:"sharedFieldNo" comment:"共享字段编号" `
}
