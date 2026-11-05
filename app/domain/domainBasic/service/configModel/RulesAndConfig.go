package configModel

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/basic/modBasicConfigModel"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumCommonPg/configModelPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumCommonPg/typeSysPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/state/yesNoPg/yesNoIntPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityBasic"
	"github.com/pangu-2/go-tools/tools/cryptPg"
	"github.com/pangu-2/go-tools/tools/noPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"go-spring.org/log"
	"gorm.io/datatypes"
)

// normalizeStringSlice 去空、去重，供 JSON 字段使用。
func normalizeStringSlice(in []string) []string {
	out := make([]string, 0, len(in))
	seen := make(map[string]bool, len(in))
	for _, item := range in {
		item = strings.TrimSpace(item)
		if item == "" || seen[item] {
			continue
		}
		seen[item] = true
		out = append(out, item)
	}
	return out
}

// saveRules 设计器保存时重建字段规则与模型级规则。
//
// 字段规则写入 basic_model_rules（value_no = 模型编号，field_no = 字段编号），
// 模型级规则写入 basic_config_model_rules（model_no = 模型编号）。
func (c *CreateUpdate) saveRules(ctx *gin.Context, tenantNo string) error {
	if strPg.IsBlank(c.model.No) {
		return nil
	}
	// 重建：先清理旧规则
	if err := c.Sp.repFieldRules.DeleteAllByValueNo(ctx, c.model.No); err != nil {
		log.Error(ctx, log.TagAppDef, err, "delete field rules error")
	}
	if err := c.Sp.repRule.DeleteAllByModelNo(ctx, c.model.No); err != nil {
		log.Error(ctx, log.TagAppDef, err, "delete model rules error")
	}
	//
	bodyByField := make(map[string]modBasicConfigModel.ItemVo, len(c.ct.Body))
	for _, item := range c.ct.Body {
		bodyByField[strings.TrimSpace(item.Field)] = item
	}
	//
	fieldRules := make([]*entityBasic.BasicModelRulesEntity, 0)
	for _, field := range c.fields {
		item, ok := bodyByField[field.Field]
		if !ok {
			continue
		}
		ruleItems := item.RuleItems
		if len(ruleItems) == 0 && len(item.Rules) > 0 {
			for _, mode := range normalizeStringSlice(item.Rules) {
				ruleItems = append(ruleItems, modBasicConfigModel.RuleItemVo{
					Name:     field.Name + "-" + mode,
					RuleMode: mode,
					Show:     yesNoIntPg.Yes.IndexPg(),
				})
			}
		}
		for idx, rule := range ruleItems {
			if strPg.IsBlank(rule.RuleMode) {
				continue
			}
			rule.Name = strings.TrimSpace(rule.Name)
			if strPg.IsBlank(rule.Name) {
				rule.Name = field.Name + "-" + strings.TrimSpace(rule.RuleMode)
			}
			show := rule.Show.ToInt8()
			if show < 1 {
				show = yesNoIntPg.Yes.IndexInt8()
			}
			fieldRules = append(fieldRules, &entityBasic.BasicModelRulesEntity{
				No:           noPg.No(),
				Name:         rule.Name,
				TenantNo:     tenantNo,
				State:        enumStatePg.ENABLE.IndexInt8(),
				Sort:         int64(idx),
				Show:         show,
				Description:  field.Description,
				ValueType:    field.ValueType,
				ExtraData:    rule.ExtraData,
				RuleMode:     strings.TrimSpace(rule.RuleMode),
				Coding:       rule.Coding,
				Condition:    rule.Condition,
				ErrorMessage: rule.ErrorMessage,
				Structure:    rule.Structure,
				RuleTarget:   datatypes.NewJSONType(normalizeStringSlice(rule.RuleTarget)),
				TypeSys:      typeSysPg.General.String(),
				TypeModel:    "model",
				ValueNo:      c.model.No,
				FieldNo:      field.No,
			})
		}
	}
	if len(fieldRules) > 0 {
		tx := c.Sp.repFieldRules.DbModel().WithContext(ctx).CreateInBatches(fieldRules, 1000)
		if tx.Error != nil {
			log.Error(ctx, log.TagAppDef, tx.Error, "save field rules error")
			return tx.Error
		}
	}
	//
	modelRuleItems := c.ct.Header.RuleItems
	if len(modelRuleItems) == 0 && len(c.ct.Header.Rules) > 0 {
		for _, mode := range normalizeStringSlice(c.ct.Header.Rules) {
			modelRuleItems = append(modelRuleItems, modBasicConfigModel.RuleItemVo{
				Name:     c.model.Name + "-" + mode,
				RuleMode: mode,
				Show:     yesNoIntPg.Yes.IndexPg(),
			})
		}
	}
	modelRules := make([]*entityBasic.BasicConfigModelRulesEntity, 0, len(modelRuleItems))
	for idx, rule := range modelRuleItems {
		if strPg.IsBlank(rule.RuleMode) {
			continue
		}
		rule.Name = strings.TrimSpace(rule.Name)
		if strPg.IsBlank(rule.Name) {
			rule.Name = c.model.Name + "-" + strings.TrimSpace(rule.RuleMode)
		}
		show := rule.Show.ToInt8()
		if show < 1 {
			show = yesNoIntPg.Yes.IndexInt8()
		}
		modelRules = append(modelRules, &entityBasic.BasicConfigModelRulesEntity{
			No:           noPg.No(),
			Name:         rule.Name,
			TenantNo:     tenantNo,
			State:        enumStatePg.ENABLE.IndexInt8(),
			Sort:         int64(idx),
			Show:         show,
			ModelNo:      c.model.No,
			Description:  c.model.Description,
			ExtraData:    rule.ExtraData,
			RuleMode:     strings.TrimSpace(rule.RuleMode),
			Coding:       rule.Coding,
			Condition:    rule.Condition,
			ErrorMessage: rule.ErrorMessage,
			Structure:    rule.Structure,
			RuleTarget:   datatypes.NewJSONType(normalizeStringSlice(rule.RuleTarget)),
			SharedScope:  datatypes.NewJSONType([]string{}),
			TypeSys:      typeSysPg.General.String(),
			TypeCategory: "model",
		})
	}
	if len(modelRules) > 0 {
		tx := c.Sp.repRule.DbModel().WithContext(ctx).CreateInBatches(modelRules, 1000)
		if tx.Error != nil {
			log.Error(ctx, log.TagAppDef, tx.Error, "save model rules error")
			return tx.Error
		}
	}
	return nil
}

// syncConfig 模型保存后同步【配置】行：
//   - 每个模型字段对应一条 BasicConfig 记录，关联 model_no/model/field；
//   - 已存在的配置保留 value/value_binary，只同步展示与字段属性；
//   - 字段被删除时同步删除其配置行；
//   - 最后按 tenant_no + model 刷新 bigcache。
func (c *CreateUpdate) syncConfig(ctx *gin.Context, tenantNo string) error {
	if strPg.IsBlank(c.model.No) || strPg.IsBlank(c.model.Model) {
		return nil
	}
	existing, _ := c.Sp.repConfig.FindAllByModelNo(ctx, c.model.No)
	byField := make(map[string]*entityBasic.BasicConfigEntity, len(existing))
	for _, item := range existing {
		byField[item.Field] = item
	}
	categoryNo := strings.TrimSpace(c.model.Category)
	if categoryNo == "" {
		categoryNo = configModelPg.ModelCategoryConfig.Index()
	}
	seen := make(map[string]bool, len(c.fields))
	add := make([]*entityBasic.BasicConfigEntity, 0)
	for _, field := range c.fields {
		seen[field.Field] = true
		kindUnique := cryptPg.Md5(tenantNo + c.model.No + field.Field)
		old, exists := byField[field.Field]
		if exists {
			updates := map[string]any{
				"name":          field.Name,
				"field_path":    field.FieldPath,
				"description":   field.Description,
				"show":          field.Show,
				"binary":        field.Binary,
				"value_type":    field.ValueType,
				"form_code":     field.FormCode,
				"default_value": field.DefaultValue,
				"model_no":      c.model.No,
				"model":         c.model.Model,
				"module":        c.model.Module,
				"module_sub":    c.model.ModuleSub,
				"tenant_no":     tenantNo,
				"category_no":   categoryNo,
				"sort":          field.Sort,
				"kind_unique":   kindUnique,
			}
			if err := c.Sp.repConfig.UpdateMap(ctx, updates, old.ID); err != nil {
				log.Error(ctx, log.TagAppDef, err, "sync config update error")
				return err
			}
			continue
		}
		value := field.Value
		if value == "" {
			value = field.DefaultValue
		}
		add = append(add, &entityBasic.BasicConfigEntity{
			No:           noPg.No(),
			TenantNo:     tenantNo,
			State:        enumStatePg.ENABLE.IndexInt8(),
			Sort:         field.Sort,
			ModelNo:      c.model.No,
			Model:        c.model.Model,
			Module:       c.model.Module,
			ModuleSub:    c.model.ModuleSub,
			Field:        field.Field,
			FieldPath:    field.FieldPath,
			Name:         field.Name,
			Description:  field.Description,
			Show:         field.Show,
			Binary:       field.Binary,
			ValueType:    field.ValueType,
			FormCode:     field.FormCode,
			DefaultValue: field.DefaultValue,
			Value:        value,
			ValueBinary:  field.ValueBinary,
			CategoryNo:   categoryNo,
			KindUnique:   kindUnique,
		})
	}
	if len(add) > 0 {
		tx := c.Sp.repConfig.DbModel().WithContext(ctx).CreateInBatches(add, 1000)
		if tx.Error != nil {
			log.Error(ctx, log.TagAppDef, tx.Error, "sync config create error")
			return tx.Error
		}
	}
	// 字段已删除 -> 配置行删除
	delIds := make([]int64, 0)
	for field, item := range byField {
		if !seen[field] {
			delIds = append(delIds, item.ID)
		}
	}
	if len(delIds) > 0 {
		if err := c.Sp.repConfig.DeleteByIds(ctx, delIds); err != nil {
			log.Error(ctx, log.TagAppDef, err, "sync config delete error")
			return err
		}
	}
	// 缓存重建（tenant_no + model）
	if c.Sp.cache != nil {
		if err := c.Sp.cache.Refresh(ctx, tenantNo, c.model.Model); err != nil {
			log.Error(ctx, log.TagAppDef, err, "refresh config cache error")
			return err
		}
	}
	return nil
}
