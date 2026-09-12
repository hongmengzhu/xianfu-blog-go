package cacheAuthPubPrivPg

import (
	cache "github.com/go-pkgz/expirable-cache/v3"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityRam"
)

var dipl = cache.NewCache[string, entityRam.RamAsaJsonPrivatePublicKey]()

func Get(key string) (entityRam.RamAsaJsonPrivatePublicKey, bool) {
	return dipl.Get(key)
}

func Set(key string, value entityRam.RamAsaJsonPrivatePublicKey) {
	dipl.Set(key, value, 0)
}

func Remove(key string) {
	dipl.Remove(key)
}
