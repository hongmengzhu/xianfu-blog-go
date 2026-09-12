package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryBlog"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(TestController)).Export(gs.As[routerPg.RouteRegistrar]())
}

// TestController test
type TestController struct {
	routerPg.RouteRegistrar
	sv *repositoryBlog.BlogArticleRepository `autowire:"?"`
}

// RegisterRoutes
//
//	@Description: 注册路由
//	@receiver c
//	@param e
func (c *TestController) RegisterRoutes(e *gin.Engine) {
	e.GET("/test/cache", c.Cache)
}

func (c *TestController) Cache(ctx *gin.Context) {
	//err := articleBlogEvent.NewStartInit(c.log).Processor(context.Background())
	//if err != nil {
	//	log.Errorf(ctx, log.TagAppDef,"error:", err)
	//}
	// 模版
	ctx.JSON(200, gin.H{"data": "ok"})
}

func (c *TestController) FindAllPage(ctx *gin.Context) {

	// 模版
	ctx.JSON(200, gin.H{"data": "ok"})
}
