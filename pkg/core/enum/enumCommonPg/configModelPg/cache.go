package configModelPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// Cache 缓存
type Cache string

const (
	CacheMemory Cache = "memory" //内存
	CacheRedis  Cache = "redis"  //redis
	CacheL2     Cache = "l2"     //2级缓存
)

// Name 名称
func (this Cache) Name() string {
	switch this {
	case "memory":
		return "内存"
	case "redis":
		return "redis"
	default:
		return "未知"
	}
}

// 值
func (this Cache) String() string {
	return string(this)
}

// Index 值
func (this Cache) Index() string {
	return string(this)
}

// IsEqual 值是否相等
func (this Cache) IsEqual(id string) bool {
	return string(this) == id
}

var CacheMap = map[string]enumBasePg.EnumString{
	CacheMemory.String(): enumBasePg.EnumString{CacheMemory.String(), CacheMemory.Name()},
	CacheRedis.String():  enumBasePg.EnumString{CacheRedis.String(), CacheRedis.Name()},
}

func IsExistCache(id string) (Cache, bool) {
	_, ok := CacheMap[id]
	return Cache(id), ok
}
