package modEventBasicEvent

type EventDto struct {
	TenantNo        string   `label:"租户编号"`
	IsThisTenantAll bool     `label:"当前租户下全部"`
	IsAll           bool     `label:"全部"`
	Nos             []string `label:"指定编号"`
}
