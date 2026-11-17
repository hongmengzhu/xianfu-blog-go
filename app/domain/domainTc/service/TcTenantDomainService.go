package service

import (
	"context"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/tc/modTcTenantDomain"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/consts/automatedPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/state/yesNoPg/yesNoIntPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/model"
	entityTc "github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityTc"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/repositoryRam"
	repositoryTc "github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/repositoryTc"
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
	"golang.org/x/exp/slices"
)

func init() {
	gs.Provide(new(TcTenantDomainService))
}

// TcTenantDomainService 租户
// @Description:
type TcTenantDomainService struct {
	sv      *repositoryTc.TcTenantDomainRepository `autowire:"?"`
	acc     *repositoryRam.RamAccountRepository    `autowire:"?"`
	ten     *repositoryTc.TcTenantRepository       `autowire:"?"`
	caCache *TcTenantDomainCacheService            `autowire:"?"`
}

// domainCodeRegexp 域名标识（即 Host）允许的字符：小写字母、数字、点、连字符、下划线
var domainCodeRegexp = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]*$`)

// normalizeDomainCode 规范化域名标识并校验格式。
//
//	@Description: 域名 code 会被当作完整 Host 存入缓存，因此必须支持点号并做规范化
//	@param code
//	@return string 规范化后的 code
//	@return bool 是否合法
func normalizeDomainCode(code string) (string, bool) {
	code = strings.ToLower(strings.TrimSpace(code))
	code = strings.Trim(code, ".")
	if code == "" || len(code) > 80 {
		return "", false
	}
	if !domainCodeRegexp.MatchString(code) {
		return "", false
	}
	return code, true
}

// reloadCache 域名增删改后刷新域名缓存，避免"重启前不生效/删除后仍可解析"
func (c *TcTenantDomainService) reloadCache(ctx context.Context) {
	if nil == c.caCache {
		return
	}
	c.caCache.Reload(ctx)
}

// Create 新增
//
//	@Description:
//	@receiver c
//	@param ct
//	@return rt
func (c *TcTenantDomainService) Create(ctx *gin.Context, ct modTcTenantDomain.CreateUpdateCt) (rt rg.Rs[string]) {
	log.Info(ctx, log.TagAppDef, "ct", log.Any("ct", ct))
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
	// 域名 code 即 Host：允许字母/数字/点/连字符，且必须全局唯一（跨租户）
	if !automatedPg.IsCreateCode(ct.Code) {
		code, valid := normalizeDomainCode(ct.Code)
		if !valid {
			return rt.ErrorMessage("域名格式不正确")
		}
		info.Code = code
		if exist, found := r.FindByCodeGlobal(ctx, code); found && nil != exist {
			return rt.ErrorMessage("域名已存在")
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
	log.Info(ctx, log.TagAppDef, "info", log.Any("info", info))
	err, _ := r.Create(ctx, &info)
	if err != nil {
		return rt.ErrorMessage("保存失败 ")
	}
	log.Info(ctx, log.TagAppDef, "save", log.Any("info", info))
	// 设置默认
	c.setDefaulted(ctx, ten.No, info.ID)
	c.reloadCache(ctx)
	return rt.OkData(numberPg.Int64ToString(info.ID))
}

// Update 更新
//
//	@Description:
//	@receiver c
//	@param ct
//	@return rt
func (c *TcTenantDomainService) Update(ctx *gin.Context, ct modTcTenantDomain.CreateUpdateCt) (rt rg.Rs[string]) {
	log.Info(ctx, log.TagAppDef, "ct", log.Any("ct", ct))
	//
	var info entityTc.TcTenantDomainEntity
	copier.Copy(&info, &ct)
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
	if strPg.IsBlank(ct.TenantNo) {
		return rt.ErrorMessage("租户编号不能为空")
	}
	//
	ten, result := c.ten.FindByNo(ctx, ct.TenantNo)
	if !result {
		return rt.ErrorMessage("租户 不存在")
	}
	r := c.sv
	code, valid := normalizeDomainCode(ct.Code)
	if !valid {
		return rt.ErrorMessage("域名格式不正确")
	}
	if exist, found := r.FindByCodeGlobal(ctx, code); found && nil != exist && exist.ID != ct.ID.ToInt64() {
		return rt.ErrorMessage("域名已存在")
	}
	info.Code = code
	_, b := r.FindById(ctx, ct.ID.ToInt64())
	if !b {
		return rt.ErrorMessage("数据不存在")
	}

	//编号，不参与更新
	info.No = ""
	info.TenantNo = ten.No
	log.Info(ctx, log.TagAppDef, "save", log.Any("info", info))
	err := r.Update(ctx, info, info.ID)
	if err != nil {
		return rt.ErrorMessage("更新失败")
	}
	// 设置默认
	c.setDefaulted(ctx, ten.No, info.ID)
	c.reloadCache(ctx)
	return rt.Ok()
}

// UpdateByTenant 更新
//
//	@Description:
//	@receiver c
//	@param ct
//	@return rt
func (c *TcTenantDomainService) UpdateByTenant(ctx *gin.Context, ct modTcTenantDomain.CreateUpdateCt) (rt rg.Rs[string]) {
	log.Info(ctx, log.TagAppDef, "ct", log.Any("ct", ct))
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
	_, result := c.ten.FindByNo(ctx, holder.GetTenantNo())
	if !result {
		return rt.ErrorMessage("租户 不存在")
	}
	r := c.sv
	code, valid := normalizeDomainCode(ct.Code)
	if !valid {
		return rt.ErrorMessage("域名格式不正确")
	}
	if exist, found := r.FindByCodeGlobal(ctx, code); found && nil != exist && exist.ID != ct.ID.ToInt64() {
		return rt.ErrorMessage("域名已存在")
	}
	find, b := r.FindById(ctx, ct.ID.ToInt64())
	if !b || nil == find {
		return rt.ErrorMessage("数据不存在")
	}
	// 只能修改自己租户的域名，禁止通过 id 夺走其它租户的绑定
	if find.TenantNo != holder.GetTenantNo() {
		return rt.ErrorMessage("数据不存在")
	}
	info.Code = code

	//编号，不参与更新
	info.No = ""
	info.TenantNo = find.TenantNo
	log.Info(ctx, log.TagAppDef, "save", log.Any("info", info))
	err := r.Update(ctx, info, info.ID)
	if err != nil {
		log.Error(ctx, log.TagAppDef, err, "更新域名失败", log.Any("id", info.ID))
		return rt.ErrorMessage("更新失败")
	}
	// 设置默认
	c.setDefaulted(ctx, find.TenantNo, info.ID)
	c.reloadCache(ctx)
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
	if strPg.IsBlank(tenantNo) {
		return rt.Ok()
	}
	// 不变量：同一租户最多一个默认域名。
	// 原实现只在 defaultCount>1 或 ==0 时才动手，defaultCount==1 且要切换默认时会把两个都留在"是"。
	// 这里统一先清零再置位，保证任何调用顺序都收敛到唯一默认。
	c.sv.SetDefaultedByTenantNo(ctx, yesNoIntPg.No.IndexInt8(), tenantNo)
	if newId > 0 {
		if err := c.sv.Update(ctx, entityTc.TcTenantDomainEntity{Defaulted: yesNoIntPg.Yes.IndexInt8()}, newId); err != nil {
			log.Error(ctx, log.TagAppDef, err, "设置默认域名失败", log.Any("id", newId))
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
	log.Info(ctx, log.TagAppDef, "ct", log.Any("ct", ct))
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
	log.Info(ctx, log.TagAppDef, "ct", log.Any("ct", ct))
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
	log.Info(ctx, log.TagAppDef, "ct", log.Any("ct", ct))
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
	changed := false
	for _, info := range finds {
		if info.State != state.IndexInt8() {
			err := r.Update(ctx, entityTc.TcTenantDomainEntity{State: state.IndexInt8()}, info.ID)
			if err != nil {
				log.Error(ctx, log.TagAppDef, err, "更新失败", log.Any("id", info.ID))
				continue
			}
			changed = true
		}
	}
	if changed {
		// 启停会影响 Host 能否解析，必须刷新缓存
		c.reloadCache(ctx)
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
	// 只使用通过租户校验的行，禁止用原始 ids 越权删除
	idsFind := make([]string, 0)
	for _, info := range finds {
		idsFind = append(idsFind, numberPg.Int64ToString(info.ID))
	}
	if c.sv.Config().Data.Delete {
		for _, info := range finds {
			log.Info(ctx, log.TagAppDef, "LogicalDeletion", log.Any("id", info.ID), log.Any("TenantId", info.TenantNo))
		}
		repository.DeleteByIdsString(ctx, idsFind)
	} else {
		for _, info := range finds {
			enum := enumStatePg.State(info.State)
			// 有效 停用，反转 为对应的 取消 弃置
			if ok, reverse := enum.ReverseEnableDisable(); ok {
				if err := repository.Update(ctx, entityTc.TcTenantDomainEntity{State: reverse.IndexInt8()}, info.ID); nil != err {
					log.Error(ctx, log.TagAppDef, err, "删除域名失败", log.Any("id", info.ID))
				}
			}
		}
	}
	c.reloadCache(ctx)

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
			if err := repository.Update(ctx, entityTc.TcTenantDomainEntity{State: reverse.IndexInt8()}, info.ID); nil != err {
				log.Error(ctx, log.TagAppDef, err, "恢复域名失败", log.Any("id", info.ID))
			}
		}
	}
	c.reloadCache(ctx)
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
		log.Info(ctx, log.TagAppDef, "id=%v,TenantId=%v", log.Any("id", info.ID), log.Any("TenantId", info.TenantNo))
		idsNew = append(idsNew, info.ID)
	}
	if len(idsNew) > 0 {
		cn.DeleteByIds(ctx, idsNew)
	}
	c.reloadCache(ctx)
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
		ct.PageNum, ct.PageSize = optionsPg.ClampPage(ct.PageNum, ct.PageSize)
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
	log.Info(ctx, log.TagAppDef, "ct", log.Any("ct", ct))
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
	// 默认标记是"唯一"的：批量设置会产生多个默认域名，因此只允许一次一个
	if len(ids) != 1 {
		return rt.ErrorMessage("只能设置一个默认域名")
	}
	r := c.sv
	finds, b := r.FindAllByIdStringIn(ctx, ids, optionsPg.WithCtx(ctx))
	if !b || len(finds) != 1 {
		return rt.ErrorMessage("数据不存在")
	}
	info := finds[0]
	if state.IsEqual(yesNoIntPg.Yes.IndexInt8()) {
		// 设为默认：清掉该租户其它默认，再把当前置为默认
		c.setDefaulted(ctx, info.TenantNo, info.ID)
	} else {
		if err := r.Update(ctx, entityTc.TcTenantDomainEntity{Defaulted: yesNoIntPg.No.IndexInt8()}, info.ID); err != nil {
			log.Error(ctx, log.TagAppDef, err, "取消默认域名失败", log.Any("id", info.ID))
		}
	}
	c.reloadCache(ctx)
	return rt.Ok()
}
