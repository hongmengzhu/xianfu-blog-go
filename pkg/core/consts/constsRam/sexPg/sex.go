package sexPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// Sex 性别
type Sex string

const (
	MALE    Sex = "male"   //male
	FEMALE  Sex = "female" //female
	UNKNOWN Sex = "unknown"
)

func (this Sex) Name() string {
	switch this {
	case "male":
		return "男"
	case "female":
		return "女"
	case "unknown":
		return "未知"
	default:
		return ""
	}
}
func (this Sex) String() string {
	return string(this)
}

func (this Sex) Index() string {
	return string(this)
}

var Map = map[string]enumBasePg.EnumString{
	MALE.String():    enumBasePg.EnumString{MALE.String(), MALE.Name()},
	FEMALE.String():  enumBasePg.EnumString{FEMALE.String(), FEMALE.Name()},
	UNKNOWN.String(): enumBasePg.EnumString{UNKNOWN.String(), UNKNOWN.Name()},
}

func IsExist(id string) bool {
	_, ok := Map[id]
	return ok
}
