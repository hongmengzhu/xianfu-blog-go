package enumOrderDeliveryPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// SelectType 订单发货-选择类型
type SelectType string

const (
	SelectTypeDispatch   SelectType = "dispatch"   //发货
	SelectTypeDeliver    SelectType = "deliver"    //送货
	SelectTypeNotDeliver SelectType = "notDeliver" //无需配送
)

// Name 名称
func (this SelectType) Name() string {
	switch this {
	case "dispatch":
		return "发货"
	case "deliver":
		return "送货"
	case "notDeliver":
		return "无需配送"
	default:
		return "未知"
	}
}

// 值
func (this SelectType) String() string {
	return string(this)
}

// 值
func (this SelectType) Index() string {
	return string(this)
}

// IsEqual 值是否相等
func (this SelectType) IsEqual(id string) bool {
	return string(this) == id
}

var SelectTypeMap = map[string]enumBasePg.EnumString{
	SelectTypeDispatch.String():   enumBasePg.EnumString{SelectTypeDispatch.String(), SelectTypeDispatch.Name()},
	SelectTypeDeliver.String():    enumBasePg.EnumString{SelectTypeDeliver.String(), SelectTypeDeliver.Name()},
	SelectTypeNotDeliver.String(): enumBasePg.EnumString{SelectTypeNotDeliver.String(), SelectTypeNotDeliver.Name()},
}

func IsExistSelectType(id string) (SelectType, bool) {
	_, ok := SelectTypeMap[id]
	return SelectType(id), ok
}
