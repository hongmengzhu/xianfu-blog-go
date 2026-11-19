package modBasicOperationLog

import (
	"time"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

// QueryCt 操作日志查询条件。
type QueryCt struct {
	model.BaseQueryCt
	ID              typePg.Int64String `json:"id" label:"ID" `
	No              string             `json:"no" label:"日志编号" `
	TenantNo        string             `json:"tenantNo" label:"租户编号" `
	TypeDomain      string             `json:"typeDomain" label:"域:manage/system/customer" `
	Module          string             `json:"module" label:"模块" `
	SubModule       string             `json:"subModule" label:"子模块" `
	Action          string             `json:"action" label:"动作" `
	OperationType   typePg.Int8        `json:"operationType" label:"操作类型:1新增;2修改;3删除;4启用;5停用;6恢复;7授权;8导出;9其他" `
	RiskLevel       typePg.Int8        `json:"riskLevel" label:"风险等级:1低;2中;3高" `
	OperatorNo      string             `json:"operatorNo" label:"操作人编号" `
	OperatorAccount string             `json:"operatorAccount" label:"操作人账号" `
	TargetType      string             `json:"targetType" label:"对象类型" `
	TargetNo        string             `json:"targetNo" label:"对象编号" `
	Result          typePg.Int8        `json:"result" label:"结果:0未知;1成功;2失败;3部分成功" `
	Source          typePg.Int8        `json:"source" label:"来源:1接口;2事件;3定时;4手动" `
	StartAt         *time.Time         `json:"startAt" label:"开始时间" `
	EndAt           *time.Time         `json:"endAt" label:"结束时间" `
}
