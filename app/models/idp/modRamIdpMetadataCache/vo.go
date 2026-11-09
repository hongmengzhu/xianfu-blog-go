package modRamIdpMetadataCache

import (
	"time"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type Vo struct {
	ID          typePg.Uint64String `json:"id" label:"id" `
	IdpNo       string              `json:"idpNo" label:"身份提供商编号" `
	SourceNo    string              `json:"sourceNo" label:"认证源编号" `
	Protocol    string              `json:"protocol" label:"协议" `
	EntityID    string              `json:"entityId" label:"EntityID" `
	MetadataUrl string              `json:"metadataUrl" label:"元数据地址" `
	MetaFormat  string              `json:"metaFormat" label:"元数据格式" `
	State       typePg.Int8         `json:"state" label:"状态:1启用;2禁用" ` // 状态:1启用;2禁用
	CacheTtl    int64               `json:"cacheTtl" label:"缓存TTL(秒)" `
	LastFetchAt *time.Time          `json:"lastFetchAt" label:"上次拉取时间" `
	ExpireAt    *time.Time          `json:"expireAt" label:"过期时间" `
	CreateAt    *time.Time          `json:"createAt" label:"创建时间" ` // 创建时间
	UpdateAt    *time.Time          `json:"updateAt" label:"更新时间" ` // 更新时间
}
