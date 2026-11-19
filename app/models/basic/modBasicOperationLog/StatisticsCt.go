package modBasicOperationLog

import "time"

// StatisticsCt 统计条件。
type StatisticsCt struct {
	TenantNo   string     `json:"tenantNo" label:"租户编号" `
	TypeDomain string     `json:"typeDomain" label:"域" `
	Module     string     `json:"module" label:"模块" `
	Action     string     `json:"action" label:"动作" `
	OperatorNo string     `json:"operatorNo" label:"操作人编号" `
	StartAt    *time.Time `json:"startAt" label:"开始时间" `
	EndAt      *time.Time `json:"endAt" label:"结束时间" `
}

// Point 趋势点。
type Point struct {
	Day     string `json:"day" label:"日期" `
	Total   int64  `json:"total" label:"总数" `
	Success int64  `json:"success" label:"成功" `
	Fail    int64  `json:"fail" label:"失败" `
	Partial int64  `json:"partial" label:"部分成功" `
	Unknown int64  `json:"unknown" label:"未知" `
}

// GroupItem 按维度聚合的行。
type GroupItem struct {
	Key   string `json:"key" label:"编码" `
	Name  string `json:"name" label:"名称" `
	Count int64  `json:"count" label:"数量" `
}

// StatisticsVo 统计结果。
type StatisticsVo struct {
	Total   int64       `json:"total" label:"总数" `
	Success int64       `json:"success" label:"成功" `
	Fail    int64       `json:"fail" label:"失败" `
	Partial int64       `json:"partial" label:"部分成功" `
	Unknown int64       `json:"unknown" label:"未知" `
	Trend   []Point     `json:"trend" label:"趋势" `
	Modules []GroupItem `json:"modules" label:"模块分布" `
	Actions []GroupItem `json:"actions" label:"动作分布" `
	Risks   []GroupItem `json:"risks" label:"风险分布" `
	Results []GroupItem `json:"results" label:"结果分布" `
}
