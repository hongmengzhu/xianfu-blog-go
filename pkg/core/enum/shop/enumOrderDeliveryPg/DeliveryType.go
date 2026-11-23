package enumOrderDeliveryPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// DeliveryType 订单发货类型
type DeliveryType string

const (
	DeliveryTypeWaybillNumber   DeliveryType = "waybillNumber"   //录入单号
	DeliveryTypeBusinessSender  DeliveryType = "businessSender"  //商家寄件
	DeliveryTypeWaybillPrinting DeliveryType = "waybillPrinting" //电子面单打印
)

// Name 名称
func (this DeliveryType) Name() string {
	switch this {
	case "waybillNumber":
		return "录入单号"
	case "businessSender":
		return "商家寄件"
	case "waybillPrinting":
		return "电子面单打印"
	default:
		return "未知"
	}
}

// 值
func (this DeliveryType) String() string {
	return string(this)
}

// 值
func (this DeliveryType) Index() string {
	return string(this)
}

// IsEqual 值是否相等
func (this DeliveryType) IsEqual(id string) bool {
	return string(this) == id
}

var DeliveryTypeMap = map[string]enumBasePg.EnumString{
	DeliveryTypeWaybillNumber.String():   enumBasePg.EnumString{Id: DeliveryTypeWaybillNumber.String(), Name: DeliveryTypeWaybillNumber.Name()},
	DeliveryTypeBusinessSender.String():  enumBasePg.EnumString{Id: DeliveryTypeBusinessSender.String(), Name: DeliveryTypeBusinessSender.Name()},
	DeliveryTypeWaybillPrinting.String(): enumBasePg.EnumString{Id: DeliveryTypeWaybillPrinting.String(), Name: DeliveryTypeWaybillPrinting.Name()},
}

func IsExistDeliveryType(id string) (DeliveryType, bool) {
	_, ok := DeliveryTypeMap[id]
	return DeliveryType(id), ok
}
