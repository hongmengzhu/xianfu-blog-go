package service

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/tc/modTcTenantDomain"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityTc"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryRam"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryTc"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/holderPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/automatedPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/yesNoPg/yesNoIntPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg/optionsPg"
	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/dbPg/pagePg"
	"github.com/pangu-2/go-tools/tools/noPg"
	"github.com/pangu-2/go-tools/tools/numberPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
	"golang.org/x/exp/slices"
)

func init() {
	gs.Provide(new(TcTenantDomainService))
}

// TcTenantDomainService 租户
// @Description:
type TcTenantDomainService struct {
	sv *repositoryTc.TcTenantDomainRepository `autowire:"?"`

	acc *repositoryRam.RamAccountRepository `autowire:"?"`
	ten *repositoryTc.TcTenantRepository    `autowire:"?"`
}

// Create 新增
//
//	@Description:
//	@receiver c
//	@param ct
//	@return rt
func (c *TcTenantDomainService) Create(ctx *gin.Context, ct modTcTenantDomain.CreateCt) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	//
	var info entityTc.TcTenantDomainEntity
	copier.Copy(&info, &ct)
	holder := holderPg.GetContextAccount(ctx)
	info.TenantNo = holder.GetTenantNo()
	//
	if "" == ct.Name {
		return rt.ErrorMessage("名称不能为空")
	}
	if "" == ct.Code {
		return rt.ErrorMessage("编号/标记不能为空")
	}
	if strPg.IsBlank(info.TenantNo) {
		return rt.ErrorMessage("租户编号不能为空")
	}
	ten, result := c.ten.FindByNo(ctx, info.TenantNo)
	if !result {
		return rt.ErrorMessage("租户 不存在")
	}
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
	info.ID = strPg.GenerateNumberId18()
	//租户编号生成
	info.No = noPg.No()
	//自动设置编号
	if automatedPg.IsCreateCode(info.Code) {
		info.Code = info.No
	}
	info.CreateBy = holder.GetAccountNo()
	info.TenantNo = ten.No
	log.Infof(ctx, log.TagAppDef, "info=%+v", info)
	err, _ := r.Create(ctx, &info)
	if err != nil {
		return rt.ErrorMessage("保存失败 " + err.Error())
	}
	log.Infof(ctx, log.TagAppDef, "save=%+v", info)
	// 设置默认
	c.setDefaulted(ctx, ten.No, info.ID)
	return rt.OkData(numberPg.Int64ToString(info.ID))
}

// Update 更新
//
//	@Description:
//	@receiver c
//	@param ct
//	@return rt
func (c *TcTenantDomainService) Update(ctx *gin.Context, ct modTcTenantDomain.UpdateCt) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	//
	var info entityTc.TcTenantDomainEntity
	copier.Copy(&info, &ct)
	holder := holderPg.GetContextAccount(ctx)
	//
	if ct.ID < 1 {
		return rt.ErrorMessage("id错误")
	}
	if "" == ct.Name {
		return rt.ErrorMessage("名称不能为空")
	}
	if "" == ct.Code {
		return rt.ErrorMessage("编号/标记不能为空")
	}
	//
	ten, result := c.ten.FindByNo(ctx, holder.GetTenantNo())
	if !result {
		return rt.ErrorMessage("租户 不存在")
	}
	r := c.sv
	{
		_, result := r.FindByCodeAndIdNot(ctx, ct.Code, ct.ID.ToString(), optionsPg.WithCtx(ctx))
		if result {
			return rt.ErrorMessage("编号已存在")
		}
	}
	_, b := r.FindById(ctx, ct.ID.ToInt64())
	if !b {
		return rt.ErrorMessage("数据不存在")
	}

	//编号，不参与更新
	info.No = ""
	info.TenantNo = ten.No
	log.Infof(ctx, log.TagAppDef, "save=%+v", info)
	err := r.Update(ctx, info, info.ID)
	if err != nil {
		return rt.ErrorMessage("更新失败:" + err.Error())
	}
	// 设置默认
	c.setDefaulted(ctx, ten.No, info.ID)
	return rt.Ok()
}

// setDefaulted
//
//	@Description: 设置默认
//	@receiver c
//	@param ctx
//	@param tenantNo
//	@param newId
//	@return rt
func (c *TcTenantDomainService) setDefaulted(ctx context.Context, tenantNo string, newId int64) (rt rg.Rs[string]) {
	infos, query := c.sv.FindAllByTenantNo(ctx, tenantNo)
	if query {
		defaultCount := 0
		for _, item := range infos {
			if yesNoIntPg.Yes.IsExistInt8(item.Defaulted) {
				defaultCount++
			}
		}
		// 默认值只能有一个
		if defaultCount > 1 {
			//设置所有默认值为否
			c.sv.SetDefaultedByTenantNo(ctx, yesNoIntPg.No.IndexInt8(), tenantNo)
			//设置当前默认值是是
			if newId > 0 {
				c.sv.Update(ctx, entityTc.TcTenantDomainEntity{Defaulted: yesNoIntPg.Yes.IndexInt8()}, newId)
			}
		} else if defaultCount == 0 && newId > 0 {
			//如果没有默认值，则设置当前默认值是是
			c.sv.Update(ctx, entityTc.TcTenantDomainEntity{Defaulted: yesNoIntPg.Yes.IndexInt8()}, newId)
		}
	}
	return rt.Ok()
}

// Detail 详情
//
//	@Description:
//	@receiver c
//	@param id
func (c *TcTenantDomainService) Detail(ctx *gin.Context, id int64) (rt rg.Rs[modTcTenantDomain.Vo]) {
	if id < 1 {
		return rt.ErrorMessage("id错误")
	}
	find, b := c.sv.FindById(ctx, id)
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	var info modTcTenantDomain.Vo
	copier.Copy(&info, &find)
	//
	if strPg.IsNotBlank(find.TenantNo) {
		tmp, result := c.ten.FindByNo(ctx, find.TenantNo)
		if result {
			info.TenantName = tmp.Name
		}
	}
	//
	return rt.OkData(info)
}

// Enable 启用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *TcTenantDomainService) Enable(ctx *gin.Context, ct model.BaseIdsCt[string]) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	ct2 := model.BaseStateIdsCt[string]{
		Ids: ct.Ids,
	}
	ct2.State = enumStatePg.ENABLE.IndexInt64()
	return c.State(ctx, ct2)
}

// Disable 禁用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *TcTenantDomainService) Disable(ctx *gin.Context, ct model.BaseIdsCt[string]) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	ct2 := model.BaseStateIdsCt[string]{
		Ids: ct.Ids,
	}
	ct2.State = enumStatePg.DISABLE.IndexInt64()
	return c.State(ctx, ct2)
}

// State 状态 启用/禁用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *TcTenantDomainService) State(ctx *gin.Context, ct model.BaseStateIdsCt[string]) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	ids := ct.Ids
	if len(ids) < 1 {
		return rt.ErrorMessage("id错误")
	}
	state, ok := enumStatePg.IsExistInt64(ct.State)
	if !ok {
		return rt.ErrorMessage("类型不正确")
	}
	if !state.IsEnableDisable() {
		return rt.ErrorMessage("状态错误")
	}
	r := c.sv
	finds, b := r.FindAllByIdStringIn(ctx, ids, optionsPg.WithCtx(ctx))
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	for _, info := range finds {
		if info.State != state.IndexInt8() {
			err := r.Update(ctx, entityTc.TcTenantDomainEntity{State: state.IndexInt8()}, info.ID)
			if err != nil {
				log.Errorf(ctx, log.TagAppDef, "更新失败", "id:", info.ID, "err", err)
			}
		}
	}
	return rt.Ok()
}

// StateEnableDisable 状态 设置 有效 停用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *TcTenantDomainService) StateEnableDisable(ctx *gin.Context, ct model.BaseStateIdsCt[string]) (rt rg.Rs[string]) {
	return c.State(ctx, ct)
}

// LogicalDeletion 逻辑删除
//
//	@Description:
//	@receiver c
//	@param ct
func (c *TcTenantDomainService) LogicalDeletion(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
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
			log.Infof(ctx, log.TagAppDef, "id=%v,TenantId=%v", info.ID, "info.TenantId")
		}
		repository.DeleteByIdsString(ctx, ids)
	} else {
		for _, info := range finds {
			enum := enumStatePg.State(info.State)
			// 有效 停用，反转 为对应的 取消 弃置
			if ok, reverse := enum.ReverseEnableDisable(); ok {
				repository.Update(ctx, entityTc.TcTenantDomainEntity{State: reverse.IndexInt8()}, info.ID)
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
func (c *TcTenantDomainService) LogicalRecovery(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
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
		//  取消 弃置 批量删除，反转 为对应的 有效 停用 停用
		if ok, reverse := enum.ReverseCancelLayAside(); ok {
			repository.Update(ctx, entityTc.TcTenantDomainEntity{State: reverse.IndexInt8()}, info.ID)
		}
	}
	return rt.Ok()
}

// PhysicalDeletion 物理删除
//
//	@Description:
//	@receiver c
//	@param ct
func (c *TcTenantDomainService) PhysicalDeletion(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
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
func (c *TcTenantDomainService) Query(ctx *gin.Context, ct modTcTenantDomain.QueryCt) (rt rg.Rs[pagePg.Paginator[modTcTenantDomain.Vo]]) {
	var query entityTc.TcTenantDomainEntity
	copier.Copy(&query, &ct)
	slice := make([]modTcTenantDomain.Vo, 0)
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
	}), optionsPg.WithCtx(ctx))
	if nil != err {
		return rt.Ok()
	}

	if page.Total > 0 && page.Data != nil && len(page.Data) > 0 {
		pg := pagePg.NewPaginatorByPageable[modTcTenantDomain.Vo](page.Pageable)
		mapTenant := make(map[string]*entityTc.TcTenantEntity)
		idsTenant := make([]string, 0)
		for _, item := range page.Data {
			if strPg.IsNotBlank(item.TenantNo) && !slices.Contains(idsTenant, item.TenantNo) {
				idsTenant = append(idsTenant, item.TenantNo)
			}
		}
		if len(idsTenant) > 0 {
			info, result := c.ten.FindAllByNoIn(ctx, idsTenant)
			if result {
				for _, item := range info {
					mapTenant[item.No] = item
				}
			}
		}
		//字段赋值
		for _, item := range page.Data {
			var vo modTcTenantDomain.Vo
			copier.Copy(&vo, &item)
			vo.Ext = make(map[string]interface{})
			//
			if strPg.IsNotBlank(item.TenantNo) {
				if info, ok := mapTenant[item.TenantNo]; ok {
					vo.TenantName = info.Name
					vo.Ext["tenantCode"] = info.Code
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
func (c *TcTenantDomainService) SelectNodePublic(ctx *gin.Context, ct modTcTenantDomain.QueryCt) (rt rg.Rs[[]model.BaseNode]) {
	var query entityTc.TcTenantDomainEntity
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
func (c *TcTenantDomainService) SelectNodeAllPublic(ctx *gin.Context, ct modTcTenantDomain.QueryCt) (rt rg.Rs[[]model.BaseNode]) {
	var query entityTc.TcTenantDomainEntity
	copier.Copy(&query, &ct)
	slice := make([]model.BaseNode, 0)
	rt.Data = slice
	infos := c.sv.FindAll(ctx, query)
	if len(infos) > 0 {

		for _, item := range infos {
			var vo modTcTenantDomain.Vo
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
func (c *TcTenantDomainService) SelectPublic(ctx *gin.Context, ct modTcTenantDomain.QueryCt) (rt rg.Rs[[]modTcTenantDomain.Vo]) {
	var query entityTc.TcTenantDomainEntity
	copier.Copy(&query, &ct)
	rt.Data = []modTcTenantDomain.Vo{}
	infos := c.sv.FindAll(ctx, query)
	if len(infos) > 0 {
		slice := make([]modTcTenantDomain.Vo, 0)
		for _, item := range infos {
			var vo modTcTenantDomain.Vo
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
func (c *TcTenantDomainService) ExistName(ctx *gin.Context, ct model.BaseExistWdCt[string]) (rt rg.Rs[string]) {
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

// ExistCode 查重
//
//	@Description:
//	@receiver c
//	@param ct
func (c *TcTenantDomainService) ExistCode(ctx *gin.Context, ct model.BaseExistWdCt[string]) (rt rg.Rs[string]) {
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

// ExistNo 查重
//
//	@Description:
//	@receiver c
//	@param ct
func (c *TcTenantDomainService) ExistNo(ctx *gin.Context, ct model.BaseExistWdCt[string]) (rt rg.Rs[string]) {
	if "" == ct.Wd {
		return rt.ErrorMessage("查询内容不能为空")
	}
	id := "0"
	if strPg.IsNotBlank(ct.Id) {
		id = ct.Id
	}
	_, result := c.sv.FindByNoAndIdNot(ctx, ct.Wd, id, optionsPg.WithCtx(ctx))
	if result {
		return rt.ErrorMessage("重复，已存在")
	}
	return rt.OkMessage("可以使用")
}

// SetDefaulted 状态 默认是否
//
//	@Description:
//	@receiver c
//	@param ct
func (c *TcTenantDomainService) SetDefaulted(ctx *gin.Context, ct model.BaseStateIdsCt[string]) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	ids := ct.Ids
	if len(ids) < 1 {
		return rt.ErrorMessage("id错误")
	}
	state, ok := yesNoIntPg.IsExistInt64(ct.State)
	if !ok {
		return rt.ErrorMessage("类型不正确")
	}
	if !state.IsEnableDisable() {
		return rt.ErrorMessage("状态错误")
	}
	r := c.sv
	finds, b := r.FindAllByIdStringIn(ctx, ids, optionsPg.WithCtx(ctx))
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	//如果是 默认启用，那么要把 当前租户下的所有都设置为 否
	if state.IsEqual(yesNoIntPg.Yes.IndexInt8()) {
		c.setDefaulted(ctx, finds[0].TenantNo, 0)
	}
	for _, info := range finds {
		err := r.Update(ctx, entityTc.TcTenantDomainEntity{Defaulted: state.IndexInt8()}, info.ID)
		if err != nil {
			log.Errorf(ctx, log.TagAppDef, "更新失败", "id:", info.ID, "err", err)
		}
	}
	return rt.Ok()
}
