package yesNoString

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

type String string

// 验证状态
const (
	Yes String = "yes" //是
	No  String = "no"  //否
)

func (this String) Name() string {
	switch this {
	case "yes":
		return "是"
	case "no":
		return "否"
	default:
		return "未知"
	}
}

// 值
func (this String) String() string {
	return string(this)
}

func (this String) Index() string {
	return string(this)
}

// IsEqual 值是否相等
func (this String) IsEqual(id string) bool {
	return string(this) == id
}

// IsEnableDisable 是否是 有效 停用
func (this String) IsEnableDisable() bool {
	if Yes == this {
		return true
	}
	if No == this {
		return true
	}
	return false
}

var StringMap = map[string]enumBasePg.EnumString{
	Yes.Index(): enumBasePg.EnumString{Yes.Index(), Yes.String()},
	No.Index():  enumBasePg.EnumString{Yes.Index(), Yes.String()},
}

func IsExistString(id string) (enumBasePg.EnumString, bool) {
	enumInt8, ok := StringMap[id]
	return enumInt8, ok
}
