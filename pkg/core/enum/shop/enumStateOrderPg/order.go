package enumStateOrderPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// StateOrder 订单状态
type StateOrder string

const (
	DELETE          StateOrder = "delete"          //删除
	CREATE          StateOrder = "create"          //创建
	PendingApproval StateOrder = "pendingApproval" //待审核
	InProcess       StateOrder = "inProcess"       //处理中
	PendingDispatch StateOrder = "pendingDispatch" //待发货
	Shipped         StateOrder = "shipped"         //已发货
	Received        StateOrder = "received"        //已收货 Signed
	Completed       StateOrder = "completed"       //完成
	Cancel          StateOrder = "cancel"          //交易取消
)

// Name 名称
func (this StateOrder) Name() string {
	switch this {
	case "delete":
		return "删除"
	case "create":
		return "待付款"
	case "pendingApproval":
		return "待审核"
	case "inProcess":
		return "待处理"
	case "pendingDispatch":
		return "待发货"
	case "shipped":
		return "已发货"
	case "received":
		return "已收货"
	case "completed":
		return "已完成"
	case "cancel":
		return "交易取消"
	default:
		return "未知"
	}
}

// 值
func (this StateOrder) String() string {
	return string(this)
}

// 值
func (this StateOrder) Index() string {
	return string(this)
}

// IsEqual 值是否相等
func (this StateOrder) IsEqual(id string) bool {
	return string(this) == id
}

var StateOrderMap = map[string]enumBasePg.EnumString{
	DELETE.String():          enumBasePg.EnumString{DELETE.String(), DELETE.Name()},
	CREATE.String():          enumBasePg.EnumString{CREATE.String(), CREATE.Name()},
	PendingApproval.String(): enumBasePg.EnumString{PendingApproval.String(), PendingApproval.Name()},
	InProcess.String():       enumBasePg.EnumString{InProcess.String(), InProcess.Name()},
	PendingDispatch.String(): enumBasePg.EnumString{PendingDispatch.String(), PendingDispatch.Name()},
	Shipped.String():         enumBasePg.EnumString{Shipped.String(), Shipped.Name()},
	Received.String():        enumBasePg.EnumString{Received.String(), Received.Name()},
	Completed.String():       enumBasePg.EnumString{Completed.String(), Completed.Name()},
	Cancel.String():          enumBasePg.EnumString{Cancel.String(), Cancel.Name()},
}

func IsExistStateOrder(id string) (StateOrder, bool) {
	_, ok := StateOrderMap[id]
	return StateOrder(id), ok
}

// 可以发货的几种状态
var deliveryStatusMap = map[string]enumBasePg.EnumString{
	PendingApproval.String(): enumBasePg.EnumString{PendingApproval.String(), PendingApproval.Name()},
	InProcess.String():       enumBasePg.EnumString{InProcess.String(), InProcess.Name()},
	PendingDispatch.String(): enumBasePg.EnumString{PendingDispatch.String(), PendingDispatch.Name()},
}

// 可以发起售后的几种状态
var startAfterSalesStateMap = map[string]enumBasePg.EnumString{
	PendingDispatch.String(): enumBasePg.EnumString{PendingDispatch.String(), PendingDispatch.Name()},
	Shipped.String():         enumBasePg.EnumString{Shipped.String(), Shipped.Name()},
	Received.String():        enumBasePg.EnumString{Received.String(), Received.Name()},
}

// VerifyDeliveryStatus
//
//	@Description: 可以发货的状态
//	@param tp
//	@return bool
func VerifyDeliveryStatus(tp string) bool {
	if _, ok := deliveryStatusMap[tp]; ok {
		return true
	}
	return false
}

// VerifyStartAfterSalesState
//
//	@Description: 可以发起售后的状态
//	@param tp
//	@return bool
func VerifyStartAfterSalesState(tp string) bool {
	if _, ok := startAfterSalesStateMap[tp]; ok {
		return true
	}
	return false
}
