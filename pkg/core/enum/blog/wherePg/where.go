package wherePg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// Where 类型源|采集|手写
type Where string

const (
	ALL           Where = "all"
	ORG           Where = "org"
	ORG_SPECIFIED Where = "orgSpecified"
	FANS          Where = "fans"
	VIP           Where = "vip"
)

func (this Where) Name() string {
	switch this {
	case "all":
		return "全网"
	case "org":
		return "本组织"
	case "orgSpecified":
		return "指定组织"
	case "fans":
		return "粉丝"
	default:
		return "未知"
	}
}
func (this Where) String() string {
	return string(this)
}

func (this Where) Index() string {
	return string(this)
}

var MapWhere = map[string]enumBasePg.EnumString{
	ALL.String():           enumBasePg.EnumString{Id: ALL.String(), Name: ALL.Name()},
	ORG.String():           enumBasePg.EnumString{Id: ORG.String(), Name: ORG.Name()},
	ORG_SPECIFIED.String(): enumBasePg.EnumString{Id: ORG_SPECIFIED.String(), Name: ORG_SPECIFIED.Name()},
	FANS.String():          enumBasePg.EnumString{Id: FANS.String(), Name: FANS.Name()},
}

func IsExistWhere(id string) bool {
	_, ok := MapWhere[id]
	return ok
}
