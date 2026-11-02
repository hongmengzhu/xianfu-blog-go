package configModelPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// FormAttribute 表单FormAttribute 属性
type FormAttribute string

const (
	FormAttributeAutocomplete FormAttribute = "autocomplete" //表单自动填充特性提示   off/on
	FormAttributeDisabled     FormAttribute = "disabled"     //禁用
	FormAttributeMax          FormAttribute = "max"          //最大
	FormAttributeMaxlength    FormAttribute = "maxlength"    //最大字符数（以 UTF-16 码点为单位）
	FormAttributeMin          FormAttribute = "max"          //最小
	FormAttributeMinlength    FormAttribute = "maxlength"    //最小字符数（以 UTF-16 码点为单位）
	FormAttributeMultiple     FormAttribute = "multiple"     //多个
	FormAttributePattern      FormAttribute = "pattern"      //正则表达式
	FormAttributePlaceholder  FormAttribute = "placeholder"  //简短提示
	FormAttributeReadonly     FormAttribute = "readonly"     //只读
	FormAttributeRequired     FormAttribute = "required"     //指定一个非空值
	FormAttributeSize         FormAttribute = "size"         //内容字体大小
	FormAttributeStep         FormAttribute = "step"         //步进值
)

// Name 名称
func (this FormAttribute) Name() string {
	switch this {
	case "memory":
		return "内存"
	case "redis":
		return "redis"
	default:
		return "未知"
	}
}

// 值
func (this FormAttribute) String() string {
	return string(this)
}

// 值
func (this FormAttribute) Index() string {
	return string(this)
}

// IsEqual 值是否相等
func (this FormAttribute) IsEqual(id string) bool {
	return string(this) == id
}

var FormAttributeMap = map[string]enumBasePg.EnumString{
	FormAttributeAutocomplete.String(): enumBasePg.EnumString{FormAttributeAutocomplete.String(), FormAttributeAutocomplete.Name()},
	FormAttributeDisabled.String():     enumBasePg.EnumString{FormAttributeDisabled.String(), FormAttributeDisabled.Name()},
	FormAttributePlaceholder.String():  enumBasePg.EnumString{FormAttributePlaceholder.String(), FormAttributePlaceholder.Name()},
	FormAttributeReadonly.String():     enumBasePg.EnumString{FormAttributeReadonly.String(), FormAttributeReadonly.Name()},
	FormAttributeRequired.String():     enumBasePg.EnumString{FormAttributeRequired.String(), FormAttributeRequired.Name()},
}

func IsExistFormAttribute(id string) (FormAttribute, bool) {
	_, ok := FormAttributeMap[id]
	return FormAttribute(id), ok
}
