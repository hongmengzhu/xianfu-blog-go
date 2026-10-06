package cacheBlogTopicCategoryPg

import (
	"context"
	"encoding/json"
	"errors"
	"slices"
	"strings"

	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBlog"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryBlog"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/cachePg/rdsPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constRedisPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/redis/go-redis/v9"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(Cache))
}

// Cache  专题分类缓存
type Cache struct {
	sv  *repositoryBlog.BlogTopicCategoryRepository `autowire:"?"`
	rdb *rdsPg.BatchString                          `autowire:"?"`
}

func (n *Cache) KeyNo(no string) string {
	return "cac_blog_topic_category:" + no
}

// GetMapByNo 多个
func (n *Cache) GetMapByNo(ctx context.Context, list []string) (maps map[string]*entityBlog.BlogTopicCategoryEntity) {
	maps = make(map[string]*entityBlog.BlogTopicCategoryEntity)
	//
	if nil == list || len(list) == 0 {
		return maps
	}
	idsKey := make([]string, 0)
	ids := make([]string, 0)
	str := ""
	for _, id := range list {
		if strPg.IsNotBlank(id) {
			str = strings.TrimSpace(id)
			ids = append(ids, str)
			//
			idsKey = append(idsKey, n.KeyNo(str))
		}
	}
	if len(idsKey) == 0 {
		return maps
	}
	keys := make([]string, 0)
	keys = append(keys, strings.Join(idsKey, ","))
	//
	resp, err := n.rdb.GetRdb().Eval(ctx, constRedisPg.BATCH_MGET, keys).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return maps
		}
		log.Errorf(ctx, log.TagAppDef, err, "获取缓存失败:%+v", err)
		return maps
	}
	// 解析返回结果
	result, ok := resp.(string)
	if !ok {
		log.Errorf(ctx, log.TagAppDef, err, "获取缓存失败:返回结果格式错误，预期为数组类型:%+v", err)
		return maps
	}
	ret := strings.TrimSpace(result)
	if strPg.IsNotBlank(ret) && ret != "[false]" {
		idsNot := make([]string, 0)
		valueList := make([]string, 0)
		err2 := json.Unmarshal([]byte(ret), &valueList)
		if err2 != nil {
			//解析失败，重新生成
			return maps
		}
		if nil != valueList && len(valueList) > 0 {
			idsFind := make([]string, 0)
			for _, item := range valueList {
				if strPg.IsBlank(item) {
					continue
				}
				if "false" == item {
					continue
				}
				var vo entityBlog.BlogTopicCategoryEntity
				err3 := json.Unmarshal([]byte(item), &vo)
				if err3 != nil {
					log.Warnf(ctx, log.TagAppDef, "解析json失败:%+v,err:%+v", item, err3)
					//解析失败，重新生成
					return maps
				}

				if vo.ID > 0 {
					maps[vo.No] = &vo
					idsFind = append(idsFind, vo.No)
				}
			}
			//
			for _, id := range ids {
				if !slices.Contains(idsFind, id) {
					idsNot = append(idsNot, id)
				}
			}
		}
		//
		if len(idsNot) > 0 {
			info, b := n.sv.FindAllByNoIn(ctx, idsNot)
			if b {
				mapTmp := make(map[string]any)
				for _, item := range info {
					str1, err4 := json.Marshal(item)
					if err4 == nil {
						mapTmp[n.KeyNo(item.No)] = string(str1)
						maps[item.No] = item
					}
				}
				//
				n.rdb.SetPipelineTimeDuration(ctx, mapTmp, 30*24*60*60+10)
			}
		}
		return maps
	} else {
		//先读 数据库
		info, b := n.sv.FindAllByNoIn(ctx, ids)
		if b {
			mapTmp := make(map[string]any)
			for _, item := range info {
				str1, err4 := json.Marshal(item)
				if err4 == nil {
					mapTmp[n.KeyNo(item.No)] = string(str1)
					maps[item.No] = item
				}
			}
			//
			n.rdb.SetPipelineTimeDuration(ctx, mapTmp, 30*24*60*60+10)
		}
	}
	//
	return maps
}

// GetNo 单个
func (n *Cache) GetNo(ctx context.Context, no string) (*entityBlog.BlogTopicCategoryEntity, bool) {
	key := n.KeyNo(no)
	result, err := n.rdb.GetRdb().Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		//先读 数据库
		info, b := n.sv.FindByNo(ctx, no)
		if b {
			str1, err4 := json.Marshal(info)
			if err4 == nil {
				n.rdb.GetRdb().Set(ctx, key, str1, 30*24*60*60+10)
			}
			return info, b
		}
	}
	if err != nil {
		return nil, false
	}
	if strPg.IsNotBlank(result) {
		var vo entityBlog.BlogTopicCategoryEntity
		err4 := json.Unmarshal([]byte(result), &vo)
		if err4 == nil {
			return &vo, true
		}
	}
	return nil, false
}

// UpdateNo 更新
func (n *Cache) UpdateNo(ctx context.Context, no string) {
	key := n.KeyNo(no)
	//先读 数据库
	info, b := n.sv.FindByNo(ctx, no)
	if b {
		str1, err4 := json.Marshal(info)
		if err4 == nil {
			n.rdb.GetRdb().Set(ctx, key, str1, 30*24*60*60+10)
		}
	}
}

// DeleteNos 删除
func (n *Cache) DeleteNos(ctx context.Context, nos []string) {
	if nil == nos || len(nos) <= 0 {
		return
	}
	keys := make([]string, 0)
	for _, no := range nos {
		if strPg.IsBlank(no) {
			continue
		}
		key := n.KeyNo(no)
		keys = append(keys, key)
	}
	if len(keys) > 0 {
		n.rdb.GetRdb().Del(ctx, keys...)
	}
}
