package identityPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// IdentityType 身份类型
type IdentityType string

const (
	GENERAL           IdentityType = "general"
	MANAGER           IdentityType = "manager"
	ASSISTANT_MANAGER IdentityType = "assistant_manager"
)

func (this IdentityType) Name() string {
	switch this {
	case "general":
		return "普通"
	case "manager":
		return "经理"
	case "assistant_manager":
		return "副经理"
	default:
		return "未知"
	}
}
func (this IdentityType) String() string {
	return string(this)
}

func (this IdentityType) Index() string {
	return string(this)
}

var Map = map[string]enumBasePg.EnumString{
	GENERAL.String():           enumBasePg.EnumString{GENERAL.String(), GENERAL.Name()},
	MANAGER.String():           enumBasePg.EnumString{MANAGER.String(), MANAGER.Name()},
	ASSISTANT_MANAGER.String(): enumBasePg.EnumString{ASSISTANT_MANAGER.String(), ASSISTANT_MANAGER.Name()},
}

func IsExist(id string) bool {
	_, ok := Map[id]
	return ok
}
