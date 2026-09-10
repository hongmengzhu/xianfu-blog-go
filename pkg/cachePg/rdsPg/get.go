package rdsPg

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/go-viper/mapstructure/v2"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constBlogPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"github.com/redis/go-redis/v9"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(Get))
}

// Get
// @Description: 获取缓存
type Get struct {
	rdb *redis.Client `autowire:""`
}

func NewGet(rdb *redis.Client) *Get {
	return &Get{rdb: rdb}
}

// GetJson
//
//	@Description:
//	@receiver c
//	@param ctx
//	@param key
//	@param v
//	@return rt
func (c *Get) GetJson(ctx context.Context, key string) (rt rg.Rs[string]) {
	result, err := c.rdb.JSONGet(ctx, key).Result()
	if err != nil {
		log.Warnf(ctx, log.TagAppDef, "获取缓存失败: %+v", err.Error())
		if errors.Is(err, redis.Nil) {
			return rt.ErrorMessage("缓存不存在")
		}
		return rt.ErrorMessage("获取缓存错误")
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
func (c *Get) GetJsonStruct(ctx context.Context, key string, v interface{}) (rt rg.Rs[string]) {
	result, err := c.rdb.JSONGet(ctx, key).Result()
	if err == redis.Nil {
		return rt.ErrorMessage("缓存不存在")
	}
	if err != nil {
		log.Warnf(ctx, log.TagAppDef, "获取缓存失败: %+v", err.Error())
		if errors.Is(err, redis.Nil) {
			return rt.ErrorMessage("缓存不存在")
		}
		return rt.ErrorMessage("获取缓存错误")
	}
	log.Debugf(ctx, log.TagAppDef, "缓存[key]=%+v,[data]=%+v", key, result)
	err = json.Unmarshal([]byte(result), &v)
	if err != nil {
		log.Errorf(ctx, log.TagAppDef, "缓存序列化失败: %+v", err.Error())
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
	result, err := c.rdb.JSONGet(ctx, key).Result()
	if err == redis.Nil {
		return rt.ErrorMessage("缓存不存在")
	}
	if err != nil {
		log.Warnf(ctx, log.TagAppDef, "获取缓存失败: %+v", err.Error())
		if errors.Is(err, redis.Nil) {
			return rt.ErrorMessage("缓存不存在")
		}
		return rt.ErrorMessage("获取缓存错误")
	}
	log.Debugf(ctx, log.TagAppDef, "缓存[key]=%+v,[data]=%+v", key, result)
	var tmp map[string]interface{}
	err = json.Unmarshal([]byte(result), &tmp)
	if err != nil {
		log.Errorf(ctx, log.TagAppDef, "缓存序列化失败: %+v", err.Error())
		return rt.ErrorMessage("缓存序列化失败")
	}
	log.Debugf(ctx, log.TagAppDef, "缓存[key]=%+v,[map[string]interface{}]=%+v", key, tmp)
	//map 转为 struct
	if err = mapstructure.Decode(tmp, &v); err != nil {
		log.Errorf(ctx, log.TagAppDef, "map 转 struct err=%+v", err)
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
	result, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		log.Warnf(ctx, log.TagAppDef, "获取缓存失败: %+v", err.Error())
		if errors.Is(err, redis.Nil) {
			return rt.ErrorMessage("缓存不存在或已过期")
		}
		return rt.ErrorMessage("获取缓存错误")
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
	result, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		log.Warnf(ctx, log.TagAppDef, "获取缓存失败: %+v", err.Error())
		if errors.Is(err, redis.Nil) {
			return rt.ErrorMessage("缓存不存在或已过期")
		}
		return rt.ErrorMessage("获取缓存错误")
	}
	log.Debugf(ctx, log.TagAppDef, "缓存[key]=%+v,[data]=%+v", key, result)
	err = json.Unmarshal([]byte(result), v)
	if err != nil {
		log.Errorf(ctx, log.TagAppDef, "缓存序列化失败: %+v", err.Error())
		return rt.ErrorMessage("缓存序列化失败")
	}
	log.Debugf(ctx, log.TagAppDef, "缓存[key]=%+v,[data]=%+v", key, v)
	return rt.OkData(result)
}

func (t *Get) GetAllEvalByLua(ctx context.Context, key []string) ([]interface{}, bool) {
	resp, err := t.rdb.Eval(ctx, constBlogPg.ArticleCategoryLua, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, false
		}
		log.Errorf(context.Background(), log.TagAppDef, "获取缓存失败:", err)
		return nil, false
	}
	// 解析返回结果
	result, ok := resp.([]interface{})
	if !ok {
		log.Errorf(context.Background(), log.TagAppDef, "获取缓存失败:返回结果格式错误，预期为数组类型:", err)
		return nil, false
	}
	return result, true
}
