package service

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/manage/domainBlog/service/blogArticle"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/web/api/model/modBlogArticle"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBlog"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryBlog"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/holderPg"
	"github.com/pangu-2/go-tools/tools/noPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(ArticleService))
}

type ArticleService struct {
	sv    *repositoryBlog.BlogArticleRepository           `autowire:"?"`
	catDb *repositoryBlog.BlogArticleCategoryRepository   `autowire:"?"`
	sata  *repositoryBlog.BlogArticleStatisticsRepository `autowire:"?"`
	sp    *blogArticle.Sp                                 `autowire:"?"`
}

// Push
//
//	@Description: 推送文章连接
//	@receiver c
func (c *ArticleService) Push(ctx *gin.Context, ct modBlogArticle.PushCt) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	if strPg.IsBlank(ct.CategoryNo) {
		return rt.ErrorMessage("请选择分类")
	}
	if strPg.IsBlank(ct.Title) {
		return rt.ErrorMessage("标题不能为空")
	}
	if strPg.IsBlank(ct.Url) {
		return rt.ErrorMessage("url地址不能为空")
	}
	info, result := c.catDb.FindByNo(ctx, ct.CategoryNo)
	if !result {
		return rt.ErrorMessage("分类不存在")
	}
	holder := holderPg.GetContextAccount(ctx)
	save := entityBlog.BlogArticleEntity{
		CategoryNo:  info.No,
		Name:        ct.Title,
		UrlSource:   ct.Url,
		Description: ct.Description,
	}
	save.No = noPg.No()
	save.Code = save.No
	save.TenantNo = holder.GetTenantNo()
	save.Ano = holder.GetAccountNo()
	err, _ := c.sv.Create(ctx, &save)
	if err != nil {
		log.Debugf(ctx, log.TagAppDef, "save err=%+v", err)
		return rt.ErrorMessage("保存失败：" + err.Error())
	}
	err, _ = c.sata.Create(ctx, &entityBlog.BlogArticleStatisticsEntity{ID: save.ID, ArticleNo: save.No})
	if err != nil {
		log.Debugf(ctx, log.TagAppDef, "save err=%+v", err)
	}
	return rt.Ok()
}
