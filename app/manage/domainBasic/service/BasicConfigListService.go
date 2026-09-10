package service

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/manage/domainBasic/service/configBasic"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/basic/modBasicConfigList"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBasic"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/holderPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg/optionsPg"
	"github.com/pangu-2/go-tools/tools/noPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"

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
	sv *repositoryBasic.BasicConfigListRepository `autowire:"?"`

	Sp            *configBasic.Sp                                   `autowire:"?"`
	repModel      *repositoryBasic.BasicConfigModelRepository       `autowire:"?"`
	repEvent      *repositoryBasic.BasicConfigEventRepository       `autowire:"?"`
	repEventField *repositoryBasic.BasicConfigEventFieldsRepository `autowire:"?"`
	repConfig     *repositoryBasic.BasicConfigRepository            `autowire:"?"`
	repConfigList *repositoryBasic.BasicConfigListRepository        `autowire:"?"`
}

// CreateUpdate 创建/更新
//
//	@Description:
//	@receiver c
//	@param ct
//	@return rt
func (c *BasicConfigListService) CreateUpdate(ctx *gin.Context, ct modBasicConfigList.CreateUpdateCt) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	var info entityBasic.BasicConfigListEntity
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
		log.Infof(ctx, log.TagAppDef, "info.save=%+v", info)
		err, _ := r.Create(ctx, &info)
		if err != nil {
			log.Errorf(ctx, log.TagAppDef, "update error=%+v", err)
			return rt.ErrorMessage(err.Error())
		}
	} else {
		find, b := r.FindById(ctx, ct.ID.ToInt64())
		if !b {
			return rt.ErrorMessage("数据不存在")
		}
		info.ID = 0
		info.No = ""
		log.Infof(ctx, log.TagAppDef, "info.save=%+v", info)
		err := r.Update(ctx, info, find.ID)
		if err != nil {
			log.Errorf(ctx, log.TagAppDef, "update error=%+v", err)
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
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	return c.State(ctx, ct.Ids, enumStatePg.ENABLE)
}

// Disable 禁用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicConfigListService) Disable(ctx *gin.Context, ct model.BaseIdsCt[string]) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
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
			r.Update(ctx, entityBasic.BasicConfigListEntity{State: state.IndexInt8()}, info.ID)
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
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ids)
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
			log.Infof(ctx, log.TagAppDef, "id=%v,TenantId=%v", info.ID, info.TenantNo)
		}
		repository.DeleteByIdsString(ctx, ids)
	} else {
		for _, info := range finds {
			enum := enumStatePg.State(info.State)
			// 有效 停用，反转 为对应的 取消 弃置
			if ok, reverse := enum.ReverseEnableDisable(); ok {
				repository.Update(ctx, entityBasic.BasicConfigListEntity{State: reverse.IndexInt8()}, info.ID)
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
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ids)
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
			repository.Update(ctx, entityBasic.BasicConfigListEntity{State: reverse.IndexInt8()}, info.ID)
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
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ids)
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
		log.Infof(ctx, log.TagAppDef, "id=%v,TenantId=%v", info.ID, info.TenantNo)
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
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	var query entityBasic.BasicConfigListEntity
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
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	var query entityBasic.BasicConfigListEntity
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
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
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
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
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

// DetailForm 查重
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicConfigListService) DetailForm(ctx *gin.Context, ct modBasicConfigList.DetailFormCt) (rt rg.Rs[modBasicConfigList.DetailFormVo]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	return configBasic.NewDetailForm(c.Sp).Process(ctx, ct)
}

// ConfigUpdate 查重
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicConfigListService) ConfigUpdate(ctx *gin.Context, ct modBasicConfigList.ConfigUpdateCt) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	return configBasic.NewConfigUpdate(c.Sp).Process(ctx, ct)
}
