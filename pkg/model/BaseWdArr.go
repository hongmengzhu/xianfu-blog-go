package model

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/interfaces"
)

// BaseWdArr 基础 详情
type BaseWdArr struct {
	Wd     []string             `json:"wd"`
	Holder interfaces.IHolderPg `json:"holder"` // 会话信息
}
