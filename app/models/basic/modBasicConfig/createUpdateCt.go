package modBasicConfig

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

// CreateUpdateCt 配置项创建/更新参数
// ModelNo、Model 二选一即可，优先使用 ModelNo；两者都为空时按当前租户 + model 查找。
type CreateUpdateCt struct {
	ID typePg.Uint64String `json:"id" form:"id" label:"id" `

	ModelNo string `json:"modelNo" label:"模型编号"`
	Model   string `json:"model" label:"模型标识"`

	Field       string      `json:"field" binding:"required" label:"字段标识"`
	Name        string      `json:"name" label:"名称"`
	FieldPath   string      `json:"fieldPath" label:"路径字段"`
	Description string      `json:"description" label:"描述"`
	Show        typePg.Int8 `json:"show" label:"1显示2隐藏"`
	Binary      typePg.Int8 `json:"binary" label:"值二进制1是2否"`
	ValueType   string      `json:"valueType" label:"字段值类型"`
	FormCode    string      `json:"formCode" label:"表单控件"`

	DefaultValue string `json:"defaultValue" label:"默认值"`
	Value        string `json:"value" label:"值"`
	ValueBinary  string `json:"valueBinary" label:"值二进制"`

	CategoryNo string             `json:"categoryNo" label:"配置分类编号"`
	Sort       typePg.Int64String `json:"sort" label:"排序"`
}
