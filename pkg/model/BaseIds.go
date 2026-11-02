package model

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/components/idpPg/auth/interfaces"
)

// BaseIds 基础 详情
type BaseIds[ID any] struct {
	Ids    []ID                 `json:"ids"`
	Holder interfaces.IHolderPg `json:"holder"` // 会话信息
}
