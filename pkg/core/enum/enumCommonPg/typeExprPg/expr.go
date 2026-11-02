package typeExprPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// TypeExpr 表达式类型;普通;正则
type TypeExpr string

const (
	TypeExprGeneral TypeExpr = "general" //普通
	TypeExprRegular TypeExpr = "regular" //正则
)

// Name 名称
func (this TypeExpr) Name() string {
	switch this {
	case "general":
		return "普通"
	case "regular":
		return "正则"
	default:
		return "未知"
	}
}

// 值
func (this TypeExpr) String() string {
	return string(this)
}

// 值
func (this TypeExpr) Index() string {
	return string(this)
}

// IsEqual 值是否相等
func (this TypeExpr) IsEqual(id string) bool {
	return string(this) == id
}

var TypeExprMap = map[string]enumBasePg.EnumString{
	TypeExprGeneral.String(): enumBasePg.EnumString{TypeExprGeneral.String(), TypeExprGeneral.Name()},
	TypeExprRegular.String(): enumBasePg.EnumString{TypeExprRegular.String(), TypeExprRegular.Name()},
}

func IsExistTypeExpr(id string) (TypeExpr, bool) {
	_, ok := TypeExprMap[id]
	return TypeExpr(id), ok
}
