package opLogPg

// MatchActionMap 匹配动作映射。
// key 精确对应路由最后一段；Prefix 用于同义/衍生动作用名，如 exportExcel -> export。
var MatchActionMap = map[string]Match{
	"createUpdate": {
		Action: "createUpdate", ActionName: "创建更新",
		TargetKey: "id", RiskLevel: 2,
		Prefix: []string{"createUpdateSimple", "createUpdateData", "createUpdateCategory", "createUpdateByCategory"},
	},
	"create": {
		Action: "create", ActionName: "创建",
		RiskLevel: 2,
	},
	"update": {
		Action: "update", ActionName: "更新",
		TargetKey: "id", RiskLevel: 2,
	},
	"delete": {
		Action: "delete", ActionName: "删除",
		TargetKey: "ids", RiskLevel: 3,
	},
	"enable": {
		Action: "enable", ActionName: "启用",
		TargetKey: "ids", RiskLevel: 2,
	},
	"disable": {
		Action: "disable", ActionName: "禁用",
		TargetKey: "ids", RiskLevel: 2,
	},
	"state": {
		Action: "state", ActionName: "状态变更",
		TargetKey: "ids", RiskLevel: 2,
	},
	"recovery": {
		Action: "recovery", ActionName: "恢复",
		TargetKey: "ids", RiskLevel: 2,
		Prefix: []string{"recover"},
	},
	"physicalDeletion": {
		Action: "physicalDeletion", ActionName: "删除",
		TargetKey: "ids", RiskLevel: 3,
		Prefix: []string{"physicalDelete"},
	},
	"export": {
		Action: "export", ActionName: "导出",
		RiskLevel: 2,
		Prefix:    []string{"exportExcel", "download", "downloadExcel"},
	},
	"import": {
		Action: "import", ActionName: "导入",
		RiskLevel: 2,
		Prefix:    []string{"importExcel", "uploadExcel"},
	},
	"reset": {
		Action: "reset", ActionName: "重置",
		TargetKey: "ids", RiskLevel: 2,
	},
	"audit": {
		Action: "audit", ActionName: "审核",
		TargetKey: "ids", RiskLevel: 2,
	},
	"querySensitive": {
		Action: "querySensitive", ActionName: "查询敏感数据",
		TargetKey: "ids", RiskLevel: 2,
	},
	"grant": {
		Action: "grant", ActionName: "授权",
		TargetKey: "ids", RiskLevel: 2,
	},
	"set": {
		Action: "set", ActionName: "设置",
		TargetKey: "ids", RiskLevel: 2,
	},
	"query": {
		Action: "query", ActionName: "查询",
		TargetKey: "ids", RiskLevel: 2,
		Prefix: []string{"list"},
	},
}
