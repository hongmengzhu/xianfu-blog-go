package service

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/basic/modBasicConfigEventFields"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBasic"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/enumCommonPg/configModelPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg/optionsPg"
	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/cryptPg"
	"github.com/pangu-2/go-tools/tools/dbPg/pagePg"
	"github.com/pangu-2/go-tools/tools/noPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
	"gorm.io/datatypes"
)

func init() {
	gs.Provide(new(BasicConfigEventFieldsService))
}

// BasicConfigEventFieldsService 省市区
// @Description:
type BasicConfigEventFieldsService struct {
	sv    *repositoryBasic.BasicConfigEventFieldsRepository `autowire:"?"`
	event *repositoryBasic.BasicConfigEventRepository       `autowire:"?"`
}

// CreateUpdate 新增
//
//	@Description:
//	@receiver c
//	@param ct
//	@return rt
func (c *BasicConfigEventFieldsService) CreateUpdate(ctx *gin.Context, ct modBasicConfigEventFields.CreateUpdateCt) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%#v", ct)
	if nil == ct.Body || len(ct.Body) < 1 {
		return rt.ErrorMessage("表体字段不能为空")
	}
	if strPg.IsBlank(ct.EventNo) {
		return rt.ErrorMessage("事件编号不能为空")
	}
	event, result := c.event.FindByNo(ctx, ct.EventNo)
	if !result {
		return rt.ErrorMessage("事件 不存在")
	}
	if nil != ct.BodyDelIds || len(ct.BodyDelIds) > 0 {
		delIds := make([]string, 0)
		for _, id := range ct.BodyDelIds {
			if strPg.IsNotBlank(id) {
				delIds = append(delIds, strings.TrimSpace(id))
			}
		}
		if len(delIds) > 0 {
			c.sv.DeleteAllByModelNoAndIds(ctx, event.No, delIds)
		}

	}
	dataAdd := make([]*entityBasic.BasicConfigEventFieldsEntity, 0)
	dataUpdate := make([]*entityBasic.BasicConfigEventFieldsEntity, 0)
	findField := make(map[string]bool)
	for _, item := range ct.Body {
		//
		if strPg.IsBlank(item.Name) {
			return rt.ErrorMessage("名称不能为空")
		}
		if strPg.IsBlank(item.Field) {
			return rt.ErrorMessage(item.Name + " 字段不能为空")
		}
		if _, ok := findField[strings.TrimSpace(item.Field)]; ok {
			return rt.ErrorMessage(item.Field + " 字段重复")
		}
		findField[strings.TrimSpace(item.Field)] = true
	}
	for _, item := range ct.Body {
		parameterConfig := entityBasic.BasicConfigEventFieldsJsonParameterConfig{}
		obj := entityBasic.BasicConfigEventFieldsEntity{}
		err := copier.Copy(&obj, &item)
		if err != nil {
			log.Infof(ctx, log.TagAppDef, "copier.Copy error: %+v", err)
		}
		obj.Name = strings.TrimSpace(obj.Name)
		obj.Field = strings.TrimSpace(obj.Field)
		//
		tags := make([]string, 0)
		if nil != item.Rules && len(item.Rules) > 0 {
			for _, v := range item.Rules {
				if strPg.IsNotBlank(v) {
					tags = append(tags, strings.TrimSpace(v))
				}
			}
		}
		obj.Rules = datatypes.NewJSONType[[]string](tags)
		//数据字典
		if configModelPg.ParameterSourceDataDictionary.Index() == item.ParameterSource {
			if strPg.IsNotBlank(item.ParameterConfigDataDictionary) {
				parameterConfig.DataDictionary = strings.TrimSpace(item.ParameterConfigDataDictionary)
			}
		}
		obj.ParameterConfig = datatypes.NewJSONType(parameterConfig)
		//
		if item.Id.ToInt64() > 0 {
			obj.No = ""
			obj.EventNo = event.No
			obj.ModelNo = event.ModelNo
			obj.Model = event.Model
			obj.Module = event.Module
			obj.ModuleSub = event.ModuleSub
			obj.TenantNo = event.TenantNo
			obj.KindUnique = cryptPg.Md5(obj.EventNo + obj.Field)
			dataUpdate = append(dataUpdate, &obj)
		} else {
			obj.ID = 0
			obj.No = noPg.No()
			obj.State = enumStatePg.ENABLE.IndexInt8()
			obj.EventNo = event.No
			obj.ModelNo = event.ModelNo
			obj.Model = event.Model
			obj.Module = event.Module
			obj.ModuleSub = event.ModuleSub
			obj.TenantNo = event.TenantNo
			obj.Sort = 0
			obj.KindUnique = cryptPg.Md5(obj.EventNo + obj.Field)
			dataAdd = append(dataAdd, &obj)
		}
	}
	if len(dataAdd) > 0 {
		tx := c.sv.DbModel().CreateInBatches(dataAdd, 1000000)
		if tx.Error != nil {
			log.Errorf(ctx, log.TagAppDef, "save err=%+v", tx.Error)
			return rt.ErrorMessage("保存失败：")
		}
	}
	if len(dataUpdate) > 0 {
		for _, entity := range dataUpdate {
			c.sv.Update(ctx, *entity, entity.ID)
		}
	}
	return rg.OkData("成功")
}

// Detail 详情
//
//	@Description:
//	@receiver c
//	@param id
func (c *BasicConfigEventFieldsService) Detail(ctx *gin.Context, id int64) (rt rg.Rs[modBasicConfigEventFields.Vo]) {
	if id < 1 {
		return rt.ErrorMessage("id错误")
	}
	find, b := c.sv.FindById(ctx, id)
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	var info modBasicConfigEventFields.Vo
	copier.Copy(&info, &find)
	return rt.OkData(info)
}

// Enable 启用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicConfigEventFieldsService) Enable(ctx *gin.Context, ct model.BaseIdsCt[string]) (rt rg.Rs[string]) {
	return c.State(ctx, ct.Ids, enumStatePg.ENABLE)
}

// Disable 禁用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicConfigEventFieldsService) Disable(ctx *gin.Context, ct model.BaseIdsCt[string]) (rt rg.Rs[string]) {
	return c.State(ctx, ct.Ids, enumStatePg.GetType(enumStatePg.DISABLE))
}

// State 状态 启用/禁用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicConfigEventFieldsService) State(ctx *gin.Context, ids []string, state enumStatePg.State) (rt rg.Rs[string]) {
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
			r.Update(ctx, entityBasic.BasicConfigEventFieldsEntity{State: state.IndexInt8()}, info.ID)
		}
	}
	return rt.Ok()
}

// StateEnableDisable 状态 设置 有效 停用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicConfigEventFieldsService) StateEnableDisable(ctx *gin.Context, ids []string, state enumStatePg.State) (rt rg.Rs[string]) {
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
func (c *BasicConfigEventFieldsService) LogicalDeletion(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
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
			log.Infof(ctx, log.TagAppDef, "id=%v", info.ID)
		}
		repository.DeleteByIdsString(ctx, ids)
	} else {
		for _, info := range finds {
			enum := enumStatePg.State(info.State)
			// 有效 停用，反转 为对应的 取消 弃置
			if ok, reverse := enum.ReverseEnableDisable(); ok {
				repository.Update(ctx, entityBasic.BasicConfigEventFieldsEntity{State: reverse.IndexInt8()}, info.ID)
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
func (c *BasicConfigEventFieldsService) LogicalRecovery(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
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
			repository.Update(ctx, entityBasic.BasicConfigEventFieldsEntity{State: reverse.IndexInt8()}, info.ID)
		}
	}
	return rt.Ok()
}

// PhysicalDeletion 物理删除
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicConfigEventFieldsService) PhysicalDeletion(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
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
func (c *BasicConfigEventFieldsService) Query(ctx *gin.Context, ct modBasicConfigEventFields.QueryCt) (rt rg.Rs[pagePg.Paginator[modBasicConfigEventFields.Vo]]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	var query entityBasic.BasicConfigEventFieldsEntity
	copier.Copy(&query, &ct)
	slice := make([]modBasicConfigEventFields.Vo, 0)
	rt.Data.Data = slice
	r := c.sv
	page, err := r.FindAllPage(ctx, query, optionsPg.WithOption(func(arg *optionsPg.OptionParams) {
		if ct.PageSize < 1 {
			ct.PageSize = 20
		}
		arg.Pageable = new(pagePg.PageablePageSize(0, ct.PageNum, ct.PageSize))
		//排序
		arg.Db = arg.Db.Order("create_at asc")
		if strPg.IsNotBlank(ct.Wd) {
			arg.Db = arg.Db.Where("name like ?", "%"+ct.Wd+"%")
		}
	}), optionsPg.WithCtx(ctx))
	if nil != err {
		return rt.Ok()
	}

	if page.Total > 0 && page.Data != nil && len(page.Data) > 0 {
		pg := pagePg.NewPaginatorByPageable[modBasicConfigEventFields.Vo](page.Pageable)
		//字段赋值
		for _, item := range page.Data {
			var vo modBasicConfigEventFields.Vo
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

// AllByEventNo 查询
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicConfigEventFieldsService) AllByEventNo(ctx *gin.Context, ct modBasicConfigEventFields.QueryPublicCt) (rt rg.Rs[[]modBasicConfigEventFields.Vo]) {

	var query entityBasic.BasicConfigEventFieldsEntity
	copier.Copy(&query, &ct)
	slice := make([]modBasicConfigEventFields.Vo, 0)
	rt.Data = slice
	//
	if strPg.IsBlank(ct.EventNo) {
		return rt.ErrorMessage("参数错误")
	}
	//
	infos := c.sv.FindAll(ctx, query)
	if len(infos) > 0 {
		for _, item := range infos {
			var vo modBasicConfigEventFields.Vo
			copier.Copy(&vo, &item)
			//
			obj := item.ParameterConfig.Data()
			vo.ParameterConfigDataDictionary = obj.DataDictionary
			//
			slice = append(slice, vo)
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
func (c *BasicConfigEventFieldsService) SelectNodeAllPublic(ctx *gin.Context, ct modBasicConfigEventFields.QueryPublicCt) (rt rg.Rs[[]model.BaseNodeNo]) {
	var query entityBasic.BasicConfigEventFieldsEntity
	copier.Copy(&query, &ct)
	slice := make([]model.BaseNodeNo, 0)
	rt.Data = slice
	infos := c.sv.FindAll(ctx, query)
	if len(infos) > 0 {
		for _, item := range infos {
			var vo modBasicConfigEventFields.Vo
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
func (c *BasicConfigEventFieldsService) ExistName(ctx *gin.Context, ct model.BaseExistWdCt[string]) (rt rg.Rs[string]) {
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
func (c *BasicConfigEventFieldsService) ExistCode(ctx *gin.Context, ct model.BaseExistWdCt[string]) (rt rg.Rs[string]) {
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
