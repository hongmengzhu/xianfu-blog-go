package auditPg

// Match 匹配
type Match struct {
	Prefix        []string   `json:"prefix" label:"前缀匹配"`
	Action        string     `json:"action" label:"动作"`
	ActionName    string     `json:"actionName" label:"动作名称"`
	OperationType int8       `json:"operationType" label:"操作类型"` // 见 entityBasic.OperationLogTypeXxx
	TargetKey     string     `json:"targetKey" label:"目标键"`
	TargetFrom    TargetFrom `json:"targetFrom" label:"目标来源｜对象 ID 来源"`
	RiskLevel     int8       `json:"riskLevel" label:"风险等级"`
}
