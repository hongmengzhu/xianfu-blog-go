package enumParameterPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// DataSourceBy 数据源
type DataSourceBy string

const (
	DataSourceByGeneral DataSourceBy = "general" //普通
	DataSourceByCache   DataSourceBy = "cache"   //缓存
)

// Name 名称
func (this DataSourceBy) Name() string {
	switch this {
	case "general":
		return "普通"
	case "cache":
		return "缓存"
	default:
		return "未知"
	}
}

// 值
func (this DataSourceBy) String() string {
	return string(this)
}

// 值
func (this DataSourceBy) Index() string {
	return string(this)
}

// IsEqual 值是否相等
func (this DataSourceBy) IsEqual(id string) bool {
	return string(this) == id
}

var DataSourceByMap = map[string]enumBasePg.EnumString{
	DataSourceByGeneral.String(): enumBasePg.EnumString{DataSourceByGeneral.String(), DataSourceByGeneral.Name()},
	DataSourceByCache.String():   enumBasePg.EnumString{DataSourceByCache.String(), DataSourceByCache.Name()},
}

func IsExistDataSourceBy(id string) (DataSourceBy, bool) {
	_, ok := DataSourceByMap[id]
	return DataSourceBy(id), ok
}
