package modBasicOperationLog

import (
	"time"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

// Vo 列表行。
type Vo struct {
	ID                typePg.Int64String `json:"id" label:"ID" `
	No                string             `json:"no" label:"日志编号" `
	TenantNo          string             `json:"tenantNo" label:"租户编号" `
	OrgNo             string             `json:"orgNo" label:"组织编号" `
	TypeDomain        string             `json:"typeDomain" label:"域" `
	Module            string             `json:"module" label:"模块" `
	ModuleName        string             `json:"moduleName" label:"模块名称" `
	SubModule         string             `json:"subModule" label:"子模块" `
	Action            string             `json:"action" label:"动作" `
	ActionName        string             `json:"actionName" label:"动作名称" `
	OperationType     typePg.Int8        `json:"operationType" label:"操作类型" `
	OperationTypeName string             `json:"operationTypeName" label:"操作类型名称" `
	RiskLevel         typePg.Int8        `json:"riskLevel" label:"风险等级" `
	RiskLevelName     string             `json:"riskLevelName" label:"风险等级名称" `
	Description       string             `json:"description" label:"描述" `
	OperatorNo        string             `json:"operatorNo" label:"操作人编号" `
	OperatorAccount   string             `json:"operatorAccount" label:"操作人账号" `
	OperatorName      string             `json:"operatorName" label:"操作人姓名" `
	RoleName          string             `json:"roleName" label:"角色" `
	TargetType        string             `json:"targetType" label:"对象类型" `
	TargetID          string             `json:"targetId" label:"对象ID" `
	TargetNo          string             `json:"targetNo" label:"对象编号" `
	TargetName        string             `json:"targetName" label:"对象名称" `
	TargetCount       int                `json:"targetCount" label:"对象数量" `
	Result            typePg.Int8        `json:"result" label:"结果" `
	ResultName        string             `json:"resultName" label:"结果名称" `
	ResultCode        string             `json:"resultCode" label:"结果码" `
	ErrorMessage      string             `json:"errorMessage" label:"失败原因" `
	DurationMs        int64              `json:"durationMs" label:"耗时(毫秒)" `
	IPAddress         string             `json:"ipAddress" label:"IP" `
	RequestID         string             `json:"requestId" label:"请求ID" `
	TraceID           string             `json:"traceId" label:"链路ID" `
	Source            typePg.Int8        `json:"source" label:"来源" `
	SourceName        string             `json:"sourceName" label:"来源名称" `
	State             typePg.Int8        `json:"state" label:"状态" `
	OperateAt         *time.Time         `json:"operateAt" label:"操作时间" `
	CreateAt          *time.Time         `json:"createAt" label:"创建时间" `
}

// DetailVo 详情，包含请求/响应与数据快照。
type DetailVo struct {
	Vo
	StoreNo          string      `json:"storeNo" label:"门店编号" `
	AppNo            string      `json:"appNo" label:"应用编号" `
	OperatorType     typePg.Int8 `json:"operatorType" label:"操作人类型" `
	OperatorTypeName string      `json:"operatorTypeName" label:"操作人类型名称" `
	RoleNo           string      `json:"roleNo" label:"角色编号" `
	BeforeData       string      `json:"beforeData" label:"变更前" `
	AfterData        string      `json:"afterData" label:"变更后" `
	DiffData         string      `json:"diffData" label:"差异" `
	RequestMethod    string      `json:"requestMethod" label:"请求方法" `
	RequestPath      string      `json:"requestPath" label:"请求路径" `
	RequestQuery     string      `json:"requestQuery" label:"查询参数" `
	RequestBody      string      `json:"requestBody" label:"请求体" `
	ResponseBody     string      `json:"responseBody" label:"响应体" `
	UserAgent        string      `json:"userAgent" label:"User-Agent" `
	ExtraData        string      `json:"extraData" label:"扩展数据" `
}
