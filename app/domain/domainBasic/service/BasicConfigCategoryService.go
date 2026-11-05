package service

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/basic/modBasicConfigCategory"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/consts/automatedPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/consts/constNodePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumCommonPg/typeSysPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/repositoryBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/sdk/auth/holderPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg/optionsPg"
	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/dbPg/pagePg"
	"github.com/pangu-2/go-tools/tools/noPg"
	"github.com/pangu-2/go-tools/tools/numberPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(BasicConfigCategoryService))
}

// BasicConfigCategoryService 配置分类
type BasicConfigCategoryService struct {
	sv *repositoryBasic.BasicConfigCategoryRepository `autowire:"?"`
}

// Create 新增配置分类
func (c *BasicConfigCategoryService) Create(ctx *gin.Context, ct modBasicConfigCategory.CreateUpdateCt) (rt rg.Rs[string]) {
	log.Info(ctx, log.TagAppDef, "ct", log.Any("ct", ct))
	var info entityBasic.BasicConfigCategoryEntity
	if err := copier.Copy(&info, &ct); err != nil {
		log.Info(ctx, log.TagAppDef, "copier.Copy error:", log.Any("err", err.Error()))
	}
	info.Name = strings.TrimSpace(info.Name)
	if strPg.IsBlank(info.Name) {
		return rt.ErrorMessage("名称不能为空")
	}
	holder := holderPg.GetContextAccount(ctx)
	info.TenantNo = holder.GetTenantNo()
	if strPg.IsBlank(info.Code) {
		info.Code = automatedPg.CREATE_CODE
	}
	r := c.sv
	if !automatedPg.IsCreateCode(info.Code) {
		if !automatedPg.FormatVerify(info.Code) {
			return rt.ErrorMessage("标志格式错误")
		}
		exist, ok := r.FindByCode(ctx, info.Code)
		if ok && exist.TenantNo == info.TenantNo {
			return rt.ErrorMessage("标志已存在")
		}
	}
	parent := &entityBasic.BasicConfigCategoryEntity{}
	if strPg.IsBlank(info.ParentNo) && strPg.IsNotBlank(info.ParentId) {
		var ok bool
		parent, ok = r.FindByIdString(ctx, info.ParentId)
		if !ok {
			return rt.ErrorMessage("上级不存在")
		}
		info.ParentNo = parent.No
	}
	if strPg.IsNotBlank(info.ParentNo) {
		var ok bool
		parent, ok = r.FindByNo(ctx, info.ParentNo)
		if !ok {
			return rt.ErrorMessage("上级不存在")
		}
	}
	info.No = noPg.No()
	if automatedPg.IsCreateCode(info.Code) {
		info.Code = strPg.GenerateNumberId22()
	}
	if info.State == 0 {
		info.State = enumStatePg.ENABLE.Index()
	}
	info.TypeSys = typeSysPg.General.String()
	info.Ano = holder.GetAccount().Account
	if err, _ := r.Create(ctx, &info); err != nil {
		return rt.ErrorMessage("保存失败 " + err.Error())
	}
	if strPg.IsNotBlank(info.ParentNo) {
		info.IdLink = constNodePg.NoLinkAssemble(parent.IdLink, numberPg.Int64ToString(info.ID))
		info.NoLink = constNodePg.NoLinkAssemble(parent.NoLink, info.No)
		info.ParentNo = parent.No
		info.ParentId = numberPg.Int64ToString(parent.ID)
	} else {
		info.IdLink = constNodePg.NoLinkDefault(numberPg.Int64ToString(info.ID))
		info.NoLink = constNodePg.NoLinkDefault(info.No)
		info.ParentNo = ""
		info.ParentId = ""
	}
	if err := r.Update(ctx, info, info.ID); err != nil {
		return rt.ErrorMessage(err.Error())
	}
	return rg.OkData(numberPg.Int64ToString(info.ID))
}

// Update 更新配置分类
func (c *BasicConfigCategoryService) Update(ctx *gin.Context, ct modBasicConfigCategory.CreateUpdateCt) (rt rg.Rs[string]) {
	log.Info(ctx, log.TagAppDef, "ct", log.Any("ct", ct))
	if ct.ID.ToInt64() < 1 {
		return rt.ErrorMessage("id错误")
	}
	find, ok := c.sv.FindById(ctx, ct.ID.ToInt64())
	if !ok {
		return rt.ErrorMessage("数据不存在")
	}
	holder := holderPg.GetContextAccount(ctx)
	if strPg.IsNotBlank(find.TenantNo) && find.TenantNo != holder.GetTenantNo() {
		return rt.ErrorMessage("无权操作其他租户的分类")
	}
	var info entityBasic.BasicConfigCategoryEntity
	if err := copier.Copy(&info, &ct); err != nil {
		log.Info(ctx, log.TagAppDef, "copier.Copy error:", log.Any("err", err.Error()))
	}
	info.Name = strings.TrimSpace(info.Name)
	if strPg.IsBlank(info.Name) {
		return rt.ErrorMessage("名称不能为空")
	}
	parent := &entityBasic.BasicConfigCategoryEntity{}
	parentMoved := false
	var descendants []*entityBasic.BasicConfigCategoryEntity
	if strPg.IsBlank(info.ParentNo) && strPg.IsNotBlank(info.ParentId) {
		parent, ok = c.sv.FindByIdString(ctx, info.ParentId)
		if !ok {
			return rt.ErrorMessage("上级不存在")
		}
		info.ParentNo = parent.No
	}
	if strPg.IsNotBlank(info.ParentNo) {
		parent, ok = c.sv.FindByNo(ctx, info.ParentNo)
		if !ok {
			return rt.ErrorMessage("上级不存在")
		}
		if parent.ID == find.ID {
			return rt.ErrorMessage("上级不能等于自己")
		}
		parentMoved = parent.No != find.ParentNo
	} else {
		parentMoved = strPg.IsNotBlank(find.ParentNo)
	}
	if strPg.IsBlank(info.Code) {
		info.Code = find.Code
	} else {
		exist, exists := c.sv.FindByCodeAndIdNot(ctx, info.Code, ct.ID.ToString())
		if exists && exist.TenantNo == holder.GetTenantNo() {
			return rt.ErrorMessage("标志已存在")
		}
	}
	if parentMoved {
		descendants, _ = c.sv.FindAllByNoLink(ctx, find.No)
		if strPg.IsNotBlank(info.ParentNo) {
			for _, item := range descendants {
				if item.No == parent.No {
					return rt.ErrorMessage("无法保存，不能设置为自己的子集")
				}
			}
		}
	}
	info.TypeSys = typeSysPg.General.String()
	info.No = ""
	if strPg.IsNotBlank(info.ParentNo) {
		info.IdLink = constNodePg.NoLinkAssemble(parent.IdLink, numberPg.Int64ToString(find.ID))
		info.NoLink = constNodePg.NoLinkAssemble(parent.NoLink, find.No)
		info.ParentNo = parent.No
		info.ParentId = numberPg.Int64ToString(parent.ID)
	} else {
		info.IdLink = constNodePg.NoLinkDefault(numberPg.Int64ToString(find.ID))
		info.NoLink = constNodePg.NoLinkDefault(find.No)
		info.ParentNo = ""
		info.ParentId = ""
	}
	updates := map[string]any{
		"name":        info.Name,
		"name_fl":     info.NameFl,
		"name_full":   info.NameFull,
		"code":        info.Code,
		"description": info.Description,
		"sort":        info.Sort,
		"parent_no":   info.ParentNo,
		"parent_id":   info.ParentId,
		"id_link":     info.IdLink,
		"no_link":     info.NoLink,
		"type_sys":    info.TypeSys,
	}
	if err := c.sv.UpdateMap(ctx, updates, find.ID); err != nil {
		return rt.ErrorMessage(err.Error())
	}
	if parentMoved {
		for _, child := range descendants {
			if child.ID == find.ID {
				continue
			}
			if strPg.IsNotBlank(find.NoLink) && strings.HasPrefix(child.NoLink, find.NoLink) {
				child.NoLink = strings.Replace(child.NoLink, find.NoLink, info.NoLink, 1)
			}
			if strPg.IsNotBlank(find.IdLink) && strings.HasPrefix(child.IdLink, find.IdLink) {
				child.IdLink = strings.Replace(child.IdLink, find.IdLink, info.IdLink, 1)
			}
			if err := c.sv.Update(ctx, entityBasic.BasicConfigCategoryEntity{IdLink: child.IdLink, NoLink: child.NoLink}, child.ID); err != nil {
				return rt.ErrorMessage(err.Error())
			}
		}
	}
	return rt.Ok()
}

// Detail 详情
func (c *BasicConfigCategoryService) Detail(ctx *gin.Context, id int64) (rt rg.Rs[modBasicConfigCategory.Vo]) {
	if id < 1 {
		return rt.ErrorMessage("id错误")
	}
	find, ok := c.sv.FindById(ctx, id)
	if !ok {
		return rt.ErrorMessage("数据不存在")
	}
	var vo modBasicConfigCategory.Vo
	copier.Copy(&vo, find)
	return rt.OkData(vo)
}

// Query 分页查询
func (c *BasicConfigCategoryService) Query(ctx *gin.Context, ct modBasicConfigCategory.QueryCt) (rt rg.Rs[pagePg.Paginator[modBasicConfigCategory.Vo]]) {
	var query entityBasic.BasicConfigCategoryEntity
	copier.Copy(&query, &ct)
	holder := holderPg.GetContextAccount(ctx)
	slice := make([]modBasicConfigCategory.Vo, 0)
	rt.Data.Data = slice
	page, err := c.sv.FindAllPage(ctx, query, optionsPg.WithOption(func(arg *optionsPg.OptionParams) {
		if ct.PageSize < 1 {
			ct.PageSize = 20
		}
		arg.Pageable = new(pagePg.PageablePageSize(0, ct.PageNum, ct.PageSize))
		arg.Db = arg.Db.Where("tenant_no=?", holder.GetTenantNo()).Order("sort asc, create_at")
		if strPg.IsNotBlank(ct.Wd) {
			arg.Db = arg.Db.Where("name like ? or code like ? or description like ?", "%"+ct.Wd+"%", "%"+ct.Wd+"%", "%"+ct.Wd+"%")
		}
	}), optionsPg.WithCtx(ctx))
	if err != nil {
		log.Error(ctx, log.TagAppDef, err, "query config category error")
		return rt.Ok()
	}
	if page.Total > 0 && page.Data != nil && len(page.Data) > 0 {
		pg := pagePg.NewPaginatorByPageable[modBasicConfigCategory.Vo](page.Pageable)
		for _, item := range page.Data {
			var vo modBasicConfigCategory.Vo
			copier.Copy(&vo, item)
			slice = append(slice, vo)
		}
		pg.Data = slice
		pg.Pageable = page.Pageable
		rt.Data = pg
	}
	return rt.Ok()
}

// Enable 启用
func (c *BasicConfigCategoryService) Enable(ctx *gin.Context, ct model.BaseIdsCt[string]) (rt rg.Rs[string]) {
	return c.State(ctx, ct.Ids, enumStatePg.ENABLE)
}

// Disable 禁用
func (c *BasicConfigCategoryService) Disable(ctx *gin.Context, ct model.BaseIdsCt[string]) (rt rg.Rs[string]) {
	return c.State(ctx, ct.Ids, enumStatePg.GetType(enumStatePg.DISABLE))
}

// State 状态 启用/禁用
func (c *BasicConfigCategoryService) State(ctx *gin.Context, ids []string, state enumStatePg.State) (rt rg.Rs[string]) {
	if len(ids) < 1 {
		return rt.ErrorMessage("id错误")
	}
	finds, ok := c.sv.FindAllByIdStringIn(ctx, ids, optionsPg.WithCtx(ctx))
	if !ok {
		return rt.ErrorMessage("数据不存在")
	}
	for _, info := range finds {
		if info.State != state.IndexInt8() {
			c.sv.Update(ctx, entityBasic.BasicConfigCategoryEntity{State: state.IndexInt8()}, info.ID)
		}
	}
	return rt.Ok()
}

// StateEnableDisable 状态 设置 有效 停用
func (c *BasicConfigCategoryService) StateEnableDisable(ctx *gin.Context, ids []string, state enumStatePg.State) (rt rg.Rs[string]) {
	if !state.IsEnableDisable() {
		return rt.ErrorMessage("状态错误")
	}
	return c.State(ctx, ids, state)
}

// LogicalDeletion 逻辑删除
func (c *BasicConfigCategoryService) LogicalDeletion(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
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
			if ok, reverse := enumStatePg.State(info.State).ReverseEnableDisable(); ok {
				c.sv.Update(ctx, entityBasic.BasicConfigCategoryEntity{State: reverse.IndexInt8()}, info.ID)
			}
		}
	}
	return rt.Ok()
}

// LogicalRecovery 逻辑删除恢复
func (c *BasicConfigCategoryService) LogicalRecovery(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
	if len(ids) < 1 {
		return rt.ErrorMessage("id错误")
	}
	finds, ok := c.sv.FindAllByIdStringIn(ctx, ids, optionsPg.WithCtx(ctx))
	if !ok {
		return rt.ErrorMessage("数据不存在")
	}
	for _, info := range finds {
		if ok, reverse := enumStatePg.State(info.State).ReverseCancelLayAside(); ok {
			c.sv.Update(ctx, entityBasic.BasicConfigCategoryEntity{State: reverse.IndexInt8()}, info.ID)
		}
	}
	return rt.Ok()
}

// PhysicalDeletion 物理删除
func (c *BasicConfigCategoryService) PhysicalDeletion(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
	if len(ids) < 1 {
		return rt.ErrorMessage("id错误")
	}
	if err := c.sv.DeleteByIdsString(ctx, ids); err != nil {
		return rt.ErrorMessage(err.Error())
	}
	return rt.Ok()
}

// ExistName 查重
func (c *BasicConfigCategoryService) ExistName(ctx *gin.Context, ct model.BaseExistWdCt[string]) (rt rg.Rs[string]) {
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

// ExistCode 查重
func (c *BasicConfigCategoryService) ExistCode(ctx *gin.Context, ct model.BaseExistWdCt[string]) (rt rg.Rs[string]) {
	if strPg.IsBlank(ct.Wd) {
		return rt.ErrorMessage("关键词不能为空")
	}
	id := "0"
	if strPg.IsNotBlank(ct.Id) {
		id = ct.Id
	}
	if _, exists := c.sv.FindByCodeAndIdNot(ctx, ct.Wd, id); exists {
		return rt.ErrorMessage("重复，已存在")
	}
	return rt.OkMessage("可以使用")
}
