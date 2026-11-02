package modelBasePg

type ItemResult struct {
	Msg   string `json:"msg" label:"消息"`
	Row   int64  `json:"row" label:"行数"`
	Col   int64  `json:"col" label:"列数"`
	Field string `json:"field" label:"字段"`
}
