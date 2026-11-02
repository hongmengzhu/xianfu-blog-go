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
	ALL.String():           enumBasePg.EnumString{ALL.String(), ALL.Name()},
	ORG.String():           enumBasePg.EnumString{ORG.String(), ORG.Name()},
	ORG_SPECIFIED.String(): enumBasePg.EnumString{ORG_SPECIFIED.String(), ORG_SPECIFIED.Name()},
	FANS.String():          enumBasePg.EnumString{FANS.String(), FANS.Name()},
}

func IsExistWhere(id string) bool {
	_, ok := MapWhere[id]
	return ok
}
