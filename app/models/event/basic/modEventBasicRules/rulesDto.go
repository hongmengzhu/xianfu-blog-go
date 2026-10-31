package modEventBasicRules

type RulesDto struct {
	TenantNo        string   `label:"租户编号"`
	IsThisTenantAll bool     `label:"当前租户下全部"`
	IsAll           bool     `label:"全部"`
	FieldNo         []string `label:"指定字段编号"`
}
