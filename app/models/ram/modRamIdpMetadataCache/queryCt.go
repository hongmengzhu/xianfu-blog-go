package modRamIdpMetadataCache

import (
	"time"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type QueryCt struct {
	model.BaseQueryCt
	ID       typePg.Uint64String `json:"id" label:"" `
	IdpNo    string              `json:"idpNo" label:"身份提供商编号" `
	SourceNo string              `json:"sourceNo" label:"认证源编号" `
	Protocol string              `json:"protocol" label:"协议" `
	EntityID string              `json:"entityId" label:"EntityID" `
	State    typePg.Int8         `json:"state" label:"状态:1启用;2禁用" ` // 状态:1启用;2禁用
	CreateAt *time.Time          `json:"createAt" label:"创建时间" `    // 创建时间
}
