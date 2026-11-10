package service

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/basic/modBasicTagsRelation"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/consts/constTags"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/repositoryBasic"
	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/strPg"
)

// loadTagStyles 批量加载标签显示数据（含 attributeMap 颜色对象）
func loadTagStyles(ctx *gin.Context, tagsRelat *repositoryBasic.BasicTagsRelationRepository, codes []string) map[string]modBasicTagsRelation.AllVo {
	result := make(map[string]modBasicTagsRelation.AllVo)
	if tagsRelat == nil || len(codes) < 1 {
		return result
	}
	infos, found := tagsRelat.FindAllByTagNoInAndCategoryRoot(ctx, codes, constTags.ArticleInfo.String())
	if !found {
		return result
	}
	for _, item := range infos {
		var vo modBasicTagsRelation.AllVo
		copier.Copy(&vo, &item)
		vo.AttributeMap = make(map[string]interface{})
		if strPg.IsNotBlank(item.Attribute) {
			if err := json.Unmarshal([]byte(item.Attribute), &vo.AttributeMap); err == nil {
				if colorObj, exists := vo.AttributeMap["color"]; exists {
					color := make(map[string]interface{})
					if value, isString := colorObj.(string); isString && strPg.IsNotBlank(value) {
						_ = json.Unmarshal([]byte(value), &color)
					}
					vo.AttributeMap["color"] = color
				}
			}
		} else {
			vo.AttributeMap = defaultTagAttributeMap()
		}
		result[item.TagNo] = vo
	}
	return result
}

// tagStylesOf 按顺序组装标签显示数据，查不到的标签用默认样式兜底
func tagStylesOf(ctx *gin.Context, tagsRelat *repositoryBasic.BasicTagsRelationRepository, codes []string) []modBasicTagsRelation.AllVo {
	styles := make([]modBasicTagsRelation.AllVo, 0)
	if len(codes) < 1 {
		return styles
	}
	result := loadTagStyles(ctx, tagsRelat, codes)
	for _, code := range codes {
		if obj, ok := result[code]; ok {
			styles = append(styles, obj)
		} else {
			styles = append(styles, defaultTagStyle(code))
		}
	}
	return styles
}

// defaultTagStyle 默认标签显示数据
func defaultTagStyle(code string) modBasicTagsRelation.AllVo {
	return modBasicTagsRelation.AllVo{
		Code:         code,
		Name:         code,
		NameFl:       code,
		TagNo:        code,
		AttributeMap: defaultTagAttributeMap(),
	}
}

// defaultTagAttributeMap 默认标签显示属性
func defaultTagAttributeMap() map[string]interface{} {
	return map[string]interface{}{
		"type":     "default",
		"bordered": true,
		"color":    map[string]interface{}{},
		"size":     "medium",
		"strong":   false,
		"round":    false,
	}
}
