package model

// BaseIdMarkCt 基础 详情
type BaseIdsMarkCt[ID, MARK any] struct {
	Ids  []ID `json:"ids"`
	Mark MARK `json:"mark"`
}
