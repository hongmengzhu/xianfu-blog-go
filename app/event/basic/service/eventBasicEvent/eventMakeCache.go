package eventBasicEvent

import (
	"context"

	"github.com/goccy/go-json"
	"github.com/hongmengzhu/xianfu-blog-go/app/core/basic/model/modCacheBasicEvent"
	"github.com/hongmengzhu/xianfu-blog-go/app/event/basic/model/modEventBasicEvent"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/sdk/basic/key/basicEventKey"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg/optionsPg"
	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/strPg"
	"gorm.io/gorm"
)

type EventMakeCache struct {
	Sp *Sp `autowire:"?"`
	ct modEventBasicEvent.EventDto
}

func NewEventMakeCache(sp *Sp, ct modEventBasicEvent.EventDto) *EventMakeCache {
	return &EventMakeCache{
		Sp: sp,
		ct: ct,
	}
}

func (c *EventMakeCache) Processor(ctx context.Context) error {
	if c.ct.IsAll {
		return c.All(ctx)
	}
	if c.ct.IsThisTenantAll {
		return c.ThisTenantAll(ctx)
	}
	if strPg.IsNotBlank(c.ct.TenantNo) && !c.ct.IsThisTenantAll && len(c.ct.Nos) > 0 {
		return c.Nos(ctx)
	}
	return nil
}

func (c *EventMakeCache) ThisTenantAll(ctx context.Context) error {
	var query entityBasic.BasicConfigEventEntity
	if strPg.IsNotBlank(c.ct.TenantNo) {
		query.TenantNo = c.ct.TenantNo
		query.State = enumStatePg.ENABLE.Index()
		infos := c.Sp.repEvent.FindAll(ctx, query)
		if infos != nil {
			data := make(map[string]interface{})
			//keysAdd := make([]string, 0)
			for _, info := range infos {
				var obj modCacheBasicEvent.EventCache
				copier.Copy(&obj, info)
				//
				//key := basicEventKey.EventTenantNo(info.TenantNo, info.No)
				//data[key] = info.No
				//keysAdd = append(keysAdd, key)

				key2 := basicEventKey.EventTenantNo(info.TenantNo, info.No)
				str, err := json.Marshal(obj)
				if err == nil {
					data[key2] = str
				}
				// 字段 = 字段编号
				key3 := basicEventKey.EventTenantNoByEventNo(info.TenantNo, info.Field)
				data[key3] = info.No
			}
			if len(data) > 0 {
				c.Sp.rdt.SetPipeline(ctx, data)
			}
			//if len(keysAdd) > 0 {
			//	keysAll := basicEventKey.EventTenantNoKeys(c.ct.TenantNo)
			//	err := c.Sp.rdt.SAdd(ctx, keysAll, keysAdd)
			//	if err != nil {
			//		log.Error(ctx,log.TagAppDef,"缓存失败:", err)
			//	}
			//}
		}
	}
	return nil
}

func (c *EventMakeCache) All(ctx context.Context) error {
	var query entityBasic.BasicConfigEventEntity
	if c.ct.IsAll {
		query.State = enumStatePg.ENABLE.Index()
		infos := c.Sp.repEvent.FindAll(ctx, query)
		if infos != nil {
			data := make(map[string]interface{})
			//mapKeysAdd := make(map[string][]string, 0)
			for _, info := range infos {
				var obj modCacheBasicEvent.EventCache
				copier.Copy(&obj, info)
				//
				key := basicEventKey.EventTenantNo(info.TenantNo, info.No)
				str, err := json.Marshal(obj)
				if err == nil {
					data[key] = str
				}

				// 字段 = 字段编号
				key3 := basicEventKey.EventTenantNoByEventNo(info.TenantNo, info.Field)
				data[key3] = info.No
			}
			if len(data) > 0 {
				c.Sp.rdt.SetPipeline(ctx, data)
			}
			//if len(mapKeysAdd) > 0 {
			//	for tenantNo, keys := range mapKeysAdd {
			//		if nil == keys || len(keys) <= 0 {
			//			continue
			//		}
			//		keysAll := basicEventKey.EventTenantNoKeys(tenantNo)
			//		err := c.Sp.rdt.SAdd(ctx, keysAll, keys)
			//		if err != nil {
			//			log.Error(ctx,log.TagAppDef,"缓存失败:", err)
			//		}
			//	}
			//}
		}
	}
	return nil
}

func (c *EventMakeCache) Nos(ctx context.Context) error {
	var query entityBasic.BasicConfigEventEntity
	if strPg.IsNotBlank(c.ct.TenantNo) && !c.ct.IsThisTenantAll && len(c.ct.Nos) > 0 {
		nos := make([]string, 0)
		for _, item := range c.ct.Nos {
			if strPg.IsNotBlank(item) {
				nos = append(nos, item)
			}
		}
		if len(nos) <= 0 {
			return nil
		}
		query.TenantNo = c.ct.TenantNo
		query.State = enumStatePg.ENABLE.Index()
		infos := c.Sp.repEvent.FindAll(ctx, query, optionsPg.WithCondition(func(db *gorm.DB) *gorm.DB {
			db = db.Order("create_at desc")
			db.Where("no in ?", nos)
			return db
		}))
		if infos != nil {
			//keysAdd := make([]string, 0)
			data := make(map[string]interface{})
			for _, info := range infos {
				var obj modCacheBasicEvent.EventCache
				copier.Copy(&obj, info)
				//
				key := basicEventKey.EventTenantNo(info.TenantNo, info.No)
				str, err := json.Marshal(obj)
				if err == nil {
					data[key] = str
				}
				//keysAdd = append(keysAdd, key)

				// 字段 = 字段编号
				key3 := basicEventKey.EventTenantNoByEventNo(info.TenantNo, info.Field)
				data[key3] = info.No
			}
			if len(data) > 0 {
				c.Sp.rdt.SetPipeline(ctx, data)
			}
			//if len(keysAdd) > 0 {
			//	keysAll := basicEventKey.EventTenantNoKeys(c.ct.TenantNo)
			//	err := c.Sp.rdt.SAdd(ctx, keysAll, keysAdd)
			//	if err != nil {
			//		log.Error(ctx,log.TagAppDef,"缓存失败:", err)
			//	}
			//}
		}
	}
	return nil
}
