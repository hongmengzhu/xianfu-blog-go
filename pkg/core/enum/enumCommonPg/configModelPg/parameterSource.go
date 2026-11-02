package configModelPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// ParameterSource 参数源
type ParameterSource string

const (
	ParameterSourceDataDictionary ParameterSource = "dataDictionary" //数据字典
	ParameterSourceCustom         ParameterSource = "custom"         //自定义
	ParameterSourceUrl            ParameterSource = "url"            //自定义url
	ParameterSourceCustomCreate   ParameterSource = "customCreate"   //自定义创建
)

// Name 名称
func (this ParameterSource) Name() string {
	switch this {
	case "dataDictionary":
		return "数据字典"
	case "custom":
		return "自定义"
	case "url":
		return "自定义url"
	case "customCreate":
		return "自定义创建"
	default:
		return "未知"
	}
}

// 值
func (this ParameterSource) String() string {
	return string(this)
}

// 值
func (this ParameterSource) Index() string {
	return string(this)
}

// IsEqual 值是否相等
func (this ParameterSource) IsEqual(id string) bool {
	return string(this) == id
}

var ParameterSourceMap = map[string]enumBasePg.EnumString{
	ParameterSourceDataDictionary.String(): enumBasePg.EnumString{ParameterSourceDataDictionary.String(), ParameterSourceDataDictionary.Name()},
	ParameterSourceCustom.String():         enumBasePg.EnumString{ParameterSourceCustom.String(), ParameterSourceCustom.Name()},
	ParameterSourceUrl.String():            enumBasePg.EnumString{ParameterSourceUrl.String(), ParameterSourceUrl.Name()},
	ParameterSourceCustomCreate.String():   enumBasePg.EnumString{ParameterSourceCustomCreate.String(), ParameterSourceCustomCreate.Name()},
}

func IsExistParameterSource(id string) (ParameterSource, bool) {
	_, ok := ParameterSourceMap[id]
	return ParameterSource(id), ok
}
