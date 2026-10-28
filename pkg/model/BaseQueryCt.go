package model

import "encoding/json"

type BaseQueryCt struct {
	PageNum  int64    `json:"pageNum"  form:"pageNum"`
	PageSize int64    `json:"pageSize"  form:"pageSize"`
	Wd       string   `json:"wd"  form:"wd"`
	OrderBy  []string `json:"orderBy"  form:"orderBy" label:"排序"`
}

func (c BaseQueryCt) ToJsonString() (string, error) {
	marshal, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(marshal), nil
}
