package model

// BaseIdMarkCt 基础 详情
type BaseIdMarkCt[ID, MARK any] struct {
	Id   ID   `json:"id"`
	Mark MARK `json:"mark"`
}
