package model

// BaseIdMarksCt 基础 详情
type BaseIdsMarksCt[ID, MARK any] struct {
	Ids   []ID   `json:"ids"`
	Marks []MARK `json:"marks"`
}
