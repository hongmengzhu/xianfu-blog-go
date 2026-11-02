package ram

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/domainRam/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamAccountLoginLog"
	routerPg2 "github.com/hongmengzhu/xianfu-blog-go/pkg/core/routerPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(AccountLoginLogController)).Name("ManageAccountLoginLogController").Export(gs.As[routerPg2.RouteRegistrar]())
}

// AccountLoginLogController 团队
// @Description:
type AccountLoginLogController struct {
	routerPg2.RouteRegistrar
	Sp *authPg.GroupManageMiddlewareSp    `autowire:"?"`
	sv *service.RamAccountLoginLogService `autowire:"?"`
}

// RegisterRoutes 注册路由
//
//	@Description:
//	@receiver c
//	@param e
func (c *AccountLoginLogController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/manage/ram/account-login-log", authPg.GroupManageMiddleware(c.Sp))
	group.POST("/physicalDeletion", c.PhysicalDeletion)
	group.POST("/query", c.Query)
}

// PhysicalDeletion 物理删除
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *AccountLoginLogController) PhysicalDeletion(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg2.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.PhysicalDeletion(ctx, ct.Ids))
}

// Query 查询列表
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *AccountLoginLogController) Query(ctx *gin.Context) {
	var ct modRamAccountLoginLog.QueryCt
	if !routerPg2.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Query(ctx, ct))
}
