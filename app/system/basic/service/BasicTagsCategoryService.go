package service

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/basic/modBasicTagsCategory"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBasic"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/holderPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/automatedPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constNodePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/enumCommonPg/typeSysPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg/optionsPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/excelPg"
	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/dbPg/pagePg"
	"github.com/pangu-2/go-tools/tools/noPg"
	"github.com/pangu-2/go-tools/tools/numberPg"
	"github.com/pangu-2/go-tools/tools/slicePg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(BasicTagsCategoryService))
}

// BasicTagsCategoryService 模块
// @Description:
type BasicTagsCategoryService struct {
	sv *repositoryBasic.BasicTagsCategoryRepository `autowire:"?"`
}

// Create 新增
//
//	@Description:
//	@receiver c
//	@param ct
//	@return rt
func (c *BasicTagsCategoryService) Create(ctx *gin.Context, ct modBasicTagsCategory.CreateUpdateCt) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%#v", ct)
	var info entityBasic.BasicTagsCategoryEntity
	err := copier.Copy(&info, &ct)
	if err != nil {
		log.Infof(ctx, log.TagAppDef, "copier.Copy error: %+v", err)
	}
	if "" == ct.Name {
		return rt.ErrorMessage("名称不能为空")
	}
	if strPg.IsBlank(ct.Code) {
		info.Code = automatedPg.CREATE_CODE
	}
	holder := holderPg.GetContextAccount(ctx)
	parent := &entityBasic.BasicTagsCategoryEntity{}
	r := c.sv
	//判断是否是自动,不是自动
	if !automatedPg.IsCreateCode(info.Code) {
		//判断格式是否满足要求
		if !automatedPg.FormatVerify(info.Code) {
			return rt.ErrorMessage("标志格式不能为空")
		}
		//不是自动
		_, result := r.FindByCode(ctx, info.Code)
		if result {
			return rt.ErrorMessage("标志已存在")
		}
	}
	result := false
	if strPg.IsNotBlank(info.ParentNo) {
		parent, result = r.FindByNo(ctx, info.ParentNo)
		if !result {
			return rt.ErrorMessage("上级不存在")
		}
	}
	info.No = noPg.No()
	//自动设置编号
	if automatedPg.IsCreateCode(info.Code) {
		info.Code = strPg.GenerateNumberId22()
	}
	info.TenantNo = holder.GetTenantNo()
	log.Infof(ctx, log.TagAppDef, "info=%+v", info)
	err, _ = r.Create(ctx, &info)
	if err != nil {
		return rt.ErrorMessage("保存失败 " + err.Error())
	}
	//设置上级 link
	if strPg.IsNotBlank(info.ParentNo) {
		info.IdLink = constNodePg.NoLinkAssemble(parent.IdLink, numberPg.Int64ToString(info.ID))
		info.NoLink = constNodePg.NoLinkAssemble(parent.NoLink, info.No)
		info.ParentNo = parent.No
		info.ParentId = numberPg.Int64ToString(parent.ID)
	} else {
		info.IdLink = constNodePg.NoLinkDefault(numberPg.Int64ToString(info.ID))
		info.NoLink = constNodePg.NoLinkDefault(info.No)
		info.ParentId = ""
		info.ParentNo = ""
	}
	err = r.Update(ctx, info, info.ID)
	if err != nil {
		return rt.ErrorMessage(err.Error())
	}
	return rg.OkData(numberPg.Int64ToString(info.ID))
}

// Update 更新
//
//	@Description:
//	@receiver c
//	@param ct
//	@return rt
func (c *BasicTagsCategoryService) Update(ctx *gin.Context, ct modBasicTagsCategory.CreateUpdateCt) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%#v", ct)
	var info entityBasic.BasicTagsCategoryEntity
	copier.Copy(&info, &ct)
	if ct.ID < 1 {
		return rt.ErrorMessage("id错误")
	}
	if "" == ct.Name {
		return rt.ErrorMessage("名称不能为空")
	}
	if strPg.IsBlank(ct.Code) {
		info.Code = ""
	}
	r := c.sv
	if strPg.IsBlank(ct.Code) {
		info.Code = ""
	} else {
		_, result := r.FindByCodeAndIdNot(ctx, ct.Code, ct.ID.ToString())
		if result {
			return rt.ErrorMessage("标志已存在")
		}
	}
	find, b := r.FindById(ctx, ct.ID.ToInt64())
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	//上级
	parent := &entityBasic.BasicTagsCategoryEntity{}
	var childData []*entityBasic.BasicTagsCategoryEntity
	if strPg.IsNotBlank(info.ParentNo) {
		result := false
		parent, result = r.FindByNo(ctx, info.ParentNo)
		if !result {
			return rt.ErrorMessage("上级不存在")
		}
		if parent.ID == ct.ID.ToInt64() {
			return rt.ErrorMessage("上级不能等于自己")
		}
		//新的ID 不等于 旧的上级时,检测是否已经 在新的子集已存在
		if parent.No != find.ParentNo {
			result2 := false
			childData, result2 = r.FindAllByNoLink(ctx, find.No)
			if result2 {
				//log.Infof(ctx, log.TagAppDef,"data=%+v \n", childData)
				for _, item := range childData {
					if item.No == parent.No {
						return rt.ErrorMessage("无法保存，不能设置为自己的子集")
					}
				}
			}
		}
	}

	info.TypeSys = typeSysPg.General.Index()
	//设置上级 link
	if strPg.IsNotBlank(ct.ParentNo) {
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
	info.No = ""
	log.Infof(ctx, log.TagAppDef, "info.IdLink=%+v", info.IdLink)
	err := r.Update(ctx, info, info.ID)
	if err != nil {
		log.Errorf(ctx, log.TagAppDef, "update error=%+v", err)
		return rt.ErrorMessage(err.Error())
	}
	log.Infof(ctx, log.TagAppDef, "save.info=%+v", info)
	//更改上级后，相关子集修改
	if strPg.IsNotBlank(ct.ParentNo) && nil != childData {
		maps := slicePg.ToMapArray(childData, func(t *entityBasic.BasicTagsCategoryEntity) (string, *entityBasic.BasicTagsCategoryEntity) {
			if strPg.IsBlank(t.ParentNo) {
				return constNodePg.ROOT, t
			}
			return t.ParentNo, t
		})
		if strPg.IsBlank(info.ParentNo) {
			info.ParentNo = constNodePg.ROOT
		}
		for _, item := range maps[info.ParentNo] {
			item.IdLink = constNodePg.NoLinkAssemble(info.IdLink, numberPg.Int64ToString(find.ID))
			item.NoLink = constNodePg.NoLinkAssemble(info.NoLink, item.No)
			c.childParentIdLink(maps, item)
		}
		log.Infof(ctx, log.TagAppDef, "maps=%+v", maps)
		for _, val := range maps {
			for _, item := range val {
				if item.ID == find.ID {
					continue
				}
				r.Update(ctx, entityBasic.BasicTagsCategoryEntity{IdLink: item.IdLink,
					NoLink: item.NoLink}, item.ID)
			}
		}
		maps = nil
	}
	return rt.Ok()
}

// ChildParentIdLink 子集 上级 link更新
//
//	@Description:
//	@receiver c
//	@param id
func (c *BasicTagsCategoryService) childParentIdLink(maps map[string][]*entityBasic.BasicTagsCategoryEntity, parent *entityBasic.BasicTagsCategoryEntity) {
	key := parent.ParentNo
	if strPg.IsBlank(parent.ParentNo) {
		key = constNodePg.ROOT
	}
	entities := maps[key]
	for _, item := range entities {
		item.IdLink = constNodePg.NoLinkAssemble(parent.IdLink, numberPg.Int64ToString(item.ID))
		item.NoLink = constNodePg.NoLinkAssemble(parent.NoLink, item.No)
	}
}

// CacheOverride 缓存重载
//
//	@Description:
//	@receiver c
func (c *BasicTagsCategoryService) CacheOverride(ctx *gin.Context) {
	r := c.sv
	infos, b := r.FindAllData(ctx)
	if !b {
		return
	}
	maps := slicePg.ToMapArray(infos, func(t *entityBasic.BasicTagsCategoryEntity) (string, *entityBasic.BasicTagsCategoryEntity) {
		if strPg.IsBlank(t.ParentNo) {
			return constNodePg.ROOT, t
		}
		return t.ParentNo, t
	})
	for _, item := range maps[constNodePg.ROOT] {
		item.IdLink = constNodePg.NoLinkDefault(numberPg.Int64ToString(item.ID))
		item.NoLink = constNodePg.NoLinkDefault(item.No)
		c.childParentIdLink(maps, item)
	}
	log.Infof(ctx, log.TagAppDef, "maps=%+v", maps)
	for _, val := range maps {
		for _, item := range val {
			r.Update(ctx, entityBasic.BasicTagsCategoryEntity{
				IdLink: item.IdLink,
				NoLink: item.NoLink},
				item.ID)
		}
	}
	maps = nil
}

// Detail 详情
//
//	@Description:
//	@receiver c
//	@param id
func (c *BasicTagsCategoryService) Detail(ctx *gin.Context, id int64) (rt rg.Rs[modBasicTagsCategory.Vo]) {
	if id < 1 {
		return rt.ErrorMessage("id错误")
	}
	find, b := c.sv.FindById(ctx, id)
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	var info modBasicTagsCategory.Vo
	copier.Copy(&info, &find)
	return rt.OkData(info)
}

// Enable 启用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicTagsCategoryService) Enable(ctx *gin.Context, ct model.BaseIdsCt[string]) (rt rg.Rs[string]) {
	return c.State(ctx, ct.Ids, enumStatePg.ENABLE)
}

// Disable 禁用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicTagsCategoryService) Disable(ctx *gin.Context, ct model.BaseIdsCt[string]) (rt rg.Rs[string]) {
	return c.State(ctx, ct.Ids, enumStatePg.GetType(enumStatePg.DISABLE))
}

// State 状态 启用/禁用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicTagsCategoryService) State(ctx *gin.Context, ids []string, state enumStatePg.State) (rt rg.Rs[string]) {
	if len(ids) < 1 {
		return rt.ErrorMessage("id错误")
	}
	ctxR := c.sv
	finds, b := ctxR.FindAllByIdStringIn(ctx, ids)
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	for _, info := range finds {
		if info.State != state.IndexInt8() {
			ctxR.Update(ctx, entityBasic.BasicTagsCategoryEntity{State: state.IndexInt8()}, info.ID)
		}
	}
	return rt.Ok()
}

// StateEnableDisable 状态 设置 有效 停用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicTagsCategoryService) StateEnableDisable(ctx *gin.Context, ids []string, state enumStatePg.State) (rt rg.Rs[string]) {
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
func (c *BasicTagsCategoryService) LogicalDeletion(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ids)
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
			log.Infof(ctx, log.TagAppDef, "id=%v,TenantId=%v", info.ID, "")
		}
		repository.DeleteByIdsString(ctx, ids)
	} else {
		for _, info := range finds {
			enum := enumStatePg.State(info.State)
			// 有效 停用，反转 为对应的 取消 弃置
			if ok, reverse := enum.ReverseEnableDisable(); ok {
				repository.Update(ctx, entityBasic.BasicTagsCategoryEntity{State: reverse.IndexInt8()}, info.ID)
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
func (c *BasicTagsCategoryService) LogicalRecovery(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
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
			repository.Update(ctx, entityBasic.BasicTagsCategoryEntity{State: reverse.IndexInt8()}, info.ID)
		}
	}
	return rt.Ok()
}

// PhysicalDeletion 物理删除
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicTagsCategoryService) PhysicalDeletion(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
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
		log.Infof(ctx, log.TagAppDef, "id=%v,TenantId=%v", info.ID, "")
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
func (c *BasicTagsCategoryService) Query(ctx *gin.Context, ct modBasicTagsCategory.QueryCt) (rt rg.Rs[pagePg.Paginator[modBasicTagsCategory.Vo]]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	var query entityBasic.BasicTagsCategoryEntity
	copier.Copy(&query, &ct)
	slice := make([]modBasicTagsCategory.Vo, 0)
	rt.Data.Data = slice
	r := c.sv
	page, err := r.FindAllPage(ctx, query, optionsPg.WithOption(func(arg *optionsPg.OptionParams) {
		if ct.PageSize < 1 {
			ct.PageSize = 20
		}
		arg.Pageable = new(pagePg.PageablePageSize(0, ct.PageNum, ct.PageSize))
		//自定义查询
		arg.Db = arg.Db.Order("create_at desc")
		//自定义查询
		if strPg.IsNotBlank(ct.Wd) {
			arg.Db = arg.Db.Where("name like ?", "%"+ct.Wd+"%")
		}
	}))
	if nil != err {
		return rt.Ok()
	}

	if page.Total > 0 && page.Data != nil && len(page.Data) > 0 {
		pg := pagePg.NewPaginatorByPageable[modBasicTagsCategory.Vo](page.Pageable)
		//字段赋值
		for _, item := range page.Data {
			var vo modBasicTagsCategory.Vo
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

// QueryPublic 查询
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicTagsCategoryService) QueryPublic(ctx *gin.Context, ct modBasicTagsCategory.QueryCt) (rt rg.Rs[pagePg.Paginator[modBasicTagsCategory.Vo]]) {
	var query entityBasic.BasicTagsCategoryEntity
	copier.Copy(&query, &ct)
	slice := make([]modBasicTagsCategory.Vo, 0)
	rt.Data.Data = slice
	r := c.sv
	page, err := r.FindAllPage(ctx, query, optionsPg.WithOption(func(arg *optionsPg.OptionParams) {
		if ct.PageSize < 1 {
			ct.PageSize = 20
		}
		arg.Pageable = new(pagePg.PageablePageSize(0, ct.PageNum, ct.PageSize))
		//自定义查询
		arg.Db = arg.Db.Order("create_at desc")
		if strPg.IsNotBlank(ct.Wd) {
			arg.Db = arg.Db.Where("name like ?", "%"+ct.Wd+"%")
		}
	}))
	if nil != err {
		return rt.Ok()
	}

	if page.Total > 0 && page.Data != nil && len(page.Data) > 0 {

		pg := pagePg.NewPaginatorByPageable[modBasicTagsCategory.Vo](page.Pageable)
		//字段赋值
		for _, item := range page.Data {
			var vo modBasicTagsCategory.Vo
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

// SelectNodePublic 查询
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicTagsCategoryService) SelectNodePublic(ctx *gin.Context, ct modBasicTagsCategory.QueryPublicCt) (rt rg.Rs[[]model.BaseNodeNo]) {
	var query entityBasic.BasicTagsCategoryEntity
	copier.Copy(&query, &ct)
	slice := make([]model.BaseNodeNo, 0)
	rt.Data = slice
	infos := c.sv.FindAll(ctx, query)
	if len(infos) > 0 {
		for _, item := range infos {
			var vo modBasicTagsCategory.Vo
			copier.Copy(&vo, &item)
			code := model.BaseNodeNo{
				Value:    item.No,
				Label:    item.Name,
				ParentNo: item.ParentNo,
				Extend:   vo,
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
func (c *BasicTagsCategoryService) SelectNodeAllPublic(ctx *gin.Context, ct modBasicTagsCategory.QueryPublicCt) (rt rg.Rs[[]model.BaseNodeNo]) {
	var query entityBasic.BasicTagsCategoryEntity
	copier.Copy(&query, &ct)
	slice := make([]model.BaseNodeNo, 0)
	rt.Data = slice
	infos := c.sv.FindAll(ctx, query)
	if len(infos) > 0 {
		for _, item := range infos {
			var vo modBasicTagsCategory.Vo
			copier.Copy(&vo, &item)
			code := model.BaseNodeNo{
				Value:    item.No,
				Label:    item.Name,
				ParentNo: item.ParentNo,
				Extend:   vo,
			}

			slice = append(slice, code)
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
func (c *BasicTagsCategoryService) SelectPublic(ctx *gin.Context, ct modBasicTagsCategory.QueryCt) (rt rg.Rs[[]modBasicTagsCategory.Vo]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	var query entityBasic.BasicTagsCategoryEntity
	copier.Copy(&query, &ct)
	slice := make([]modBasicTagsCategory.Vo, 0)
	rt.Data = slice
	infos := c.sv.FindAll(ctx, query)
	if len(infos) > 0 {
		for _, item := range infos {
			var vo modBasicTagsCategory.Vo
			copier.Copy(&vo, &item)
			slice = append(slice, vo)
		}
		rt.Data = slice
	}
	return rt.Ok()
}

// ExportExcel 导出
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicTagsCategoryService) ExportExcel(ctx *gin.Context, ct modBasicTagsCategory.QueryCt) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	var query entityBasic.BasicTagsCategoryEntity
	copier.Copy(&query, &ct)
	infos := c.sv.FindAll(ctx, query)
	if len(infos) > 0 {
		slice := make([]interface{}, 0)
		for _, item := range infos {
			var vo modBasicTagsCategory.Vo
			copier.Copy(&vo, &item)
			slice = append(slice, vo)
		}
		log.Infof(ctx, log.TagAppDef, "导出数据 %+v", slice)
		strings := []string{"ID", "名称", "名称外文",
			"编号代号",
			"全称",
			"状态:1启用;2禁用",
			"删除:1是;2否",
			"描述",
			"创建时间",
			"更新时间",
			"创建人",
			"更新人", "组织id"}
		err := excelPg.ExportExcelByStruct(ctx, strings, slice, "department", "Sheet1")
		if nil != err {
			r := rg.Rs[string]{}
			ctx.JSON(200, r.ErrorMessage(err.Error()))
		}
	} else {
		r := rg.Rs[string]{}
		ctx.JSON(200, r.ErrorMessage("没有任何数据"))
	}

}

// ExistName 查重
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicTagsCategoryService) ExistName(ctx *gin.Context, ct model.BaseExistWdCt[string]) (rt rg.Rs[string]) {
	if "" == ct.Wd {
		return rt.ErrorMessage("名称不能为空")
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
func (c *BasicTagsCategoryService) ExistCode(ctx *gin.Context, ct model.BaseExistWdCt[string]) (rt rg.Rs[string]) {
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
