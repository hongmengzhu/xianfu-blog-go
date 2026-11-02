package accessMethodPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// AccessMethod 类型源|采集|手写
type AccessMethod string

const (
	CONSOLE AccessMethod = "console"
	API     AccessMethod = "api"
)

func (this AccessMethod) Name() string {
	switch this {
	case "console":
		return "控制台"
	case "api":
		return "Api"
	default:
		return "未知"
	}
}
func (this AccessMethod) String() string {
	return string(this)
}

func (this AccessMethod) Index() string {
	return string(this)
}

var MapAccessMethod = map[string]enumBasePg.EnumString{
	CONSOLE.String(): enumBasePg.EnumString{CONSOLE.String(), CONSOLE.Name()},
	API.String():     enumBasePg.EnumString{API.String(), API.Name()},
}

func IsExistAccessMethod(id string) bool {
	_, ok := MapAccessMethod[id]
	return ok
}
