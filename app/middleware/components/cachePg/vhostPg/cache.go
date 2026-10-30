package vhostPg

import (
	cache "github.com/go-pkgz/expirable-cache/v3"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityTc"
)

var dipl = cache.NewCache[string, entityTc.TcTenantDomainEntity]()

func Get(key string) (entityTc.TcTenantDomainEntity, bool) {
	return dipl.Get(key)
}

func Set(key string, value entityTc.TcTenantDomainEntity) {
	dipl.Set(key, value, 0)
}

func Remove(key string) {
	dipl.Remove(key)
}
