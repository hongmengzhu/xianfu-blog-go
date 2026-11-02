package enumPaymentPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// PaymentCode 支付公司
type PaymentCode string

const (
	PaymentCodeAlipay   PaymentCode = "alipay"   //支付宝
	PaymentCodeWeChat   PaymentCode = "wechat"   //微信
	PaymentCodeUnionPay PaymentCode = "UPQP"     //银联 unionPay
	PaymentCodePayPal   PaymentCode = "payPal"   //payPal
	PaymentCodeAllinPay PaymentCode = "allinPay" //通联
)

// Name 名称
func (this PaymentCode) Name() string {
	switch this {
	case "alipay":
		return "支付宝"
	case "wechat":
		return "微信"
	case "allinPay":
		return "通联"
	case "UPQP":
		return "银联"
	case "payPal":
		return "payPal"
	default:
		return "未知"
	}
}

// 值
func (this PaymentCode) String() string {
	return string(this)
}

// 值
func (this PaymentCode) Index() string {
	return string(this)
}

// IsEqual 值是否相等
func (this PaymentCode) IsEqual(id string) bool {
	return string(this) == id
}

var PaymentCodeMap = map[string]enumBasePg.EnumString{
	PaymentCodeAlipay.String():   enumBasePg.EnumString{PaymentCodeAlipay.String(), PaymentCodeAlipay.Name()},
	PaymentCodeWeChat.String():   enumBasePg.EnumString{PaymentCodeWeChat.String(), PaymentCodeWeChat.Name()},
	PaymentCodeUnionPay.String(): enumBasePg.EnumString{PaymentCodeUnionPay.String(), PaymentCodeUnionPay.Name()},
	PaymentCodePayPal.String():   enumBasePg.EnumString{PaymentCodePayPal.String(), PaymentCodePayPal.Name()},
	PaymentCodeAllinPay.String(): enumBasePg.EnumString{PaymentCodeAllinPay.String(), PaymentCodeAllinPay.Name()},
}

func IsExistPaymentCodeMap(id string) (PaymentCode, bool) {
	_, ok := PaymentCodeMap[id]
	return PaymentCode(id), ok
}
