package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/web/api/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(BookmarkApiController)).Export(gs.As[routerPg.RouteRegistrar]())
}

type BookmarkApiController struct {
	routerPg.RouteRegistrar
	Sp *authPg.GroupApiMiddlewareSp `autowire:"?"`
	sv *service.BookmarkService     `autowire:"?"`
}

// RegisterRoutes
//
//	@Description: 注册路由
//	@receiver c
//	@param e
func (c *BookmarkApiController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/api/bookmark", authPg.GroupApiMiddleware(c.Sp))
	group.POST("/getAll", c.GetAll)
	group.POST("/getMy", c.GetMy)
}

// GetAll 获取所有
func (c *BookmarkApiController) GetAll(ctx *gin.Context) {
	ctx.JSON(200, c.sv.GetAll(ctx))
}

// GetMy 获取所有
func (c *BookmarkApiController) GetMy(ctx *gin.Context) {
	ctx.JSON(200, c.sv.GetMy(ctx))
}
