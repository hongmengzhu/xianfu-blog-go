package service

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryRam"
	holderPg2 "github.com/hongmengzhu/xianfu-blog-go/pkg/auth/holderPg"
	multiTenantPg2 "github.com/hongmengzhu/xianfu-blog-go/pkg/auth/holderPg/multiTenantPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constsRam/typeDomainPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"

	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
)

func init() {
	gs.Provide(NewRamAccountMiddlewareService)
}

// RamAccountMiddlewareService 账户公共动作
// @Description:
type RamAccountMiddlewareService struct {
	sv    *repositoryRam.RamAccountRepository              `autowire:"?"`
	aAuth *repositoryRam.RamAccountAuthorizationRepository `autowire:"?"`
}

func NewRamAccountMiddlewareService() *RamAccountMiddlewareService {
	return new(RamAccountMiddlewareService)
}

// FindByLoginNo
//
//	@Description: 登录信息
//	@receiver c
//	@param jwt
//	@return rt
func (c *RamAccountMiddlewareService) FindByLoginNo(ctx *gin.Context, no string) (rt rg.Rs[holderPg2.HolderPg]) {
	log.Debugf(ctx, log.TagAppDef, "jwt=%+v", no)
	if strPg.IsBlank(no) {
		return rt.ErrorMessage("账号登陆失败")
	}
	info, b := c.sv.FindByNo(ctx, no)
	if !b {
		return rt.ErrorMessage("账号不存在")
	}
	pg := rt.Data
	rule := multiTenantPg2.NewMultiRuleDefaultBySystem()
	//
	rule.Tenant = true
	rule.Merchant = false
	rule.Org = false
	rule.Owner = false
	rule.MultiOwner = false
	//
	var accountHolder2 holderPg2.AccountHolder
	copier.Copy(&accountHolder2, &info)
	//
	accountHolder2.Os = holderPg2.NewAccountHolderOs()
	accountHolder2.Os.Departments = info.Os.Data().Departments
	accountHolder2.Os.Roles = info.Os.Data().Roles
	accountHolder2.Os.Orgs = info.Os.Data().Orgs
	accountHolder2.Os.Tenants = info.Os.Data().Tenants
	accountHolder2.Os.Teams = info.Os.Data().Teams
	accountHolder2.Os.Shops = info.Os.Data().Shops
	accountHolder2.Os.Stores = info.Os.Data().Stores
	//
	pg.MultiTenant = multiTenantPg2.MultiTenantPg{
		//TenantNo: accountHolder2.Os.Tenants,
		TenantNo: []string{info.TenantNo},
	}
	pg.HolderData = accountHolder2
	pg.TypeDomain = typeDomainPg.System.Code()
	pg.Rule = &rule
	rt.Data = pg
	return rt.Ok()
}
