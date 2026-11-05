package modBasicConfigModel

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
	"gorm.io/datatypes"
)

type ItemVo struct {
	Id               typePg.Uint64String `json:"id" form:"id" label:"id" `
	No               string              `json:"no" label:"编号"`
	Name             string              `json:"name" label:"字段中文名称"`
	Field            string              `json:"field" label:"字段英文名称"`
	FieldPath        string              `json:"fieldPath" label:"路径字段名称"`
	Description      string              `json:"description" label:"描述"`
	Show             typePg.Int8         `json:"show" label:"显示隐藏"`
	Independent      typePg.Int8         `json:"independent" label:"独立1是2否"`
	Binary           typePg.Int8         `json:"binary" label:"是否二进制"`
	DefaultValue     string              `json:"defaultValue" label:"默认值"`
	Value            string              `json:"value" label:"值"`
	ValueBinary      string              `json:"valueBinary" label:"值二进制"`
	ValueType        string              `json:"valueType" label:"字段类型"`
	ValueAttr        string              `json:"valueAttr" label:"字段值属性"`
	FormCode         string              `json:"formCode" label:"表单类型"`
	FormAttr         datatypes.JSON      `json:"formAttr" label:"表单属性"`
	FormGroup        string              `json:"formGroup" label:"表单分组"`
	ParameterSource  string              `json:"parameterSource" label:"参数源"`
	ParameterContent string              `json:"parameterContent" label:"参数内容"`
	ParameterConfig  string              `json:"parameterConfig" label:"参数配置"`
	SharedScope      []string            `json:"sharedScope" label:"共享范围"`
	SharedFieldNo    string              `json:"sharedFieldNo" label:"共享字段编号"`
	Sort             typePg.Int64String  `json:"sort" label:"排序"`
	// Rules 字段验证规则的简写集合（ruleMode 列表），用于快速保存。
	Rules []string `json:"rules" label:"验证模式类型"`
	// RuleItems 字段验证规则的完整定义，保存到 basic_model_rules。
	RuleItems []RuleItemVo `json:"ruleItems" label:"验证规则明细"`
}
