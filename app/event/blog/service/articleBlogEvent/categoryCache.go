package articleBlogEvent

import (
	"context"
	"strings"

	"github.com/goccy/go-json"
	"github.com/hongmengzhu/xianfu-blog-go/app/event/blog/model/modEventBlogArticleCategory"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/blog/modBlogArticleCategory"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBlog"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/sdk/blog/key/blogKeyPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg/optionsPg"
	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/strPg"
	"go-spring.org/log"
	"gorm.io/gorm"
)

// CategoryCache 分类缓存处理
type CategoryCache struct {
	sp  *Sp `autowire:"?"`
	dto modEventBlogArticleCategory.CacheDto
}

func NewCategoryCache(sp *Sp, dto modEventBlogArticleCategory.CacheDto) *CategoryCache {
	return &CategoryCache{
		sp:  sp,
		dto: dto,
	}
}

// Processor 处理
//
//	@Description:
//	@receiver c
//	@param ctx
//	@return error
func (c *CategoryCache) Processor(ctx context.Context) error {
	if c.dto.IsAll {
		return c.all(ctx)
	}
	if c.dto.IsThisTenantAll {
		return c.thisAll(ctx)
	}
	if strPg.IsNotBlank(c.dto.TenantNo) && !c.dto.IsThisTenantAll && len(c.dto.Nos) > 0 {
		return c.custom(ctx)
	}
	return nil
}

func (c *CategoryCache) thisAll(ctx context.Context) error {
	var query entityBlog.BlogArticleCategoryEntity
	if strPg.IsNotBlank(c.dto.TenantNo) {
		query.TenantNo = c.dto.TenantNo
		query.State = enumStatePg.ENABLE.Index()
		infos := c.sp.catRep.FindAll(ctx, query)
		if infos != nil {
			mapTmp := make(map[string]modBlogArticleCategory.Cache)
			data := make(map[string]interface{})
			for _, info := range infos {
				//
				var obj modBlogArticleCategory.Cache
				copier.Copy(&obj, info)
				//
				mapTmp[info.No] = obj
			}
			//
			keysAdd := make([]string, 0)
			code_link := make([]string, 0)
			for _, info := range infos {
				code_link = make([]string, 0)
				if strPg.IsNotBlank(info.NoLink) {
					split := strings.Split(info.NoLink, "|")
					for _, item := range split {
						if strPg.IsNotBlank(item) {
							obj := mapTmp[item]
							code_link = append(code_link, obj.Code)
						}
					}
				}
				//
				if obj, ok := mapTmp[info.No]; ok {
					obj.CodeLink = strings.Join(code_link, "/")
					//
					key := blogKeyPg.ArticleCategoryTenantNo(info.TenantNo, info.No)
					str, err := json.Marshal(obj)
					if err == nil {
						data[key] = str
					}
					keysAdd = append(keysAdd, key)
					//

					// 缓存  code = no
					key2 := blogKeyPg.ArticleCategoryTenantNoAndNoByCode(info.TenantNo, info.Code)
					data[key2] = info.No
				}
			}
			//存入所有集合
			if len(keysAdd) > 0 {
				keysAll := blogKeyPg.ArticleCategoryTenantNoKeys(c.dto.TenantNo)
				err := c.sp.rdt.GetRdb().SAdd(ctx, keysAll, keysAdd).Err()
				if err != nil {
					log.Errorf(ctx, log.TagAppDef, err, "缓存失败:", err)
				}
			}
			//
			if len(data) > 0 {
				c.sp.rdt.SetPipeline(ctx, data)
			}
		}
	}
	return nil
}

func (c *CategoryCache) all(ctx context.Context) error {
	var query entityBlog.BlogArticleCategoryEntity
	if c.dto.IsAll {
		query.State = enumStatePg.ENABLE.Index()
		infos := c.sp.catRep.FindAll(ctx, query)
		if infos != nil {
			mapTmp := make(map[string]modBlogArticleCategory.Cache)
			data := make(map[string]interface{})
			for _, info := range infos {
				//
				var obj modBlogArticleCategory.Cache
				copier.Copy(&obj, info)
				//
				mapTmp[info.No] = obj
				//
			}
			//
			mapKeysAdd := make(map[string][]string, 0)
			code_link := make([]string, 0)
			for _, info := range infos {
				//
				//
				code_link = make([]string, 0)
				if strPg.IsNotBlank(info.NoLink) {
					split := strings.Split(info.NoLink, "|")
					for _, item := range split {
						if strPg.IsNotBlank(item) {
							obj := mapTmp[item]
							code_link = append(code_link, obj.Code)
						}
					}
				}
				//
				if obj, ok := mapTmp[info.No]; ok {
					obj.CodeLink = strings.Join(code_link, "/")
					//
					key := blogKeyPg.ArticleCategoryTenantNo(info.TenantNo, info.No)
					str, err := json.Marshal(obj)
					if err == nil {
						data[key] = str
					}
					//
					if _, ok2 := mapKeysAdd[info.TenantNo]; !ok2 {
						mapKeysAdd[info.TenantNo] = make([]string, 0)
					}
					mapKeysAdd[info.TenantNo] = append(mapKeysAdd[info.TenantNo], key)
					//

					// 缓存  code = no
					key2 := blogKeyPg.ArticleCategoryTenantNoAndNoByCode(info.TenantNo, info.Code)
					data[key2] = info.No
				}
			}
			if len(data) > 0 {
				c.sp.rdt.SetPipeline(ctx, data)
			}
			if len(mapKeysAdd) > 0 {
				for tenantNo, keys := range mapKeysAdd {
					if nil == keys || len(keys) <= 0 {
						continue
					}
					//存入所有集合
					keysAll := blogKeyPg.ArticleCategoryTenantNoKeys(tenantNo)
					err := c.sp.rdt.GetRdb().SAdd(ctx, keysAll, keys).Err()
					if err != nil {
						log.Errorf(ctx, log.TagAppDef, err, "缓存失败:", err)
					}
				}
			}
		}
	}
	return nil
}

func (c *CategoryCache) custom(ctx context.Context) error {
	var query entityBlog.BlogArticleCategoryEntity
	if strPg.IsNotBlank(c.dto.TenantNo) && !c.dto.IsThisTenantAll && len(c.dto.Nos) > 0 {
		nos := make([]string, 0)
		for _, item := range c.dto.Nos {
			if strPg.IsNotBlank(item) {
				nos = append(nos, item)
			}
		}
		if len(nos) <= 0 {
			return nil
		}
		query.TenantNo = c.dto.TenantNo
		query.State = enumStatePg.ENABLE.Index()
		infos := c.sp.catRep.FindAll(ctx, query, optionsPg.WithCondition(func(db *gorm.DB) *gorm.DB {
			db = db.Order("create_at desc")
			db.Where("no in ?", nos)
			return db
		}))
		if infos != nil {
			keysAdd := make([]string, 0)
			mapTmp := make(map[string]modBlogArticleCategory.Cache)
			data := make(map[string]interface{})
			for _, info := range infos {
				//
				var obj modBlogArticleCategory.Cache
				copier.Copy(&obj, info)
				//
				mapTmp[info.No] = obj
			}
			//
			code_link := make([]string, 0)
			for _, info := range infos {
				code_link = make([]string, 0)
				if strPg.IsNotBlank(info.NoLink) {
					split := strings.Split(info.NoLink, "|")
					for _, item := range split {
						if strPg.IsNotBlank(item) {
							obj := mapTmp[item]
							code_link = append(code_link, obj.Code)
						}
					}
				}
				//
				if obj, ok := mapTmp[info.No]; ok {
					obj.CodeLink = strings.Join(code_link, "/")
					//
					key := blogKeyPg.ArticleCategoryTenantNo(info.TenantNo, info.No)
					str, err := json.Marshal(obj)
					if err == nil {
						data[key] = str
					}
					keysAdd = append(keysAdd, key)
					//

					// 缓存  code = no
					key2 := blogKeyPg.ArticleCategoryTenantNoAndNoByCode(info.TenantNo, info.Code)
					data[key2] = info.No
				}
			}
			if len(data) > 0 {
				c.sp.rdt.SetPipeline(ctx, data)
			}
			if len(keysAdd) > 0 {
				//存入所有集合
				keysAll := blogKeyPg.ArticleCategoryTenantNoKeys(c.dto.TenantNo)
				err := c.sp.rdt.GetRdb().SAdd(ctx, keysAll, keysAdd).Err()
				if err != nil {
					log.Errorf(ctx, log.TagAppDef, err, "缓存失败:", err)
				}
			}

		}
	}
	return nil
}
