package modBasicOperationLog

// Option 枚举选项。
type Option struct {
	Value int8   `json:"value" label:"值" `
	Label string `json:"label" label:"名称" `
}

// OptionsVo 前端字典。
type OptionsVo struct {
	OperationTypes []Option    `json:"operationTypes" label:"操作类型" `
	RiskLevels     []Option    `json:"riskLevels" label:"风险等级" `
	Results        []Option    `json:"results" label:"结果" `
	Sources        []Option    `json:"sources" label:"来源" `
	States         []Option    `json:"states" label:"状态" `
	OperatorTypes  []Option    `json:"operatorTypes" label:"操作人类型" `
	Modules        []GroupItem `json:"modules" label:"模块" `
}
