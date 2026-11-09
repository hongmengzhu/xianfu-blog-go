package modRamIdpCredential

import (
	"time"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type Vo struct {
	ID          typePg.Uint64String `json:"id" label:"id" `
	Description string              `json:"description" label:"描述" `
	State       typePg.Int8         `json:"state" label:"状态:1启用;2禁用" ` // 状态:1启用;2禁用
	Idp         string              `json:"idp" label:"身份提供商" `
	SourceNo    string              `json:"sourceNo" label:"认证源编号" `
	CredType    string              `json:"credType" label:"凭证类型" `
	ExpireAt    *time.Time          `json:"expireAt" label:"过期时间" `
	CreateAt    *time.Time          `json:"createAt" label:"创建时间" ` // 创建时间
	CreateBy    string              `json:"createBy" label:"创建人" `  // 创建人
}
