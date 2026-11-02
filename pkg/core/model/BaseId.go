package model

import "github.com/hongmengzhu/xianfu-blog-go/pkg/sdk/auth/holderPg/holderIfPg"

// BaseId 基础 详情
type BaseId[ID any] struct {
	Id     ID                   `json:"id"`
	Holder holderIfPg.IHolderPg `json:"holder"` // 会话信息
}
