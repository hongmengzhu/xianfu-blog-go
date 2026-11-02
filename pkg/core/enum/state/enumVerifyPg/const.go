package enumVerifyPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
	"github.com/pangu-2/go-tools/tools/numberPg"
)

// 验证状态
const (
	YES VerifyEnum = 1 //是
	NO  VerifyEnum = 2 //否
)

type VerifyEnum int8

func (this VerifyEnum) String() string {
	switch this {
	case 1:
		return "是"
	case 2:
		return "否"
	default:
		return "未知"
	}
}
func (this VerifyEnum) Index() int8 {
	return int8(this)
}
func (this VerifyEnum) IndexInt8() int8 {
	return int8(this)
}
func (this VerifyEnum) IndexInt64() int64 {
	return int64(this)
}
func (this VerifyEnum) IndexPg() typePg.Int8 {
	return typePg.Int8(this)
}
func (this VerifyEnum) IndexString() string {
	return numberPg.Int8ToString(this.Index())
}

func (this VerifyEnum) IsExistInt64(id int64) bool {
	return int64(this) == id
}

func (this VerifyEnum) IsExistInt8(id int8) bool {
	return int8(this) == id
}

// IsEnableDisable 是否是 有效 停用
func (this VerifyEnum) IsEnableDisable() bool {
	if YES == this {
		return true
	}
	if NO == this {
		return true
	}
	return false
}

var Map = map[int8]enumBasePg.EnumInt8{
	YES.Index(): enumBasePg.EnumInt8{YES.Index(), YES.String()},
	NO.Index():  enumBasePg.EnumInt8{YES.Index(), YES.String()},
}

func GetInt64(id int64) (enumBasePg.EnumInt8, bool) {
	enumInt8, ok := Map[int8(id)]
	return enumInt8, ok
}

func IsExistInt64(id int64) (VerifyEnum, bool) {
	_, ok := Map[int8(id)]
	return VerifyEnum(id), ok
}

func GetType(id VerifyEnum) VerifyEnum {
	return id
}
