package ram

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/domainRam/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamAccountSession"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/model"
	routerPg2 "github.com/hongmengzhu/xianfu-blog-go/pkg/core/routerPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/authPg/controllerPg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(AccountSessionController)).Name("SystemAccountSessionController").Export(gs.As[routerPg2.RouteRegistrar]())
}

type AccountSessionController struct {
	routerPg2.RouteRegistrar
	controllerPg.SpSystemAuth
	sv *service.RamAccountSessionService `autowire:"?"`
}

func (c *AccountSessionController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/sys/ram/account-session", authPg.GroupSystemMiddleware(c.Sp))
	group.POST("/physicalDeletion", c.PhysicalDeletion)
	group.POST("/query", c.Query)
}

func (c *AccountSessionController) PhysicalDeletion(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg2.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.PhysicalDeletion(ctx, ct.Ids))
}

func (c *AccountSessionController) Query(ctx *gin.Context) {
	var ct modRamAccountSession.QueryCt
	if !routerPg2.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Query(ctx, ct))
}
