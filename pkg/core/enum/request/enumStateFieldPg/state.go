package enumStateFieldPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// Field 数据源
type Field string

const (
	Id Field = "id" //
	No Field = "no" //
)

// Name 名称
func (this Field) Name() string {
	switch this {
	case "id":
		return "id"
	case "no":
		return "no"
	default:
		return "未知"
	}
}

// 值
func (this Field) String() string {
	return string(this)
}

// 值
func (this Field) Index() string {
	return string(this)
}

// IsEqual 值是否相等
func (this Field) IsEqual(id string) bool {
	return string(this) == id
}

var FieldMap = map[string]enumBasePg.EnumString{
	Id.String(): enumBasePg.EnumString{Id.String(), Id.Name()},
	No.String(): enumBasePg.EnumString{No.String(), No.Name()},
}

func IsExistField(id string) (Field, bool) {
	_, ok := FieldMap[id]
	return Field(id), ok
}
