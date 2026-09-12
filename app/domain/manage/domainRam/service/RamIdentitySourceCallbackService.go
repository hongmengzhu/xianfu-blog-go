package service

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamIdentitySourceCallback"
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
	"github.com/pangu-2/go-tools/tools/numberPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
)

func init() {
	gs.Provide(new(RamIdentitySourceCallbackService))
}

// RamIdentitySourceCallbackService 认证源回调白名单
// @Description:
type RamIdentitySourceCallbackService struct {
	sv *repositoryRam.RamIdentitySourceCallbackRepository `autowire:"?"`
}

// Create 新增
func (c *RamIdentitySourceCallbackService) Create(ctx *gin.Context, ct modRamIdentitySourceCallback.CreateUpdateCt) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	var info entityRam.RamIdentitySourceCallbackEntity
	copier.Copy(&info, &ct)
	if "" == ct.Name {
		return rt.ErrorMessage("名称不能为空")
	}
	r := c.sv
	holder := holderPg.GetContextAccount(ctx)
	info.TenantNo = holder.GetTenantNo()
	info.No = noPg.No()
	log.Infof(ctx, log.TagAppDef, "info%+v", info)
	err, _ := r.Create(ctx, &info)
	if err != nil {
		return rt.ErrorMessage("保存失败 " + err.Error())
	}
	log.Infof(ctx, log.TagAppDef, "save=%+v", info)
	return rg.OkData(numberPg.Int64ToString(info.ID))
}

// Update 更新
func (c *RamIdentitySourceCallbackService) Update(ctx *gin.Context, ct modRamIdentitySourceCallback.CreateUpdateCt) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	var info entityRam.RamIdentitySourceCallbackEntity
	copier.Copy(&info, &ct)
	r := c.sv
	if ct.ID < 1 {
		return rt.ErrorMessage("id错误")
	}
	if "" == ct.Name {
		return rt.ErrorMessage("名称不能为空")
	}
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
	return rt.Ok()
}

// Detail 详情
func (c *RamIdentitySourceCallbackService) Detail(ctx *gin.Context, id int64) (rt rg.Rs[modRamIdentitySourceCallback.Vo]) {
	if id < 1 {
		return rt.ErrorMessage("id错误")
	}
	find, b := c.sv.FindById(ctx, id)
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	var info modRamIdentitySourceCallback.Vo
	copier.Copy(&info, &find)
	return rt.OkData(info)
}

// Enable 启用
func (c *RamIdentitySourceCallbackService) Enable(ctx *gin.Context, ct model.BaseIdsCt[string]) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	return c.State(ctx, ct.Ids, enumStatePg.ENABLE)
}

// Disable 禁用
func (c *RamIdentitySourceCallbackService) Disable(ctx *gin.Context, ct model.BaseIdsCt[string]) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	return c.State(ctx, ct.Ids, enumStatePg.GetType(enumStatePg.DISABLE))
}

// State 状态 启用/禁用
func (c *RamIdentitySourceCallbackService) State(ctx *gin.Context, ids []string, state enumStatePg.State) (rt rg.Rs[string]) {
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
			r.Update(ctx, entityRam.RamIdentitySourceCallbackEntity{State: state.IndexInt8()}, info.ID)
		}
	}
	return rt.Ok()
}

// StateEnableDisable 状态 设置 有效 停用
func (c *RamIdentitySourceCallbackService) StateEnableDisable(ctx *gin.Context, ids []string, state enumStatePg.State) (rt rg.Rs[string]) {
	if !state.IsEnableDisable() {
		return rt.ErrorMessage("状态错误")
	}
	return c.State(ctx, ids, state)
}

// LogicalDeletion 逻辑删除
func (c *RamIdentitySourceCallbackService) LogicalDeletion(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
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
			if ok, reverse := enum.ReverseEnableDisable(); ok {
				repository.Update(ctx, entityRam.RamIdentitySourceCallbackEntity{State: reverse.IndexInt8()}, info.ID)
			}
		}
	}
	return rt.Ok()
}

// LogicalRecovery 逻辑删除恢复
func (c *RamIdentitySourceCallbackService) LogicalRecovery(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
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
		if ok, reverse := enum.ReverseCancelLayAside(); ok {
			repository.Update(ctx, entityRam.RamIdentitySourceCallbackEntity{State: reverse.IndexInt8()}, info.ID)
		}
	}
	return rt.Ok()
}

// PhysicalDeletion 物理删除
func (c *RamIdentitySourceCallbackService) PhysicalDeletion(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
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
func (c *RamIdentitySourceCallbackService) Query(ctx *gin.Context, ct modRamIdentitySourceCallback.QueryCt) (rt rg.Rs[pagePg.Paginator[modRamIdentitySourceCallback.Vo]]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	var query entityRam.RamIdentitySourceCallbackEntity
	copier.Copy(&query, &ct)
	slice := make([]modRamIdentitySourceCallback.Vo, 0)
	rt.Data.Data = slice
	r := c.sv
	page, err := r.FindAllPage(ctx, query, optionsPg.WithOption(func(arg *optionsPg.OptionParams) {
		if ct.PageSize < 1 {
			ct.PageSize = 20
		}
		arg.Pageable = new(pagePg.PageablePageSize(0, ct.PageNum, ct.PageSize))
		arg.Db = arg.Db.Order("create_at desc")
		if strPg.IsNotBlank(ct.Wd) {
			arg.Db = arg.Db.Where("name like ?", "%"+ct.Wd+"%")
		}
	}), optionsPg.WithCtx(ctx))
	if nil != err {
		return rt.Ok()
	}

	if page.Total > 0 && page.Data != nil && len(page.Data) > 0 {
		pg := pagePg.NewPaginatorByPageable[modRamIdentitySourceCallback.Vo](page.Pageable)
		for _, item := range page.Data {
			var vo modRamIdentitySourceCallback.Vo
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

// SelectPublic 查询
func (c *RamIdentitySourceCallbackService) SelectPublic(ctx *gin.Context, ct modRamIdentitySourceCallback.QueryCt) (rt rg.Rs[[]modRamIdentitySourceCallback.Vo]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	var query entityRam.RamIdentitySourceCallbackEntity
	copier.Copy(&query, &ct)
	rt.Data = []modRamIdentitySourceCallback.Vo{}
	infos := c.sv.FindAll(ctx, query)
	if len(infos) > 0 {
		slice := make([]modRamIdentitySourceCallback.Vo, 0)
		for _, item := range infos {
			var vo modRamIdentitySourceCallback.Vo
			copier.Copy(&vo, &item)
			slice = append(slice, vo)
		}
		rt.Data = slice
	}
	return rt.Ok()
}

// ExistName 查重
func (c *RamIdentitySourceCallbackService) ExistName(ctx *gin.Context, ct model.BaseExistWdCt[string]) (rt rg.Rs[string]) {
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
