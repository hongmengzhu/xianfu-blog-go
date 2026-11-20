package modBasicOperationLog

import (
	"time"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

// Vo 列表行。
type Vo struct {
	ID              typePg.Int64String `json:"id" label:"ID" `
	No              string             `json:"no" label:"日志编号" `
	TraceID         string             `json:"traceId" label:"链路ID" `
	StatusCode      int                `json:"statusCode" label:"HTTP状态码" `
	Method          string             `json:"method" label:"请求方法" `
	Path            string             `json:"path" label:"请求路径" `
	TypeDomain      string             `json:"typeDomain" label:"域类型" `
	Module          string             `json:"module" label:"模块" `
	Action          string             `json:"action" label:"动作" `
	OperatorNo      string             `json:"operatorNo" label:"操作人编号" `
	OperatorAccount string             `json:"operatorAccount" label:"操作人账号" `
	OperatorName    string             `json:"operatorName" label:"操作人姓名" `
	IP              string             `json:"ip" label:"客户端IP" `
	Param           string             `json:"param" label:"query参数" `
	RequestBody     string             `json:"requestBody" label:"请求体" `
	CreateAt        *time.Time         `json:"createAt" label:"创建时间" `
}
