package modRamAppAccessKey

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
	"time"
)

type Vo struct {
	ID          typePg.Uint64String `json:"id" label:"id" `
	Name        string              `json:"name" label:"名称" ` // 名称
	No          string              `json:"no" label:"编号代号"`
	Code        string              `json:"code" label:"码值" `
	State       typePg.Int8         `json:"state" label:"状态:1启用;2禁用" ` // 状态:1启用;2禁用
	Description string              `json:"description" label:"描述" `   // 描述
	CreateAt    *time.Time          `json:"createAt" label:"创建时间" `    // 创建时间
	UpdateAt    *time.Time          `json:"updateAt" label:"更新时间" `    // 更新时间
	Type        string              `json:"type" label:"类型" binding:"required,min=1,max=255" validate:"required,min=1,max=255"`
	Key         string              `json:"key" label:"键" `
	Secret      string              `json:"secret" label:"密钥" `
	ExpiryDate  *time.Time          `json:"expiryDate" label:"有效期" `
	AppNo       string              `json:"appNo" label:"应用编号" `
}
