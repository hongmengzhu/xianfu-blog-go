package domain

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/core/cache/cacheRam"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/pub/model/modSysPub"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/clientPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
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
func (c *PublicAuthPubPrivKeyController) config(ctx *gin.Context) {
	key := c.sv.ManageLoginKey(ctx, clientPg.Browser, "1000")
	vo := modSysPub.NewAuthConfigPublic(key)
	vo.LoginEncrypt = true
	ctx.JSON(200, rg.OkData(vo))
}
