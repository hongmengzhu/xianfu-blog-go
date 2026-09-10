package service

import (
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/basic/modBasicTagsCategory"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/basic/modBasicTagsRelation"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBasic"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/holderPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/enumCommonPg/typeSysPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg/optionsPg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"

	"strings"

	"github.com/go-viper/mapstructure/v2"
	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/dbPg/pagePg"
	"github.com/pangu-2/go-tools/tools/jsonPg"
	"github.com/pangu-2/go-tools/tools/numberPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
)

func init() {
	gs.Provide(new(BasicTagsRelationService))
}

// BasicTagsRelationService 团队
// @Description:
type BasicTagsRelationService struct {
	sv          *repositoryBasic.BasicTagsRelationRepository `autowire:"?"`
	categoryRep *repositoryBasic.BasicTagsCategoryRepository `autowire:"?"`
}

// Create 新增
//
//	@Description:
//	@receiver c
//	@param ct
//	@return rt
func (c *BasicTagsRelationService) Create(ctx *gin.Context, ct modBasicTagsRelation.CreateCt) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	ct.No = strings.TrimSpace(ct.No)
	ct.Name = strings.TrimSpace(ct.Name)
	ct.NameShort = strings.TrimSpace(ct.NameShort)
	ct.Category = strings.TrimSpace(ct.Category)
	if strPg.IsBlank(ct.Name) {
		return rt.ErrorMessage("名称不能为空")
	}
	if strPg.IsBlank(ct.No) {
		return rt.ErrorMessage("标签不能为空")
	}
	if strPg.IsBlank(ct.Category) {
		return rt.ErrorMessage("分类不能为空")
	}
	if strPg.IsBlank(ct.TypeSys) {
		return rt.ErrorMessage("类型不能为空")
	}
	if _, ok := typeSysPg.IsExistTypeSys(ct.TypeSys); !ok {
		return rt.ErrorMessage("类型错误")
	}
	if _, result := c.categoryRep.FindByNo(ctx, ct.Category); !result {
		return rt.ErrorMessage("分类不存在")
	}
	if strPg.IsBlank(ct.NameShort) {
		ct.NameShort = ct.Name
	}
	var attributeVo modBasicTagsRelation.AttributeVo
	if nil != ct.AttributeMap {
		//map 转为 struct
		if err := mapstructure.Decode(ct.AttributeMap, &attributeVo); err != nil {
			log.Errorf(ctx, log.TagAppDef, "map 转 struct err=%+v", err)
		}
	}

	r := c.sv
	//不是自动
	_, result := r.FindByCodeAndIdNotAndCategoryNot(ctx, ct.No, 0, ct.Category)
	if result {
		return rt.ErrorMessage("标签已存在")
	}

	holder := holderPg.GetContextAccount(ctx)
	var info entityBasic.BasicTagsRelationEntity
	copier.Copy(&info, &ct)
	toJson, err2 := jsonPg.ObjToJson(attributeVo)
	if nil != err2 {
		log.Infof(ctx, log.TagAppDef, "jsonPg.ObjToJson err=%+v", err2)
		info.Attribute = "{}"
	} else {
		info.Attribute = toJson
	}

	log.Infof(ctx, log.TagAppDef, "info%+v", info)
	info.TenantNo = holder.GetTenantNo()
	err, _ := r.Create(ctx, &info)
	if nil != err {
		log.Errorf(ctx, log.TagAppDef, "save err=%+v", err)
		return rt.ErrorMessage("保存失败")
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
func (c *BasicTagsRelationService) Update(ctx *gin.Context, ct modBasicTagsRelation.UpdateCt) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	if ct.ID < 1 {
		return rt.ErrorMessage("id错误")
	}
	ct.No = strings.TrimSpace(ct.No)
	ct.Name = strings.TrimSpace(ct.Name)
	ct.NameShort = strings.TrimSpace(ct.NameShort)
	ct.Category = strings.TrimSpace(ct.Category)
	if strPg.IsBlank(ct.Name) {
		return rt.ErrorMessage("名称不能为空")
	}
	if strPg.IsBlank(ct.No) {
		return rt.ErrorMessage("标签不能为空")
	}
	if strPg.IsBlank(ct.Category) {
		return rt.ErrorMessage("分类不能为空")
	}
	if strPg.IsBlank(ct.TypeSys) {
		return rt.ErrorMessage("类型不能为空")
	}
	if _, ok := typeSysPg.IsExistTypeSys(ct.TypeSys); !ok {
		return rt.ErrorMessage("类型错误")
	}
	if _, result := c.categoryRep.FindByNo(ctx, ct.Category); !result {
		return rt.ErrorMessage("分类不存在")
	}
	if strPg.IsBlank(ct.NameShort) {
		ct.NameShort = ct.Name
	}
	r := c.sv

	_, result := r.FindByCodeAndIdNotAndCategoryNot(ctx, ct.No, ct.ID.ToInt64(), ct.Category)
	if result {
		return rt.ErrorMessage("标签已存在")
	}
	_, b := r.FindById(ctx, ct.ID.ToInt64())
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	log.Infof(ctx, log.TagAppDef, "AttributeMap=%+v", ct.AttributeMap)
	var attributeVo modBasicTagsRelation.AttributeVo
	if nil != ct.AttributeMap {
		//map 转为 struct
		if err := mapstructure.Decode(ct.AttributeMap, &attributeVo); err != nil {
			log.Errorf(ctx, log.TagAppDef, "map 转 struct err=%+v", err)
		}
		log.Infof(ctx, log.TagAppDef, "attributeVo=%+v", attributeVo)
	}
	var info entityBasic.BasicTagsRelationEntity
	copier.Copy(&info, &ct)
	toJson, err2 := jsonPg.ObjToJson(attributeVo)
	if nil != err2 {
		log.Infof(ctx, log.TagAppDef, "jsonPg.ObjToJson err=%+v", err2)
		info.Attribute = "{}"
	} else {
		info.Attribute = toJson
	}
	r.Update(ctx, info, info.ID)
	return rt.Ok()
}

// Detail 详情
//
//	@Description:
//	@receiver c
//	@param id
func (c *BasicTagsRelationService) Detail(ctx *gin.Context, id int64) (rt rg.Rs[modBasicTagsRelation.Vo]) {
	if id < 1 {
		return rt.ErrorMessage("id错误")
	}
	find, b := c.sv.FindById(ctx, id)
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	var info modBasicTagsRelation.Vo
	copier.Copy(&info, &find)
	return rt.OkData(info)
}

// Enable 启用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicTagsRelationService) Enable(ctx *gin.Context, ct model.BaseIdsCt[string]) (rt rg.Rs[string]) {
	return c.State(ctx, ct.Ids, enumStatePg.ENABLE)
}

// Disable 禁用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicTagsRelationService) Disable(ctx *gin.Context, ct model.BaseIdsCt[string]) (rt rg.Rs[string]) {
	return c.State(ctx, ct.Ids, enumStatePg.GetType(enumStatePg.DISABLE))
}

// State 状态 启用/禁用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicTagsRelationService) State(ctx *gin.Context, ids []string, state enumStatePg.State) (rt rg.Rs[string]) {
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
			r.Update(ctx, entityBasic.BasicTagsRelationEntity{State: state.IndexInt8()}, info.ID)
		}
	}
	return rt.Ok()
}

// StateEnableDisable 状态 设置 有效 停用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicTagsRelationService) StateEnableDisable(ctx *gin.Context, ids []string, state enumStatePg.State) (rt rg.Rs[string]) {
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
func (c *BasicTagsRelationService) LogicalDeletion(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
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
			log.Infof(ctx, log.TagAppDef, "id=%v,TenantId=%v", info.ID, info.TenantNo)
		}
		//不是系统类型情况下可以删除
		repository.DeleteByIdsStringAndTypeSysNot(ctx, ids, typeSysPg.System.String())
	} else {
		for _, info := range finds {
			enum := enumStatePg.State(info.State)
			// 有效 停用，反转 为对应的 取消 弃置
			if ok, reverse := enum.ReverseEnableDisable(); ok {
				repository.Update(ctx, entityBasic.BasicTagsRelationEntity{State: reverse.IndexInt8()}, info.ID)
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
func (c *BasicTagsRelationService) LogicalRecovery(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
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
			repository.Update(ctx, entityBasic.BasicTagsRelationEntity{State: reverse.IndexInt8()}, info.ID)
		}
	}
	return rt.Ok()
}

// PhysicalDeletion 物理删除
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicTagsRelationService) PhysicalDeletion(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
	if len(ids) < 1 {
		return rt.ErrorMessage("id错误")
	}
	cn := c.sv
	finds, b := cn.FindAllByIdStringIn(ctx, ids)
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	idsNew := make([]string, 0)
	for _, info := range finds {
		log.Infof(ctx, log.TagAppDef, "id=%v,TenantId=%v", info.ID, info.TenantNo)
		idsNew = append(idsNew, numberPg.Int64ToString(info.ID))
	}
	if len(idsNew) > 0 {
		//不是系统类型情况下可以删除
		cn.DeleteByIdsStringAndTypeSysNot(ctx, idsNew, typeSysPg.System.String())
	}
	return rt.Ok()
}

// Query 查询
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicTagsRelationService) Query(ctx *gin.Context, ct modBasicTagsRelation.QueryCt) (rt rg.Rs[pagePg.Paginator[modBasicTagsRelation.Vo]]) {
	var query entityBasic.BasicTagsRelationEntity
	copier.Copy(&query, &ct)
	slice := make([]modBasicTagsRelation.Vo, 0)
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
			arg.Db = arg.Db.Where("name like ?", "%"+ct.Wd+"%")
		}
	}), optionsPg.WithCtx(ctx))
	if nil != err {
		return rt.Ok()
	}

	if page.Total > 0 && page.Data != nil && len(page.Data) > 0 {

		pg := pagePg.NewPaginatorByPageable[modBasicTagsRelation.Vo](page.Pageable)
		//字段赋值
		for _, item := range page.Data {
			var vo modBasicTagsRelation.Vo
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
func (c *BasicTagsRelationService) SelectNodePublic(ctx *gin.Context, ct modBasicTagsRelation.QueryCt) (rt rg.Rs[[]model.BaseNode]) {
	var query entityBasic.BasicTagsRelationEntity
	copier.Copy(&query, &ct)
	slice := make([]model.BaseNode, 0)
	rt.Data = slice
	infos := c.sv.FindAll(ctx, query)
	if len(infos) > 0 {

		for _, item := range infos {
			slice = append(slice, model.BaseNode{Key: numberPg.Int64ToString(item.ID), Label: item.Name})
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
func (c *BasicTagsRelationService) SelectNodeAllPublic(ctx *gin.Context, ct modBasicTagsRelation.QueryCt) (rt rg.Rs[[]model.BaseNode]) {
	var query entityBasic.BasicTagsRelationEntity
	copier.Copy(&query, &ct)
	slice := make([]model.BaseNode, 0)
	rt.Data = slice
	infos := c.sv.FindAll(ctx, query)
	if len(infos) > 0 {

		for _, item := range infos {
			var vo modBasicTagsRelation.Vo
			copier.Copy(&vo, &item)
			slice = append(slice, model.BaseNode{Key: numberPg.Int64ToString(item.ID), Label: item.Name, Extend: vo})
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
func (c *BasicTagsRelationService) SelectPublic(ctx *gin.Context, ct modBasicTagsRelation.QueryCt) (rt rg.Rs[[]modBasicTagsRelation.Vo]) {
	var query entityBasic.BasicTagsRelationEntity
	copier.Copy(&query, &ct)
	rt.Data = []modBasicTagsRelation.Vo{}
	infos := c.sv.FindAll(ctx, query)
	if len(infos) > 0 {
		slice := make([]modBasicTagsRelation.Vo, 0)
		for _, item := range infos {
			var vo modBasicTagsRelation.Vo
			copier.Copy(&vo, &item)
			slice = append(slice, vo)
		}
		rt.Data = slice
	}
	return rt.Ok()
}

// All 查询
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicTagsRelationService) All(ctx *gin.Context, ct modBasicTagsRelation.AllCt) (rt rg.Rs[[]modBasicTagsRelation.AllVo]) {
	var query entityBasic.BasicTagsRelationEntity
	copier.Copy(&query, &ct)
	slice := make([]modBasicTagsRelation.AllVo, 0)
	rt.Data = slice
	infos := c.sv.FindAll(ctx, query)
	if len(infos) > 0 {
		for _, item := range infos {
			var vo modBasicTagsRelation.AllVo
			copier.Copy(&vo, &item)
			vo.AttributeMap = make(map[string]interface{})
			vo.Show = true
			//log.Infof(ctx, log.TagAppDef,"item.AttributeMap=%+v", item.Attribute)
			if strPg.IsNotBlank(item.Attribute) {
				err := json.Unmarshal([]byte(item.Attribute), &vo.AttributeMap)
				if err != nil {
					log.Errorf(ctx, log.TagAppDef, "json解析失败 %+v", err)
				}
				if obj, ok := vo.AttributeMap["color"]; ok {
					color := make(map[string]interface{})
					if strPg.IsNotBlank(obj.(string)) {
						err := json.Unmarshal([]byte(obj.(string)), &color)
						if err != nil {
							log.Errorf(ctx, log.TagAppDef, "json解析失败 %+v", err)
						}
					}
					vo.AttributeMap["color"] = color
				}
			}

			slice = append(slice, vo)
		}
		rt.Data = slice
	}
	return rt.Ok()
}

// AllByLink 查询包括所有子类
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicTagsRelationService) AllByLink(ctx *gin.Context, ct modBasicTagsRelation.AllCt) (rt rg.Rs[[]modBasicTagsRelation.AllVo]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	category := make([]string, 0)
	if strPg.IsNotBlank(ct.CategoryNo) {
		data, result := c.categoryRep.FindAllByNoLink(ctx, ct.CategoryNo)
		if result {
			for _, item := range data {
				category = append(category, item.No)
			}
		} else {
			category = append(category, ct.CategoryNo)
		}
		ct.CategoryNo = ""
	}
	var query entityBasic.BasicTagsRelationEntity
	copier.Copy(&query, &ct)
	slice := make([]modBasicTagsRelation.AllVo, 0)
	rt.Data = slice
	infos, b := c.sv.FindAllByCategoryNoIn(ctx, query, category)
	if b {
		for _, item := range infos {
			var vo modBasicTagsRelation.AllVo
			copier.Copy(&vo, &item)
			vo.AttributeMap = make(map[string]interface{})
			vo.Show = true
			//log.Infof(ctx, log.TagAppDef,"item.AttributeMap=%+v", item.Attribute)
			if strPg.IsNotBlank(item.Attribute) {
				err := json.Unmarshal([]byte(item.Attribute), &vo.AttributeMap)
				if err != nil {
					log.Errorf(ctx, log.TagAppDef, "json解析失败 %+v", err)
				}
				if obj, ok := vo.AttributeMap["color"]; ok {
					color := make(map[string]interface{})
					if strPg.IsNotBlank(obj.(string)) {
						err := json.Unmarshal([]byte(obj.(string)), &color)
						if err != nil {
							log.Errorf(ctx, log.TagAppDef, "json解析失败 %+v", err)
						}
					}
					vo.AttributeMap["color"] = color
				}
			}

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
func (c *BasicTagsRelationService) ExistName(ctx *gin.Context, ct modBasicTagsRelation.ExistWdCt) (rt rg.Rs[string]) {
	if "" == ct.Wd {
		return rt.ErrorMessage("关键词不能为空")
	}
	if strPg.IsBlank(ct.Category) {
		return rt.ErrorMessage("分类不能为空")
	}
	_, result := c.sv.FindByNameAndIdNotAndCategoryNot(ctx, ct.Wd, numberPg.StrToInt64(ct.Id), ct.Category)
	if result {
		return rt.ErrorMessage("重复，已存在")
	}
	if strPg.IsBlank(ct.Category) {
		return rt.ErrorMessage("分类不能为空")
	}
	return rt.OkMessage("可以使用")
}

// ExistCode 查重
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicTagsRelationService) ExistCode(ctx *gin.Context, ct modBasicTagsRelation.ExistWdCt) (rt rg.Rs[string]) {
	if "" == ct.Wd {
		return rt.ErrorMessage("关键词不能为空")
	}
	if strPg.IsBlank(ct.Category) {
		return rt.ErrorMessage("分类不能为空")
	}
	_, result := c.sv.FindByCodeAndIdNotAndCategoryNot(ctx, ct.Wd, numberPg.StrToInt64(ct.Id), ct.Category)
	if result {
		return rt.ErrorMessage("重复，已存在")
	}
	return rt.OkMessage("可以使用")
}

// GetCategory
//
//	@Description: 获取下 所有 分类
//	@receiver c
func (c *BasicTagsRelationService) GetCategory(ctx *gin.Context, category string) (rt rg.Rs[modBasicTagsRelation.GroupVo]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", category)
	var groupVo modBasicTagsRelation.GroupVo
	groupVo.General = make([]modBasicTagsCategory.Vo, 0)
	groupVo.Sys = make([]modBasicTagsCategory.Vo, 0)
	// 获取所有分类
	infos, result := c.categoryRep.FindAllByNoLink(ctx, category)
	if result {
		//整理 所有普通分类
		for _, item := range infos {
			var vo modBasicTagsCategory.Vo
			copier.Copy(&vo, &item)
			//整理 所有普通分类
			if item.TypeSys == typeSysPg.System.String() {
				//系统分类
				groupVo.Sys = append(groupVo.Sys, vo)
				//如果这个分类是 一级分类 那么 添加到 slice
				if item.ParentNo == category {
					//系统分类
					groupVo.General = append(groupVo.General, vo)
				} else if item.Code == category {
					//系统分类
					groupVo.General = append(groupVo.General, vo)
				}
			} else {
				//整理 所有普通分类
				groupVo.General = append(groupVo.General, vo)
			}
		}
	}
	return rt.OkData(groupVo)
}

// GetCategoryTagsAll
//
//	@Description: 获取指定分类下标签
//	@receiver c
//	@param ctx
//	@param categoryRoot 根分类
//	@param category 分类
//	@return rt
func (c *BasicTagsRelationService) GetCategoryTagsAll(ctx *gin.Context, categoryRoot string, ct modBasicTagsRelation.QueryCt) (rt rg.Rs[[]modBasicTagsRelation.AllVo]) {
	slice := make([]modBasicTagsRelation.AllVo, 0)
	rt.Data = slice
	var infos []*entityBasic.BasicTagsRelationEntity
	//根分类条件
	tx := c.sv.DbModel().Where("category_root = ?", categoryRoot)
	//指定分类
	if strPg.IsNotBlank(ct.CategoryNo) {
		tx = tx.Where("category_no = ?", ct.CategoryNo)
	}
	//租户条件
	holder := holderPg.GetContextAccount(ctx)
	tx.Where("tenant_no = ?", holder.GetTenantNo())
	tx = tx.Find(&infos)
	if tx.Error != nil {
		return rt.Ok()
	}
	if 0 == tx.RowsAffected {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return rt.Ok()
	}
	for _, item := range infos {
		var vo modBasicTagsRelation.AllVo
		copier.Copy(&vo, &item)
		vo.AttributeMap = make(map[string]interface{})
		vo.Show = true
		if strPg.IsNotBlank(item.Attribute) {
			err := json.Unmarshal([]byte(item.Attribute), &vo.AttributeMap)
			if err != nil {
				log.Errorf(ctx, log.TagAppDef, "json解析失败 %+v", err)
			}
			if obj, ok := vo.AttributeMap["color"]; ok {
				color := make(map[string]interface{})
				if strPg.IsNotBlank(obj.(string)) {
					err := json.Unmarshal([]byte(obj.(string)), &color)
					if err != nil {
						log.Errorf(ctx, log.TagAppDef, "json解析失败 %+v", err)
					}
				}
				vo.AttributeMap["color"] = color
			}
		}

		slice = append(slice, vo)
	}
	return rt.OkData(slice)
}

// GetCategoryTags
//
//	@Description: 获取指定分类下标签
//	@receiver c
//	@param ctx
//	@param categoryRoot 根分类
//	@param category 分类
//	@return rt
func (c *BasicTagsRelationService) GetCategoryTags(ctx *gin.Context, categoryRoot string, ct modBasicTagsRelation.QueryCt) (rt rg.Rs[pagePg.Paginator[modBasicTagsRelation.AllVo]]) {
	var query entityBasic.BasicTagsRelationEntity
	copier.Copy(&query, &ct)
	//根分类
	query.CategoryRoot = categoryRoot
	//商户条件
	holder := holderPg.GetContextAccount(ctx)
	//
	slice := make([]modBasicTagsRelation.AllVo, 0)
	rt.Data.Data = slice
	r := c.sv
	page, err := r.FindAllPage(ctx, query, optionsPg.WithOption(func(arg *optionsPg.OptionParams) {
		if ct.PageSize < 1 {
			ct.PageSize = 20
		}
		arg.Pageable = new(pagePg.PageablePageSize(0, ct.PageNum, ct.PageSize))
		//自定义查询
		arg.Db = arg.Db.Order("category_no,name ASC,create_at desc")
		if strPg.IsNotBlank(ct.Wd) {
			arg.Db = arg.Db.Where("name like ?", "%"+ct.Wd+"%").
				Or("name_short like ?", "%"+ct.Wd+"%").
				Or("code like ?", "%"+ct.Wd+"%").
				Or("name_full like ?", "%"+ct.Wd+"%")
		}
		//名称
		if strPg.IsNotBlank(ct.Name) {
			arg.Db = arg.Db.Where("name like ?", "%"+ct.Name+"%")
		}
		//商户条件 或者系统分类
		arg.Db = arg.Db.Where(
			r.Db().Where("tenant_no = ?", holder.GetTenantNo()).
				Or("type_sys = ?", typeSysPg.System.String()))

	}), optionsPg.WithCtx(ctx))
	if nil != err {
		return rt.Ok()
	}
	if page.Total > 0 && page.Data != nil && len(page.Data) > 0 {
		pg := pagePg.NewPaginatorByPageable[modBasicTagsRelation.AllVo](page.Pageable)
		//字段赋值
		for _, item := range page.Data {
			var vo modBasicTagsRelation.AllVo
			copier.Copy(&vo, &item)
			vo.Show = true
			vo.AttributeMap = make(map[string]interface{})
			vo.AttributeMap["bordered"] = true
			vo.AttributeMap["type"] = "default"
			vo.AttributeMap["color"] = struct {
			}{}
			vo.AttributeMap["strong"] = false
			vo.AttributeMap["round"] = false
			if strPg.IsNotBlank(item.Attribute) {
				err := json.Unmarshal([]byte(item.Attribute), &vo.AttributeMap)
				if err != nil {
					log.Errorf(ctx, log.TagAppDef, "json解析失败 %+v", err)
				}
				if obj, ok := vo.AttributeMap["color"]; ok {
					color := make(map[string]interface{})
					if strPg.IsNotBlank(obj.(string)) {
						err := json.Unmarshal([]byte(obj.(string)), &color)
						if err != nil {
							log.Errorf(ctx, log.TagAppDef, "json解析失败 %+v", err)
						}
					}
					vo.AttributeMap["color"] = color
				}
			}
			if strPg.IsBlank(vo.AttributeMap["type"].(string)) {
				vo.AttributeMap["type"] = "default"
			}
			//
			//
			slice = append(slice, vo)
		}
		pg.Data = slice
		pg.Pageable = page.Pageable
		rt.Data = pg
	}
	return rt.Ok()
}
