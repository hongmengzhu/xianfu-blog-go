package model

type BaseNodeDataTree struct {
	No       string `json:"no" label:"键"`
	ParentNo string `json:"parentNo" label:"上级"`
	Name     string `json:"name" label:"名称"`
	Code     string `json:"code" label:"名称"`
	Extend   any    `json:"extend"`
	Data     any    `json:"Data" label:"下级"`
}
