package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/manage/domainRam/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modPublic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/holderPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(PublicController)).Name("ManagePublicController").Export(gs.As[routerPg.RouteRegistrar]())
}

// PublicController 用户公共动作
// @Description:
type PublicController struct {
	routerPg.RouteRegistrar
	Sp *authPg.GroupManageMiddlewareSp  `autowire:"?"`
	sv *service.RamAccountPublicService `autowire:"?"`
}

// RegisterRoutes 注册路由
//
//	@Description:
//	@receiver c
//	@param e
func (c *PublicController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/manage/public", authPg.GroupManageMiddleware(c.Sp))
	group.GET("/info", c.Public)
	group.GET("/infoPublic", c.Public)
	group.POST("/password", c.UpdatePassword)
}

// Public 用户详情
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *PublicController) Public(ctx *gin.Context) {
	ctx.JSON(200, c.sv.Public(ctx, holderPg.GetContextAccount(ctx)))
}

// UpdatePassword 修改密码
func (c *PublicController) UpdatePassword(ctx *gin.Context) {
	var ct modPublic.PasswordCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.UpdatePassword(ctx, ct))
}
