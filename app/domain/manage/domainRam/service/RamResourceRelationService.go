package service

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamResourceRelation"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityRam"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/holderPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg/optionsPg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"

	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/dbPg/pagePg"
	"github.com/pangu-2/go-tools/tools/numberPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
)

func init() {
	gs.Provide(new(RamResourceRelationService))
}

// RamResourceRelationService 资源关系
// @Description:
type RamResourceRelationService struct {
	sv *repositoryRam.RamResourceRelationRepository `autowire:"?"`
}

// Create 新增
//
//	@Description:
//	@receiver c
//	@param ct
//	@return rt
func (c *RamResourceRelationService) Create(ctx *gin.Context, ct modRamResourceRelation.CreateCt) (rt rg.Rs[string]) {
	if "" == ct.Name {
		return rt.ErrorMessage("名称不能为空")
	}
	holder := holderPg.GetContextAccount(ctx)

	var info entityRam.RamResourceRelationEntity
	copier.Copy(&info, &ct)
	log.Infof(ctx, log.TagAppDef, "info%+v", info)
	info.TenantNo = holder.GetTenantNo()
	c.sv.Create(ctx, &info)
	log.Infof(ctx, log.TagAppDef, "save=%+v", info)
	return rg.OkData(numberPg.Int64ToString(info.ID))
}

// Update 更新
//
//	@Description:
//	@receiver c
//	@param ct
//	@return rt
func (c *RamResourceRelationService) Update(ctx *gin.Context, ct modRamResourceRelation.UpdateCt) (rt rg.Rs[string]) {
	if ct.ID < 1 {
		return rt.ErrorMessage("id错误")
	}
	if "" == ct.Name {
		return rt.ErrorMessage("名称不能为空")
	}
	r := c.sv
	_, b := r.FindById(ctx, ct.ID.ToInt64())
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	var info entityRam.RamResourceRelationEntity
	copier.Copy(&info, &ct)
	r.Update(ctx, info, info.ID)
	return rt.Ok()
}

// Detail 详情
//
//	@Description:
//	@receiver c
//	@param id
func (c *RamResourceRelationService) Detail(ctx *gin.Context, id int64) (rt rg.Rs[modRamResourceRelation.Vo]) {
	if id < 1 {
		return rt.ErrorMessage("id错误")
	}
	find, b := c.sv.FindById(ctx, id)
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	var info modRamResourceRelation.Vo
	copier.Copy(&info, find)
	return rt.OkData(info)
}

// Delete 删除
//
//	@Description:
//	@receiver c
//	@param ct
func (c *RamResourceRelationService) Delete(ctx *gin.Context, ct model.BaseIdsCt[string]) (rt rg.Rs[string]) {
	if len(ct.Ids) < 1 {
		return rt.ErrorMessage("id错误")
	}
	r := c.sv
	finds, b := r.FindAllByIdStringIn(ctx, ct.Ids)
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	for _, info := range finds {
		log.Infof(ctx, log.TagAppDef, "id=%v,TenantId=%v", info.ID, info.TenantNo)
		r.DeleteById(ctx, info.ID)
	}
	return rt.Ok()
}

// Enable 启用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *RamResourceRelationService) Enable(ctx *gin.Context, ct model.BaseIdsCt[string]) (rt rg.Rs[string]) {
	return c.State(ctx, ct.Ids, enumStatePg.ENABLE)
}

// Disable 禁用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *RamResourceRelationService) Disable(ctx *gin.Context, ct model.BaseIdsCt[string]) (rt rg.Rs[string]) {
	return c.State(ctx, ct.Ids, enumStatePg.GetType(enumStatePg.DISABLE))
}

// State 状态 启用/禁用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *RamResourceRelationService) State(ctx *gin.Context, ids []string, state enumStatePg.State) (rt rg.Rs[string]) {
	if len(ids) < 1 {
		return rt.ErrorMessage("id错误")
	}
	r := c.sv
	_, b := r.FindAllByIdStringIn(ctx, ids)
	if b {
		return rt.ErrorMessage("数据不存在")
	}
	//for _, info := range finds {
	//	if info.StateOrder != state.IndexInt8() {
	//		r.Update(ctx, entityRam.RamResourceRelationEntity{StateOrder: state.IndexInt8()}, info.ID)
	//	}
	//}
	return rt.Ok()
}

// StateEnableDisable 状态 设置 有效 停用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *RamResourceRelationService) StateEnableDisable(ctx *gin.Context, ids []string, state enumStatePg.State) (rt rg.Rs[string]) {
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
func (c *RamResourceRelationService) LogicalDeletion(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
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
		//for _, info := range finds {
		//	enum := enumStatePg.State(info.StateOrder)
		//	// 有效 停用，反转 为对应的 取消 弃置
		//	if ok, reverse := enum.ReverseEnableDisable(); ok {
		//		repository.Update(ctx, entityRam.RamResourceRelationEntity{StateOrder: reverse.IndexInt8()}, info.ID)
		//	}
		//}
	}

	return rt.Ok()
}

// LogicalRecovery 逻辑删除恢复
//
//	@Description:
//	@receiver c
//	@param ct
func (c *RamResourceRelationService) LogicalRecovery(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ids)
	if len(ids) < 1 {
		return rt.ErrorMessage("id错误")
	}
	repository := c.sv
	_, b := repository.FindAllByIdStringIn(ctx, ids)
	if b {
		return rt.ErrorMessage("数据不存在")
	}
	//for _, info := range finds {
	//	enum := enumStatePg.State(info.StateOrder)
	//	//  取消 弃置 批量删除，反转 为对应的 有效 停用 停用
	//	if ok, reverse := enum.ReverseCancelLayAside(); ok {
	//		repository.Update(ctx, entityRam.RamResourceRelationEntity{StateOrder: reverse.IndexInt8()}, info.ID)
	//	}
	//}
	return rt.Ok()
}

// PhysicalDeletion 物理删除
//
//	@Description:
//	@receiver c
//	@param ct
func (c *RamResourceRelationService) PhysicalDeletion(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
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
func (c *RamResourceRelationService) Query(ctx *gin.Context, ct modRamResourceRelation.QueryCt) (rt rg.Rs[pagePg.Paginator[modRamResourceRelation.Vo]]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	var query entityRam.RamResourceRelationEntity
	copier.Copy(&query, &ct)
	//
	slice := make([]modRamResourceRelation.Vo, 0)
	//
	page, err := c.sv.FindAllPage(ctx, query, optionsPg.WithOption(func(arg *optionsPg.OptionParams) {
		arg.Pageable = new(pagePg.PageablePageSize(0, ct.PageNum, ct.PageSize))
	}), optionsPg.WithOption(func(arg *optionsPg.OptionParams) {
		arg.Db = arg.Db.Order("create_at desc")
	}))
	if nil != err {
		return rt.Ok()
	}

	if page.Total > 0 && page.Data != nil && len(page.Data) > 0 {
		pg := pagePg.NewPaginatorByPageable[modRamResourceRelation.Vo](page.Pageable)
		//字段赋值
		for _, item := range page.Data {
			var vo modRamResourceRelation.Vo
			copier.Copy(&vo, &item)
			//
			slice = append(slice, vo)
		}
		pg.Data = slice
		pg.Pageable = page.Pageable
		rt.Data = pg
		return rt.Ok()
	}
	return rt.Ok()
}

// SelectNodePublic 查询
//
//	@Description:
//	@receiver c
//	@param ct
func (c *RamResourceRelationService) SelectNodePublic(ctx *gin.Context, ct modRamResourceRelation.QueryCt) (rt rg.Rs[[]model.BaseNode]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	var query entityRam.RamResourceRelationEntity
	copier.Copy(&query, &ct)
	slice := make([]model.BaseNode, 0)
	rt.Data = slice
	infos := c.sv.FindAll(ctx, query)
	if len(infos) > 0 {
		//for _, item := range infos {
		//	slice = append(slice, model.BaseNode{Key: numberPg.Int64ToString(item.ID),
		//		Label: item.Name,
		//	})
		//}
		rt.Data = slice
	}
	return rt.Ok()
}

// SelectNodeAllPublic 查询
//
//	@Description:
//	@receiver c
//	@param ct
func (c *RamResourceRelationService) SelectNodeAllPublic(ctx *gin.Context, ct modRamResourceRelation.QueryCt) (rt rg.Rs[[]model.BaseNode]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	var query entityRam.RamResourceRelationEntity
	copier.Copy(&query, &ct)
	slice := make([]model.BaseNode, 0)
	rt.Data = slice
	infos := c.sv.FindAll(ctx, query)
	if len(infos) > 0 {
		for _, item := range infos {
			var vo modRamResourceRelation.Vo
			copier.Copy(&vo, &item)
			slice = append(slice, model.BaseNode{Key: numberPg.Int64ToString(item.ID),
				//Label:  item.Name,
				Extend: vo,
			})
		}
		rt.Data = slice
	}
	return rt.Ok()
}

// SelectPublic 查询
//
//	@Description:
//	@receiver c
//	@param ct
func (c *RamResourceRelationService) SelectPublic(ctx *gin.Context, ct modRamResourceRelation.QueryCt) (rt rg.Rs[[]modRamResourceRelation.Vo]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	var query entityRam.RamResourceRelationEntity
	copier.Copy(&query, &ct)
	rt.Data = []modRamResourceRelation.Vo{}
	infos := c.sv.FindAll(ctx, query)
	if len(infos) > 0 {
		slice := make([]modRamResourceRelation.Vo, 0)
		for _, item := range infos {
			var vo modRamResourceRelation.Vo
			copier.Copy(&vo, &item)
			slice = append(slice, vo)
		}
		rt.Data = slice
	}
	return rt.Ok()
}

// Selected 查询已选中的
//
//	@Description:
//	@receiver c
//	@param ct
func (c *RamResourceRelationService) Selected(ctx *gin.Context, code string) (rt rg.Rs[[]string]) {
	var query entityRam.RamResourceRelationEntity
	query.Mark = code
	slice := make([]string, 0)
	rt.Data = slice
	infos := c.sv.FindAll(ctx, query)
	if len(infos) > 0 {
		for _, item := range infos {
			slice = append(slice, numberPg.Int64ToString(item.ResourceId))
		}
		rt.Data = slice
	}
	return rt.Ok()
}
