package model

// BaseStateIdsCt 基础 详情
type BaseStateIdsCt[ID any] struct {
	Ids    []ID                   `json:"ids" binding:"required" validate:"required"`
	State  int64                  `json:"state" binding:"required,min=1,max=2" validate:"required,min=1,max=2"`
	Extend map[string]interface{} `json:"extend,omitempty"`
	Field  string                 `json:"field,omitempty" label:"字段"`
}
