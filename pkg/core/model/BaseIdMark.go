package model

import "github.com/hongmengzhu/xianfu-blog-go/pkg/sdk/auth/holderPg/holderIfPg"

// BaseIdMark 基础 详情
type BaseIdMark[ID, MARK any] struct {
	Id     ID                   `json:"id"`
	Mark   MARK                 `json:"mark"`
	Holder holderIfPg.IHolderPg `json:"holder"` // 会话信息
}
