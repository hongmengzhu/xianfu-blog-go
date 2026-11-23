package typeIdentityPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// TypeIdentity 身份类型;普通;经理;副经理
type TypeIdentity string

const (
	General   TypeIdentity = "general"   //普通
	System    TypeIdentity = "system"    //系统
	Manager   TypeIdentity = "manager"   //经理
	Assistant TypeIdentity = "assistant" //助理
)

// Name 名称
func (this TypeIdentity) Name() string {
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
func (this TypeIdentity) String() string {
	return string(this)
}

// 值
func (this TypeIdentity) Index() string {
	return string(this)
}

// IsEqual 值是否相等
func (this TypeIdentity) IsEqual(id string) bool {
	return string(this) == id
}

var TypeIdentityMap = map[string]enumBasePg.EnumString{
	General.String():   enumBasePg.EnumString{Id: General.String(), Name: General.Name()},
	System.String():    enumBasePg.EnumString{Id: System.String(), Name: System.Name()},
	Manager.String():   enumBasePg.EnumString{Id: Manager.String(), Name: Manager.Name()},
	Assistant.String(): enumBasePg.EnumString{Id: Assistant.String(), Name: Manager.Name()},
}

func IsExistTypeIdentity(id string) (TypeIdentity, bool) {
	_, ok := TypeIdentityMap[id]
	return TypeIdentity(id), ok
}
