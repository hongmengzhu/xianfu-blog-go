package eventBasicRules

import (
	"context"

	"github.com/goccy/go-json"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/cache/basic/modCacheBasicRules"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/event/basic/modEventBasicRules"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/sdk/basic/key/basicEventKey"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg/optionsPg"
	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/strPg"
	"gorm.io/gorm"
)

type RulesMakeCache struct {
	Sp *Sp `autowire:"?"`
	ct modEventBasicRules.RulesDto
}

func NewRulesMakeCache(sp *Sp, ct modEventBasicRules.RulesDto) *RulesMakeCache {
	return &RulesMakeCache{
		Sp: sp,
		ct: ct,
	}
}

func (c *RulesMakeCache) Processor(ctx context.Context) error {
	if c.ct.IsAll {
		return c.All(ctx)
	}
	if c.ct.IsThisTenantAll {
		return c.ThisTenantAll(ctx)
	}
	if strPg.IsNotBlank(c.ct.TenantNo) && !c.ct.IsThisTenantAll && len(c.ct.FieldNo) > 0 {
		return c.FieldNos(ctx)
	}
	return nil
}

func (c *RulesMakeCache) ThisTenantAll(ctx context.Context) error {
	var query entityBasic.BasicModelRulesEntity
	if strPg.IsNotBlank(c.ct.TenantNo) {
		query.TenantNo = c.ct.TenantNo
		query.State = enumStatePg.ENABLE.Index()
		infos := c.Sp.repRules.FindAll(ctx, query)
		if infos != nil {
			data := make(map[string]map[string]interface{})
			for _, info := range infos {
				var obj modCacheBasicRules.RulesCache
				copier.Copy(&obj, info)

				key := basicEventKey.RulesByEventFieldTenantNo(info.TenantNo, info.ValueNo, info.No)
				if _, ok := data[key]; !ok {
					data[key] = make(map[string]interface{})
				}
				str, err := json.Marshal(obj)
				if err == nil {
					data[key][obj.No] = str
				}
			}
			if len(data) > 0 {
				c.Sp.rdt.HSetPipelineMapAll(ctx, data)
			}
		}
	}
	return nil
}

func (c *RulesMakeCache) All(ctx context.Context) error {
	var query entityBasic.BasicModelRulesEntity
	if c.ct.IsAll {
		query.State = enumStatePg.ENABLE.Index()
		infos := c.Sp.repRules.FindAll(ctx, query)
		if infos != nil {
			data := make(map[string]map[string]interface{})
			for _, info := range infos {
				var obj modCacheBasicRules.RulesCache
				copier.Copy(&obj, info)

				key := basicEventKey.RulesByFieldTenantNo(info.TenantNo, info.ValueNo)
				if _, ok := data[key]; !ok {
					data[key] = make(map[string]interface{})
				}
				str, err := json.Marshal(obj)
				if err == nil {
					data[key][obj.No] = str
				}
			}
			if len(data) > 0 {
				c.Sp.rdt.HSetPipelineMapAll(ctx, data)
			}
		}
	}
	return nil
}

func (c *RulesMakeCache) FieldNos(ctx context.Context) error {
	var query entityBasic.BasicModelRulesEntity
	if strPg.IsNotBlank(c.ct.TenantNo) && !c.ct.IsThisTenantAll && len(c.ct.FieldNo) > 0 {
		fieldNos := make([]string, 0)
		for _, item := range c.ct.FieldNo {
			if strPg.IsNotBlank(item) {
				fieldNos = append(fieldNos, item)
			}
		}
		if len(fieldNos) <= 0 {
			return nil
		}
		query.TenantNo = c.ct.TenantNo
		query.State = enumStatePg.ENABLE.Index()
		infos := c.Sp.repRules.FindAll(ctx, query, optionsPg.WithCondition(func(db *gorm.DB) *gorm.DB {
			db = db.Order("create_at desc")
			db.Where("field_no in ?", fieldNos)
			return db
		}))
		if infos != nil {
			data := make(map[string]map[string]interface{})
			for _, info := range infos {
				var obj modCacheBasicRules.RulesCache
				copier.Copy(&obj, info)

				key := basicEventKey.RulesByEventFieldTenantNo(info.TenantNo, info.ValueNo, info.No)
				if _, ok := data[key]; !ok {
					data[key] = make(map[string]interface{})
				}
				str, err := json.Marshal(obj)
				if err == nil {
					data[key][obj.No] = str
				}
			}
			if len(data) > 0 {
				c.Sp.rdt.HSetPipelineMapAll(ctx, data)
			}
		}
	}
	return nil
}
