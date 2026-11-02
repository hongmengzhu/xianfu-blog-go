package enumPaymentPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// PayProduct 支付产品
type PayProduct string

const (
	PayProductAlipay   PayProduct = "alipay"   //支付宝
	PayProductWeChat   PayProduct = "wechat"   //微信
	PayProductAllinPay PayProduct = "allinPay" //通联
	PayProductUnionPay PayProduct = "unionPay" //银联
	PayProductPayPal   PayProduct = "payPal"   //payPal
)

// Name 名称
func (this PayProduct) Name() string {
	switch this {
	case "alipay":
		return "支付宝"
	case "wechat":
		return "微信"
	case "allinPay":
		return "通联"
	case "unionPay":
		return "银联"
	case "payPal":
		return "payPal"
	default:
		return "未知"
	}
}

// 值
func (this PayProduct) String() string {
	return string(this)
}

// 值
func (this PayProduct) Index() string {
	return string(this)
}

// IsEqual 值是否相等
func (this PayProduct) IsEqual(id string) bool {
	return string(this) == id
}

var PayProductMap = map[string]enumBasePg.EnumString{
	PayProductAlipay.String():   enumBasePg.EnumString{PayProductAlipay.String(), PayProductAlipay.Name()},
	PayProductWeChat.String():   enumBasePg.EnumString{PayProductWeChat.String(), PayProductWeChat.Name()},
	PayProductAllinPay.String(): enumBasePg.EnumString{PayProductAllinPay.String(), PayProductAllinPay.Name()},
	PayProductUnionPay.String(): enumBasePg.EnumString{PayProductUnionPay.String(), PayProductUnionPay.Name()},
	PayProductPayPal.String():   enumBasePg.EnumString{PayProductPayPal.String(), PayProductPayPal.Name()},
}

func IsExistPayProductMap(id string) (PayProduct, bool) {
	_, ok := PayProductMap[id]
	return PayProduct(id), ok
}
