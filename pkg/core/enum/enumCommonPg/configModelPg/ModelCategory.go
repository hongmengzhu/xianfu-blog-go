package configModelPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// ModelCategory 模型分类
type ModelCategory string

const (
	ModelCategoryTable  ModelCategory = "table"
	ModelCategoryConfig ModelCategory = "config"
)

// Name 名称
func (this ModelCategory) Name() string {
	switch this {
	case "table":
		return "库表模型"
	case "config":
		return "配置模型"
	default:
		return "未知"
	}
}

// 值
func (this ModelCategory) String() string {
	return string(this)
}

// Index 值
func (this ModelCategory) Index() string {
	return string(this)
}

// IsEqual 值是否相等
func (this ModelCategory) IsEqual(id string) bool {
	return string(this) == id
}

var ModelCategoryMap = map[string]enumBasePg.EnumString{
	ModelCategoryTable.String():  enumBasePg.EnumString{ModelCategoryTable.String(), ModelCategoryTable.Name()},
	ModelCategoryConfig.String(): enumBasePg.EnumString{ModelCategoryConfig.String(), ModelCategoryConfig.Name()},
}

func IsExistModelCategory(id string) (ModelCategory, bool) {
	_, ok := ModelCategoryMap[id]
	return ModelCategory(id), ok
}
