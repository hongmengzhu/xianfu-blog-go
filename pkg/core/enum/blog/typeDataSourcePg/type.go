package typeDataSourcePg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// TypeDataSource 类型源|采集|手写
type TypeDataSource string

const (
	PLATFORM TypeDataSource = "platform"
	EXTERNAL TypeDataSource = "external"
)

func (this TypeDataSource) Name() string {
	switch this {
	case "platform":
		return "平台"
	case "external":
		return "外部"
	default:
		return "未知"
	}
}
func (this TypeDataSource) String() string {
	return string(this)
}

func (this TypeDataSource) Index() string {
	return string(this)
}

var MapTypeDataSource = map[string]enumBasePg.EnumString{
	PLATFORM.String(): enumBasePg.EnumString{PLATFORM.String(), PLATFORM.Name()},
	EXTERNAL.String(): enumBasePg.EnumString{EXTERNAL.String(), EXTERNAL.Name()},
}

func IsExistTypeDataSource(id string) bool {
	_, ok := MapTypeDataSource[id]
	return ok
}
