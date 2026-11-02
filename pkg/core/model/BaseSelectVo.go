package model

// BaseSelectVo 基础 详情
type BaseSelectVo[ID any] struct {
	Value  ID          `json:"value"`
	Label  string      `json:"label"`
	Name   string      `json:"name"`
	Extend interface{} `json:"extend"`
}
