package bigCachePg

import (
	"github.com/allegro/bigcache/v3"
	StarterBigCache "go-spring.org/starter-bigcache"
)

// newBigCacheConfig projects the starter's Config onto bigcache.Config. All
// three custom drivers start from it; each then applies its own policy or adds
// a field the starter does not bind.
func newBigCacheConfig(c StarterBigCache.Config) bigcache.Config {
	conf := bigcache.DefaultConfig(c.LifeWindow)
	conf.Shards = c.Shards
	conf.CleanWindow = c.CleanWindow
	conf.MaxEntriesInWindow = c.MaxEntriesInWindow
	conf.MaxEntrySize = c.MaxEntrySize
	conf.HardMaxCacheSize = c.HardMaxCacheSize
	conf.StatsEnabled = c.StatsEnabled
	return conf
}
