package service

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	cacheDiplPg2 "github.com/hongmengzhu/xianfu-blog-go/app/middleware/components/cachePg/cacheDiplPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/api/modApiDiplAccessKey"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityApi"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryApi"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/holderPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg/optionsPg"
	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/dbPg/pagePg"
	"github.com/pangu-2/go-tools/tools/noPg"
	"github.com/pangu-2/go-tools/tools/numberPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/userPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
	"gorm.io/gorm"
)

func init() {
	gs.Provide(new(ApiDiplAccessKeyService))
}

// ApiDiplAccessKeyService 团队
// @Description:
type ApiDiplAccessKeyService struct {
	sv  *repositoryApi.ApiDiplAccessKeyRepository `autowire:"?"`
	app *repositoryApi.ApiDiplRepository          `autowire:"?"`
}

// MakeNewRecord
//
//	@Description: 生成新记录
//	@receiver c
//	@param ctx
//	@param ct
//	@return rt
func (c *ApiDiplAccessKeyService) MakeNewRecord(ctx *gin.Context, ct model.BaseIdCt[string]) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	save := entityApi.ApiDiplAccessKeyEntity{}
	save.DiplNo = strings.TrimSpace(ct.Id)
	if strPg.IsBlank(save.DiplNo) {
		return rt.ErrorMessage("应用 编号 错误")
	}
	//当前时间增加 5年
	add := time.Now().AddDate(10, 0, 0)
	holder := holderPg.GetContextAccount(ctx)
	save.State = enumStatePg.ENABLE.IndexInt8()
	save.Name = noPg.No()
	save.TenantNo = holder.GetTenantNo()
	save.Key = strPg.GetNanoIdNumber(20)
	save.Secret = strPg.GetNanoid(20)
	save.ExpiryDate = &add
	save.KindUnique = userPg.SaltMake(save.Key, save.Secret+save.ExpiryDate.String())
	err, _ := c.sv.Create(ctx, &save)
	if err != nil {
		log.Errorf(ctx, log.TagAppDef, "", err)
		return rt.ErrorMessage("保存失败")
	}
	//加入缓存
	sha := cacheDiplPg2.HashSha(save.Key, save.Secret)
	obj := cacheDiplPg2.DiplCo{
		HashSha:  sha,
		No:       save.DiplNo,
		TenantNo: save.TenantNo,
		Key:      save.Key,
		Secret:   save.Secret,
		Ano:      save.Ano,
	}
	cacheDiplPg2.Set(save.Key, obj)
	return rt.Ok()
}

// Enable 启用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *ApiDiplAccessKeyService) Enable(ctx *gin.Context, ct model.BaseIdsCt[string]) (rt rg.Rs[string]) {
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
func (c *ApiDiplAccessKeyService) Disable(ctx *gin.Context, ct model.BaseIdsCt[string]) (rt rg.Rs[string]) {
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
func (c *ApiDiplAccessKeyService) State(ctx *gin.Context, ct model.BaseStateIdsCt[string]) (rt rg.Rs[string]) {
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
	if nil == ct.Extend {
		return rt.ErrorMessage("扩展参数错误")
	}
	no := ""
	{
		if obj, ok := ct.Extend["no"]; ok {
			no = obj.(string)
		} else {
			return rt.ErrorMessage("扩展参数错误")
		}
	}
	no = strings.TrimSpace(no)
	if strPg.IsBlank(no) {
		return rt.ErrorMessage("扩展参数错误")
	}
	for _, info := range finds {
		if info.State != state.IndexInt8() {
			r.UpdateAllByDiplNoAndNoSetState(ctx, no, numberPg.Int64ToString(info.ID), state.IndexInt8())
		}
	}
	return rt.Ok()
}

// StateEnableDisable 状态 设置 有效 停用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *ApiDiplAccessKeyService) StateEnableDisable(ctx *gin.Context, ct model.BaseStateIdsCt[string]) (rt rg.Rs[string]) {
	return c.State(ctx, ct)
}

// LogicalDeletion 逻辑删除
//
//	@Description:
//	@receiver c
//	@param ct
func (c *ApiDiplAccessKeyService) LogicalDeletion(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ids)
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
			log.Infof(ctx, log.TagAppDef, "id=%v,TenantId=%v", info.ID, info.TenantNo)
		}
		repository.DeleteByIdsString(ctx, ids)
	} else {
		for _, info := range finds {
			enum := enumStatePg.State(info.State)
			// 有效 停用，反转 为对应的 取消 弃置
			if ok, reverse := enum.ReverseEnableDisable(); ok {
				repository.Update(ctx, entityApi.ApiDiplAccessKeyEntity{State: reverse.IndexInt8()}, info.ID)
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
func (c *ApiDiplAccessKeyService) LogicalRecovery(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ids)
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
			repository.Update(ctx, entityApi.ApiDiplAccessKeyEntity{State: reverse.IndexInt8()}, info.ID)
		}
	}
	return rt.Ok()
}

// PhysicalDeletion 物理删除
//
//	@Description:
//	@receiver c
//	@param ct
func (c *ApiDiplAccessKeyService) PhysicalDeletion(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ids)
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
		log.Infof(ctx, log.TagAppDef, "id=%v,TenantId=%v", info.ID, info.TenantNo)
		idsNew = append(idsNew, info.ID)
		//删除缓存
		cacheDiplPg2.Remove(info.Key)
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
func (c *ApiDiplAccessKeyService) Query(ctx *gin.Context, ct modApiDiplAccessKey.QueryCt) (rt rg.Rs[pagePg.Paginator[modApiDiplAccessKey.Vo]]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	var query entityApi.ApiDiplAccessKeyEntity
	copier.Copy(&query, &ct)
	slice := make([]modApiDiplAccessKey.Vo, 0)
	rt.Data.Data = slice
	r := c.sv
	page, err := r.FindAllPage(ctx, query, optionsPg.WithOption(func(arg *optionsPg.OptionParams) {
		if ct.PageSize < 1 {
			ct.PageSize = 20
		}
		arg.Pageable = new(pagePg.PageablePageSize(0, ct.PageNum, ct.PageSize))
		//排序
		arg.Db = arg.Db.Order("create_at desc")
		if strPg.IsNotBlank(ct.Wd) {
			arg.Db = arg.Db.Where("name like ?", "%"+ct.Wd+"%")
		}
	}), optionsPg.WithCtx(ctx))
	if nil != err {
		return rt.Ok()
	}

	if page.Total > 0 && page.Data != nil && len(page.Data) > 0 {

		pg := pagePg.NewPaginatorByPageable[modApiDiplAccessKey.Vo](page.Pageable)
		//字段赋值
		for _, item := range page.Data {
			var vo modApiDiplAccessKey.Vo
			copier.Copy(&vo, &item)
			//
			vo.Hash = item.Key + ":" + cacheDiplPg2.HashSha(item.Key, item.Secret)
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

// SelectPublic 查询
//
//	@Description:
//	@receiver c
//	@param ct
func (c *ApiDiplAccessKeyService) SelectPublic(ctx *gin.Context, ct modApiDiplAccessKey.QueryCt) (rt rg.Rs[[]modApiDiplAccessKey.Vo]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	var query entityApi.ApiDiplAccessKeyEntity
	copier.Copy(&query, &ct)
	if strPg.IsBlank(ct.DiplNo) {
		query.DiplNo = "-1"
	}

	slice := make([]modApiDiplAccessKey.Vo, 0)
	rt.Data = slice
	infos := c.sv.FindAll(ctx, query, optionsPg.WithCondition(func(db *gorm.DB) *gorm.DB {
		db = db.Order("create_at desc")
		return db
	}))
	if len(infos) > 0 {
		for _, item := range infos {
			var vo modApiDiplAccessKey.Vo
			copier.Copy(&vo, &item)
			//
			vo.Hash = item.Key + ":" + cacheDiplPg2.HashSha(item.Key, item.Secret)
			//log.Debugf(ctx, log.TagAppDef, "vo.Hash=%+v", vo.Hash)
			//
			slice = append(slice, vo)
		}
		rt.Data = slice
	}
	return rt.Ok()
}
