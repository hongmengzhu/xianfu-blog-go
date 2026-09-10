package model

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/interfaces"
)

// BaseId 基础 详情
type BaseId[ID any] struct {
	Id     ID                   `json:"id"`
	Holder interfaces.IHolderPg `json:"holder"` // 会话信息
}
