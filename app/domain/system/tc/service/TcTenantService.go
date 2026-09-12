package service

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/tc/modTcTenant"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityRam"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityTc"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryRam"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryTc"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/holderPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/automatedPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg/optionsPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/noPg2"
	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/dbPg/pagePg"
	"github.com/pangu-2/go-tools/tools/numberPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(TcTenantService))
}

// TcTenantService 租户
// @Description:
type TcTenantService struct {
	sv  *repositoryTc.TcTenantRepository    `autowire:"?"`
	acc *repositoryRam.RamAccountRepository `autowire:"?"`
}

// Create 新增
//
//	@Description:
//	@receiver c
//	@param ct
//	@return rt
func (c *TcTenantService) Create(ctx *gin.Context, ct modTcTenant.CreateUpdateCt) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	if "" == ct.Name {
		return rt.ErrorMessage("名称不能为空")
	}
	if "" == ct.Code {
		return rt.ErrorMessage("编号/标记不能为空")
	}
	holder := holderPg.GetContextAccount(ctx)
	r := c.sv
	//判断是否是自动,不是自动
	if !automatedPg.IsCreateCode(ct.Code) {
		//判断格式是否满足要求
		if !automatedPg.FormatVerify(ct.Code) {
			return rt.ErrorMessage("编号格式不能为空")
		}
		//不是自动
		_, result := r.FindByNo(ctx, ct.Code)
		if result {
			return rt.ErrorMessage("编号已存在")
		}
	}
	if strPg.IsNotBlank(ct.Founder) {
		{
			_, result := c.acc.FindByNo(ctx, ct.Founder)
			if !result {
				return rt.ErrorMessage("创始人不存在")
			}
		}
		// 判断创始人是否已经绑定过租户
		{
			_, result := r.FindByFounder(ctx, ct.Founder)
			if result {
				return rt.ErrorMessage("该帐户已绑定租户")
			}
		}
	}
	var info entityTc.TcTenantEntity
	copier.Copy(&info, &ct)
	info.ID = strPg.GenerateNumberId18()
	//租户编号生成
	info.No = noPg2.TenantNo()
	//自动设置编号
	if automatedPg.IsCreateCode(info.Code) {
		info.Code = info.No
	}
	info.CreateBy = holder.GetAccountNo()
	log.Infof(ctx, log.TagAppDef, "info=%+v", info)
	err, _ := r.Create(ctx, &info)
	if err != nil {
		return rt.ErrorMessage("保存失败 " + err.Error())
	}
	log.Infof(ctx, log.TagAppDef, "save=%+v", info)
	return rt.OkData(numberPg.Int64ToString(info.ID))
}

// Update 更新
//
//	@Description:
//	@receiver c
//	@param ct
//	@return rt
func (c *TcTenantService) Update(ctx *gin.Context, ct modTcTenant.CreateUpdateCt) (rt rg.Rs[string]) {
	if ct.ID < 1 {
		return rt.ErrorMessage("id错误")
	}
	if "" == ct.Name {
		return rt.ErrorMessage("名称不能为空")
	}
	if "" == ct.Code {
		return rt.ErrorMessage("编号/标记不能为空")
	}
	r := c.sv
	{
		_, result := r.FindByCodeAndIdNot(ctx, ct.Code, ct.ID.ToString())
		if result {
			return rt.ErrorMessage("编号已存在")
		}
	}
	_, b := r.FindById(ctx, ct.ID.ToInt64())
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	if strPg.IsNotBlank(ct.Founder) {
		{
			_, result := c.acc.FindByNo(ctx, ct.Founder)
			if !result {
				return rt.ErrorMessage("创始人不存在")
			}
		}
		// 判断创始人是否已经绑定过租户
		{
			_, result := r.FindByFounderAndNotIdString(ctx, ct.Founder, ct.ID.ToString())
			if result {
				return rt.ErrorMessage("该帐户已绑定租户")
			}
		}
	}
	var info entityTc.TcTenantEntity
	copier.Copy(&info, &ct)
	//编号，不参与更新
	info.No = ""
	log.Infof(ctx, log.TagAppDef, "save=%+v", info)
	err := r.Update(ctx, info, info.ID)
	if err != nil {
		return rt.ErrorMessage("更新失败:" + err.Error())
	}
	return rt.Ok()
}

// Detail 详情
//
//	@Description:
//	@receiver c
//	@param id
func (c *TcTenantService) Detail(ctx *gin.Context, id int64) (rt rg.Rs[modTcTenant.Vo]) {
	if id < 1 {
		return rt.ErrorMessage("id错误")
	}
	find, b := c.sv.FindById(ctx, id)
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	var info modTcTenant.Vo
	copier.Copy(&info, &find)
	return rt.OkData(info)
}

// Enable 启用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *TcTenantService) Enable(ctx *gin.Context, ct model.BaseIdsCt[string]) (rt rg.Rs[string]) {
	return c.State(ctx, ct.Ids, enumStatePg.ENABLE)
}

// Disable 禁用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *TcTenantService) Disable(ctx *gin.Context, ct model.BaseIdsCt[string]) (rt rg.Rs[string]) {
	return c.State(ctx, ct.Ids, enumStatePg.GetType(enumStatePg.DISABLE))
}

// State 状态 启用/禁用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *TcTenantService) State(ctx *gin.Context, ids []string, state enumStatePg.State) (rt rg.Rs[string]) {
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
			r.Update(ctx, entityTc.TcTenantEntity{State: state.IndexInt8()}, info.ID)
		}
	}
	return rt.Ok()
}

// StateEnableDisable 状态 设置 有效 停用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *TcTenantService) StateEnableDisable(ctx *gin.Context, ids []string, state enumStatePg.State) (rt rg.Rs[string]) {
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
func (c *TcTenantService) LogicalDeletion(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
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
			log.Infof(ctx, log.TagAppDef, "id=%v,TenantId=%v", info.ID, "info.TenantId")
		}
		repository.DeleteByIdsString(ctx, ids)
	} else {
		for _, info := range finds {
			enum := enumStatePg.State(info.State)
			// 有效 停用，反转 为对应的 取消 弃置
			if ok, reverse := enum.ReverseEnableDisable(); ok {
				repository.Update(ctx, entityTc.TcTenantEntity{State: reverse.IndexInt8()}, info.ID)
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
func (c *TcTenantService) LogicalRecovery(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
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
			repository.Update(ctx, entityTc.TcTenantEntity{State: reverse.IndexInt8()}, info.ID)
		}
	}
	return rt.Ok()
}

// PhysicalDeletion 物理删除
//
//	@Description:
//	@receiver c
//	@param ct
func (c *TcTenantService) PhysicalDeletion(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
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
		log.Infof(ctx, log.TagAppDef, "id=%v,TenantId=%v", info.ID, 0)
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
func (c *TcTenantService) Query(ctx *gin.Context, ct modTcTenant.QueryCt) (rt rg.Rs[pagePg.Paginator[modTcTenant.Vo]]) {
	var query entityTc.TcTenantEntity
	copier.Copy(&query, &ct)
	slice := make([]modTcTenant.Vo, 0)
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
			arg.Db = arg.Db.Where("name like ?", "%"+ct.Wd+"%")
		}
	}))
	if nil != err {
		return rt.Ok()
	}

	if page.Total > 0 && page.Data != nil && len(page.Data) > 0 {
		pg := pagePg.NewPaginatorByPageable[modTcTenant.Vo](page.Pageable)
		mapAccount := make(map[string]*entityRam.RamAccountEntity)
		idsFounder := make([]string, 0)
		for _, item := range page.Data {
			if strPg.IsNotBlank(item.Founder) {
				idsFounder = append(idsFounder, item.Founder)
			}
		}
		if len(idsFounder) > 0 {
			info, result := c.acc.FindAllByNoIn(ctx, idsFounder)
			if result {
				for _, item := range info {
					mapAccount[item.No] = item
				}
			}
		}
		//字段赋值
		for _, item := range page.Data {
			var vo modTcTenant.Vo
			copier.Copy(&vo, &item)
			vo.Ext = make(map[string]interface{})
			//
			if strPg.IsNotBlank(item.Founder) {
				if info, ok := mapAccount[item.Founder]; ok {
					vo.FounderName = info.Name
					vo.Ext["founderPhone"] = info.Phone
					vo.Ext["founderRealName"] = info.RealName
				}
			}
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
func (c *TcTenantService) SelectNodePublic(ctx *gin.Context, ct modTcTenant.QueryCt) (rt rg.Rs[[]model.BaseNode]) {
	var query entityTc.TcTenantEntity
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
func (c *TcTenantService) SelectNodeAllPublic(ctx *gin.Context, ct modTcTenant.QueryCt) (rt rg.Rs[[]model.BaseNode]) {
	var query entityTc.TcTenantEntity
	copier.Copy(&query, &ct)
	slice := make([]model.BaseNode, 0)
	rt.Data = slice
	infos := c.sv.FindAll(ctx, query)
	if len(infos) > 0 {

		for _, item := range infos {
			var vo modTcTenant.Vo
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
func (c *TcTenantService) SelectPublic(ctx *gin.Context, ct modTcTenant.QueryCt) (rt rg.Rs[[]modTcTenant.Vo]) {
	var query entityTc.TcTenantEntity
	copier.Copy(&query, &ct)
	rt.Data = []modTcTenant.Vo{}
	infos := c.sv.FindAll(ctx, query)
	if len(infos) > 0 {
		slice := make([]modTcTenant.Vo, 0)
		for _, item := range infos {
			var vo modTcTenant.Vo
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
func (c *TcTenantService) ExistName(ctx *gin.Context, ct model.BaseExistWdCt[string]) (rt rg.Rs[string]) {
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
func (c *TcTenantService) ExistCode(ctx *gin.Context, ct model.BaseExistWdCt[string]) (rt rg.Rs[string]) {
	if "" == ct.Wd {
		return rt.ErrorMessage("查询内容不能为空")
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

// ExistNo 查重
//
//	@Description:
//	@receiver c
//	@param ct
func (c *TcTenantService) ExistNo(ctx *gin.Context, ct model.BaseExistWdCt[string]) (rt rg.Rs[string]) {
	if "" == ct.Wd {
		return rt.ErrorMessage("查询内容不能为空")
	}
	id := "0"
	if strPg.IsNotBlank(ct.Id) {
		id = ct.Id
	}
	_, result := c.sv.FindByNoAndIdNot(ctx, ct.Wd, id)
	if result {
		return rt.ErrorMessage("重复，已存在")
	}
	return rt.OkMessage("可以使用")
}
