package modBasicOperationLog

import "time"

// ExportVo 导出列；字段顺序即 Excel 列顺序。
type ExportVo struct {
	No                string     `json:"no" label:"日志编号" `
	OperateAt         *time.Time `json:"operateAt" label:"操作时间" `
	TypeDomain        string     `json:"typeDomain" label:"域" `
	ModuleName        string     `json:"moduleName" label:"模块" `
	ActionName        string     `json:"actionName" label:"动作" `
	OperationTypeName string     `json:"operationTypeName" label:"操作类型" `
	RiskLevelName     string     `json:"riskLevelName" label:"风险等级" `
	Description       string     `json:"description" label:"描述" `
	OperatorNo        string     `json:"operatorNo" label:"操作人编号" `
	OperatorAccount   string     `json:"operatorAccount" label:"操作人账号" `
	OperatorName      string     `json:"operatorName" label:"操作人姓名" `
	TargetType        string     `json:"targetType" label:"对象类型" `
	TargetNo          string     `json:"targetNo" label:"对象编号" `
	TargetName        string     `json:"targetName" label:"对象名称" `
	ResultName        string     `json:"resultName" label:"结果" `
	ResultCode        string     `json:"resultCode" label:"结果码" `
	ErrorMessage      string     `json:"errorMessage" label:"失败原因" `
	DurationMs        int64      `json:"durationMs" label:"耗时(毫秒)" `
	IPAddress         string     `json:"ipAddress" label:"IP" `
	RequestID         string     `json:"requestId" label:"请求ID" `
	TraceID           string     `json:"traceId" label:"链路ID" `
	TenantNo          string     `json:"tenantNo" label:"租户编号" `
}
