package enumOrderPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// Payment 订单支付状态
type Payment string

const (
	PaymentUnpaid   Payment = "unpaid"   //未支付
	PaymentPaid     Payment = "paid"     //已支付[补]
	PaymentPaidAuto Payment = "paidAuto" //已支付自动(支付公司反馈)
)

// Name 名称
func (this Payment) Name() string {
	switch this {
	case "unpaid":
		return "未支付"
	case "paid":
		return "已支付[补]"
	case "paidAuto":
		return "已支付[自动]"
	default:
		return "未知"
	}
}

// 值
func (this Payment) String() string {
	return string(this)
}

// 值
func (this Payment) Index() string {
	return string(this)
}

// IsEqual 值是否相等
func (this Payment) IsEqual(id string) bool {
	return string(this) == id
}

var PaymentMap = map[string]enumBasePg.EnumString{
	PaymentUnpaid.String():   enumBasePg.EnumString{PaymentUnpaid.String(), PaymentUnpaid.Name()},
	PaymentPaid.String():     enumBasePg.EnumString{PaymentPaid.String(), PaymentPaid.Name()},
	PaymentPaidAuto.String(): enumBasePg.EnumString{PaymentPaidAuto.String(), PaymentPaidAuto.Name()},
}

func IsExistPayment(id string) (Payment, bool) {
	_, ok := PaymentMap[id]
	return Payment(id), ok
}

// PaymentIsPaid
//
//	@Description: 是否已支付
//	@param id
//	@return bool
func PaymentIsPaid(id string) bool {
	switch id {
	case PaymentPaid.String():
	case PaymentPaidAuto.String():
		return true
	}
	return false
}
