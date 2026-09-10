package service

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/manage/domainBasic/service/configEvent"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/basic/modBasicConfigEvent"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBasic"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/holderPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg/optionsPg"
	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/dbPg/pagePg"
	"github.com/pangu-2/go-tools/tools/noPg"
	"github.com/pangu-2/go-tools/tools/numberPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
	"gorm.io/gorm"
)

func init() {
	gs.Provide(new(BasicConfigEventService))
}

type BasicConfigEventService struct {
	sv        *repositoryBasic.BasicConfigEventRepository       `autowire:"?"`
	repoField *repositoryBasic.BasicConfigEventFieldsRepository `autowire:"?"`
	sp        *configEvent.Sp                                   `autowire:"?"`
}

// Create 创建事件
func (s *BasicConfigEventService) Create(ctx *gin.Context, ct modBasicConfigEvent.CreateCt) (rt rg.Rs[string]) {
	if strPg.IsBlank(ct.Name) {
		return rt.ErrorMessage("中文名称不能为空")
	}
	if strPg.IsBlank(ct.Field) {
		return rt.ErrorMessage("英文标识不能为空")
	}
	if strPg.IsBlank(ct.Model) {
		return rt.ErrorMessage("模型英文标识不能为空")
	}

	// 检查事件是否已存在 (EventNo = Field)
	var existEvent entityBasic.BasicConfigEventEntity
	// Field 对应 EventNo 或者 Field 本身，这里假设 BasicConfigEventEntity.Field 存储事件的英文标识
	// 根据用户描述：Field 英文标识
	tx := s.sv.Db().Where("field = ?", ct.Field).First(&existEvent)
	if tx.RowsAffected > 0 {
		return rt.ErrorMessage("事件标识已存在")
	}

	holder := holderPg.GetContextAccount(ctx)

	// 1. 创建事件头
	eventEntity := entityBasic.BasicConfigEventEntity{
		Name:      ct.Name,
		Field:     ct.Field, // 事件标识
		Model:     ct.Model, // 关联的模型
		ModelNo:   ct.Model, // 关联的模型编号
		ModuleSub: ct.ModuleSub,
		// TypeSys:      ct.TypeSys, // Entity中没有该字段
		Description: ct.Description,
		No:          noPg.No(), // 生成唯一编号
		TenantNo:    holder.GetTenantNo(),
		CreateBy:    holder.GetAccount().Account,
	}

	// Entity 中没有 TypeCategory, 只有 TypeSys 和 ModuleSub
	// KindUnique = ModelNo + Field (EventNo)
	eventEntity.KindUnique = eventEntity.ModelNo + eventEntity.Field

	err, _ := s.sv.Create(ctx, &eventEntity)
	if err != nil {
		log.Errorf(ctx, log.TagAppDef, "Create event error: %v", err)
		return rt.ErrorMessage("创建事件失败")
	}

	// 2. 创建事件字段
	for _, fieldCt := range ct.Fields {
		fieldEntity := entityBasic.BasicConfigEventFieldsEntity{
			EventNo:         eventEntity.Field,
			ModelNo:         eventEntity.ModelNo,
			Name:            fieldCt.Name,
			Field:           fieldCt.Field,
			Show:            fieldCt.Show,
			Binary:          fieldCt.Binary,
			DefaultValue:    fieldCt.DefaultValue,
			ValueType:       fieldCt.ValueType,
			FormCode:        fieldCt.FormCode,
			ParameterSource: fieldCt.ParameterSource,
			TenantNo:        holder.GetTenantNo(),
			CreateBy:        holder.GetAccount().Account,
		}
		// FieldKindUnique = EventNo + Field
		fieldEntity.KindUnique = fieldEntity.EventNo + fieldEntity.Field

		err, _ := s.repoField.Create(ctx, &fieldEntity)
		if err != nil {
			log.Errorf(ctx, log.TagAppDef, "Create event field error: %v", err)
			continue
		}
	}

	return rg.OkData(numberPg.Int64ToString(eventEntity.ID))
}

func (c *BasicConfigEventService) CreateUpdate(ctx *gin.Context, ct modBasicConfigEvent.CreateUpdateCt) (rt rg.Rs[string]) {
	return configEvent.NewCreateUpdate(c.sp, ct, true).Process(ctx)
}

// Detail 详情
func (c *BasicConfigEventService) Detail(ctx *gin.Context, id string) (rt rg.Rs[modBasicConfigEvent.CreateUpdateCt]) {
	return configEvent.NewDetail(c.sp).Process(ctx, id)
}

// Enable 启用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicConfigEventService) Enable(ctx *gin.Context, ct model.BaseIdsCt[string]) (rt rg.Rs[string]) {
	return c.State(ctx, ct.Ids, enumStatePg.ENABLE)
}

// Disable 禁用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicConfigEventService) Disable(ctx *gin.Context, ct model.BaseIdsCt[string]) (rt rg.Rs[string]) {
	return c.State(ctx, ct.Ids, enumStatePg.GetType(enumStatePg.DISABLE))
}

// State 状态 启用/禁用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicConfigEventService) State(ctx *gin.Context, ids []string, state enumStatePg.State) (rt rg.Rs[string]) {
	if len(ids) < 1 {
		return rt.ErrorMessage("id错误")
	}
	r := c.sv
	finds, b := r.FindAllByIdStringIn(ctx, ids)
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	for _, info := range finds {
		if info.State != state.IndexInt8() {
			r.Update(ctx, entityBasic.BasicConfigEventEntity{State: state.IndexInt8()}, info.ID)
		}
	}
	return rt.Ok()
}

// StateEnableDisable 状态 设置 有效 停用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicConfigEventService) StateEnableDisable(ctx *gin.Context, ids []string, state enumStatePg.State) (rt rg.Rs[string]) {
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
func (c *BasicConfigEventService) LogicalDeletion(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
	if len(ids) < 1 {
		return rt.ErrorMessage("id错误")
	}
	repository := c.sv
	finds, b := repository.FindAllByIdStringIn(ctx, ids)
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	if c.sv.Config().Data.Delete {
		for _, info := range finds {
			log.Infof(ctx, log.TagAppDef, "id=%v", info.ID)
		}
		repository.DeleteByIdsString(ctx, ids)
	} else {
		for _, info := range finds {
			enum := enumStatePg.State(info.State)
			// 有效 停用，反转 为对应的 取消 弃置
			if ok, reverse := enum.ReverseEnableDisable(); ok {
				repository.Update(ctx, entityBasic.BasicConfigEventEntity{State: reverse.IndexInt8()}, info.ID)
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
func (c *BasicConfigEventService) LogicalRecovery(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ids)
	if len(ids) < 1 {
		return rt.ErrorMessage("id错误")
	}
	repository := c.sv
	finds, b := repository.FindAllByIdStringIn(ctx, ids)
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	for _, info := range finds {
		enum := enumStatePg.State(info.State)
		//  取消 弃置 批量删除，反转 为对应的 有效 停用 停用
		if ok, reverse := enum.ReverseCancelLayAside(); ok {
			repository.Update(ctx, entityBasic.BasicConfigEventEntity{State: reverse.IndexInt8()}, info.ID)
		}
	}
	return rt.Ok()
}

// PhysicalDeletion 物理删除
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicConfigEventService) PhysicalDeletion(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ids)
	if len(ids) < 1 {
		return rt.ErrorMessage("id错误")
	}
	cn := c.sv
	finds, b := cn.FindAllByIdStringIn(ctx, ids)
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	idsNew := make([]int64, 0)
	for _, info := range finds {
		log.Infof(ctx, log.TagAppDef, "id=%v", info.ID)
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
func (c *BasicConfigEventService) Query(ctx *gin.Context, ct modBasicConfigEvent.QueryCt) (rt rg.Rs[pagePg.Paginator[modBasicConfigEvent.Vo]]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	var query entityBasic.BasicConfigEventEntity
	copier.Copy(&query, &ct)
	slice := make([]modBasicConfigEvent.Vo, 0)
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
			arg.Db = arg.Db.Where("name like ?", "%"+ct.Wd+"%").Or("field like ?", "%"+ct.Wd+"%").Or("description like ?", "%"+ct.Wd+"%")
		}
	}), optionsPg.WithCtx(ctx))
	if nil != err {
		return rt.Ok()
	}

	if page.Total > 0 && page.Data != nil && len(page.Data) > 0 {
		pg := pagePg.NewPaginatorByPageable[modBasicConfigEvent.Vo](page.Pageable)
		//字段赋值
		for _, item := range page.Data {
			var vo modBasicConfigEvent.Vo
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
func (c *BasicConfigEventService) SelectNodeAllPublic(ctx *gin.Context, ct modBasicConfigEvent.QueryPublicCt) (rt rg.Rs[[]model.BaseNodeNo]) {
	var query entityBasic.BasicConfigEventEntity
	copier.Copy(&query, &ct)
	slice := make([]model.BaseNodeNo, 0)
	rt.Data = slice
	holder := holderPg.GetContextAccount(ctx)
	infos := c.sv.FindAll(ctx, query, optionsPg.WithCondition(func(db *gorm.DB) *gorm.DB {
		return db.Where("tenant_no=?", holder.GetTenantNo())
	}))
	if len(infos) > 0 {
		for _, item := range infos {
			var vo modBasicConfigEvent.Vo
			copier.Copy(&vo, &item)
			//
			code := model.BaseNodeNo{
				Value:    item.No,
				Label:    item.Name,
				ParentNo: "",
				Extend:   vo,
			}
			slice = append(slice, code)
		}
		rt.Data = slice
	}
	return rt.Ok()
}

// AllByModel 查询
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicConfigEventService) AllByModel(ctx *gin.Context, ct modBasicConfigEvent.QueryCt) (rt rg.Rs[[]modBasicConfigEvent.Vo]) {
	var query entityBasic.BasicConfigEventEntity
	copier.Copy(&query, &ct)
	//
	slice := make([]modBasicConfigEvent.Vo, 0)
	rt.Data = slice
	holder := holderPg.GetContextAccount(ctx)
	//
	if strPg.IsBlank(ct.ModelNo) {
		return rt.ErrorMessage("模型编号错误")
	}
	infos := c.sv.FindAll(ctx, query, optionsPg.WithCondition(func(db *gorm.DB) *gorm.DB {
		return db.Where("tenant_no=?", holder.GetTenantNo())
	}))
	if len(infos) > 0 {
		for _, item := range infos {
			var vo modBasicConfigEvent.Vo
			copier.Copy(&vo, &item)
			slice = append(slice, vo)
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
func (c *BasicConfigEventService) ExistName(ctx *gin.Context, ct model.BaseExistWdCt[string]) (rt rg.Rs[string]) {
	if "" == ct.Wd {
		return rt.ErrorMessage("关键词不能为空")
	}
	id := "0"
	if strPg.IsNotBlank(ct.Id) {
		id = ct.Id
	}
	_, result := c.sv.FindByNameAndIdNot(ctx, ct.Wd, id)
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
func (c *BasicConfigEventService) ExistCode(ctx *gin.Context, ct model.BaseExistWdCt[string]) (rt rg.Rs[string]) {
	if "" == ct.Wd {
		return rt.ErrorMessage("关键词不能为空")
	}
	id := "0"
	if strPg.IsNotBlank(ct.Id) {
		id = ct.Id
	}
	_, result := c.sv.FindByCodeAndIdNot(ctx, ct.Wd, id)
	if result {
		return rt.ErrorMessage("重复，已存在")
	}
	return rt.OkMessage("可以使用")
}
