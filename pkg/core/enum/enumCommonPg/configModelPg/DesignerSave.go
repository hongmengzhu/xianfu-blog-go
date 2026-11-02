package configModelPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// DesignerSave 设计保存条件
type DesignerSave string

const (
	DesignerSaveNoDeletion DesignerSave = "noDeletion" //无删除
)

// Name 名称
func (this DesignerSave) Name() string {
	switch this {
	case "noDeletion":
		return "无删除"
	default:
		return "未知"
	}
}

// 值
func (this DesignerSave) String() string {
	return string(this)
}

// Index 值
func (this DesignerSave) Index() string {
	return string(this)
}

// IsEqual 值是否相等
func (this DesignerSave) IsEqual(id string) bool {
	return string(this) == id
}

var DesignerSaveMap = map[string]enumBasePg.EnumString{
	DesignerSaveNoDeletion.String(): enumBasePg.EnumString{DesignerSaveNoDeletion.String(), DesignerSaveNoDeletion.Name()},
}

func IsExistDesignerSave(id string) (DesignerSave, bool) {
	_, ok := DesignerSaveMap[id]
	return DesignerSave(id), ok
}
