package service

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/domainBasic/service/configCache"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/basic/modBasicConfig"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/cache/basic/modCacheBasicConfig"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumCommonPg/configModelPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/state/yesNoPg/yesNoIntPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/repositoryBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/sdk/auth/holderPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg/optionsPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/cryptPg"
	"github.com/pangu-2/go-tools/tools/dbPg/pagePg"
	"github.com/pangu-2/go-tools/tools/noPg"
	"github.com/pangu-2/go-tools/tools/numberPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(BasicConfigService))
}

// BasicConfigService 配置（BasicConfigEntity）增删改查与缓存。
//
// 配置行与模型字段通过 model_no / model / field 关联；缓存统一使用
// tenant_no + model 作为 key，其他服务通过 GetByTenantNoAndModel /
// ValuesByTenantNoAndModel 读取即可，命中时不会访问数据库。
type BasicConfigService struct {
	sv       *repositoryBasic.BasicConfigRepository            `autowire:"?"`
	repModel *repositoryBasic.BasicConfigModelRepository       `autowire:"?"`
	repField *repositoryBasic.BasicConfigModelFieldsRepository `autowire:"?"`
	repRules *repositoryBasic.BasicModelRulesRepository        `autowire:"?"`
	sp       *configCache.Sp                                   `autowire:"?"`
}

// CreateUpdate 新增/更新配置项
func (c *BasicConfigService) CreateUpdate(ctx *gin.Context, ct modBasicConfig.CreateUpdateCt) (rt rg.Rs[string]) {
	log.Info(ctx, log.TagAppDef, "ct", log.Any("ct", ct))
	holder := holderPg.GetContextAccount(ctx)
	tenantNo := holder.GetTenantNo()
	ct.Field = strings.TrimSpace(ct.Field)
	if strPg.IsBlank(ct.Field) {
		return rt.ErrorMessage("字段标识不能为空")
	}
	modelInfo, ok := c.findModel(ctx, tenantNo, strings.TrimSpace(ct.ModelNo), strings.TrimSpace(ct.Model))
	if !ok {
		return rt.ErrorMessage("模型不存在")
	}
	fieldInfo, ok := c.repField.FindByModelNoAndField(ctx, modelInfo.No, ct.Field)
	if !ok {
		return rt.ErrorMessage("模型字段不存在，请先在设计器保存字段")
	}
	show := normalizeYesNo(ct.Show.ToInt8(), fieldInfo.Show)
	binary := normalizeYesNo(ct.Binary.ToInt8(), fieldInfo.Binary)
	name := strings.TrimSpace(ct.Name)
	if strPg.IsBlank(name) {
		name = fieldInfo.Name
	}
	fieldPath := strings.TrimSpace(ct.FieldPath)
	if strPg.IsBlank(fieldPath) {
		fieldPath = fieldInfo.FieldPath
	}
	if strPg.IsBlank(fieldPath) {
		fieldPath = fieldInfo.Field
	}
	valueType := strings.TrimSpace(ct.ValueType)
	if strPg.IsBlank(valueType) {
		valueType = fieldInfo.ValueType
	}
	defaultValue := strings.TrimSpace(ct.DefaultValue)
	if strPg.IsBlank(defaultValue) {
		defaultValue = fieldInfo.DefaultValue
	}
	categoryNo := strings.TrimSpace(ct.CategoryNo)
	if strPg.IsBlank(categoryNo) {
		categoryNo = modelInfo.Category
	}
	if strPg.IsBlank(categoryNo) {
		categoryNo = configModelPg.ModelCategoryConfig.Index()
	}
	sort := ct.Sort.ToInt64()
	if sort == 0 {
		sort = fieldInfo.Sort
	}
	//
	if ct.ID.ToInt64() <= 0 {
		if _, exists := c.sv.FindByTenantNoAndModelNoAndField(ctx, tenantNo, modelInfo.No, ct.Field); exists {
			return rt.ErrorMessage("该模型字段的配置已存在")
		}
		value := ct.Value
		if value == "" {
			value = defaultValue
		}
		info := entityBasic.BasicConfigEntity{
			No:           noPg.No(),
			TenantNo:     tenantNo,
			State:        enumStatePg.ENABLE.IndexInt8(),
			Sort:         sort,
			ModelNo:      modelInfo.No,
			Model:        modelInfo.Model,
			Module:       modelInfo.Module,
			ModuleSub:    modelInfo.ModuleSub,
			Field:        ct.Field,
			FieldPath:    fieldPath,
			Name:         name,
			Description:  strings.TrimSpace(ct.Description),
			Show:         show,
			Binary:       binary,
			ValueType:    valueType,
			FormCode:     ct.FormCode,
			DefaultValue: defaultValue,
			Value:        value,
			ValueBinary:  ct.ValueBinary,
			CategoryNo:   categoryNo,
			KindUnique:   cryptPg.Md5(tenantNo + modelInfo.No + ct.Field),
		}
		err, _ := c.sv.Create(ctx, &info)
		if err != nil {
			log.Error(ctx, log.TagAppDef, err, "create config error")
			return rt.ErrorMessage("保存失败 " + err.Error())
		}
		c.refresh(ctx, tenantNo, modelInfo.Model)
		return rg.OkData(numberPg.Int64ToString(info.ID))
	}
	// 更新
	find, ok := c.sv.FindById(ctx, ct.ID.ToInt64())
	if !ok {
		return rt.ErrorMessage("配置不存在")
	}
	if strPg.IsNotBlank(find.TenantNo) && find.TenantNo != tenantNo {
		return rt.ErrorMessage("无权操作其他租户的配置")
	}
	if dup, exists := c.sv.FindByTenantNoAndModelNoAndField(ctx, tenantNo, modelInfo.No, ct.Field); exists && dup.ID != find.ID {
		return rt.ErrorMessage("该模型字段的配置已存在")
	}
	updates := map[string]any{
		"model_no":      modelInfo.No,
		"model":         modelInfo.Model,
		"module":        modelInfo.Module,
		"module_sub":    modelInfo.ModuleSub,
		"field":         ct.Field,
		"field_path":    fieldPath,
		"name":          name,
		"description":   strings.TrimSpace(ct.Description),
		"show":          show,
		"binary":        binary,
		"value_type":    valueType,
		"form_code":     ct.FormCode,
		"default_value": defaultValue,
		"value":         ct.Value,
		"value_binary":  ct.ValueBinary,
		"category_no":   categoryNo,
		"sort":          sort,
		"kind_unique":   cryptPg.Md5(tenantNo + modelInfo.No + ct.Field),
	}
	if err := c.sv.UpdateMap(ctx, updates, find.ID); err != nil {
		log.Error(ctx, log.TagAppDef, err, "update config error")
		return rt.ErrorMessage("保存失败 " + err.Error())
	}
	// model 可能已变化，旧 model 缓存也要刷新
	c.refresh(ctx, tenantNo, find.Model)
	c.refresh(ctx, tenantNo, modelInfo.Model)
	return rg.OkData(numberPg.Int64ToString(find.ID))
}

// Detail 配置详情
func (c *BasicConfigService) Detail(ctx *gin.Context, id int64) (rt rg.Rs[modBasicConfig.Vo]) {
	if id < 1 {
		return rt.ErrorMessage("id错误")
	}
	find, ok := c.sv.FindById(ctx, id)
	if !ok {
		return rt.ErrorMessage("数据不存在")
	}
	var vo modBasicConfig.Vo
	if err := copier.Copy(&vo, find); err != nil {
		log.Info(ctx, log.TagAppDef, "copier.Copy error:", log.Any("err", err.Error()))
	}
	return rt.OkData(vo)
}

// Query 配置分页查询
func (c *BasicConfigService) Query(ctx *gin.Context, ct modBasicConfig.QueryCt) (rt rg.Rs[pagePg.Paginator[modBasicConfig.Vo]]) {
	log.Info(ctx, log.TagAppDef, "ct", log.Any("ct", ct))
	var query entityBasic.BasicConfigEntity
	copier.Copy(&query, &ct)
	slice := make([]modBasicConfig.Vo, 0)
	rt.Data.Data = slice
	holder := holderPg.GetContextAccount(ctx)
	page, err := c.sv.FindAllPage(ctx, query, optionsPg.WithOption(func(arg *optionsPg.OptionParams) {
		if ct.PageSize < 1 {
			ct.PageSize = 20
		}
		arg.Pageable = new(pagePg.PageablePageSize(0, ct.PageNum, ct.PageSize))
		arg.Db = arg.Db.Where("tenant_no=?", holder.GetTenantNo()).Order("model, sort asc, create_at")
		if strPg.IsNotBlank(ct.Wd) {
			arg.Db = arg.Db.Where("name like ? or field like ? or description like ?", "%"+ct.Wd+"%", "%"+ct.Wd+"%", "%"+ct.Wd+"%")
		}
	}), optionsPg.WithCtx(ctx))
	if err != nil {
		log.Error(ctx, log.TagAppDef, err, "query config error")
		return rt.Ok()
	}
	if page.Total > 0 && page.Data != nil && len(page.Data) > 0 {
		pg := pagePg.NewPaginatorByPageable[modBasicConfig.Vo](page.Pageable)
		for _, item := range page.Data {
			var vo modBasicConfig.Vo
			copier.Copy(&vo, item)
			slice = append(slice, vo)
		}
		pg.Data = slice
		pg.Pageable = page.Pageable
		rt.Data = pg
	}
	return rt.Ok()
}

// AllByModel 查询当前租户指定模型的全部启用配置（走缓存）。
func (c *BasicConfigService) AllByModel(ctx *gin.Context, ct modBasicConfig.QueryCt) (rt rg.Rs[[]modBasicConfig.Vo]) {
	holder := holderPg.GetContextAccount(ctx)
	model := strings.TrimSpace(ct.Model)
	if strPg.IsBlank(model) && strPg.IsNotBlank(ct.ModelNo) {
		if info, ok := c.repModel.FindByTenantNoAndNo(ctx, holder.GetTenantNo(), strings.TrimSpace(ct.ModelNo)); ok {
			model = info.Model
		}
	}
	if strPg.IsBlank(model) {
		return rt.ErrorMessage("模型不能为空")
	}
	items, ok := c.GetByTenantNoAndModel(ctx, holder.GetTenantNo(), model)
	if !ok {
		rt.Data = make([]modBasicConfig.Vo, 0)
		return rt.Ok()
	}
	slice := make([]modBasicConfig.Vo, 0, len(items))
	for _, item := range items {
		slice = append(slice, modBasicConfig.Vo{
			ID:           typePg.Int64String(item.ID),
			TenantNo:     item.TenantNo,
			No:           item.No,
			Name:         item.Name,
			ModelNo:      item.ModelNo,
			Model:        item.Model,
			Module:       item.Module,
			ModuleSub:    item.ModuleSub,
			Field:        item.Field,
			FieldPath:    item.FieldPath,
			Description:  item.Description,
			Show:         typePg.Int8(item.Show),
			Binary:       typePg.Int8(item.Binary),
			ValueType:    item.ValueType,
			FormCode:     item.FormCode,
			DefaultValue: item.DefaultValue,
			Value:        item.Value,
			ValueBinary:  item.ValueBinary,
			CategoryNo:   item.CategoryNo,
			Sort:         typePg.Int64String(item.Sort),
		})
	}
	rt.Data = slice
	return rt.Ok()
}

// GetByTenantNoAndModel 按 tenant_no + model 读取配置（优先 bigcache）。
// 这是其他服务读取配置的统一入口。
func (c *BasicConfigService) GetByTenantNoAndModel(ctx context.Context, tenantNo, model string) ([]modCacheBasicConfig.Item, bool) {
	if strPg.IsBlank(tenantNo) || strPg.IsBlank(model) {
		return nil, false
	}
	if items, ok := c.sp.Get(ctx, tenantNo, model); ok {
		return items, true
	}
	infos, ok := c.sv.FindAllEnabledByTenantNoAndModel(ctx, tenantNo, model)
	if !ok {
		_ = c.sp.Set(ctx, tenantNo, model, []modCacheBasicConfig.Item{})
		return nil, false
	}
	items := make([]modCacheBasicConfig.Item, 0, len(infos))
	for _, info := range infos {
		items = append(items, configCache.CopyItem(info))
	}
	_ = c.sp.Set(ctx, tenantNo, model, items)
	return items, true
}

// ValuesByTenantNoAndModel 按 tenant_no + model 读取 field -> value 的配置值。
// 二进制字段返回 value_binary。
func (c *BasicConfigService) ValuesByTenantNoAndModel(ctx context.Context, tenantNo, model string) (map[string]string, bool) {
	items, ok := c.GetByTenantNoAndModel(ctx, tenantNo, model)
	if !ok {
		return nil, false
	}
	values := make(map[string]string, len(items))
	for _, item := range items {
		values[item.Field] = item.ValueMap()
	}
	return values, true
}

// UpdateValues 按模型批量保存配置值，并按字段规则校验。
func (c *BasicConfigService) UpdateValues(ctx *gin.Context, ct modBasicConfig.UpdateValuesCt) (rt rg.Rs[string]) {
	holder := holderPg.GetContextAccount(ctx)
	tenantNo := holder.GetTenantNo()
	modelInfo, ok := c.findModel(ctx, tenantNo, strings.TrimSpace(ct.ModelNo), strings.TrimSpace(ct.Model))
	if !ok {
		return rt.ErrorMessage("模型不存在")
	}
	if len(ct.Data) == 0 {
		return rt.ErrorMessage("配置数据不能为空")
	}
	infos, _ := c.sv.FindAllByModelNo(ctx, modelInfo.No)
	byField := make(map[string]*entityBasic.BasicConfigEntity, len(infos))
	for _, item := range infos {
		byField[item.Field] = item
	}
	fields, _ := c.repField.FindAllByModelNo(ctx, modelInfo.No)
	fieldByField := make(map[string]*entityBasic.BasicConfigModelFieldsEntity, len(fields))
	for _, item := range fields {
		fieldByField[item.Field] = item
	}
	rules, _ := c.repRules.FindAllByValueNo(ctx, modelInfo.No)
	rulesByField := make(map[string][]*entityBasic.BasicModelRulesEntity)
	for _, item := range rules {
		rulesByField[item.FieldNo] = append(rulesByField[item.FieldNo], item)
	}
	updates := make([]*entityBasic.BasicConfigEntity, 0, len(ct.Data))
	updateFields := make([]map[string]any, 0, len(ct.Data))
	updateIds := make([]int64, 0, len(ct.Data))
	for field, value := range ct.Data {
		field = strings.TrimSpace(field)
		if strPg.IsBlank(field) {
			continue
		}
		info, exists := byField[field]
		if !exists {
			return rt.ErrorMessage("配置字段不存在：" + field)
		}
		fieldInfo := fieldByField[info.Field]
		fieldNo := ""
		if fieldInfo != nil {
			fieldNo = fieldInfo.No
		}
		if msg := validateRules(value, rulesByField[fieldNo], fieldInfo); msg != "" {
			return rt.ErrorMessage(msg)
		}
		// 二进制字段写 value_binary，其余写 value
		if info.Binary == yesNoIntPg.Yes.IndexInt8() {
			updateFields = append(updateFields, map[string]any{"value_binary": value, "value": ""})
		} else {
			updateFields = append(updateFields, map[string]any{"value": value})
		}
		updates = append(updates, info)
		updateIds = append(updateIds, info.ID)
	}
	for i := range updates {
		if err := c.sv.UpdateMap(ctx, updateFields[i], updateIds[i]); err != nil {
			log.Error(ctx, log.TagAppDef, err, "update config value error")
			return rt.ErrorMessage("保存失败 " + err.Error())
		}
	}
	c.refresh(ctx, tenantNo, modelInfo.Model)
	return rt.Ok()
}

// CacheOverride 全量重建配置缓存。
func (c *BasicConfigService) CacheOverride(ctx *gin.Context) {
	if err := c.sp.Override(ctx); err != nil {
		log.Error(ctx, log.TagAppDef, err, "override config cache error")
	}
}

// Enable 启用
func (c *BasicConfigService) Enable(ctx *gin.Context, ct model.BaseIdsCt[string]) (rt rg.Rs[string]) {
	return c.State(ctx, ct.Ids, enumStatePg.ENABLE)
}

// Disable 禁用
func (c *BasicConfigService) Disable(ctx *gin.Context, ct model.BaseIdsCt[string]) (rt rg.Rs[string]) {
	return c.State(ctx, ct.Ids, enumStatePg.GetType(enumStatePg.DISABLE))
}

// State 状态 启用/禁用
func (c *BasicConfigService) State(ctx *gin.Context, ids []string, state enumStatePg.State) (rt rg.Rs[string]) {
	if len(ids) < 1 {
		return rt.ErrorMessage("id错误")
	}
	finds, ok := c.sv.FindAllByIdStringIn(ctx, ids, optionsPg.WithCtx(ctx))
	if !ok {
		return rt.ErrorMessage("数据不存在")
	}
	for _, info := range finds {
		if info.State != state.IndexInt8() {
			c.sv.Update(ctx, entityBasic.BasicConfigEntity{State: state.IndexInt8()}, info.ID)
		}
	}
	c.refreshByConfigs(ctx, finds)
	return rt.Ok()
}

// StateEnableDisable 状态 设置 有效 停用
func (c *BasicConfigService) StateEnableDisable(ctx *gin.Context, ids []string, state enumStatePg.State) (rt rg.Rs[string]) {
	if !state.IsEnableDisable() {
		return rt.ErrorMessage("状态错误")
	}
	return c.State(ctx, ids, state)
}

// LogicalDeletion 逻辑删除
func (c *BasicConfigService) LogicalDeletion(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
	if len(ids) < 1 {
		return rt.ErrorMessage("id错误")
	}
	finds, ok := c.sv.FindAllByIdStringIn(ctx, ids, optionsPg.WithCtx(ctx))
	if !ok {
		return rt.ErrorMessage("数据不存在")
	}
	if c.sv.Config().Data.Delete {
		if err := c.sv.DeleteByIdsString(ctx, ids); err != nil {
			return rt.ErrorMessage(err.Error())
		}
	} else {
		for _, info := range finds {
			enum := enumStatePg.State(info.State)
			if ok, reverse := enum.ReverseEnableDisable(); ok {
				c.sv.Update(ctx, entityBasic.BasicConfigEntity{State: reverse.IndexInt8()}, info.ID)
			}
		}
	}
	c.refreshByConfigs(ctx, finds)
	return rt.Ok()
}

// LogicalRecovery 逻辑删除恢复
func (c *BasicConfigService) LogicalRecovery(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
	if len(ids) < 1 {
		return rt.ErrorMessage("id错误")
	}
	finds, ok := c.sv.FindAllByIdStringIn(ctx, ids, optionsPg.WithCtx(ctx))
	if !ok {
		return rt.ErrorMessage("数据不存在")
	}
	for _, info := range finds {
		enum := enumStatePg.State(info.State)
		if ok, reverse := enum.ReverseCancelLayAside(); ok {
			c.sv.Update(ctx, entityBasic.BasicConfigEntity{State: reverse.IndexInt8()}, info.ID)
		}
	}
	c.refreshByConfigs(ctx, finds)
	return rt.Ok()
}

// PhysicalDeletion 物理删除
func (c *BasicConfigService) PhysicalDeletion(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
	if len(ids) < 1 {
		return rt.ErrorMessage("id错误")
	}
	finds, ok := c.sv.FindAllByIdStringIn(ctx, ids, optionsPg.WithCtx(ctx))
	if !ok {
		return rt.ErrorMessage("数据不存在")
	}
	idsInt := make([]int64, 0, len(finds))
	for _, info := range finds {
		idsInt = append(idsInt, info.ID)
	}
	if err := c.sv.DeleteByIds(ctx, idsInt); err != nil {
		return rt.ErrorMessage(err.Error())
	}
	c.refreshByConfigs(ctx, finds)
	return rt.Ok()
}

// ExistName 查重
func (c *BasicConfigService) ExistName(ctx *gin.Context, ct model.BaseExistWdCt[string]) (rt rg.Rs[string]) {
	if strPg.IsBlank(ct.Wd) {
		return rt.ErrorMessage("关键词不能为空")
	}
	id := "0"
	if strPg.IsNotBlank(ct.Id) {
		id = ct.Id
	}
	if _, exists := c.sv.FindByNameAndIdNot(ctx, ct.Wd, id); exists {
		return rt.ErrorMessage("重复，已存在")
	}
	return rt.OkMessage("可以使用")
}

// ---- helpers ----

func (c *BasicConfigService) findModel(ctx *gin.Context, tenantNo, modelNo, model string) (*entityBasic.BasicConfigModelEntity, bool) {
	if strPg.IsNotBlank(modelNo) {
		return c.repModel.FindByTenantNoAndNo(ctx, tenantNo, modelNo)
	}
	if strPg.IsNotBlank(model) {
		return c.repModel.FindByTenantNoAndModel(ctx, tenantNo, model)
	}
	return nil, false
}

func (c *BasicConfigService) refresh(ctx context.Context, tenantNo, model string) {
	if strPg.IsBlank(tenantNo) || strPg.IsBlank(model) {
		return
	}
	if err := c.sp.Refresh(ctx, tenantNo, model); err != nil {
		log.Error(ctx, log.TagAppDef, err, "refresh config cache error")
	}
}

func (c *BasicConfigService) refreshByConfigs(ctx context.Context, infos []*entityBasic.BasicConfigEntity) {
	pairs := make(map[string]bool, len(infos))
	for _, info := range infos {
		if strPg.IsBlank(info.TenantNo) || strPg.IsBlank(info.Model) {
			continue
		}
		pairs[info.TenantNo+"\x00"+info.Model] = true
	}
	for pair := range pairs {
		for i := 0; i < len(pair); i++ {
			if pair[i] == '\x00' {
				c.refresh(ctx, pair[:i], pair[i+1:])
				break
			}
		}
	}
}

func normalizeYesNo(value, fallback int8) int8 {
	if yesNoIntPg.Yes.IsEqual(value) || yesNoIntPg.No.IsEqual(value) {
		return value
	}
	if fallback == 0 {
		return yesNoIntPg.No.IndexInt8()
	}
	return fallback
}

func validateRules(value string, rules []*entityBasic.BasicModelRulesEntity, field *entityBasic.BasicConfigModelFieldsEntity) string {
	if len(rules) == 0 {
		return ""
	}
	fieldName := ""
	if field != nil {
		fieldName = field.Name
	}
	for _, rule := range rules {
		if rule == nil || rule.State != enumStatePg.ENABLE.IndexInt8() {
			continue
		}
		switch {
		case configModelPg.RuleModeRequired.IsEqual(rule.RuleMode):
			if strPg.IsBlank(value) {
				return ruleMessage(rule, fieldName, "不能为空")
			}
		case configModelPg.RuleModeLength.IsEqual(rule.RuleMode):
			min, max, exact := parseLengthCondition(rule.Condition)
			n := utf8.RuneCountInString(value)
			if exact > 0 && n != exact {
				return ruleMessage(rule, fieldName, fmt.Sprintf("长度必须为 %d", exact))
			}
			if min > 0 && n < min {
				return ruleMessage(rule, fieldName, fmt.Sprintf("长度不能小于 %d", min))
			}
			if max > 0 && n > max {
				return ruleMessage(rule, fieldName, fmt.Sprintf("长度不能大于 %d", max))
			}
		case configModelPg.RuleModePattern.IsEqual(rule.RuleMode) && strPg.IsNotBlank(rule.Condition):
			match, err := regexp.MatchString(strings.TrimSpace(rule.Condition), value)
			if err != nil || !match {
				return ruleMessage(rule, fieldName, "格式不正确")
			}
		case configModelPg.RuleModeGreaterThan0.IsEqual(rule.RuleMode):
			num, err := strconv.ParseFloat(value, 64)
			if err != nil || num <= 0 {
				return ruleMessage(rule, fieldName, "必须大于0")
			}
		}
	}
	return ""
}

func ruleMessage(rule *entityBasic.BasicModelRulesEntity, fieldName, msg string) string {
	if strPg.IsNotBlank(rule.ErrorMessage) {
		return rule.ErrorMessage
	}
	if strPg.IsBlank(fieldName) {
		fieldName = rule.Name
	}
	if strPg.IsBlank(fieldName) {
		fieldName = "字段"
	}
	return fieldName + " " + msg
}

func parseLengthCondition(condition string) (min, max, exact int) {
	condition = strings.TrimSpace(condition)
	if condition == "" {
		return
	}
	parts := strings.Split(condition, ",")
	if len(parts) == 1 {
		if v, err := strconv.Atoi(strings.TrimSpace(parts[0])); err == nil {
			exact = v
		}
		return
	}
	if len(parts) == 2 {
		min, _ = strconv.Atoi(strings.TrimSpace(parts[0]))
		max, _ = strconv.Atoi(strings.TrimSpace(parts[1]))
	}
	return
}

func typeInt8(v int8) typePg.Int8 {
	return typePg.Int8(v)
}
