package typePubPrivePg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// PubPrive 密钥类型
type PubPrive string

const (
	Login PubPrive = "pg:login"
)

// Name 名称
func (this PubPrive) Name() string {
	switch this {
	case "general":
		return "普通"
	case "system":
		return "系统"
	case "manager":
		return "经理"
	case "assistant":
		return "助理"
	default:
		return "未知"
	}
}

// 值
func (this PubPrive) String() string {
	return string(this)
}

// 值
func (this PubPrive) Index() string {
	return string(this)
}

// IsEqual 值是否相等
func (this PubPrive) IsEqual(id string) bool {
	return string(this) == id
}

var PubPriveMap = map[string]enumBasePg.EnumString{
	Login.String(): enumBasePg.EnumString{Login.String(), Login.Name()},
}

func IsExistPubPrive(id string) (PubPrive, bool) {
	_, ok := PubPriveMap[id]
	return PubPrive(id), ok
}
