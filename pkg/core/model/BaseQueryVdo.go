package model

import "encoding/json"

type BaseQueryVdo struct {
	PageNum  int64 `json:"pageNum"`
	PageSize int64 `json:"pageSize"`
	Deleted  int64 `json:"deleted"`
	Opened   int64 `json:"opened"`
}

func (c BaseQueryVdo) ToJsonString() (string, error) {
	marshal, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(marshal), nil
}
