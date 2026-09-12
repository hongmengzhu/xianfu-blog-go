package service

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamResourceGroup"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityRam"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/automatedPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constNodePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constsRam/resourceTypeCategoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constsRam/typeAttrPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg/optionsPg"
	"github.com/pangu-2/go-tools/tools/noPg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"

	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/dbPg/pagePg"
	"github.com/pangu-2/go-tools/tools/numberPg"
	"github.com/pangu-2/go-tools/tools/slicePg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
)

func init() {
	gs.Provide(new(RamResourceGroupService))
}

// RamResourceGroupService 资源组
// @Description:
type RamResourceGroupService struct {
	sv     *repositoryRam.RamResourceGroupRepository     `autowire:"?"`
	authDb *repositoryRam.RamResourceAuthorityRepository `autowire:"?"`
}

// Create 新增
//
//	@Description:
//	@receiver c
//	@param ct
//	@return rt
func (c *RamResourceGroupService) Create(ctx *gin.Context, ct modRamResourceGroup.CreateCt) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	var info entityRam.RamResourceGroupEntity
	err2 := copier.Copy(&info, &ct)
	if err2 != nil {
		log.Errorf(ctx, log.TagAppDef, "copier.Copy=%+v", err2)
		return rt.ErrorMessage(err2.Error())
	}
	if "" == ct.Name {
		return rt.ErrorMessage("名称不能为空")
	}
	if strPg.IsBlank(ct.Code) {
		info.Code = automatedPg.CREATE_CODE
	}
	//holder := holderPg.GetContextAccount(ctx)
	parent := &entityRam.RamResourceGroupEntity{}
	r := c.sv
	//判断是否是自动,不是自动
	if !automatedPg.IsCreateCode(info.Code) {
		//判断格式是否满足要求
		if !automatedPg.FormatVerify(info.Code) {
			return rt.ErrorMessage("标志格式不能为空")
		}
		//不是自动
		_, result := r.FindByCode(ctx, info.Code, optionsPg.WithCtx(ctx))
		if result {
			return rt.ErrorMessage("标志已存在")
		}
	}
	//默认资源组
	if strPg.IsBlank(ct.TypeCategory) {
		ct.TypeCategory = resourceTypeCategoryPg.Group.String()
	}
	//默认 资源
	if strPg.IsBlank(ct.TypeAttr) {
		ct.TypeAttr = typeAttrPg.Resource.String()
	}
	if strPg.IsNotBlank(ct.ParentNo) {
		result := false
		parent, result = r.FindByNo(ctx, ct.ParentNo, optionsPg.WithCtx(ctx))
		if !result {
			return rt.ErrorMessage("上级不存在")
		}
		if strPg.IsNotBlank(parent.ParentNo) {
			return rt.ErrorMessage("只允许2级存在")
		}
		if !typeAttrPg.Category.IsEqual(parent.TypeAttr) {
			return rt.ErrorMessage("上级属性为分类属性时,才能有下级")
		}
		if !typeAttrPg.CategoryLast.IsEqual(parent.TypeAttr) {
			return rt.ErrorMessage("当前属性为分类末级属性时,才可以关联资源权限")
		}
		info.ParentId = numberPg.Int64ToString(parent.ID)
		info.ParentNo = parent.No
	} else {
		if !typeAttrPg.Category.IsEqual(ct.TypeAttr) {
			return rt.ErrorMessage("顶级属性为分类属性")
		}
	}
	//
	info.No = noPg.No()
	info.ID = 0
	info.State = enumStatePg.ENABLE.Index()
	log.Infof(ctx, log.TagAppDef, "info=%+v", info)
	err, _ := r.Create(ctx, &info)
	if err != nil {
		return rt.ErrorMessage("保存失败 " + err.Error())
	}
	//设置上级 link
	if strPg.IsNotBlank(ct.ParentNo) {
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
	err = r.Update(ctx, info, info.ID)
	if err != nil {
		return rt.ErrorMessage(err.Error())
	}
	log.Infof(ctx, log.TagAppDef, "save=%+v", info)
	return rg.OkData(numberPg.Int64ToString(info.ID))
}

// Update 更新
//
//	@Description:
//	@receiver c
//	@param ct
//	@return rt
func (c *RamResourceGroupService) Update(ctx *gin.Context, ct modRamResourceGroup.CreateUpdateCt) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	var info entityRam.RamResourceGroupEntity
	err := copier.Copy(&info, &ct)
	if err != nil {
		return rt.ErrorMessage(err.Error())
	}
	if ct.ID < 1 {
		return rt.ErrorMessage("id错误")
	}
	if "" == ct.Name {
		return rt.ErrorMessage("名称不能为空")
	}
	if strPg.IsBlank(ct.TypeAttr) {
		return rt.ErrorMessage("属性不能为空")
	}
	r := c.sv
	find, b := r.FindById(ctx, ct.ID.ToInt64())
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	//上级
	parent := &entityRam.RamResourceGroupEntity{}
	var childData []*entityRam.RamResourceGroupEntity
	if strPg.IsNotBlank(ct.ParentNo) {
		result := false
		parent, result = r.FindByNo(ctx, ct.ParentNo, optionsPg.WithCtx(ctx))
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
					if item.ID == parent.ID {
						return rt.ErrorMessage("无法保存，不能设置为自己的子集")
					}
				}
			}
		}
	}
	// 如果是分类
	if typeAttrPg.CategoryLast.IsEqual(ct.TypeAttr) {
		if typeAttrPg.Category.IsEqual(parent.TypeAttr) {
			return rt.ErrorMessage("上级属性应该为分类")
		}
	} else if typeAttrPg.Resource.IsEqual(ct.TypeAttr) {
		//如果是资源
		return rt.ErrorMessage("属性不能是资源属性")
	}
	//
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
	err = r.Update(ctx, info, info.ID)
	if err != nil {
		log.Errorf(ctx, log.TagAppDef, "update error=%+v", err)
		return rt.ErrorMessage(err.Error())
	}
	log.Infof(ctx, log.TagAppDef, "save.info=%+v", info)
	//更改上级后，相关子集修改
	if strPg.IsNotBlank(ct.ParentNo) && nil != childData {
		maps := slicePg.ToMapArray(childData, func(t *entityRam.RamResourceGroupEntity) (string, *entityRam.RamResourceGroupEntity) {
			if len(t.ParentId) == 0 {
				return constNodePg.ROOT, t
			}
			return t.ParentId, t
		})
		for _, item := range maps[numberPg.Int64ToString(find.ID)] {
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
				err = r.Update(ctx, entityRam.RamResourceGroupEntity{IdLink: item.IdLink}, item.ID)
				if err != nil {
					return rt.ErrorMessage(err.Error())
				}
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
func (c *RamResourceGroupService) childParentIdLink(maps map[string][]*entityRam.RamResourceGroupEntity, parent *entityRam.RamResourceGroupEntity) {
	entities := maps[numberPg.Int64ToString(parent.ID)]
	for _, item := range entities {
		item.IdLink = constNodePg.NoLinkAssemble(parent.IdLink, numberPg.Int64ToString(item.ID))
		item.NoLink = constNodePg.NoLinkAssemble(parent.NoLink, item.No)
	}
}

// CacheOverride 缓存重载
//
//	@Description:
//	@receiver c
func (c *RamResourceGroupService) CacheOverride(ctx *gin.Context) {
	r := c.sv
	infos, b := r.FindAllData(ctx)
	if !b {
		return
	}
	maps := slicePg.ToMapArray(infos, func(t *entityRam.RamResourceGroupEntity) (string, *entityRam.RamResourceGroupEntity) {
		if len(t.ParentId) == 0 {
			return constNodePg.ROOT, t
		}
		return t.ParentId, t
	})
	for _, item := range maps[constNodePg.ROOT] {
		item.IdLink = constNodePg.NoLinkDefault(numberPg.Int64ToString(item.ID))
		item.NoLink = constNodePg.NoLinkDefault(item.No)
		c.childParentIdLink(maps, item)
	}
	log.Infof(ctx, log.TagAppDef, "maps=%+v", maps)
	for _, val := range maps {
		for _, item := range val {
			r.Update(ctx, entityRam.RamResourceGroupEntity{IdLink: item.IdLink, NoLink: item.NoLink}, item.ID)
		}
	}
	maps = nil
}

// PhysicalDeletion 物理删除
//
//	@Description:
//	@receiver c
//	@param ct
func (c *RamResourceGroupService) PhysicalDeletion(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ids)
	if len(ids) < 1 {
		return rt.ErrorMessage("id错误")
	}
	r := c.sv
	authDb := c.authDb
	finds, b := r.FindAllByIdStringIn(ctx, ids, optionsPg.WithCtx(ctx))
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	idsNew := make([]int64, 0)
	idsCategory := make([]string, 0)
	for _, info := range finds {
		log.Infof(ctx, log.TagAppDef, "id=%v,TenantId=%v,TypeAttr=%+v", info.ID, 0, info.TypeAttr)
		//判断 是否是 分类
		if typeAttrPg.CategoryLast.IsEqual(info.TypeAttr) || typeAttrPg.Category.IsEqual(info.TypeAttr) {
			link, b := r.FindAllByIdLink(ctx, numberPg.Int64ToString(info.ID))
			if b {
				for _, entity := range link {
					idsCategory = append(idsCategory, numberPg.Int64ToString(entity.ID))
				}
				if len(idsCategory) > 0 {
					_, result := authDb.FindAllByTypeCategoryAndGroupIdStringIn(ctx, resourceTypeCategoryPg.Group.String(), idsCategory)
					if result {
						return rt.ErrorMessage("该分类下存在子数据，请先删除子数据")
					}
				}

			}
		}
		idsNew = append(idsNew, info.ID)
	}
	if len(idsNew) > 0 {
		r.DeleteByIds(ctx, idsNew)
	}
	return rt.Ok()
}

// Detail 详情
//
//	@Description:
//	@receiver c
//	@param id
func (c *RamResourceGroupService) Detail(ctx *gin.Context, id int64) (rt rg.Rs[modRamResourceGroup.Vo]) {
	if id < 1 {
		return rt.ErrorMessage("id错误")
	}
	find, b := c.sv.FindById(ctx, id)
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	var info modRamResourceGroup.Vo
	copier.Copy(&info, &find)
	return rt.OkData(info)
}

// Query 查询
//
//	@Description:
//	@receiver c
//	@param ct
func (c *RamResourceGroupService) Query(ctx *gin.Context, ct modRamResourceGroup.QueryCt) (rt rg.Rs[pagePg.Paginator[modRamResourceGroup.Vo]]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	var query entityRam.RamResourceGroupEntity
	copier.Copy(&query, &ct)
	slice := make([]modRamResourceGroup.Vo, 0)
	rt.Data.Data = slice
	r := c.sv
	page, err := r.FindAllPage(ctx, query, optionsPg.WithOption(func(arg *optionsPg.OptionParams) {
		if ct.PageSize < 1 {
			ct.PageSize = 20
		}
		arg.Pageable = new(pagePg.PageablePageSize(0, ct.PageNum, ct.PageSize))
		arg.Db = arg.Db.Order("create_at asc")
		//自定义查询
		if strPg.IsNotBlank(ct.Wd) {
			arg.Db = arg.Db.Where("name like ?", "%"+ct.Wd+"%")
		}
	}), optionsPg.WithCtx(ctx))
	if nil != err {
		return rt.Ok()
	}

	if page.Total > 0 && page.Data != nil && len(page.Data) > 0 {

		pg := pagePg.NewPaginatorByPageable[modRamResourceGroup.Vo](page.Pageable)
		//字段赋值
		for _, item := range page.Data {
			var vo modRamResourceGroup.Vo
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
func (c *RamResourceGroupService) SelectNodePublic(ctx *gin.Context, ct modRamResourceGroup.QueryPublicCt) (rt rg.Rs[[]model.BaseNodeNo]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	var query entityRam.RamResourceGroupEntity
	copier.Copy(&query, &ct)
	slice := make([]model.BaseNodeNo, 0)
	rt.Data = slice
	infos := c.sv.FindAll(ctx, query)
	if len(infos) > 0 {
		for _, item := range infos {
			var vo modRamResourceGroup.Vo
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
func (c *RamResourceGroupService) SelectNodeAllPublic(ctx *gin.Context, ct modRamResourceGroup.QueryPublicCt) (rt rg.Rs[[]model.BaseNodeNo]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	var query entityRam.RamResourceGroupEntity
	copier.Copy(&query, &ct)
	slice := make([]model.BaseNodeNo, 0)
	rt.Data = slice
	infos := c.sv.FindAll(ctx, query)
	if len(infos) > 0 {
		for _, item := range infos {
			var vo modRamResourceGroup.Vo
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
func (c *RamResourceGroupService) SelectPublic(ctx *gin.Context, ct modRamResourceGroup.QueryCt) (rt rg.Rs[[]modRamResourceGroup.Vo]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	var query entityRam.RamResourceGroupEntity
	copier.Copy(&query, &ct)
	slice := make([]modRamResourceGroup.Vo, 0)
	rt.Data = slice
	infos := c.sv.FindAll(ctx, query)
	if len(infos) > 0 {
		for _, item := range infos {
			var vo modRamResourceGroup.Vo
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
func (c *RamResourceGroupService) ExistName(ctx *gin.Context, ct model.BaseExistWdCt[string]) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	if "" == ct.Wd {
		return rt.ErrorMessage("查询内容不能为空")
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
func (c *RamResourceGroupService) ExistCode(ctx *gin.Context, ct model.BaseExistWdCt[string]) (rt rg.Rs[string]) {
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
