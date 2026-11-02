package model

type BaseNodeKeyValue struct {
	Value    string      `json:"value" label:"键"`
	Label    string      `json:"label" label:"值"`
	ParentNo string      `json:"parentNo" label:"上级"`
	Extend   interface{} `json:"extend"`
}
