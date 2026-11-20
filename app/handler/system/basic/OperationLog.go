package basic

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/domainBasic/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/basic/modBasicOperationLog"
	routerPg "github.com/hongmengzhu/xianfu-blog-go/pkg/core/routerPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/authPg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(OperationLogController)).Name("SystemOperationLogController").Export(gs.As[routerPg.RouteRegistrar]())
}

// OperationLogController 操作日志系统端。
type OperationLogController struct {
	routerPg.RouteRegistrar
	Sp *authPg.GroupSystemMiddlewareSp   `autowire:"?"`
	sv *service.BasicOperationLogService `autowire:"?"`
}

func (c *OperationLogController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/sys/basic/operation-log", authPg.GroupSystemMiddleware(c.Sp))
	group.POST("/query", c.Query)
}

func (c *OperationLogController) Query(ctx *gin.Context) {
	var ct modBasicOperationLog.QueryCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Query(ctx, ct))
}
