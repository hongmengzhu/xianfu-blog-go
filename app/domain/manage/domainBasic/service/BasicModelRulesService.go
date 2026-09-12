package service

import (
	"strings"

	"github.com/gin-gonic/gin"
	modelRules2 "github.com/hongmengzhu/xianfu-blog-go/app/domain/manage/domainBasic/service/modelRules"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/basic/modBasicModelRules"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBasic"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg/optionsPg"
	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/dbPg/pagePg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(BasicModelRulesService))
}

type BasicModelRulesService struct {
	sv *repositoryBasic.BasicModelRulesRepository `autowire:"?"`
	sp *modelRules2.Sp                            `autowire:"?"`
}

func (c *BasicModelRulesService) CreateUpdateData(ctx *gin.Context, ct modBasicModelRules.CreateUpdateDataCt) (rt rg.Rs[string]) {
	return modelRules2.NewCreateUpdateData(c.sp, ct, true).Process(ctx)
}
func (c *BasicModelRulesService) CreateUpdate(ctx *gin.Context, ct modBasicModelRules.CreateUpdateCt) (rt rg.Rs[string]) {
	return modelRules2.NewCreateUpdate(c.sp, ct, true).Process(ctx)
}

// Detail 详情
//func (c *BasicModelRulesService) Detail(ctx *gin.Context, id string) (rt rg.Rs[modBasicModelRules.CreateUpdateCt]) {
//	return configEvent.NewDetail(c.sp).Process(ctx, id)
//}

// Enable 启用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicModelRulesService) Enable(ctx *gin.Context, ct model.BaseIdsCt[string]) (rt rg.Rs[string]) {
	return c.State(ctx, ct.Ids, enumStatePg.ENABLE)
}

// Disable 禁用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicModelRulesService) Disable(ctx *gin.Context, ct model.BaseIdsCt[string]) (rt rg.Rs[string]) {
	return c.State(ctx, ct.Ids, enumStatePg.GetType(enumStatePg.DISABLE))
}

// State 状态 启用/禁用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicModelRulesService) State(ctx *gin.Context, ids []string, state enumStatePg.State) (rt rg.Rs[string]) {
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
			r.Update(ctx, entityBasic.BasicModelRulesEntity{State: state.IndexInt8()}, info.ID)
		}
	}
	return rt.Ok()
}

// StateEnableDisable 状态 设置 有效 停用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicModelRulesService) StateEnableDisable(ctx *gin.Context, ids []string, state enumStatePg.State) (rt rg.Rs[string]) {
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
func (c *BasicModelRulesService) LogicalDeletion(ctx *gin.Context, ct model.BaseIdsCt[string]) (rt rg.Rs[string]) {
	if nil == ct.Ids || len(ct.Ids) < 1 {
		return rt.ErrorMessage("id错误")
	}
	if ct.Extend == nil {
		return rt.ErrorMessage("参数错误")
	}
	fieldNo := ""
	if obj, ok := ct.Extend["fieldNo"]; ok {
		fieldNo = obj.(string)
	}
	if strPg.IsBlank(fieldNo) {
		return rt.ErrorMessage("字段错误")
	}
	ids := make([]string, 0)
	for _, id := range ct.Ids {
		if strPg.IsNotBlank(id) {
			ids = append(ids, strings.TrimSpace(id))
		}
	}
	if len(ids) < 1 {
		return rt.ErrorMessage("请选择数据")
	}
	repository := c.sv
	finds, b := repository.FindAllByIdStringIn(ctx, ids)
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	if c.sv.Config().Data.Delete {
		idsInt := make([]int64, 0)
		for _, info := range finds {
			log.Infof(ctx, log.TagAppDef, "id=%v", info.ID)
			if info.FieldNo == fieldNo {
				idsInt = append(idsInt, info.ID)
			}
		}
		if len(idsInt) > 0 {
			repository.DeleteByIds(ctx, idsInt)
		}
	} else {
		for _, info := range finds {
			if info.FieldNo != fieldNo {
				continue
			}
			enum := enumStatePg.State(info.State)
			// 有效 停用，反转 为对应的 取消 弃置
			if ok, reverse := enum.ReverseEnableDisable(); ok {
				repository.Update(ctx, entityBasic.BasicModelRulesEntity{State: reverse.IndexInt8()}, info.ID)
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
func (c *BasicModelRulesService) LogicalRecovery(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
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
			repository.Update(ctx, entityBasic.BasicModelRulesEntity{State: reverse.IndexInt8()}, info.ID)
		}
	}
	return rt.Ok()
}

// PhysicalDeletion 物理删除
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicModelRulesService) PhysicalDeletion(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
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
func (c *BasicModelRulesService) Query(ctx *gin.Context, ct modBasicModelRules.QueryCt) (rt rg.Rs[pagePg.Paginator[modBasicModelRules.Vo]]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	var query entityBasic.BasicModelRulesEntity
	copier.Copy(&query, &ct)
	slice := make([]modBasicModelRules.Vo, 0)
	rt.Data.Data = slice
	r := c.sv
	page, err := r.FindAllPage(ctx, query, optionsPg.WithOption(func(arg *optionsPg.OptionParams) {
		if ct.PageSize < 1 {
			ct.PageSize = 20
		}
		arg.Pageable = new(pagePg.PageablePageSize(0, ct.PageNum, ct.PageSize))
		//自定义查询
		arg.Db = arg.Db.Order("create_at asc")
		if strPg.IsNotBlank(ct.Wd) {
			arg.Db = arg.Db.Where("name like ?", "%"+ct.Wd+"%").
				Or("description like ?", "%"+ct.Wd+"%")
		}
	}), optionsPg.WithCtx(ctx))

	if nil != err {
		return rt.Ok()
	}

	if page.Total > 0 && page.Data != nil && len(page.Data) > 0 {
		pg := pagePg.NewPaginatorByPageable[modBasicModelRules.Vo](page.Pageable)
		//字段赋值
		for _, item := range page.Data {
			var vo modBasicModelRules.Vo
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
func (c *BasicModelRulesService) SelectNodeAllPublic(ctx *gin.Context, ct modBasicModelRules.QueryPublicCt) (rt rg.Rs[[]model.BaseNodeNo]) {
	var query entityBasic.BasicModelRulesEntity
	copier.Copy(&query, &ct)
	slice := make([]model.BaseNodeNo, 0)
	rt.Data = slice
	infos := c.sv.FindAll(ctx, query)
	if len(infos) > 0 {
		for _, item := range infos {
			var vo modBasicModelRules.Vo
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

// AllByValueNo 查询
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicModelRulesService) AllByValueNo(ctx *gin.Context, ct modBasicModelRules.QueryPublicCt) (rt rg.Rs[[]modBasicModelRules.Vo]) {
	var query entityBasic.BasicModelRulesEntity
	copier.Copy(&query, &ct)
	//
	slice := make([]modBasicModelRules.Vo, 0)
	rt.Data = slice
	//
	if strPg.IsBlank(ct.ValueNo) {
		return rt.ErrorMessage("模型/事件编号错误")
	}
	infos := c.sv.FindAll(ctx, query)
	if len(infos) > 0 {
		for _, item := range infos {
			var vo modBasicModelRules.Vo
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
func (c *BasicModelRulesService) ExistName(ctx *gin.Context, ct model.BaseExistWdCt[string]) (rt rg.Rs[string]) {
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
func (c *BasicModelRulesService) ExistCode(ctx *gin.Context, ct model.BaseExistWdCt[string]) (rt rg.Rs[string]) {
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
