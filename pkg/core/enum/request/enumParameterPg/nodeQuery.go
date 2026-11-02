package enumParameterPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// NodeQueryBy 输出类型
type NodeQueryBy string

const (
	NodeQueryByGeneral NodeQueryBy = "general" //普通
	NodeQueryByCode    NodeQueryBy = "code"    //编码
	NodeQueryByNo      NodeQueryBy = "no"      //编码
	NodeQueryById      NodeQueryBy = "id"      //编码
)

// Name 名称
func (this NodeQueryBy) Name() string {
	switch this {
	case "general":
		return "普通"
	case "code":
		return "编码"
	case "no":
		return "编号"
	default:
		return "未知"
	}
}

// 值
func (this NodeQueryBy) String() string {
	return string(this)
}

// 值
func (this NodeQueryBy) Index() string {
	return string(this)
}

// IsEqual 值是否相等
func (this NodeQueryBy) IsEqual(id string) bool {
	return string(this) == id
}

var NodeQueryByMap = map[string]enumBasePg.EnumString{
	NodeQueryByGeneral.String(): enumBasePg.EnumString{NodeQueryByGeneral.String(), NodeQueryByGeneral.Name()},
	NodeQueryByCode.String():    enumBasePg.EnumString{NodeQueryByCode.String(), NodeQueryByCode.Name()},
	NodeQueryByNo.String():      enumBasePg.EnumString{NodeQueryByNo.String(), NodeQueryByNo.Name()},
	NodeQueryById.String():      enumBasePg.EnumString{NodeQueryById.String(), NodeQueryById.Name()},
}

func IsExistNodeQueryBy(id string) (NodeQueryBy, bool) {
	_, ok := NodeQueryByMap[id]
	return NodeQueryBy(id), ok
}
