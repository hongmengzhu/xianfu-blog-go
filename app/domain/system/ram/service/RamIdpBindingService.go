package service

import (
	"github.com/gin-gonic/gin"
	modRamIdpBinding2 "github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamIdpBinding"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityRam"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryRam"
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
	gs.Provide(new(RamIdpBindingService))
}

// RamIdpBindingService 身份绑定
// @Description:
type RamIdpBindingService struct {
	sv *repositoryRam.RamIdpBindingRepository `autowire:"?"`
}

// CreateUpdate 新增更新
//
//	@Description:
//	@receiver c
//	@param ct
//	@return rt
func (c *RamIdpBindingService) CreateUpdate(ctx *gin.Context, ct modRamIdpBinding2.CreateUpdateCt) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	//
	holder := holderPg.GetContextAccount(ctx)
	//
	var info entityRam.RamIdpBindingEntity
	copier.Copy(&info, &ct)
	//
	r := c.sv
	//是否是更新
	isUpdate := false
	b := false
	//
	find := &entityRam.RamIdpBindingEntity{}
	//
	if ct.ID.ToInt64() > 0 {
		isUpdate = true
		///
		find, b = r.FindById(ctx, ct.ID.ToInt64())
		if !b {
			return rt.ErrorMessage("数据不存在")
		}
	}
	if strPg.IsBlank(info.No) {
		info.No = noPg.No()
	}
	if isUpdate {
		info.ID = 0
	} else {
		info.TenantNo = holder.GetTenantNo()
		info.State = enumStatePg.ENABLE.Index()
	}

	log.Infof(ctx, log.TagAppDef, "info.save=%+v", info)
	if isUpdate {
		err := r.Update(ctx, info, find.ID)
		if err != nil {
			log.Errorf(ctx, log.TagAppDef, "update error=%+v", err)
			return rt.ErrorMessage(err.Error())
		}
	} else {
		err, _ := r.Create(ctx, &info)
		if err != nil {
			return rt.ErrorMessage("保存失败 " + err.Error())
		}
		log.Infof(ctx, log.TagAppDef, "save=%+v", info)
	}

	return rt.Ok()
}

// Detail 详情
//
//	@Description:
//	@receiver c
//	@param id
func (c *RamIdpBindingService) Detail(ctx *gin.Context, id int64) (rt rg.Rs[modRamIdpBinding2.Vo]) {
	if id < 1 {
		return rt.ErrorMessage("id错误")
	}
	find, b := c.sv.FindById(ctx, id)
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	var info modRamIdpBinding2.Vo
	copier.Copy(&info, &find)
	return rt.OkData(info)
}

// Enable 启用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *RamIdpBindingService) Enable(ctx *gin.Context, ct model.BaseIdsCt[string]) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	return c.State(ctx, ct.Ids, enumStatePg.ENABLE)
}

// Disable 禁用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *RamIdpBindingService) Disable(ctx *gin.Context, ct model.BaseIdsCt[string]) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	return c.State(ctx, ct.Ids, enumStatePg.GetType(enumStatePg.DISABLE))
}

// State 状态 启用/禁用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *RamIdpBindingService) State(ctx *gin.Context, ids []string, state enumStatePg.State) (rt rg.Rs[string]) {
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
			r.Update(ctx, entityRam.RamIdpBindingEntity{State: state.IndexInt8()}, info.ID)
		}
	}
	return rt.Ok()
}

// LogicalDeletion 逻辑删除
//
//	@Description:
//	@receiver c
//	@param ct
func (c *RamIdpBindingService) LogicalDeletion(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
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
				repository.Update(ctx, entityRam.RamIdpBindingEntity{State: reverse.IndexInt8()}, info.ID)
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
func (c *RamIdpBindingService) LogicalRecovery(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
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
			repository.Update(ctx, entityRam.RamIdpBindingEntity{State: reverse.IndexInt8()}, info.ID)
		}
	}
	return rt.Ok()
}

// PhysicalDeletion 物理删除
//
//	@Description:
//	@receiver c
//	@param ct
func (c *RamIdpBindingService) PhysicalDeletion(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
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
func (c *RamIdpBindingService) Query(ctx *gin.Context, ct modRamIdpBinding2.QueryCt) (rt rg.Rs[pagePg.Paginator[modRamIdpBinding2.Vo]]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	var query entityRam.RamIdpBindingEntity
	copier.Copy(&query, &ct)
	slice := make([]modRamIdpBinding2.Vo, 0)
	rt.Data.Data = slice
	r := c.sv
	page, err := r.FindAllPage(ctx, query, optionsPg.WithOption(func(arg *optionsPg.OptionParams) {
		if ct.PageSize < 1 {
			ct.PageSize = 20
		}
		arg.Pageable = new(pagePg.PageablePageSize(0, ct.PageNum, ct.PageSize))
		arg.Db = arg.Db.Order("create_at desc")
		//自定义查询
		if strPg.IsNotBlank(ct.Wd) {
			arg.Db = arg.Db.Where("description like ?", "%"+ct.Wd+"%")
		}
	}), optionsPg.WithCtx(ctx))
	if nil != err {
		return rt.Ok()
	}

	if page.Total > 0 && page.Data != nil && len(page.Data) > 0 {

		pg := pagePg.NewPaginatorByPageable[modRamIdpBinding2.Vo](page.Pageable)
		//字段赋值
		for _, item := range page.Data {
			var vo modRamIdpBinding2.Vo
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

// SelectNodeAll 查询
//
//	@Description:
//	@receiver c
//	@param ct
func (c *RamIdpBindingService) SelectNodeAll(ctx *gin.Context, ct modRamIdpBinding2.QueryPublicCt) (rt rg.Rs[[]model.BaseNodeNo]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	var query entityRam.RamIdpBindingEntity
	copier.Copy(&query, &ct)
	slice := make([]model.BaseNodeNo, 0)
	rt.Data = slice
	infos := c.sv.FindAll(ctx, query)
	if len(infos) > 0 {
		for _, item := range infos {
			var vo modRamIdpBinding2.Vo
			copier.Copy(&vo, &item)
			code := model.BaseNodeNo{
				Value:  item.No,
				Label:  item.Description,
				Extend: vo,
			}
			slice = append(slice, code)
		}
		rt.Data = slice
	}
	return rt.Ok()
}

// SelectNodeAllPublic 查询
//
//	@Description:
//	@receiver c
//	@param ct
func (c *RamIdpBindingService) SelectNodeAllPublic(ctx *gin.Context, ct modRamIdpBinding2.QueryPublicCt) (rt rg.Rs[[]model.BaseNodeNo]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	var query entityRam.RamIdpBindingEntity
	copier.Copy(&query, &ct)
	slice := make([]model.BaseNodeNo, 0)
	rt.Data = slice
	infos := c.sv.FindAll(ctx, query)
	if len(infos) > 0 {
		for _, item := range infos {
			var vo modRamIdpBinding2.Vo
			copier.Copy(&vo, &item)
			code := model.BaseNodeNo{
				Value:  item.No,
				Label:  item.Description,
				Extend: vo,
			}
			slice = append(slice, code)
		}
		rt.Data = slice
	}
	return rt.Ok()
}
