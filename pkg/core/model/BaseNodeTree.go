package model

type BaseNodeTree struct {
	Id        string         `json:"id" label:"键"`
	ParentId  string         `json:"parentId" label:"上级"`
	Name      string         `json:"name" label:"名称"`
	Extend    interface{}    `json:"extend"`
	ChildData []BaseNodeTree `json:"childData" label:"下级"`
}
