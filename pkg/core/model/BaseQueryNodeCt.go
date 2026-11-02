package model

import "encoding/json"

type BaseQueryNodeCt struct {
	PageNum  int64  `json:"pageNum" default:"1"`
	PageSize int64  `json:"pageSize" default:"20"`
	Wd       string `json:"wd"`
	Key      string `json:"key"`
	By       string `json:"by" label:"主键类型"`
	ByData   string `json:"byData" label:"数据源类型"`
}

func (c BaseQueryNodeCt) ToJsonString() (string, error) {
	marshal, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(marshal), nil
}
