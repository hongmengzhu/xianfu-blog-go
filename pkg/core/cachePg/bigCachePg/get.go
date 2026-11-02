package bigCachePg

import (
	"context"
	"encoding/json"

	"github.com/go-viper/mapstructure/v2"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(Get))
}

// Get
// @Description: 获取缓存
//
// 原来这里直接注入 *redis.Client 并调用 RedisJSON 的 JSONGet。迁到 starter-bigcache 后，
// bigcache 只保存字节，没有 JSON 文档/路径查询，所以这些方法统一退回「取整段字符串，
// 再在 Go 侧反序列化」：调用方传的都是无 path 的 JSONGet，语义完全一致。
type Get struct {
	rdb *BatchString `autowire:"?"`
}

func NewGet(rdb *BatchString) *Get {
	return &Get{rdb: rdb}
}

// get 读取原始字符串，miss 与错误都折叠为 false。
func (c *Get) get(ctx context.Context, key string) (string, bool) {
	if c.rdb == nil {
		log.Warn(ctx, log.TagAppDef, "缓存未启用", log.String("key", key))
		return "", false
	}
	result, ok := c.rdb.Get(ctx, key)
	if !ok {
		log.Warn(ctx, log.TagAppDef, "获取缓存失败: 缓存不存在", log.String("key", key))
		return "", false
	}
	return result, true
}

// GetJson
//
//	@Description:
//	@receiver c
//	@param ctx
//	@param key
//	@return rt
func (c *Get) GetJson(ctx context.Context, key string) (rt rg.Rs[string]) {
	result, ok := c.get(ctx, key)
	if !ok {
		return rt.ErrorMessage("缓存不存在")
	}
	return rt.OkData(result)
}

// GetJsonStruct
//
//	@Description:
//	@receiver c
//	@param ctx
//	@param key
//	@param v
//	@return rt
func (c *Get) GetJsonStruct(ctx context.Context, key string, v any) (rt rg.Rs[string]) {
	result, ok := c.get(ctx, key)
	if !ok {
		return rt.ErrorMessage("缓存不存在")
	}
	log.Debug(ctx, log.TagAppDef, "缓存", func() []log.Field {
		return []log.Field{log.Any("key", key), log.Any("result", result)}
	})
	if err := json.Unmarshal([]byte(result), v); err != nil {
		log.Error(ctx, log.TagAppDef, err, "缓存序列化失败", log.Any("err", err.Error()))
		return rt.ErrorMessage("缓存序列化失败")
	}
	return rt.OkData(result)
}

// GetJsonMapStruct
//
//	@Description:
//	@receiver c
//	@param ctx
//	@param key
//	@param v
//	@return rt
func (c *Get) GetJsonMapStruct(ctx context.Context, key string, v any) (rt rg.Rs[string]) {
	result, ok := c.get(ctx, key)
	if !ok {
		return rt.ErrorMessage("缓存不存在")
	}
	log.Debug(ctx, log.TagAppDef, "缓存", func() []log.Field {
		return []log.Field{log.Any("key", key), log.Any("result", result)}
	})
	var tmp map[string]any
	if err := json.Unmarshal([]byte(result), &tmp); err != nil {
		log.Error(ctx, log.TagAppDef, err, "缓存序列化失败", log.Any("err", err.Error()))
		return rt.ErrorMessage("缓存序列化失败")
	}
	log.Debug(ctx, log.TagAppDef, "缓存", func() []log.Field {
		return []log.Field{log.Any("key", key), log.Any("tmp", tmp)}
	})
	//map 转为 struct
	if err := mapstructure.Decode(tmp, &v); err != nil {
		log.Error(ctx, log.TagAppDef, err, "map 转 struct", log.Any("err", err.Error()))
		return rt.ErrorMessage("缓存序列化失败")
	}
	return rt.OkData(result)
}

// GetString
//
//	@Description:
//	@receiver c
//	@param ctx
//	@param key
//	@return rt
func (c *Get) GetString(ctx context.Context, key string) (rt rg.Rs[string]) {
	result, ok := c.get(ctx, key)
	if !ok {
		return rt.ErrorMessage("缓存不存在或已过期")
	}
	return rt.OkData(result)
}

// GetStringToJson
//
//	@Description:
//	@receiver c
//	@param ctx
//	@param key
//	@param v   为指针对象
//	@return rt
func (c *Get) GetStringToJson(ctx context.Context, key string, v any) (rt rg.Rs[string]) {
	result, ok := c.get(ctx, key)
	if !ok {
		return rt.ErrorMessage("缓存不存在或已过期")
	}
	log.Debug(ctx, log.TagAppDef, "缓存", func() []log.Field {
		return []log.Field{log.Any("key", key), log.Any("result", result)}
	})
	if err := json.Unmarshal([]byte(result), v); err != nil {
		log.Error(ctx, log.TagAppDef, err, "缓存序列化失败", log.Any("err", err.Error()))
		return rt.ErrorMessage("缓存序列化失败")
	}
	log.Debug(ctx, log.TagAppDef, "缓存", func() []log.Field {
		return []log.Field{log.Any("key", key), log.Any("v", v)}
	})
	return rt.OkData(result)
}
