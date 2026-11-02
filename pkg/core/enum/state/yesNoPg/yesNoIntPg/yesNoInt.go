package yesNoIntPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
	"github.com/pangu-2/go-tools/tools/numberPg"
)

// 验证状态
const (
	Yes YesNoInt = 1 //是
	No  YesNoInt = 2 //否
)

type YesNoInt int8

func (this YesNoInt) String() string {
	switch this {
	case 1:
		return "是"
	case 2:
		return "否"
	default:
		return "未知"
	}
}
func (this YesNoInt) Index() int8 {
	return int8(this)
}
func (this YesNoInt) IndexInt8() int8 {
	return int8(this)
}
func (this YesNoInt) IndexInt64() int64 {
	return int64(this)
}
func (this YesNoInt) IndexPg() typePg.Int8 {
	return typePg.Int8(this)
}
func (this YesNoInt) IndexString() string {
	return numberPg.Int8ToString(this.Index())
}

func (this YesNoInt) IsExistInt64(id int64) bool {
	return int64(this) == id
}

func (this YesNoInt) IsExistInt8(id int8) bool {
	return int8(this) == id
}

// IsEqual 值是否相等
func (this YesNoInt) IsEqual(id int8) bool {
	return int8(this) == id
}

// IsEnableDisable 是否是 有效 停用
func (this YesNoInt) IsEnableDisable() bool {
	if Yes == this {
		return true
	}
	if No == this {
		return true
	}
	return false
}

var MapYesNoInt = map[int8]enumBasePg.EnumInt8{
	Yes.Index(): enumBasePg.EnumInt8{Yes.Index(), Yes.String()},
	No.Index():  enumBasePg.EnumInt8{Yes.Index(), Yes.String()},
}

func GetYesNoInt(id int64) (enumBasePg.EnumInt8, bool) {
	enumInt8, ok := MapYesNoInt[int8(id)]
	return enumInt8, ok
}

func IsExistInt64(id int64) (YesNoInt, bool) {
	_, ok := MapYesNoInt[int8(id)]
	return YesNoInt(id), ok
}

func IsExistInt8(id int8) (YesNoInt, bool) {
	_, ok := MapYesNoInt[id]
	return YesNoInt(id), ok
}
