package modRamIdpCredential

import (
	"time"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type CreateUpdateCt struct {
	ID          typePg.Uint64String `json:"id" form:"id" label:"id" `
	Description string              `json:"description" label:"描述" `
	Idp         string              `json:"idp" label:"身份提供商" `
	SourceNo    string              `json:"sourceNo" label:"认证源编号" `
	CredType    string              `json:"credType" label:"凭证类型" `
	Value       string              `json:"value" label:"凭证值" `
	ExpireAt    *time.Time          `json:"expireAt" label:"过期时间" `
	State       typePg.Int8         `json:"state" label:"状态" `
}
