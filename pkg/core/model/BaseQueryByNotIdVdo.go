package model

import (
	"encoding/json"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/sdk/auth/holderPg/holderIfPg"
)

type BaseQueryByNotIdVdo[ID any] struct {
	IdNot  ID                   `json:"idNot"`
	Holder holderIfPg.IHolderPg `json:"holder"` // 会话信息
}

func (c BaseQueryByNotIdVdo[ID]) ToJsonString() (string, error) {
	marshal, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(marshal), nil
}
