package bmToolEvent

import (
	"context"
	"strings"

	"github.com/goccy/go-json"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/bm/modBmToolCategory"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/event/bm/modEventBmToolCategory"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityBm"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/sdk/bm/key/bmKeyPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg/optionsPg"
	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/strPg"
	"go-spring.org/log"
	"gorm.io/gorm"
)

// CategoryCache 分类缓存处理
type CategoryCache struct {
	sp  *Sp `autowire:"?"`
	dto modEventBmToolCategory.CacheDto
}

func NewCategoryCache(sp *Sp, dto modEventBmToolCategory.CacheDto) *CategoryCache {
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
	var query entityBm.BmToolCategoryEntity
	if strPg.IsNotBlank(c.dto.TenantNo) {
		query.TenantNo = c.dto.TenantNo
		query.State = enumStatePg.ENABLE.Index()
		infos := c.sp.catRep.FindAll(ctx, query)
		if infos != nil {
			mapTmp := make(map[string]modBmToolCategory.Cache)
			data := make(map[string]interface{})
			for _, info := range infos {
				//
				var obj modBmToolCategory.Cache
				copier.Copy(&obj, info)
				//
				mapTmp[info.No] = obj
			}
			//
			keysAdd := make([]string, 0)
			no_link := make([]string, 0)
			for _, info := range infos {
				no_link = make([]string, 0)
				if strPg.IsNotBlank(info.NoLink) {
					split := strings.Split(info.NoLink, "|")
					for _, item := range split {
						if strPg.IsNotBlank(item) {
							obj := mapTmp[item]
							no_link = append(no_link, obj.Code)
						}
					}
				}
				//
				if obj, ok := mapTmp[info.No]; ok {
					obj.NoLink = strings.Join(no_link, "/")
					//
					key := bmKeyPg.ToolCategoryTenantNo(info.TenantNo, info.No)
					str, err := json.Marshal(obj)
					if err == nil {
						data[key] = str
					}
					keysAdd = append(keysAdd, key)
					//

					// 缓存  code = no
					key2 := bmKeyPg.ToolCategoryTenantNoAndNoByCode(info.TenantNo, info.Code)
					data[key2] = info.No
				}
			}
			//存入所有集合
			if len(keysAdd) > 0 {
				keysAll := bmKeyPg.ToolCategoryTenantNoKeys(c.dto.TenantNo)
				err := c.sp.rdt.SAdd(ctx, keysAll, keysAdd)
				if err != nil {
					log.Error(ctx, log.TagAppDef, err, "缓存失败")
				}
			}
			//
			if len(data) > 0 {
				c.sp.rdt.SetPipeline(ctx, data)
			}
			log.Debug(ctx, log.TagAppDef, "thisAll", func() []log.Field {
				return []log.Field{
					log.Any("len", len(data)),
				}
			})
		}
		log.Info(ctx, log.TagAppDef, "缓存成功")
	}
	return nil
}

func (c *CategoryCache) all(ctx context.Context) error {
	var query entityBm.BmToolCategoryEntity
	if c.dto.IsAll {
		query.State = enumStatePg.ENABLE.Index()
		infos := c.sp.catRep.FindAll(ctx, query)
		if infos != nil {
			mapTmp := make(map[string]modBmToolCategory.Cache)
			data := make(map[string]interface{})
			for _, info := range infos {
				//
				var obj modBmToolCategory.Cache
				copier.Copy(&obj, info)
				//
				mapTmp[info.No] = obj
				//
			}
			//
			mapKeysAdd := make(map[string][]string, 0)
			no_link := make([]string, 0)
			for _, info := range infos {
				//
				//
				no_link = make([]string, 0)
				if strPg.IsNotBlank(info.NoLink) {
					split := strings.Split(info.NoLink, "|")
					for _, item := range split {
						if strPg.IsNotBlank(item) {
							obj := mapTmp[item]
							no_link = append(no_link, obj.Code)
						}
					}
				}
				//
				if obj, ok := mapTmp[info.No]; ok {
					obj.NoLink = strings.Join(no_link, "/")
					//
					key := bmKeyPg.ToolCategoryTenantNo(info.TenantNo, info.No)
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
					key2 := bmKeyPg.ToolCategoryTenantNoAndNoByCode(info.TenantNo, info.Code)
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
					keysAll := bmKeyPg.ToolCategoryTenantNoKeys(tenantNo)
					err := c.sp.rdt.SAdd(ctx, keysAll, keys)
					if err != nil {
						log.Error(ctx, log.TagAppDef, err, "缓存失败")
					}
				}
			}
			log.Debug(ctx, log.TagAppDef, "all", func() []log.Field {
				return []log.Field{
					log.Any("len", len(data)),
				}
			})
		}
		log.Info(ctx, log.TagAppDef, "缓存成功")
	}
	return nil
}

func (c *CategoryCache) custom(ctx context.Context) error {
	var query entityBm.BmToolCategoryEntity
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
			mapTmp := make(map[string]modBmToolCategory.Cache)
			data := make(map[string]interface{})
			for _, info := range infos {
				//
				var obj modBmToolCategory.Cache
				copier.Copy(&obj, info)
				//
				mapTmp[info.No] = obj
			}
			//
			no_link := make([]string, 0)
			for _, info := range infos {
				no_link = make([]string, 0)
				if strPg.IsNotBlank(info.NoLink) {
					split := strings.Split(info.NoLink, "|")
					for _, item := range split {
						if strPg.IsNotBlank(item) {
							obj := mapTmp[item]
							no_link = append(no_link, obj.Code)
						}
					}
				}
				//
				if obj, ok := mapTmp[info.No]; ok {
					obj.NoLink = strings.Join(no_link, "/")
					//
					key := bmKeyPg.ToolCategoryTenantNo(info.TenantNo, info.No)
					str, err := json.Marshal(obj)
					if err == nil {
						data[key] = str
					}
					keysAdd = append(keysAdd, key)
					//

					// 缓存  code = no
					key2 := bmKeyPg.ToolCategoryTenantNoAndNoByCode(info.TenantNo, info.Code)
					data[key2] = info.No
				}
			}
			if len(data) > 0 {
				c.sp.rdt.SetPipeline(ctx, data)
			}
			if len(keysAdd) > 0 {
				//存入所有集合
				keysAll := bmKeyPg.ToolCategoryTenantNoKeys(c.dto.TenantNo)
				err := c.sp.rdt.SAdd(ctx, keysAll, keysAdd)
				if err != nil {
					log.Error(ctx, log.TagAppDef, err, "缓存失败")
				}
			}
			log.Debug(ctx, log.TagAppDef, "custom", func() []log.Field {
				return []log.Field{
					log.Any("len", len(data)),
				}
			})
		}
		log.Info(ctx, log.TagAppDef, "缓存成功")
	}
	return nil
}
