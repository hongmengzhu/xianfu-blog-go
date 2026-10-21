package cacheBlogArticleCategoryPg

import (
	"context"
	"encoding/json"
	"slices"
	"strings"

	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBlog"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryBlog"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/cachePg/bigCachePg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(Cache))
}

// Cache  文章分类缓存
type Cache struct {
	sv  *repositoryBlog.BlogArticleCategoryRepository `autowire:"?"`
	rdb *bigCachePg.BatchString                       `autowire:"?"`
}

func (n *Cache) KeyNo(no string) string {
	return "cac_blog_article_category:" + no
}

// GetMapByNo 多个
func (n *Cache) GetMapByNo(ctx context.Context, list []string) (maps map[string]*entityBlog.BlogArticleCategoryEntity) {
	maps = make(map[string]*entityBlog.BlogArticleCategoryEntity)
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
	values, ok := n.rdb.GetAllByKeys(ctx, idsKey)
	if !ok {
		return maps
	}
	idsNot := make([]string, 0)
	idsFind := make([]string, 0)
	for _, item := range values {
		valueStr, ok := item.(string)
		if !ok || strPg.IsBlank(valueStr) {
			continue
		}
		var vo entityBlog.BlogArticleCategoryEntity
		if err := json.Unmarshal([]byte(valueStr), &vo); err != nil {
			log.Warn(ctx, log.TagAppDef, "解析json失败", log.Any("item", valueStr), log.Any("err", err.Error()))
			//解析失败，重新生成
			return maps
		}
		if vo.ID > 0 {
			maps[vo.No] = &vo
			idsFind = append(idsFind, vo.No)
		}
	}
	for _, id := range ids {
		if !slices.Contains(idsFind, id) {
			idsNot = append(idsNot, id)
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
	//
	return maps
}

// GetNo 单个
func (n *Cache) GetNo(ctx context.Context, no string) (*entityBlog.BlogArticleCategoryEntity, bool) {
	key := n.KeyNo(no)
	result, ok := n.rdb.Get(ctx, key)
	if !ok {
		//先读 数据库
		info, b := n.sv.FindByNo(ctx, no)
		if b {
			str1, err4 := json.Marshal(info)
			if err4 == nil {
				_ = n.rdb.Set(ctx, key, str1, 30*24*60*60+10)
			}
			return info, b
		}
		return nil, false
	}
	if strPg.IsNotBlank(result) {
		var vo entityBlog.BlogArticleCategoryEntity
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
			_ = n.rdb.Set(ctx, key, str1, 30*24*60*60+10)
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
		_ = n.rdb.Del(ctx, keys...)
	}
}
