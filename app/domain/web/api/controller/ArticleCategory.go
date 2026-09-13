package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/web/api/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/blog/modBlogArticleCategory/modApiBlogArticleCategory"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
)

type ArticleCategoryController struct {
	Sp *authPg.GroupApiMiddlewareSp    `autowire:"?"`
	sv *service.ArticleCategoryService `autowire:"?"`
}

func (c *ArticleCategoryController) SelectNodeAllPublic(ctx *gin.Context) {
	var ct modApiBlogArticleCategory.QueryPublicCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.SelectNodeAllPublic(ctx, ct))
}
