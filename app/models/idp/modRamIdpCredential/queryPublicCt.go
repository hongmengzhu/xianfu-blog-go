package modRamIdpCredential

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type QueryPublicCt struct {
	model.BaseQueryNodeCt
	ID       typePg.Uint64String `json:"id" label:"" `
	Idp      string              `json:"idp" label:"身份提供商" `
	SourceNo string              `json:"sourceNo" label:"认证源编号" `
	CredType string              `json:"credType" label:"凭证类型" `
	State    typePg.Int8         `json:"state" label:"状态:1启用;2禁用" ` // 状态:1启用;2禁用
	CreateBy string              `json:"createBy" label:"创建人" `     // 创建人
}
