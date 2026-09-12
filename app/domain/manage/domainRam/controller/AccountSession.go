package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/manage/domainRam/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamAccountSession"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(AccountSessionController)).Name("ManageAccountSessionController").Export(gs.As[routerPg.RouteRegistrar]())
}

// AccountSessionController 团队
// @Description:
type AccountSessionController struct {
	routerPg.RouteRegistrar
	Sp *authPg.GroupManageMiddlewareSp   `autowire:""`
	sv *service.RamAccountSessionService `autowire:"?"`
}

// RegisterRoutes 注册路由
//
//	@Description:
//	@receiver c
//	@param e
func (c *AccountSessionController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/manage/ram/account-session", authPg.GroupManageMiddleware(c.Sp))
	group.POST("/physicalDeletion", c.PhysicalDeletion)
	group.POST("/query", c.Query)
}

// PhysicalDeletion 物理删除
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *AccountSessionController) PhysicalDeletion(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.PhysicalDeletion(ctx, ct.Ids))
}

// Query 查询列表
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *AccountSessionController) Query(ctx *gin.Context) {
	var ct modRamAccountSession.QueryCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Query(ctx, ct))
}
