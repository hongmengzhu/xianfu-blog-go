package auditPg

import "github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityBasic"

// MatchActionMap 匹配动作映射。
// key 精确对应路由最后一段；Prefix 用于同义/衍生动作用名，如 exportExcel -> export。
var MatchActionMap = map[string]Match{
	"createUpdate": {
		Action: "createUpdate", ActionName: "创建更新", OperationType: entityBasic.OperationLogTypeCreateUpdate,
		TargetKey: "id", TargetFrom: TargetFromBody, RiskLevel: entityBasic.OperationLogRiskMedium,
		Prefix: []string{"createUpdateSimple", "createUpdateData", "createUpdateCategory", "createUpdateByCategory"},
	},
	"create": {
		Action: "create", ActionName: "创建", OperationType: entityBasic.OperationLogTypeCreate,
		TargetFrom: TargetFromBody, RiskLevel: entityBasic.OperationLogRiskMedium,
	},
	"update": {
		Action: "update", ActionName: "更新", OperationType: entityBasic.OperationLogTypeUpdate,
		TargetKey: "id", TargetFrom: TargetFromBody, RiskLevel: entityBasic.OperationLogRiskMedium,
	},
	"delete": {
		Action: "delete", ActionName: "删除", OperationType: entityBasic.OperationLogTypeDelete,
		TargetKey: "ids", TargetFrom: TargetFromBody, RiskLevel: entityBasic.OperationLogRiskHigh,
	},
	"enable": {
		Action: "enable", ActionName: "启用", OperationType: entityBasic.OperationLogTypeStateChange,
		TargetKey: "ids", TargetFrom: TargetFromBody, RiskLevel: entityBasic.OperationLogRiskMedium,
	},
	"disable": {
		Action: "disable", ActionName: "禁用", OperationType: entityBasic.OperationLogTypeStateChange,
		TargetKey: "ids", TargetFrom: TargetFromBody, RiskLevel: entityBasic.OperationLogRiskMedium,
	},
	"state": {
		Action: "state", ActionName: "状态变更", OperationType: entityBasic.OperationLogTypeStateChange,
		TargetKey: "ids", TargetFrom: TargetFromBody, RiskLevel: entityBasic.OperationLogRiskMedium,
	},
	"recovery": {
		Action: "recovery", ActionName: "恢复", OperationType: entityBasic.OperationLogTypeStateChange,
		TargetKey: "ids", TargetFrom: TargetFromBody, RiskLevel: entityBasic.OperationLogRiskMedium,
		Prefix: []string{"recover"},
	},
	"physicalDeletion": {
		Action: "physicalDeletion", ActionName: "删除", OperationType: entityBasic.OperationLogTypeDelete,
		TargetKey: "ids", TargetFrom: TargetFromBody, RiskLevel: entityBasic.OperationLogRiskHigh,
		Prefix: []string{"physicalDelete"},
	},
	"export": {
		Action: "export", ActionName: "导出", OperationType: entityBasic.OperationLogTypeExport,
		TargetFrom: TargetFromContext, RiskLevel: entityBasic.OperationLogRiskMedium,
		Prefix: []string{"exportExcel", "download", "downloadExcel"},
	},
	"import": {
		Action: "import", ActionName: "导入", OperationType: entityBasic.OperationLogTypeImport,
		TargetFrom: TargetFromContext, RiskLevel: entityBasic.OperationLogRiskMedium,
		Prefix: []string{"importExcel", "uploadExcel"},
	},
	"reset": {
		Action: "reset", ActionName: "重置", OperationType: entityBasic.OperationLogTypeOther,
		TargetKey: "ids", TargetFrom: TargetFromBody, RiskLevel: entityBasic.OperationLogRiskMedium,
	},
	"audit": {
		Action: "audit", ActionName: "审核", OperationType: entityBasic.OperationLogTypeOther,
		TargetKey: "ids", TargetFrom: TargetFromBody, RiskLevel: entityBasic.OperationLogRiskMedium,
	},
	"querySensitive": {
		Action: "querySensitive", ActionName: "查询敏感数据", OperationType: entityBasic.OperationLogTypeOther,
		TargetKey: "ids", TargetFrom: TargetFromBody, RiskLevel: entityBasic.OperationLogRiskMedium,
	},
	"grant": {
		Action: "grant", ActionName: "授权", OperationType: entityBasic.OperationLogTypeOther,
		TargetKey: "ids", TargetFrom: TargetFromBody, RiskLevel: entityBasic.OperationLogRiskMedium,
	},
	"set": {
		Action: "set", ActionName: "设置", OperationType: entityBasic.OperationLogTypeOther,
		TargetKey: "ids", TargetFrom: TargetFromBody, RiskLevel: entityBasic.OperationLogRiskMedium,
	},
}
