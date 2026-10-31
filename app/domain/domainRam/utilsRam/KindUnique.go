package utilsRam

import (
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityRam"
)

// AuthorizationKindUniquePasswordByEntity 账户授权 种类
func AuthorizationKindUniquePasswordByEntity(entity entityRam.RamAccountAuthorizationEntity) string {
	return entity.Ano + ":" + entity.Type + ":" + entity.AppNo
}

// AuthorizationKindUniquePassword 账户授权 种类
func AuthorizationKindUniquePassword(ano, tp, appNo string) string {
	return ano + ":" + tp + ":" + appNo
}
