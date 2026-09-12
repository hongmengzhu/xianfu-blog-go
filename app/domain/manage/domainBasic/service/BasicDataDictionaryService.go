package service

import (
	"strings"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/duke-git/lancet/v2/strutil"
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/basic/modBasicDataDictionary"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBasic"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg/optionsPg"
	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/cryptPg"
	"github.com/pangu-2/go-tools/tools/dbPg/pagePg"
	"github.com/pangu-2/go-tools/tools/noPg"
	"github.com/pangu-2/go-tools/tools/numberPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
	"gorm.io/gorm"
)

func init() {
	gs.Provide(new(BasicDataDictionaryService))
}

// BasicDataDictionaryService 数据字典
// @Description:
type BasicDataDictionaryService struct {
	sv *repositoryBasic.BasicDataDictionaryRepository `autowire:"?"`
}

// CreateUpdate 新增
//
//	@Description:
//	@receiver c
//	@param ct
//	@return rt
func (c *BasicDataDictionaryService) CreateUpdate(ctx *gin.Context, ct modBasicDataDictionary.CreateUpdateCt) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%#v", ct)
	if "" == ct.Name {
		return rt.ErrorMessage("名称不能为空")
	}
	if strPg.IsBlank(ct.Code) {
		return rt.ErrorMessage("码值不能为空")
	}
	r := c.sv
	var info entityBasic.BasicDataDictionaryEntity
	copier.Copy(&info, &ct)
	//
	info.Code = strings.Trim(info.Code, "")
	info.Value = info.Code
	info.TypeUniqueMd5 = cryptPg.Md5(info.Code)
	if nil != ct.Range && len(ct.Range) > 0 {
		info.Range = slice.Join(ct.Range, ",")
	}
	info.OwnerNo = "0"
	if ct.ID < 1 {
		_, result := c.sv.FindByCodeAndIdNotAndOwnerNo(ctx, info.Code, "0", info.OwnerNo)
		if result {
			return rt.ErrorMessage("码值已存在")
		}
		log.Infof(ctx, log.TagAppDef, "info%+v", info)
		info.No = noPg.No()
		err, _ := r.Create(ctx, &info)
		if err != nil {
			return rt.ErrorMessage("保存失败")
		}
		log.Infof(ctx, log.TagAppDef, "save=%+v", info)
	} else {
		_, result := c.sv.FindByCodeAndIdNotAndOwnerNo(ctx, info.Code, ct.ID.ToString(), info.OwnerNo)
		if result {
			return rt.ErrorMessage("码值已存在")
		}
		log.Infof(ctx, log.TagAppDef, "save=%+v", info)
		r.Update(ctx, info, info.ID)
	}

	return rt.OkData(numberPg.Int64ToString(info.ID))
}

// Detail 详情
//
//	@Description:
//	@receiver c
//	@param id
func (c *BasicDataDictionaryService) Detail(ctx *gin.Context, id string) (rt rg.Rs[modBasicDataDictionary.Vo]) {
	if strPg.IsBlank(id) {
		return rt.ErrorMessage("id 错误")
	}
	r := c.sv
	find, b := r.FindByIdString(ctx, id)
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	var info modBasicDataDictionary.Vo
	copier.Copy(&info, &find)
	rt.Data = info
	return rt.Ok()
}

// Enable 启用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicDataDictionaryService) Enable(ctx *gin.Context, ct model.BaseIdsCt[string]) (rt rg.Rs[string]) {
	return c.State(ctx, ct.Ids, enumStatePg.ENABLE)
}

// Disable 禁用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicDataDictionaryService) Disable(ctx *gin.Context, ct model.BaseIdsCt[string]) (rt rg.Rs[string]) {
	return c.State(ctx, ct.Ids, enumStatePg.GetType(enumStatePg.DISABLE))
}

// State 状态 启用/禁用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicDataDictionaryService) State(ctx *gin.Context, ids []string, state enumStatePg.State) (rt rg.Rs[string]) {
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
			r.Update(ctx, entityBasic.BasicDataDictionaryEntity{State: state.IndexInt8()}, info.ID)
		}
	}
	return rt.Ok()
}

// StateEnableDisable 状态 设置 有效 停用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicDataDictionaryService) StateEnableDisable(ctx *gin.Context, ids []string, state enumStatePg.State) (rt rg.Rs[string]) {
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
func (c *BasicDataDictionaryService) LogicalDeletion(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
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
				repository.Update(ctx, entityBasic.BasicDataDictionaryEntity{State: reverse.IndexInt8()}, info.ID)
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
func (c *BasicDataDictionaryService) LogicalRecovery(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
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
			repository.Update(ctx, entityBasic.BasicDataDictionaryEntity{State: reverse.IndexInt8()}, info.ID)
		}
	}
	return rt.Ok()
}

// PhysicalDeletion 物理删除
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicDataDictionaryService) PhysicalDeletion(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
	if len(ids) < 1 {
		return rt.ErrorMessage("id错误")
	}
	cn := c.sv
	finds, b := cn.FindAllByIdStringIn(ctx, ids)
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	for _, info := range finds {
		log.Infof(ctx, log.TagAppDef, "id=%v", info.ID)
	}
	cn.DeleteByIdsString(ctx, ids)
	return rt.Ok()
}

// Query 查询
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicDataDictionaryService) Query(ctx *gin.Context, ct modBasicDataDictionary.QueryCt) (rt rg.Rs[pagePg.Paginator[modBasicDataDictionary.Vo]]) {
	var query entityBasic.BasicDataDictionaryEntity
	copier.Copy(&query, &ct)
	r := c.sv
	slice := make([]modBasicDataDictionary.Vo, 0)
	rt.Data.Data = slice
	page, err := r.FindAllPage(ctx, query, optionsPg.WithOption(func(arg *optionsPg.OptionParams) {
		if ct.PageSize < 1 {
			ct.PageSize = 20
		}
		arg.Pageable = new(pagePg.PageablePageSize(0, ct.PageNum, ct.PageSize))
		//自定义查询
		arg.Db = arg.Db.Order("create_at desc").Where("type_code is null or type_code=''")
		if strPg.IsNotBlank(ct.Wd) {
			arg.Db = arg.Db.Where("name like ?", "%"+ct.Wd+"%").
				Or("code like ?", "%"+ct.Wd+"%").
				Or("name_fl like ?", "%"+ct.Wd+"%").
				Or("name_full like ?", "%"+ct.Wd+"%")
		}
	}), optionsPg.WithCtx(ctx))

	if nil != err {
		return rt.Ok()
	}

	if page.Total > 0 && page.Data != nil && len(page.Data) > 0 {
		pg := pagePg.NewPaginatorByPageable[modBasicDataDictionary.Vo](page.Pageable)
		ids := make([]string, 0)
		//for _, item := range page.Data {
		//	if strPg.IsNotBlank(item.Code) {
		//		ids = append(ids, item.Code)
		//	}
		//}
		mapBasic := make(map[string]*entityBasic.BasicDataDictionaryEntity)
		if len(ids) > 0 {
			infos, b := r.FindAllByCodeIn(ctx, ids)
			if !b {
				for _, item := range infos {
					mapBasic[item.Code] = item
				}
			}
		}

		//字段赋值
		for _, item := range page.Data {
			var vo modBasicDataDictionary.Vo
			copier.Copy(&vo, &item)
			//
			if obj, ok := mapBasic[item.Code]; ok {
				vo.TypeCodeName = obj.Name
			}
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
func (c *BasicDataDictionaryService) SelectNodeAllPublic(ctx *gin.Context, ct modBasicDataDictionary.SelectNodeCt) (rt rg.Rs[[]model.BaseNode]) {
	log.Infof(ctx, log.TagAppDef, "ct=%v", ct)
	var query entityBasic.BasicDataDictionaryEntity
	copier.Copy(&query, &ct)
	//
	query.State = enumStatePg.ENABLE.Index()
	//
	slice := make([]model.BaseNode, 0)
	rt.Data = slice
	infos := c.sv.FindAll(ctx, query, optionsPg.WithCondition(func(db *gorm.DB) *gorm.DB {
		db = db.Order("sort,create_at asc")
		db = db.Where("type_code is null or type_code = '' ")
		return db
	}))
	if len(infos) > 0 {
		for _, item := range infos {
			var vo modBasicDataDictionary.SelectNodeVo
			copier.Copy(&vo, &item)
			//
			code := model.BaseNode{
				Key:    item.Code,
				Id:     item.Code,
				No:     item.Code,
				Label:  item.Name,
				Extend: vo,
			}
			if len(item.Range) > 0 {
				vo.Range = strutil.SplitAndTrim(item.Range, ",")
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
func (c *BasicDataDictionaryService) ExistName(ctx *gin.Context, ct model.BaseExistWdCt[string]) (rt rg.Rs[string]) {
	if "" == ct.Wd {
		return rt.ErrorMessage("关键词不能为空")
	}
	_, result := c.sv.FindByNameAndIdNot(ctx, ct.Wd, numberPg.StrToInt64(ct.Id))
	if result {
		return rt.ErrorMessage("重复，已存在")
	}
	return rt.OkMessage("可以使用")
}

// ExistCode 查重
//
//	@Description:
//	@receiver c
//	@param ctx
//	@param ct
func (c *BasicDataDictionaryService) ExistCode(ctx *gin.Context, ct model.BaseExistWdCt[string]) (rt rg.Rs[string]) {
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

// ExistValue 查重
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BasicDataDictionaryService) ExistValue(ctx *gin.Context, ct modBasicDataDictionary.ExistValue) (rt rg.Rs[string]) {
	if "" == ct.Wd {
		return rt.ErrorMessage("关键词不能为空")
	}
	id := "0"
	if strPg.IsNotBlank(ct.Id.ToString()) {
		id = ct.Id.ToString()
	}
	owner := "0"
	if strPg.IsNotBlank(ct.OwnerNo) {
		owner = ct.OwnerNo
	}
	_, result := c.sv.FindByValueAndIdNotAndOwnerNo(ctx, ct.Wd, id, owner)
	if result {
		return rt.ErrorMessage("重复，已存在")
	}
	return rt.OkMessage("可以使用")
}
