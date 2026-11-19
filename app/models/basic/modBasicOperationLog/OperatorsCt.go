package modBasicOperationLog

import "time"

// OperatorsCt 操作人下拉条件。
type OperatorsCt struct {
	Keyword string     `json:"keyword" label:"关键词" `
	StartAt *time.Time `json:"startAt" label:"开始时间" `
	EndAt   *time.Time `json:"endAt" label:"结束时间" `
}

// OperatorVo 操作人选项。
type OperatorVo struct {
	OperatorNo      string `json:"operatorNo" label:"操作人编号" `
	OperatorAccount string `json:"operatorAccount" label:"操作人账号" `
	OperatorName    string `json:"operatorName" label:"操作人姓名" `
	Count           int64  `json:"count" label:"操作次数" `
}
