package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/pub/model/modSysPub"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/cacheService/cacheRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/consts/constsPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/consts/constsRam/typeDomainPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/state/clientPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/routerPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(PublicAuthPubPrivKeyController)).Name("SystemPublicAuthPubPrivKeyController").Export(gs.As[routerPg.RouteRegistrar]())
}

// PublicAuthPubPrivKeyController 公用 配置
type PublicAuthPubPrivKeyController struct {
	sv *cacheRam.CacheSessionPubPrive `autowire:"?"`
}

func (c *PublicAuthPubPrivKeyController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/pub")
	group.GET("/config", c.config)
}

// 获取 配置信息
// 支持按域选择登录公钥：typeDomain=system 返回系统域密钥（/xianfu/auth/sys/login 使用），
// 默认 typeDomain=manage 返回管理域（租户）密钥，与旧行为保持一致。
func (c *PublicAuthPubPrivKeyController) config(ctx *gin.Context) {
	typeDomain := typeDomainPg.TypeDomain(ctx.DefaultQuery("typeDomain", typeDomainPg.Manage.Code()))
	var key model.AuthPubPriveDto
	switch typeDomain {
	case typeDomainPg.System:
		key = c.sv.SystemLoginKeyByBrowser(ctx)
	case typeDomainPg.Manage:
		tenantNo := ctx.DefaultQuery("tenantNo", constsPg.ACCOUNT_MANAGE_No)
		key = c.sv.ManageLoginKey(ctx, clientPg.Browser, tenantNo)
	default:
		ctx.JSON(200, rg.ErrorMessageData[modSysPub.AuthConfigVo]("不支持的域类型:"+typeDomain.Code()))
		return
	}
	vo := modSysPub.NewAuthConfigPublic(key)
	vo.LoginEncrypt = true
	ctx.JSON(200, rg.OkData(vo))
}
