package modBasicOperationLog

import (
	"time"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/model"
)

// QueryCt 操作日志查询条件。
type QueryCt struct {
	model.BaseQueryCt
	TraceID    string     `json:"traceId" label:"链路ID" `
	Method     string     `json:"method" label:"请求方法" `
	Path       string     `json:"path" label:"请求路径" `
	TypeDomain string     `json:"typeDomain" label:"域类型" `
	Module     string     `json:"module" label:"模块" `
	Action     string     `json:"action" label:"动作" `
	OperatorNo string     `json:"operatorNo" label:"操作人编号" `
	IP         string     `json:"ip" label:"客户端IP" `
	StatusCode int        `json:"statusCode" label:"HTTP状态码" `
	StartAt    *time.Time `json:"startAt" label:"开始时间" `
	EndAt      *time.Time `json:"endAt" label:"结束时间" `
}
