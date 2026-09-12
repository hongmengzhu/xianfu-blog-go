package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/web/api/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(BookmarkCategoryApiController)).Export(gs.As[routerPg.RouteRegistrar]())
}

type BookmarkCategoryApiController struct {
	routerPg.RouteRegistrar
	Sp *authPg.GroupApiMiddlewareSp     `autowire:"?"`
	sv *service.BookmarkCategoryService `autowire:"?"`
}

// RegisterRoutes
//
//	@Description: 注册路由
//	@receiver c
//	@param e
func (c *BookmarkCategoryApiController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/api/bookmark-category", authPg.GroupApiMiddleware(c.Sp))
	group.POST("/getAll", c.GetAll)
}

// GetAll 获取所有
func (c *BookmarkCategoryApiController) GetAll(ctx *gin.Context) {
	ctx.JSON(200, c.sv.GetAll(ctx))
}
