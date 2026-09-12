package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/web/api/model/modBlogArticle"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/web/api/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
)

type ArticleController struct {
	Sp *authPg.GroupApiMiddlewareSp `autowire:"?"`
	sv *service.ArticleService      `autowire:"?"`
}

// Push 推送
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *ArticleController) Push(ctx *gin.Context) {
	var ct modBlogArticle.PushCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Push(ctx, ct))
}
