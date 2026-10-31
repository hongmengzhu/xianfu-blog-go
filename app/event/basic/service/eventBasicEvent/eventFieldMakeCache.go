package eventBasicEvent

import (
	"context"

	"github.com/goccy/go-json"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/cache/basic/modCacheBasicEvent"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/event/basic/modEventBasicEvent"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/sdk/basic/key/basicEventKey"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg/optionsPg"
	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/strPg"
	"gorm.io/gorm"
)

type EventFieldMakeCache struct {
	Sp *Sp `autowire:"?"`
	ct modEventBasicEvent.FieldDto
}

func NewEventFieldMakeCache(sp *Sp, ct modEventBasicEvent.FieldDto) *EventFieldMakeCache {
	return &EventFieldMakeCache{
		Sp: sp,
		ct: ct,
	}
}

func (c *EventFieldMakeCache) Processor(ctx context.Context) error {
	if c.ct.IsAll {
		return c.All(ctx)
	}
	if c.ct.IsThisTenantAll {
		return c.ThisTenantAll(ctx)
	}
	if strPg.IsNotBlank(c.ct.TenantNo) && !c.ct.IsThisTenantAll && len(c.ct.EventNo) > 0 {
		return c.EventNos(ctx)
	}
	return nil
}

func (c *EventFieldMakeCache) ThisTenantAll(ctx context.Context) error {
	var query entityBasic.BasicConfigEventFieldsEntity
	if strPg.IsNotBlank(c.ct.TenantNo) {
		query.TenantNo = c.ct.TenantNo
		query.State = enumStatePg.ENABLE.Index()
		infos := c.Sp.repEventFields.FindAll(ctx, query)
		if infos != nil {
			keyNo := make(map[string]interface{})
			data := make(map[string]map[string]interface{})
			for _, info := range infos {
				var obj modCacheBasicEvent.FieldCache
				copier.Copy(&obj, info)

				key := basicEventKey.EventTenantNoAllFields(info.TenantNo, info.EventNo)
				if _, ok := data[key]; !ok {
					data[key] = make(map[string]interface{})
				}
				str, err := json.Marshal(obj)
				if err == nil {
					data[key][info.No] = str
				}
				// 字段 = 字段编号
				key3 := basicEventKey.EventFieldTenantNoByFieldNo(info.TenantNo, info.EventNo, info.Field)
				keyNo[key3] = info.No
			}
			if len(keyNo) > 0 {
				c.Sp.rdt.SetPipeline(ctx, keyNo)
			}
			if len(data) > 0 {
				c.Sp.rdt.HSetPipelineMapAll(ctx, data)
			}
		}
	}
	return nil
}

func (c *EventFieldMakeCache) All(ctx context.Context) error {
	var query entityBasic.BasicConfigEventFieldsEntity
	if c.ct.IsAll {
		query.State = enumStatePg.ENABLE.Index()
		infos := c.Sp.repEventFields.FindAll(ctx, query)
		if infos != nil {
			keyNo := make(map[string]interface{})
			data := make(map[string]map[string]interface{})
			for _, info := range infos {
				var obj modCacheBasicEvent.FieldCache
				copier.Copy(&obj, info)

				key := basicEventKey.EventTenantNoAllFields(info.TenantNo, info.EventNo)
				if _, ok := data[key]; !ok {
					data[key] = make(map[string]interface{})
				}
				str, err := json.Marshal(obj)
				if err == nil {
					data[key][info.No] = str
				}
				// 字段 = 字段编号
				key3 := basicEventKey.EventFieldTenantNoByFieldNo(info.TenantNo, info.EventNo, info.Field)
				keyNo[key3] = info.No
			}
			if len(keyNo) > 0 {
				c.Sp.rdt.SetPipeline(ctx, keyNo)
			}
			if len(data) > 0 {
				c.Sp.rdt.HSetPipelineMapAll(ctx, data)
			}
		}
	}
	return nil
}

func (c *EventFieldMakeCache) EventNos(ctx context.Context) error {
	var query entityBasic.BasicConfigEventFieldsEntity
	if strPg.IsNotBlank(c.ct.TenantNo) && !c.ct.IsThisTenantAll && len(c.ct.EventNo) > 0 {
		eventNos := make([]string, 0)
		for _, item := range c.ct.EventNo {
			if strPg.IsNotBlank(item) {
				eventNos = append(eventNos, item)
			}
		}
		if len(eventNos) <= 0 {
			return nil
		}
		query.TenantNo = c.ct.TenantNo
		query.State = enumStatePg.ENABLE.Index()
		infos := c.Sp.repEventFields.FindAll(ctx, query, optionsPg.WithCondition(func(db *gorm.DB) *gorm.DB {
			db = db.Order("create_at desc")
			db = db.Where("event_no in ?", eventNos)
			return db
		}))
		if infos != nil {
			keyNo := make(map[string]interface{})
			data := make(map[string]map[string]interface{})
			for _, info := range infos {
				var obj modCacheBasicEvent.FieldCache
				copier.Copy(&obj, info)
				//
				key := basicEventKey.EventTenantNoAllFields(info.TenantNo, info.EventNo)
				if _, ok := data[key]; !ok {
					data[key] = make(map[string]interface{})
				}
				str, err := json.Marshal(obj)
				if err == nil {
					data[key][info.No] = str
				}
				// 字段 = 字段编号
				key3 := basicEventKey.EventFieldTenantNoByFieldNo(info.TenantNo, info.EventNo, info.Field)
				keyNo[key3] = info.No
			}
			if len(keyNo) > 0 {
				c.Sp.rdt.SetPipeline(ctx, keyNo)
			}
			if len(data) > 0 {
				c.Sp.rdt.HSetPipelineMapAll(ctx, data)
			}
		}
	}
	return nil
}
