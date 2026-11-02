package cacheDiplPg

import (
	cache "github.com/go-pkgz/expirable-cache/v3"
)

var dipl = cache.NewCache[string, DiplCo]()

func Get(key string) (DiplCo, bool) {
	return dipl.Get(key)
}

func Set(key string, value DiplCo) {
	dipl.Set(key, value, 0)
}

func Remove(key string) {
	dipl.Remove(key)
}
