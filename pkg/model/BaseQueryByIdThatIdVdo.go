package model

import (
	"encoding/json"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/interfaces"
)

type BaseQueryByIdThatIdVdo[ID any] struct {
	Id     ID                   `json:"id"`
	IdThat ID                   `json:"idThat"`
	Holder interfaces.IHolderPg `json:"holder"` // 会话信息
}

func (c BaseQueryByIdThatIdVdo[ID]) ToJsonString() (string, error) {
	marshal, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(marshal), nil
}
