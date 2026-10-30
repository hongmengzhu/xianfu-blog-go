package blog

import (
	"net/http"

	"github.com/hongmengzhu/xianfu-blog-go/app/domain/host"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/staterEchoV5"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryBlog"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"github.com/labstack/echo/v5"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(TestController)).Name("BlogHostTestController").Export(gs.As[routerPg.HostRouteRegistrar]())
}

// TestController test
type TestController struct {
	staterEchoV5.RouteRegistrar
	sv *repositoryBlog.BlogArticleRepository `autowire:"?"`
}

// Hosts 模块绑定域名
func (c *TestController) Hosts() string {
	return host.Host_blog
}

// RegisterRoutes
//
//	@Description: 注册路由
//	@receiver c
//	@param g
func (c *TestController) RegisterRoutes(g *echo.Group) {
	g.GET("/test/cache", c.Cache)
}

func (c *TestController) Cache(ctx *echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]any{"data": "ok"})
}
