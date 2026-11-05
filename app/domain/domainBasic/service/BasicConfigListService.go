package service

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/domainBasic/service/configBasic"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/basic/modBasicConfig"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/basic/modBasicConfigEvent"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/basic/modBasicConfigList"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/state/yesNoPg/yesNoIntPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityB72"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/repositoryB72"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/sdk/auth/holderPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg/optionsPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
	"github.com/pangu-2/go-tools/tools/noPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"

	"github.com/goccy/go-json"
	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/dbPg/pagePg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
)

func init() {
	gs.Provide(new(BasicConfigListService))
}

// BasicConfigListService 用户组
// @Description:
type BasicConfigListService struct {
	sv *repositoryB72.B72ConfigListRepository `autowire:"?"`

	Sp            *configBasic.Sp                               `autowire:"?"`
	repModel      *repositoryB72.B72ConfigModelRepository       `autowire:"?"`
	repEvent      *repositoryB72.B72ConfigEventRepository       `autowire:"?"`
	repEventField *repositoryB72.B72ConfigEventFieldsRepository `autowire:"?"`
	repConfig     *repositoryB72.B72ConfigRepository            `autowire:"?"`
	repConfigList *repositoryB72.B72ConfigListRepository        `autowire:"?"`
	repField      *repositoryB72.B72ConfigModelFieldsRepository `autowire:"?"`

	// configService 复用按模型保存配置值的能力（规则校验、binary 字段、缓存刷新）
	configService *BasicConfigService `autowire:"?"`
}

// CreateUpdate 创建/更新
//
//	@Description:
//	@receiver c
//	@param ct
//	@return rt
func (c *BasicConfigListService) CreateUpdate(ctx *gin.Context, ct modBasicConfigList.CreateUpdateCt) (rt rg.Rs[string]) {
	log.Info(ctx, log.TagAppDef, "ct", log.Any("ct", ct))
	var info entityB72.B72ConfigListEntity
	copier.Copy(&info, &ct)
	//
	if "" == ct.Name {
		return rt.ErrorMessage("名称不能为空")
	}
	if strPg.IsBlank(ct.EventNo) {
		return rt.ErrorMessage("事件不能为空")
	}
	if strPg.IsBlank(ct.Field) {
		return rt.ErrorMessage("字段名称不能为空")
	}
	holder := holderPg.GetContextAccount(ctx)
	event, result := c.repEvent.FindByNo(ctx, ct.EventNo)
	if !result {
		return rt.ErrorMessage("事件不存在")
	}
	id := "0"
	if ct.ID.ToInt64() > 0 {
		id = ct.ID.ToString()
	}
	_, b := c.sv.FindByTenantNoAndEventNoAndIdNot(ctx, holder.GetTenantNo(), event.No, id)
	if b {
		return rt.ErrorMessage("该事件在当前租户下已存在")
	}
	//
	r := c.sv
	if ct.ID < 1 {
		info.No = noPg.No()
		info.TenantNo = holder.GetTenantNo()
		info.Model = event.Model
		info.Module = event.Module
		info.ModuleSub = event.ModuleSub
		log.Info(ctx, log.TagAppDef, "update", log.Any("info", info))
		err, _ := r.Create(ctx, &info)
		if err != nil {
			log.Error(ctx, log.TagAppDef, err, "update error")
			return rt.ErrorMessage(err.Error())
		}
	} else {
		find, b := r.FindById(ctx, ct.ID.ToInt64())
		if !b {
			return rt.ErrorMessage("数据不存在")
		}
		info.ID = 0
		info.No = ""
		log.Info(ctx, log.TagAppDef, "update", log.Any("info", info))
		err := r.Update(ctx, info, find.ID)
		if err != nil {
			log.Error(ctx, log.TagAppDef, err, "update error")
			return rt.ErrorMessage(err.Error())
		}
	}
	return rt.Ok()
}

// Detail 详情
//
//	@Description:
//	@receiver c
//	@param id
func (c *BasicConfigListService) Detail(ctx *gin.Context, id int64) (rt rg.Rs[modBasicConfigList.Vo]) {
	if id < 1 {
		return rt.ErrorMessage("id错误")
	}
	find, b := c.sv.FindById(ctx, id)
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	var info modBasicConfigList.Vo
	copier.Copy(&info, &find)
	return rt.OkData(info)
}

// Enable 启用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicConfigListService) Enable(ctx *gin.Context, ct model.BaseIdsCt[string]) (rt rg.Rs[string]) {
	log.Info(ctx, log.TagAppDef, "ct", log.Any("ct", ct))
	return c.State(ctx, ct.Ids, enumStatePg.ENABLE)
}

// Disable 禁用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicConfigListService) Disable(ctx *gin.Context, ct model.BaseIdsCt[string]) (rt rg.Rs[string]) {
	log.Info(ctx, log.TagAppDef, "ct", log.Any("ct", ct))
	return c.State(ctx, ct.Ids, enumStatePg.GetType(enumStatePg.DISABLE))
}

// State 状态 启用/禁用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicConfigListService) State(ctx *gin.Context, ids []string, state enumStatePg.State) (rt rg.Rs[string]) {
	if len(ids) < 1 {
		return rt.ErrorMessage("id错误")
	}
	r := c.sv
	finds, b := r.FindAllByIdStringIn(ctx, ids, optionsPg.WithCtx(ctx))
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	for _, info := range finds {
		if info.State != state.IndexInt8() {
			r.Update(ctx, entityB72.B72ConfigListEntity{State: state.IndexInt8()}, info.ID)
		}
	}
	return rt.Ok()
}

// StateEnableDisable 状态 设置 有效 停用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicConfigListService) StateEnableDisable(ctx *gin.Context, ids []string, state enumStatePg.State) (rt rg.Rs[string]) {
	if !state.IsEnableDisable() {
		return rt.ErrorMessage("状态错误")
	}
	return c.State(ctx, ids, state)
}

// LogicalDeletion 逻辑删除
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicConfigListService) LogicalDeletion(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
	log.Info(ctx, log.TagAppDef, "ct", log.Any("ids", ids))
	if len(ids) < 1 {
		return rt.ErrorMessage("id错误")
	}
	repository := c.sv
	finds, b := repository.FindAllByIdStringIn(ctx, ids, optionsPg.WithCtx(ctx))
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	if c.sv.Config().Data.Delete {
		for _, info := range finds {
			log.Info(ctx, log.TagAppDef, "id=%v,TenantId=%v", log.Any("id", info.ID), log.Any("TenantId", info.TenantNo))
		}
		repository.DeleteByIdsString(ctx, ids)
	} else {
		for _, info := range finds {
			enum := enumStatePg.State(info.State)
			// 有效 停用，反转 为对应的 取消 弃置
			if ok, reverse := enum.ReverseEnableDisable(); ok {
				repository.Update(ctx, entityB72.B72ConfigListEntity{State: reverse.IndexInt8()}, info.ID)
			}
		}
	}
	return rt.Ok()
}

// LogicalRecovery 逻辑删除恢复
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicConfigListService) LogicalRecovery(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
	log.Info(ctx, log.TagAppDef, "ct", log.Any("ids", ids))
	if len(ids) < 1 {
		return rt.ErrorMessage("id错误")
	}
	repository := c.sv
	finds, b := repository.FindAllByIdStringIn(ctx, ids, optionsPg.WithCtx(ctx))
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	for _, info := range finds {
		enum := enumStatePg.State(info.State)
		//  取消 弃置 批量删除，反转 为对应的 有效 停用 停用
		if ok, reverse := enum.ReverseCancelLayAside(); ok {
			repository.Update(ctx, entityB72.B72ConfigListEntity{State: reverse.IndexInt8()}, info.ID)
		}
	}
	return rt.Ok()
}

// PhysicalDeletion 物理删除
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicConfigListService) PhysicalDeletion(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
	log.Info(ctx, log.TagAppDef, "ct", log.Any("ids", ids))
	if len(ids) < 1 {
		return rt.ErrorMessage("id错误")
	}
	cn := c.sv
	finds, b := cn.FindAllByIdStringIn(ctx, ids, optionsPg.WithCtx(ctx))
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	idsNew := make([]int64, 0)
	for _, info := range finds {
		log.Info(ctx, log.TagAppDef, "id=%v,TenantId=%v", log.Any("id", info.ID), log.Any("TenantId", info.TenantNo))
		idsNew = append(idsNew, info.ID)
	}
	if len(idsNew) > 0 {
		cn.DeleteByIds(ctx, idsNew)
	}
	return rt.Ok()
}

// Query 查询
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicConfigListService) Query(ctx *gin.Context, ct modBasicConfigList.QueryCt) (rt rg.Rs[pagePg.Paginator[modBasicConfigList.Vo]]) {
	log.Info(ctx, log.TagAppDef, "ct", log.Any("ct", ct))
	var query entityB72.B72ConfigListEntity
	copier.Copy(&query, &ct)
	slice := make([]modBasicConfigList.Vo, 0)
	rt.Data.Data = slice
	r := c.sv
	holder := holderPg.GetContextAccount(ctx)
	page, err := r.FindAllPage(ctx, query, optionsPg.WithOption(func(arg *optionsPg.OptionParams) {
		if ct.PageSize < 1 {
			ct.PageSize = 20
		}
		arg.Pageable = new(pagePg.PageablePageSize(0, ct.PageNum, ct.PageSize))
		//排序
		arg.Db = arg.Db.Order("create_at asc")
		arg.Db = arg.Db.Where("tenant_no=?", holder.GetTenantNo())
		if strPg.IsNotBlank(ct.Wd) {
			arg.Db = arg.Db.Where("name like ?", "%"+ct.Wd+"%")
		}
	}), optionsPg.WithCtx(ctx))
	if nil != err {
		return rt.Ok()
	}

	if page.Total > 0 && page.Data != nil && len(page.Data) > 0 {

		pg := pagePg.NewPaginatorByPageable[modBasicConfigList.Vo](page.Pageable)
		//字段赋值
		for _, item := range page.Data {
			var vo modBasicConfigList.Vo
			copier.Copy(&vo, &item)
			slice = append(slice, vo)
		}
		pg.Data = slice
		pg.Pageable = page.Pageable
		rt.Data = pg
		return rt.Ok()
	}
	return rt.Ok()
}

// SelectNodeAllPublic 查询
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicConfigListService) SelectNodeAllPublic(ctx *gin.Context, ct modBasicConfigList.QueryPublicCt) (rt rg.Rs[[]model.BaseNodeNo]) {
	log.Info(ctx, log.TagAppDef, "ct", log.Any("ct", ct))
	var query entityB72.B72ConfigListEntity
	copier.Copy(&query, &ct)
	slice := make([]model.BaseNodeNo, 0)
	rt.Data = slice
	infos := c.sv.FindAll(ctx, query)
	if len(infos) > 0 {
		for _, item := range infos {
			var vo modBasicConfigList.Vo
			copier.Copy(&vo, &item)
			code := model.BaseNodeNo{
				Value:  item.No,
				Label:  item.Name,
				Extend: vo,
			}

			slice = append(slice, code)
		}
		rt.Data = slice
	}
	return rt.Ok()
}

// ExistName 查重
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicConfigListService) ExistName(ctx *gin.Context, ct model.BaseExistWdCt[string]) (rt rg.Rs[string]) {
	log.Info(ctx, log.TagAppDef, "ct", log.Any("ct", ct))
	if "" == ct.Wd {
		return rt.ErrorMessage("查询内容不能为空")
	}
	id := "0"
	if strPg.IsNotBlank(ct.Id) {
		id = ct.Id
	}
	_, result := c.sv.FindByNameAndIdNot(ctx, ct.Wd, id, optionsPg.WithCtx(ctx))
	if result {
		return rt.ErrorMessage("重复，已存在")
	}
	return rt.OkMessage("可以使用")
}

// ExistCode 查重
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicConfigListService) ExistCode(ctx *gin.Context, ct model.BaseExistWdCt[string]) (rt rg.Rs[string]) {
	log.Info(ctx, log.TagAppDef, "ct", log.Any("ct", ct))
	if "" == ct.Wd {
		return rt.ErrorMessage("查询内容不能为空")
	}
	id := "0"
	if strPg.IsNotBlank(ct.Id) {
		id = ct.Id
	}
	_, result := c.sv.FindByCodeAndIdNot(ctx, ct.Wd, id, optionsPg.WithCtx(ctx))
	if result {
		return rt.ErrorMessage("重复，已存在")
	}
	return rt.OkMessage("可以使用")
}

// DetailForm 配置表单：返回配置模型的表单字段（含规则/表单属性）与该模型当前的配置值。
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicConfigListService) DetailForm(ctx *gin.Context, ct modBasicConfigList.DetailFormCt) (rt rg.Rs[modBasicConfigList.DetailFormVo]) {
	log.Info(ctx, log.TagAppDef, "ct", log.Any("ct", ct))
	vo := modBasicConfigList.DetailFormVo{
		Form: modBasicConfigEvent.ModelForm{Item: make([]modBasicConfigEvent.ModelFormItem, 0)},
		Data: make(map[string]interface{}),
	}
	rt.Data = vo
	if strPg.IsBlank(ct.EventNo) {
		return rt.ErrorMessage("参数错误")
	}
	holder := holderPg.GetContextAccount(ctx)
	config, ok := c.repConfigList.FindByEventNo(ctx, ct.EventNo)
	if !ok {
		return rt.ErrorMessage("配置不存在")
	}
	if holder.GetTenantNo() != config.TenantNo {
		return rt.ErrorMessage("配置不存在")
	}
	// 配置模型：字段来自模型设计器，值来自模型配置缓存
	if modelInfo, okModel := c.findModelByConfigList(ctx, config); okModel {
		if c.repField != nil {
			if fields, okField := c.repField.FindAllByModelNo(ctx, modelInfo.No); okField {
				items := make([]modBasicConfigEvent.ModelFormItem, 0, len(fields))
				for _, field := range fields {
					items = append(items, newModelFormItemFromModelField(field))
				}
				sort.SliceStable(items, func(i, j int) bool {
					return items[i].Sort < items[j].Sort
				})
				vo.Form.Item = items
			}
		}
		if c.configService != nil {
			if items, okValue := c.configService.GetByTenantNoAndModel(ctx, config.TenantNo, modelInfo.Model); okValue {
				for _, item := range items {
					vo.Data[item.Field] = item.ValueMap()
				}
			}
		}
		rt.Data = vo
		return rt.Ok()
	}
	// 兼容事件维度：字段来自事件字段，值来自 event_no 关联的配置行
	if c.repEventField != nil {
		if fields, okField := c.repEventField.FindAllByEventNo(ctx, config.EventNo); okField {
			items := make([]modBasicConfigEvent.ModelFormItem, 0, len(fields))
			for _, field := range fields {
				items = append(items, newModelFormItemFromEventField(field))
			}
			sort.SliceStable(items, func(i, j int) bool {
				return items[i].Sort < items[j].Sort
			})
			vo.Form.Item = items
		}
	}
	if infos, okValue := c.repConfig.FindByEventNo(ctx, config.EventNo); okValue {
		for _, item := range infos {
			if item.TenantNo != config.TenantNo {
				continue
			}
			if yesNoIntPg.Yes.IsEqual(item.Binary) {
				vo.Data[item.Field] = item.ValueBinary
			} else {
				vo.Data[item.Field] = item.Value
			}
		}
	}
	rt.Data = vo
	return rt.Ok()
}

// ConfigUpdate 保存配置：按模型字段规则校验后回写配置值（binary 字段写 value_binary 并刷新缓存），
// 校验通过后把本次提交的结构体写入配置列表的 content 字段。
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicConfigListService) ConfigUpdate(ctx *gin.Context, ct modBasicConfigList.ConfigUpdateCt) (rt rg.Rs[string]) {
	log.Info(ctx, log.TagAppDef, "ct", log.Any("ct", ct))
	if strPg.IsBlank(ct.EventNo) {
		return rt.ErrorMessage("配置编号不能为空")
	}
	if len(ct.Data) < 1 {
		return rt.ErrorMessage("配置数据不能为空")
	}
	holder := holderPg.GetContextAccount(ctx)
	config, ok := c.repConfigList.FindByEventNo(ctx, ct.EventNo)
	if !ok {
		return rt.ErrorMessage("配置不存在")
	}
	if holder.GetTenantNo() != config.TenantNo {
		return rt.ErrorMessage("配置不存在")
	}
	modelInfo, okModel := c.findModelByConfigList(ctx, config)
	if !okModel {
		return rt.ErrorMessage("配置关联的模型不存在")
	}
	if c.configService == nil {
		return rt.ErrorMessage("配置服务不可用")
	}
	// 统一成字符串，避免后端按类型断言失败
	data := make(map[string]string, len(ct.Data))
	for key, value := range ct.Data {
		key = strings.TrimSpace(key)
		if strPg.IsBlank(key) {
			continue
		}
		switch v := value.(type) {
		case string:
			data[key] = v
		case nil:
			data[key] = ""
		default:
			data[key] = fmt.Sprint(v)
		}
	}
	if len(data) < 1 {
		return rt.ErrorMessage("配置数据不能为空")
	}
	// 1. 规则校验 + 回写配置值（binary 字段、配置缓存）
	res := c.configService.UpdateValues(ctx, modBasicConfig.UpdateValuesCt{
		ModelNo: modelInfo.No,
		Model:   modelInfo.Model,
		Data:    data,
	})
	if !res.SuccessIs() {
		return rt.ErrorMessage(res.Message)
	}
	// 2. 校验通过后，把本次提交的结构体写入配置列表 content
	body, err := json.Marshal(data)
	if err != nil {
		return rt.ErrorMessage("配置内容序列化失败：" + err.Error())
	}
	if err := c.repConfigList.Update(ctx, entityB72.B72ConfigListEntity{Content: string(body)}, config.ID); err != nil {
		return rt.ErrorMessage("保存配置内容失败：" + err.Error())
	}
	return rt.Ok()
}

// findModelByConfigList 解析配置列表行关联的模型：优先 model_no，其次模型标识 model。
func (c *BasicConfigListService) findModelByConfigList(ctx *gin.Context, config *entityB72.B72ConfigListEntity) (*entityB72.B72ConfigModelEntity, bool) {
	if c.repModel == nil || config == nil {
		return nil, false
	}
	if strPg.IsNotBlank(config.ModelNo) {
		if info, ok := c.repModel.FindByTenantNoAndNo(ctx, config.TenantNo, config.ModelNo); ok {
			return info, true
		}
	}
	if strPg.IsNotBlank(config.Model) {
		if info, ok := c.repModel.FindByTenantNoAndModel(ctx, config.TenantNo, config.Model); ok {
			return info, true
		}
	}
	return nil, false
}

// newModelFormItemFromModelField 模型字段 -> 表单字段
func newModelFormItemFromModelField(field *entityB72.B72ConfigModelFieldsEntity) modBasicConfigEvent.ModelFormItem {
	return modBasicConfigEvent.ModelFormItem{
		Name:             field.Name,
		Field:            field.Field,
		Show:             typePg.Int8(field.Show),
		Binary:           typePg.Int8(field.Binary),
		DefaultValue:     field.DefaultValue,
		ValueType:        field.ValueType,
		FormCode:         field.FormCode,
		ParameterSource:  field.ParameterSource,
		ParameterContent: field.ParameterContent,
		FormAttr:         string(field.FormAttr),
		FormGroup:        field.FormGroup,
		Sort:             field.Sort,
		Rules:            field.Rules.Data(),
	}
}

// newModelFormItemFromEventField 事件字段 -> 表单字段
func newModelFormItemFromEventField(field *entityB72.B72ConfigEventFieldsEntity) modBasicConfigEvent.ModelFormItem {
	return modBasicConfigEvent.ModelFormItem{
		Name:             field.Name,
		Field:            field.Field,
		Show:             typePg.Int8(field.Show),
		Binary:           typePg.Int8(field.Binary),
		DefaultValue:     field.DefaultValue,
		ValueType:        field.ValueType,
		FormCode:         field.FormCode,
		ParameterSource:  field.ParameterSource,
		ParameterContent: field.ParameterContent,
		FormAttr:         string(field.FormAttr),
		Sort:             field.Sort,
		Rules:            field.Rules.Data(),
	}
}
