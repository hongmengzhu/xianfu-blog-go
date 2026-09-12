package cacheAuthPubPrivPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constsRam/sessionKeyTypePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constsRam/typeDomainPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/clientPg"
)

const (
	// 缓存前缀
	CachePrefix = "cac:auth:pubprive:"
)

func KeySystem() string {
	return Key(sessionKeyTypePg.AccessToken, typeDomainPg.System.Code(), typeDomainPg.System, clientPg.Browser)
}

func AccessTokenKeyManage(tenantNo string) string {
	return Key(sessionKeyTypePg.AccessToken, tenantNo, typeDomainPg.Manage, clientPg.Browser)
}

func AccessTokenKey(tenantNo string, typeDomain typeDomainPg.TypeDomain, client clientPg.Client) string {
	return Key(sessionKeyTypePg.AccessToken, tenantNo, typeDomain, client)
}
func RefreshTokenKey(tenantNo string, typeDomain typeDomainPg.TypeDomain, client clientPg.Client) string {
	return Key(sessionKeyTypePg.RefreshToken, tenantNo, typeDomain, client)
}

func Key(keyType sessionKeyTypePg.SessionKeyType, tenantNo string, typeDomain typeDomainPg.TypeDomain, client clientPg.Client) string {
	return CachePrefix + keyType.Code() + ":" + tenantNo + ":" + typeDomain.Code() + ":" + client.Index()
}
