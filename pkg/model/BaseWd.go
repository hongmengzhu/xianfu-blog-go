package model

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/components/idpPg/auth/interfaces"
)

// BaseWd 基础 详情
type BaseWd struct {
	Wd     string               `json:"wd"`
	Holder interfaces.IHolderPg `json:"holder"` // 会话信息
}
