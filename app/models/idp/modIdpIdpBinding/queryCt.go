package modIdpIdpBinding

import (
	"time"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type QueryCt struct {
	model.BaseQueryCt
	ID          typePg.Uint64String `json:"id" label:"" `
	Description string              `json:"description" label:"描述" ` // 描述
	Idp         string              `json:"idp" label:"身份提供商" `
	OpenId      string              `json:"openId" label:"OpenID" `
	UnionId     string              `json:"unionId" label:"UnionID" `
	Platform    string              `json:"platform" label:"平台" `
	Protocol    string              `json:"protocol" label:"协议" `
	State       typePg.Int8         `json:"state" label:"状态:1启用;2禁用" ` // 状态:1启用;2禁用
	CreateBy    string              `json:"createBy" label:"创建人" `     // 创建人
	CreateAt    *time.Time          `json:"createAt" label:"创建时间" `    // 创建时间
}
