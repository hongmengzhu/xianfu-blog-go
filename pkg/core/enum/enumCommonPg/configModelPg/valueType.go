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
	ValueTypeString.String():       enumBasePg.EnumString{Id: ValueTypeString.String(), Name: ValueTypeString.Name()},
	ValueTypeNumber.String():       enumBasePg.EnumString{Id: ValueTypeNumber.String(), Name: ValueTypeNumber.Name()},
	ValueTypeDecimal.String():      enumBasePg.EnumString{Id: ValueTypeDecimal.String(), Name: ValueTypeDecimal.Name()},
	ValueTypeInt.String():          enumBasePg.EnumString{Id: ValueTypeInt.String(), Name: ValueTypeInt.Name()},
	ValueTypeIntString.String():    enumBasePg.EnumString{Id: ValueTypeIntString.String(), Name: ValueTypeIntString.Name()},
	ValueTypeBoolean.String():      enumBasePg.EnumString{Id: ValueTypeBoolean.String(), Name: ValueTypeBoolean.Name()},
	ValueTypeMap.String():          enumBasePg.EnumString{Id: ValueTypeMap.String(), Name: ValueTypeMap.Name()},
	ValueTypeSlice.String():        enumBasePg.EnumString{Id: ValueTypeSlice.String(), Name: ValueTypeSlice.Name()},
	ValueTypeSliceString.String():  enumBasePg.EnumString{Id: ValueTypeSliceString.String(), Name: ValueTypeSliceString.Name()},
	ValueTypeSliceInt.String():     enumBasePg.EnumString{Id: ValueTypeSliceInt.String(), Name: ValueTypeSliceInt.Name()},
	ValueTypeSliceNumber.String():  enumBasePg.EnumString{Id: ValueTypeSliceNumber.String(), Name: ValueTypeSliceNumber.Name()},
	ValueTypeSliceDecimal.String(): enumBasePg.EnumString{Id: ValueTypeSliceDecimal.String(), Name: ValueTypeSliceDecimal.Name()},
	ValueTypeSliceBoolean.String(): enumBasePg.EnumString{Id: ValueTypeSliceBoolean.String(), Name: ValueTypeSliceBoolean.Name()},
	ValueTypeSliceMap.String():     enumBasePg.EnumString{Id: ValueTypeSliceMap.String(), Name: ValueTypeSliceMap.Name()},
	ValueTypeSliceSlice.String():   enumBasePg.EnumString{Id: ValueTypeSliceSlice.String(), Name: ValueTypeSliceSlice.Name()},
	ValueTypeTime.String():         enumBasePg.EnumString{Id: ValueTypeTime.String(), Name: ValueTypeTime.Name()},
	ValueTypeDate.String():         enumBasePg.EnumString{Id: ValueTypeDate.String(), Name: ValueTypeDate.Name()},
	ValueTypeDateTime.String():     enumBasePg.EnumString{Id: ValueTypeDateTime.String(), Name: ValueTypeDateTime.Name()},
}

func IsExistValueType(id string) (ValueType, bool) {
	_, ok := ValueTypeMaps[id]
	return ValueType(id), ok
}
