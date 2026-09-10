package model

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/interfaces"
)

// BaseIdMark 基础 详情
type BaseIdsMark[ID, MARK any] struct {
	Ids    []ID                 `json:"ids"`
	Mark   MARK                 `json:"mark"`
	Holder interfaces.IHolderPg `json:"holder"` // 会话信息
}
