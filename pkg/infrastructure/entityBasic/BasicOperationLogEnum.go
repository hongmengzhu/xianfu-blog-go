package entityBasic

// 操作日志枚举：统一前缀 OperationLog*，避免与其它实体常量冲突。

// 操作类型
const (
	OperationLogTypeCreate       int8 = 1  // 新增
	OperationLogTypeUpdate       int8 = 2  // 修改
	OperationLogTypeCreateUpdate int8 = 3  // 修改
	OperationLogTypeDelete       int8 = 4  // 删除
	OperationLogTypeQuery        int8 = 5  // 查询（默认不记录）
	OperationLogTypeExport       int8 = 6  // 导出
	OperationLogTypeImport       int8 = 7  // 导入
	OperationLogTypeGrant        int8 = 8  // 授权
	OperationLogTypeStateChange  int8 = 9  // 状态变更（启用/停用/恢复）
	OperationLogTypeOther        int8 = 10 // 其他
)

// 操作结果
const (
	OperationLogResultUnknown int8 = 0 // 未知
	OperationLogResultSuccess int8 = 1 // 成功
	OperationLogResultFail    int8 = 2 // 失败
	OperationLogResultPartial int8 = 3 // 部分成功
)

// 风险级别
const (
	OperationLogRiskLow    int8 = 1 // 低
	OperationLogRiskMedium int8 = 2 // 中
	OperationLogRiskHigh   int8 = 3 // 高
)

// 操作人类型
const (
	OperationLogOperatorPlatform  int8 = 1 // 平台运营
	OperationLogOperatorTenant    int8 = 2 // 租户管理员
	OperationLogOperatorNormal    int8 = 3 // 普通用户
	OperationLogOperatorAnonymous int8 = 9 // 系统任务/匿名
)

// 来源
const (
	OperationLogSourceApi    int8 = 1 // 接口
	OperationLogSourceEvent  int8 = 2 // 事件
	OperationLogSourceCron   int8 = 3 // 定时任务
	OperationLogSourceManual int8 = 4 // 手工补录
)

// 状态
const (
	OperationLogStateNormal int8 = 1 // 正常
	OperationLogStateFrozen int8 = 2 // 冻结保留
)

// OperationLogTypeName 操作类型文案
func OperationLogTypeName(v int8) string {
	switch v {
	case OperationLogTypeCreate:
		return "新增"
	case OperationLogTypeUpdate:
		return "修改"
	case OperationLogTypeDelete:
		return "删除"
	case OperationLogTypeQuery:
		return "查询"
	case OperationLogTypeExport:
		return "导出"
	case OperationLogTypeImport:
		return "导入"
	case OperationLogTypeGrant:
		return "授权"
	case OperationLogTypeStateChange:
		return "状态变更"
	case OperationLogTypeOther:
		return "其他"
	default:
		return "未知"
	}
}

// OperationLogResultName 结果文案
func OperationLogResultName(v int8) string {
	switch v {
	case OperationLogResultSuccess:
		return "成功"
	case OperationLogResultFail:
		return "失败"
	case OperationLogResultPartial:
		return "部分成功"
	default:
		return "未知"
	}
}

// OperationLogRiskName 风险级别文案
func OperationLogRiskName(v int8) string {
	switch v {
	case OperationLogRiskMedium:
		return "中"
	case OperationLogRiskHigh:
		return "高"
	default:
		return "低"
	}
}

// OperationLogOperatorName 操作人类型文案
func OperationLogOperatorName(v int8) string {
	switch v {
	case OperationLogOperatorPlatform:
		return "平台运营"
	case OperationLogOperatorTenant:
		return "租户管理员"
	case OperationLogOperatorNormal:
		return "普通用户"
	case OperationLogOperatorAnonymous:
		return "系统任务"
	default:
		return "未知"
	}
}

// OperationLogSourceName 来源文案
func OperationLogSourceName(v int8) string {
	switch v {
	case OperationLogSourceEvent:
		return "事件"
	case OperationLogSourceCron:
		return "定时任务"
	case OperationLogSourceManual:
		return "补录"
	default:
		return "接口"
	}
}

// OperationLogStateName 状态文案
func OperationLogStateName(v int8) string {
	if v == OperationLogStateFrozen {
		return "冻结保留"
	}
	return "正常"
}
