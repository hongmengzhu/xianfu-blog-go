package bookmarkTypeOwnerPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/enumBasePg"
)

// TypeOwner 类型内容所属
type TypeOwner string

const (
	GENERAL TypeOwner = "general"
	TEAM    TypeOwner = "team"
	ORG     TypeOwner = "org"
	MY      TypeOwner = "my"
)

func (this TypeOwner) Name() string {
	switch this {
	case "general":
		return "普通"
	case "team":
		return "团队"
	case "org":
		return "组织"
	case "my":
		return "我的"
	default:
		return "未知"
	}
}
func (this TypeOwner) String() string {
	return string(this)
}

func (this TypeOwner) Code() string {
	return string(this)
}

var MapTypeOwner = map[string]enumBasePg.EnumString{
	GENERAL.String(): enumBasePg.EnumString{GENERAL.String(), GENERAL.Name()},
	TEAM.String():    enumBasePg.EnumString{TEAM.String(), TEAM.Name()},
	ORG.String():     enumBasePg.EnumString{ORG.String(), ORG.Name()},
	MY.String():      enumBasePg.EnumString{MY.String(), MY.Name()},
}

func IsExistTypeOwner(id string) bool {
	_, ok := MapTypeOwner[id]
	return ok
}
