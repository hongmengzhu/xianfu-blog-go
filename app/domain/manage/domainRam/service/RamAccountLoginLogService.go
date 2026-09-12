package service

import (
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamAccountLoginLog"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityRam"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg/optionsPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"

	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/dbPg/pagePg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
)

func init() {
	gs.Provide(new(RamAccountLoginLogService))
}

// RamAccountLoginLogService 团队
// @Description:
type RamAccountLoginLogService struct {
	sv    *repositoryRam.RamAccountLoginLogRepository `autowire:"?"`
	accDb *repositoryRam.RamAccountRepository         `autowire:"?"`
}

// PhysicalDeletion 物理删除
//
//	@Description:
//	@receiver c
//	@param ct
func (c *RamAccountLoginLogService) PhysicalDeletion(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
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
func (c *RamAccountLoginLogService) Query(ctx *gin.Context, ct modRamAccountLoginLog.QueryCt) (rt rg.Rs[pagePg.Paginator[modRamAccountLoginLog.Vo]]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	var query entityRam.RamAccountLoginLogEntity
	copier.Copy(&query, &ct)
	slice := make([]modRamAccountLoginLog.Vo, 0)
	rt.Data.Data = slice
	r := c.sv
	page, err := r.FindAllPage(ctx, query, optionsPg.WithOption(func(arg *optionsPg.OptionParams) {
		if ct.PageSize < 1 {
			ct.PageSize = 20
		}
		arg.Pageable = new(pagePg.PageablePageSize(0, ct.PageNum, ct.PageSize))
		arg.Db = arg.Db.Order("create_at desc")
		//自定义查询
		//if strPg.IsNotBlank(ct.Wd) {
		//	arg.Db.Where("ip like ?", "%"+ct.Wd+"%")
		//}
	}), optionsPg.WithCtx(ctx))
	if nil != err {
		return rt.Ok()
	}

	if page.Total > 0 && page.Data != nil && len(page.Data) > 0 {
		pg := pagePg.NewPaginatorByPageable[modRamAccountLoginLog.Vo](page.Pageable)
		mapAcc := make(map[string]*entityRam.RamAccountEntity)
		idsAcc := make([]string, 0)
		for _, item := range page.Data {
			if strPg.IsNotBlank(item.Ano) && !slices.Contains(idsAcc, item.Ano) {
				idsAcc = append(idsAcc, item.Ano)
			}
		}
		// 账号
		{
			if len(idsAcc) > 0 {
				acc, b := c.accDb.FindAllByNoIn(ctx, idsAcc)
				if b {
					for _, item := range acc {
						mapAcc[item.No] = item
					}
				}
			}
		}
		//字段赋值
		for _, item := range page.Data {
			var vo modRamAccountLoginLog.Vo
			copier.Copy(&vo, &item)
			//
			vo.ExtraData = make(map[string]any)
			if strPg.IsNotBlank(item.Ano) {
				if acc, ok := mapAcc[item.Ano]; ok {
					vo.ExtraData["account"] = acc.Account
					vo.ExtraData["cc"] = acc.Cc
					vo.ExtraData["code"] = acc.Code
					vo.ExtraData["description"] = acc.Description
					vo.ExtraData["mail"] = acc.Mail
					vo.ExtraData["phone"] = acc.Phone
					vo.ExtraData["registerTime"] = acc.RegisterTime
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
