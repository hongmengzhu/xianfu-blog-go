package modApiBlogCollect

type PushAll struct {
	Data []PushCt `json:"data"`
	Rule []string `json:"rule" label:"规则"`
}
