package model

import "github.com/hongmengzhu/xianfu-blog-go/pkg/sdk/auth/holderPg/holderIfPg"

// BaseIds 基础 详情
type BaseIds[ID any] struct {
	Ids    []ID                 `json:"ids"`
	Holder holderIfPg.IHolderPg `json:"holder"` // 会话信息
}
