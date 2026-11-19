package modBasicOperationLog

import (
	"time"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

// TargetHistoryCt 对象操作历史查询条件。
type TargetHistoryCt struct {
	model.BaseQueryCt
	TargetType string             `json:"targetType" label:"对象类型" `
	TargetNo   string             `json:"targetNo" label:"对象编号" `
	TargetID   typePg.Int64String `json:"targetId" label:"对象ID" `
	StartAt    *time.Time         `json:"startAt" label:"开始时间" `
	EndAt      *time.Time         `json:"endAt" label:"结束时间" `
}
