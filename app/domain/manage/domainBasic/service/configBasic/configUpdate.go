package configBasic

import (
	"regexp"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/hongmengzhu/xianfu-blog-go/app/core/basic/model/modCacheBasicEvent"
	"github.com/hongmengzhu/xianfu-blog-go/app/core/basic/model/modCacheBasicRules"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/basic/modBasicConfigList"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/holderPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/cachePg/rdsPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/enumCommonPg/configModelPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/sdk/basic/key/basicEventKey"
	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/log"
)

// ConfigUpdate 配置更新
// @Description:
type ConfigUpdate struct {
	Sp *Sp `autowire:"?"`

	rdt *rdsPg.BatchString `autowire:"?"`
}

// NewConfigUpdate
//
//	@Description: 配置更新
//	@param sp
//	@return *ConfigUpdateCt
func NewConfigUpdate(sp *Sp) *ConfigUpdate {
	return &ConfigUpdate{
		Sp:  sp,
		rdt: sp.rdt,
	}
}

// Process
//
//	@Description:
//	@receiver c
//	@param ctx
//	@param ct
//	@return rt
func (c *ConfigUpdate) Process(ctx *gin.Context, ct modBasicConfigList.ConfigUpdateCt) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	if strPg.IsBlank(ct.EventNo) {
		return rt.ErrorMessage("配置编号不能为空")
	}
	if nil == ct.Data || len(ct.Data) < 1 {
		return rt.ErrorMessage("配置数据不能为空")
	}
	//
	holder := holderPg.GetContextAccount(ctx)
	config, b2 := c.Sp.repConfigList.FindByEventNo(ctx, ct.EventNo)
	if !b2 {
		return rt.ErrorMessage("配置不存在")
	}
	if holder.GetTenantNo() != config.TenantNo {
		return rt.ErrorMessage("配置不存在")
	}
	hashKeys := make([]string, 0)
	hashKeysEqField := make(map[string]string)
	mapField := make(map[string]*modCacheBasicEvent.FieldCache)
	mapRules := make(map[string][]*modCacheBasicRules.RulesCache)
	key2 := basicEventKey.EventTenantNoAllFields(config.TenantNo, config.EventNo)
	all, b2 := c.rdt.HGetAll(ctx, key2)
	if b2 {
		for _, v := range all {
			var obj modCacheBasicEvent.FieldCache
			err := json.Unmarshal([]byte(v), &obj)
			if err != nil {
				log.Errorf(ctx, log.TagAppDef, "json.Unmarshal.err:%+v", err)
			} else {
				copier.Copy(&obj, v)
				//
				mapField[obj.Field] = &obj
				//
				key := basicEventKey.RulesByEventFieldTenantNo(obj.TenantNo, obj.EventNo, obj.No)
				hashKeys = append(hashKeys, key)
				//获取到 验证字段 = key
				hashKeysEqField[obj.Field] = key
			}
		}
	}
	// 获取字段验证
	{
		if len(hashKeys) > 0 {
			all2, b3 := c.rdt.HGetAllPipeline(ctx, hashKeys)
			if b3 {
				for k, v := range all2 {
					mapRules[k] = make([]*modCacheBasicRules.RulesCache, 0)
					for _, v2 := range v {
						var obj modCacheBasicRules.RulesCache
						err := json.Unmarshal([]byte(v2), &obj)
						if err != nil {
							log.Errorf(ctx, log.TagAppDef, "json.Unmarshal.err:%+v", err)
						} else {
							mapRules[k] = append(mapRules[k], &obj)
						}
					}
					//排序
					slices.SortFunc(mapRules[k], func(a, b *modCacheBasicRules.RulesCache) int {
						if a.Sort < b.Sort {
							return -1
						} else if a.Sort > b.Sort {
							return 1
						}
						return 0
					})
				}
			}
		}
	}
	// 数据库字段
	dbFields := make(map[string]string)
	//
	info, result := c.Sp.repConfig.FindByEventNo(ctx, config.EventNo)
	if result {
		for _, item := range info {
			dbFields[item.Field] = item.No
		}
	}
	dbUpdate := make(map[string]string)
	//  处理
	for k, v := range ct.Data {
		//键 不能为空
		if strPg.IsBlank(k) {
			continue
		}
		//判断 值 在数据库中是否存在，不存在 跳过
		if _, ok := dbFields[k]; !ok {
			continue
		}
		//是否 有指定 字段验证
		if keyCache, ok := hashKeysEqField[k]; ok {
			if rules, ok2 := mapRules[keyCache]; ok2 {
				value := v.(string)
				// 字段名
				fieldName := ""
				field, ok3 := mapField[k]
				if ok3 {
					fieldName = field.Name
				}
				for _, item := range rules {
					//不为空
					if configModelPg.RuleModeRequired.IsEqual(item.RuleMode) {
						if strPg.IsBlank(value) {
							//错误信息
							if strPg.IsNotBlank(item.ErrorMessage) {
								return rt.ErrorMessage(item.ErrorMessage)
							}
							return rt.ErrorMessage("字段 " + fieldName + " 不能为空")
						}
					}
					//大于0
					if configModelPg.RuleModeGreaterThan0.IsEqual(item.RuleMode) {
						if strPg.ToInt64(value) <= 0 {
							//错误信息
							if strPg.IsNotBlank(item.ErrorMessage) {
								return rt.ErrorMessage(item.ErrorMessage)
							}
							return rt.ErrorMessage("字段 " + fieldName + " 必须大于0")
						}
					}
					//正则
					if configModelPg.RuleModePattern.IsEqual(item.RuleMode) && strPg.IsNotBlank(item.Condition) {
						match, _ := regexp.MatchString(strings.TrimSpace(item.Condition), value)
						if !match {
							if strPg.IsNotBlank(item.ErrorMessage) {
								return rt.ErrorMessage(item.ErrorMessage)
							}
							return rt.ErrorMessage("字段 " + fieldName + " 正则匹配错误")
						}
					}
				}
				//数据库字段
				if ok3 {
					dbUpdate[field.Field] = value
				}
			}
		}
	}
	// 更新数据库
	log.Infof(ctx, log.TagAppDef, "dbUpdate=%+v", dbUpdate)
	if len(dbUpdate) > 0 {
		for key, val := range dbUpdate {
			if no, ok := dbFields[key]; ok {
				c.Sp.repConfig.UpdateByTenantEventNoAndNoAndValue(ctx, config.TenantNo, config.EventNo, no, val)
			}
		}
	}
	return rt.Ok()
}
