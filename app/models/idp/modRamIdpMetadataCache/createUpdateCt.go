package modRamIdpMetadataCache

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type CreateUpdateCt struct {
	ID          typePg.Uint64String `json:"id" form:"id" label:"id" `
	IdpNo       string              `json:"idpNo" label:"身份提供商编号" `
	SourceNo    string              `json:"sourceNo" label:"认证源编号" `
	Protocol    string              `json:"protocol" label:"协议" `
	EntityID    string              `json:"entityId" label:"EntityID" `
	MetadataUrl string              `json:"metadataUrl" label:"元数据地址" `
	MetaFormat  string              `json:"metaFormat" label:"元数据格式" `
	CacheTtl    int64               `json:"cacheTtl" label:"缓存TTL(秒)" `
	State       typePg.Int8         `json:"state" label:"状态" `
}
