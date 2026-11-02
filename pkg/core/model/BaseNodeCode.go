package model

type BaseNodeCode struct {
	Id       string      `json:"id" label:"键"`
	Code     string      `json:"code" label:"键"`
	Key      string      `json:"key" label:"键"`
	Label    string      `json:"label" label:"值"`
	ParentId string      `json:"parentId" label:"上级"`
	Extend   interface{} `json:"extend"`
}
