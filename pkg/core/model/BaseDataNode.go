package model

type BaseDataNode struct {
	Key      string         `json:"key"`
	Label    string         `json:"label"`
	Extend   interface{}    `json:"extend"`
	Children []BaseDataNode `json:"children"`
}
