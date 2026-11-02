package model

type BaseNode struct {
	Id       string      `json:"id" label:"键"`
	Code     string      `json:"code" label:"键"`
	Key      string      `json:"key" label:"键"`
	No       string      `json:"no" label:"编号"`
	Label    string      `json:"label" label:"值"`
	ParentId string      `json:"parentId" label:"上级"`
	ParentNo string      `json:"parentNo" label:"上级"`
	Extend   interface{} `json:"extend"`
}
