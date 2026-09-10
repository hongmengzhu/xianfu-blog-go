package model

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/interfaces"
)

// BaseIdMark 基础 详情
type BaseIdMark[ID, MARK any] struct {
	Id     ID                   `json:"id"`
	Mark   MARK                 `json:"mark"`
	Holder interfaces.IHolderPg `json:"holder"` // 会话信息
}
