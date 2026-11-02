package model

import "github.com/hongmengzhu/xianfu-blog-go/pkg/sdk/auth/holderPg/holderIfPg"

// BasePutIdState 基础 状态更新
type BasePutIdState[STATE, ID, D any] struct {
	State  STATE                `json:"state"`
	Id     ID                   `json:"id"`
	Holder holderIfPg.IHolderPg `json:"holder"` // 会话信息
	Data   D                    `json:"data"`   //数据
	Msg    string               `json:"msg"`    //消息
}
