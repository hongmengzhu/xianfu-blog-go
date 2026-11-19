package basic

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/domainBasic/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/basic/modBasicOperationLog"
	routerPg "github.com/hongmengzhu/xianfu-blog-go/pkg/core/routerPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(OperationLogController)).Name("ManageOperationLogController").Export(gs.As[routerPg.RouteRegistrar]())
}

// OperationLogController 操作日志管理端。
type OperationLogController struct {
	routerPg.RouteRegistrar
	Sp *authPg.GroupManageMiddlewareSp   `autowire:"?"`
	sv *service.BasicOperationLogService `autowire:"?"`
}

func (c *OperationLogController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/manage/basic/operation-log", authPg.GroupManageMiddleware(c.Sp))
	group.POST("/query", c.Query)
	group.GET("/detail/:id", c.Detail)
	group.POST("/statistics", c.Statistics)
	group.GET("/options", c.Options)
	group.POST("/operators", c.Operators)
	group.POST("/targetHistory", c.TargetHistory)
	group.POST("/exportExcel", c.ExportExcel)
}

func (c *OperationLogController) Query(ctx *gin.Context) {
	var ct modBasicOperationLog.QueryCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Query(ctx, ct))
}

func (c *OperationLogController) Detail(ctx *gin.Context) {
	var id typePg.Int64String
	if n, err := strconv.ParseInt(ctx.Param("id"), 10, 64); err == nil {
		id = typePg.Int64String(n)
	}
	ctx.JSON(200, c.sv.Detail(ctx, id))
}

func (c *OperationLogController) Statistics(ctx *gin.Context) {
	var ct modBasicOperationLog.StatisticsCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Statistics(ctx, ct))
}

func (c *OperationLogController) Options(ctx *gin.Context) {
	ctx.JSON(200, c.sv.Options(ctx))
}

func (c *OperationLogController) Operators(ctx *gin.Context) {
	var ct modBasicOperationLog.OperatorsCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Operators(ctx, ct))
}

func (c *OperationLogController) TargetHistory(ctx *gin.Context) {
	var ct modBasicOperationLog.TargetHistoryCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.TargetHistory(ctx, ct))
}

func (c *OperationLogController) ExportExcel(ctx *gin.Context) {
	var ct modBasicOperationLog.QueryCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	c.sv.Export(ctx, ct, 10000)
}
