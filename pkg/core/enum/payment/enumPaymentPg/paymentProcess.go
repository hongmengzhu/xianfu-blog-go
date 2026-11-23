package enumPaymentPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// PaymentProcess 处理状态
type PaymentProcess string

const (
	NotProcessed            PaymentProcess = "notProcessed"             //未处理
	Complete                PaymentProcess = "processingComplete"       //处理完成
	Processing              PaymentProcess = "processing"               //处理中
	ProcessingFailed        PaymentProcess = "processingFailed"         //处理失败
	ProcessingFailedGateway PaymentProcess = "processingFailed:gateway" //处理失败:网关
)

// Name 名称
func (this PaymentProcess) Name() string {
	switch this {
	case "notProcessed":
		return "未处理"
	case "processingComplete":
		return "处理完成"
	case "processing":
		return "处理中"
	case "processingFailed":
		return "处理失败"
	default:
		return "未知"
	}
}

// 值
func (this PaymentProcess) String() string {
	return string(this)
}

// 值
func (this PaymentProcess) Index() string {
	return string(this)
}

// IsEqual 值是否相等
func (this PaymentProcess) IsEqual(id string) bool {
	return string(this) == id
}

var PaymentMap = map[string]enumBasePg.EnumString{
	NotProcessed.String():            enumBasePg.EnumString{Id: NotProcessed.String(), Name: NotProcessed.Name()},
	Complete.String():                enumBasePg.EnumString{Id: Complete.String(), Name: Complete.Name()},
	ProcessingFailed.String():        enumBasePg.EnumString{Id: ProcessingFailed.String(), Name: ProcessingFailed.Name()},
	Processing.String():              enumBasePg.EnumString{Id: Processing.String(), Name: Processing.Name()},
	ProcessingFailedGateway.String(): enumBasePg.EnumString{Id: ProcessingFailedGateway.String(), Name: ProcessingFailedGateway.Name()},
}

func IsExistPayment(id string) (PaymentProcess, bool) {
	_, ok := PaymentMap[id]
	return PaymentProcess(id), ok
}
