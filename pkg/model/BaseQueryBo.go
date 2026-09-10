package model

import (
	"encoding/json"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/interfaces"
)

type BaseQueryBo struct {
	pageNum  int64                `json:"pageNum"`
	PageSize int64                `json:"pageSize"`
	Holder   interfaces.IHolderPg `json:"holder"` // 会话信息
}

func (c BaseQueryBo) ToJsonString() (string, error) {
	marshal, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(marshal), nil
}
