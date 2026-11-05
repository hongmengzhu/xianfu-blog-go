package modBasicConfigEvent

import "github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"

type ModelFormItem struct {
	Name            string      `json:"name" label:"字段中文名称"`
	Field           string      `json:"field" label:"字段英文名称"`
	Show            typePg.Int8 `json:"show" label:"显示隐藏"`
	Binary          typePg.Int8 `json:"binary" label:"是否二进制"`
	DefaultValue    string      `json:"defaultValue" label:"默认值"`
	ValueType       string      `json:"valueType" label:"字段类型"`
	FormCode        string      `json:"formCode" label:"表单类型"`
	ParameterSource string      `json:"parameterSource" label:"参数源"`
	Rules           []string    `json:"rules" label:"验证模式类型"`
	// 以下为渲染表单所需的模型字段属性
	ParameterContent string `json:"parameterContent" label:"参数内容"`
	FormAttr         string `json:"formAttr" label:"表单属性"`
	FormGroup        string `json:"formGroup" label:"表单分组"`
	Sort             int64  `json:"sort" label:"排序"`
}
