package enumStatePg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
	"github.com/pangu-2/go-tools/tools/numberPg"
)

// 状态 11取消(对应有效)12弃置(对应停用)13批量删除(无状态)
const (
	ENABLE       State = 1  //有效
	DISABLE      State = 2  //停用
	CANCEL       State = 11 //取消(对应有效)
	LAY_ASIDE    State = 12 //弃置(对应停用)
	BATCH_DELETE State = 13 //批量删除(无状态)
)

type State int8

func (this State) String() string {
	switch this {
	case 1:
		return "有效"
	case 2:
		return "停用"
	case 11:
		return "取消"
	case 12:
		return "弃置"
	case 13:
		return "批量删除"
	default:
		return "未知"
	}
}
func (this State) Index() int8 {
	return int8(this)
}
func (this State) IndexInt8() int8 {
	return int8(this)
}
func (this State) IndexInt64() int64 {
	return int64(this)
}
func (this State) IndexPg() typePg.Int8 {
	return typePg.Int8(this)
}
func (this State) IndexString() string {
	return numberPg.Int8ToString(this.Index())
}

func (this State) IsExistInt64(id int64) bool {
	return int64(this) == id
}

func (this State) IsExistInt8(id int8) bool {
	return int8(this) == id
}
func (this State) IsEqualInt8(id int8) bool {
	return int8(this) == id
}

// IsEnableDisable 是否是 有效 停用
func (this State) IsEnableDisable() bool {
	if ENABLE == this {
		return true
	}
	if DISABLE == this {
		return true
	}
	return false
}

// IsCancelLayAside 是否是 取消 弃置
func (this State) IsCancelLayAside() bool {
	if CANCEL == this {
		return true
	}
	if LAY_ASIDE == this {
		return true
	}
	return false
}

// ReverseEnableDisable 反转 有效 停用
func (this State) ReverseEnableDisable() (bool, State) {
	if ENABLE == this {
		return true, CANCEL
	}
	if DISABLE == this {
		return true, LAY_ASIDE
	}
	return false, State(0)
}

// ReverseCancelLayAside 反转 取消 弃置 批量删除
func (this State) ReverseCancelLayAside() (bool, State) {
	if CANCEL == this {
		return true, ENABLE
	}
	if LAY_ASIDE == this {
		return true, DISABLE
	}
	if BATCH_DELETE == this {
		return true, DISABLE
	}
	return false, State(0)
}

var MapEnableDisable = map[int8]enumBasePg.EnumInt8{
	ENABLE.Index():  enumBasePg.EnumInt8{ENABLE.Index(), ENABLE.String()},
	DISABLE.Index(): enumBasePg.EnumInt8{ENABLE.Index(), ENABLE.String()},
}
var Map = map[int8]enumBasePg.EnumInt8{
	ENABLE.Index():  enumBasePg.EnumInt8{ENABLE.Index(), ENABLE.String()},
	DISABLE.Index(): enumBasePg.EnumInt8{ENABLE.Index(), ENABLE.String()},
}

func GetInt64(id int64) (enumBasePg.EnumInt8, bool) {
	enumInt8, ok := Map[int8(id)]
	return enumInt8, ok
}

func IsExistInt64(id int64) (State, bool) {
	_, ok := Map[int8(id)]
	return State(id), ok
}

func GetType(id State) State {
	return id
}

func Get(id int8) (enumBasePg.EnumInt8, bool) {
	enumInt8, ok := Map[int8(id)]
	return enumInt8, ok
}

func IsEnableDisable(id int8) bool {
	_, ok := MapEnableDisable[int8(id)]
	return ok
}
