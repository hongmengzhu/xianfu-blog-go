package configModelPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// ValueType 值类型
type ValueType string

const (
	ValueTypeString         ValueType = "string"        //字符
	ValueTypeNumber         ValueType = "number"        //数值
	ValueTypeDecimal        ValueType = "decimal"       //数值
	ValueTypeInt            ValueType = "int"           //整数
	ValueTypeIntString      ValueType = "intString"     //整数字符
	ValueTypeBoolean        ValueType = "boolean"       //布尔值
	ValueTypeMap            ValueType = "map"           //map 对象
	ValueTypeSlice          ValueType = "slice"         //数组
	ValueTypeSliceString    ValueType = "slice|string"  //数组 字符
	ValueTypeSliceInt       ValueType = "slice|int"     //数组 整数
	ValueTypeSliceNumber    ValueType = "slice|number"  //数组 数值
	ValueTypeSliceDecimal   ValueType = "slice|decimal" //数组 数值
	ValueTypeSliceBoolean   ValueType = "slice|boolean" //数组 布尔值
	ValueTypeSliceMap       ValueType = "slice|map"     //数组 对象
	ValueTypeSliceSlice     ValueType = "slice|slice"   //数组 数组
	ValueTypeDate           ValueType = "date"          //日期
	ValueTypeDateTime       ValueType = "dateTime"      //日期时间
	ValueTypeTime           ValueType = "time"          //时间
	ValueTypeJson           ValueType = "json"          //json
	ValueTypeExtArrayString ValueType = "array|string"  //数组 字符
	ValueTypeExtArrayInt    ValueType = "array|integer" //数组 字符
	ValueTypeExtArrayNumber ValueType = "array|number"  //数组 字符
	ValueTypeExtInteger     ValueType = "integer"       //整数
)

// Name 名称
func (this ValueType) Name() string {
	switch this {
	case "string":
		return "字符"
	case "number":
		return "数值"
	default:
		return "未知"
	}
}

// 值
func (this ValueType) String() string {
	return string(this)
}

// 值
func (this ValueType) Index() string {
	return string(this)
}

// IsEqual 值是否相等
func (this ValueType) IsEqual(id string) bool {
	return string(this) == id
}

var ValueTypeMaps = map[string]enumBasePg.EnumString{
	ValueTypeString.String():       enumBasePg.EnumString{ValueTypeString.String(), ValueTypeString.Name()},
	ValueTypeNumber.String():       enumBasePg.EnumString{ValueTypeNumber.String(), ValueTypeNumber.Name()},
	ValueTypeDecimal.String():      enumBasePg.EnumString{ValueTypeDecimal.String(), ValueTypeDecimal.Name()},
	ValueTypeInt.String():          enumBasePg.EnumString{ValueTypeInt.String(), ValueTypeInt.Name()},
	ValueTypeIntString.String():    enumBasePg.EnumString{ValueTypeIntString.String(), ValueTypeIntString.Name()},
	ValueTypeBoolean.String():      enumBasePg.EnumString{ValueTypeBoolean.String(), ValueTypeBoolean.Name()},
	ValueTypeMap.String():          enumBasePg.EnumString{ValueTypeMap.String(), ValueTypeMap.Name()},
	ValueTypeSlice.String():        enumBasePg.EnumString{ValueTypeSlice.String(), ValueTypeSlice.Name()},
	ValueTypeSliceString.String():  enumBasePg.EnumString{ValueTypeSliceString.String(), ValueTypeSliceString.Name()},
	ValueTypeSliceInt.String():     enumBasePg.EnumString{ValueTypeSliceInt.String(), ValueTypeSliceInt.Name()},
	ValueTypeSliceNumber.String():  enumBasePg.EnumString{ValueTypeSliceNumber.String(), ValueTypeSliceNumber.Name()},
	ValueTypeSliceDecimal.String(): enumBasePg.EnumString{ValueTypeSliceDecimal.String(), ValueTypeSliceDecimal.Name()},
	ValueTypeSliceBoolean.String(): enumBasePg.EnumString{ValueTypeSliceBoolean.String(), ValueTypeSliceBoolean.Name()},
	ValueTypeSliceMap.String():     enumBasePg.EnumString{ValueTypeSliceMap.String(), ValueTypeSliceMap.Name()},
	ValueTypeSliceSlice.String():   enumBasePg.EnumString{ValueTypeSliceSlice.String(), ValueTypeSliceSlice.Name()},
	ValueTypeTime.String():         enumBasePg.EnumString{ValueTypeTime.String(), ValueTypeTime.Name()},
	ValueTypeDate.String():         enumBasePg.EnumString{ValueTypeDate.String(), ValueTypeDate.Name()},
	ValueTypeDateTime.String():     enumBasePg.EnumString{ValueTypeDateTime.String(), ValueTypeDateTime.Name()},
}

func IsExistValueType(id string) (ValueType, bool) {
	_, ok := ValueTypeMaps[id]
	return ValueType(id), ok
}
