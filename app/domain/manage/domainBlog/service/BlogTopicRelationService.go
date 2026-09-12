package service

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/blog/modBlogTopicRelation"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBlog"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryBlog"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/holderPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg/optionsPg"
	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/dbPg/pagePg"
	"github.com/pangu-2/go-tools/tools/numberPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(BlogTopicRelationService))
}

type BlogTopicRelationService struct {
	topic    *repositoryBlog.BlogTopicRepository         `autowire:"?"`
	article  *repositoryBlog.BlogArticleRepository       `autowire:"?"`
	relation *repositoryBlog.BlogTopicRelationRepository `autowire:"?"`
}

// AddByTopic
//
//	@Description: 加入话题
//	@receiver c
//	@param ctx
//	@param ct
//	@return rt
func (c *BlogTopicRelationService) AddByTopic(ctx *gin.Context, ct modBlogTopicRelation.AddByTopicCt) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	if strPg.IsBlank(ct.TopicNo) {
		return rt.ErrorMessage("话题编号不能为空")
	}
	if nil == ct.Nos || 0 == len(ct.Nos) {
		return rt.ErrorMessage("文章编号不能为空")
	}
	topic, result := c.topic.FindByNo(ctx, ct.TopicNo, optionsPg.WithCtx(ctx))
	if !result {
		return rt.ErrorMessage("话题 不存在")
	}
	if !enumStatePg.ENABLE.IsEqualInt8(topic.State) {
		return rt.ErrorMessage("话题 状态异常")
	}
	ids := make([]string, 0)
	for _, no := range ct.Nos {
		if strPg.IsNotBlank(no) {
			ids = append(ids, no)
		}
	}
	if len(ids) < 1 {
		return rt.ErrorMessage("文章编号不能为空")
	}
	articles, result := c.article.FindAllByNoIn(ctx, ids)
	if !result {
		return rt.ErrorMessage("文章 不存在")
	}
	holder := holderPg.GetContextAccount(ctx)
	tenantNo := holder.GetTenantNo()
	ano := holder.GetAccountNo()
	//mapArticles := make(map[string]*entityBlog.BlogArticleEntity)
	for _, article := range articles {
		//判断是否存在，如果存在则跳过
		_, r := c.relation.FindByTopicNoAndArticleNo(ctx, topic.No, article.No)
		if r {
			continue
		}
		find, b := c.article.FindByNo(ctx, article.No, optionsPg.WithCtx(ctx))
		if !b {
			continue
		}
		obj := entityBlog.BlogTopicRelationEntity{}
		obj.TenantNo = tenantNo
		obj.TopicNo = topic.No
		obj.ArticleNo = article.No
		obj.Ano = ano
		obj.Name = find.Name
		obj.Description = find.Description
		err, _ := c.relation.Create(ctx, &obj)
		if nil != err {
			log.Errorf(ctx, log.TagAppDef, "save err=%+v", err)
		}
	}
	return rt.Ok()
}

// PhysicalDeletion 物理删除
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BlogTopicRelationService) PhysicalDeletion(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ids)
	if len(ids) < 1 {
		return rt.ErrorMessage("id错误")
	}
	cn := c.relation
	finds, b := cn.FindAllByIdStringIn(ctx, ids, optionsPg.WithCtx(ctx))
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	holder := holderPg.GetContextAccount(ctx)
	tenantNo := holder.GetTenantNo()
	idsNew := make([]string, 0)
	for _, info := range finds {
		log.Infof(ctx, log.TagAppDef, "id=%v,TenantNo=%v", info.ID, info.TenantNo)
		idsNew = append(idsNew, numberPg.Int64ToString(info.ID))
	}
	if len(idsNew) > 0 {
		err := cn.DeleteAllByTenantNoAndIdsString(ctx, tenantNo, idsNew)
		if err != nil {
			log.Errorf(ctx, log.TagAppDef, "操作 err=%+v", err)
		}
	}
	return rt.Ok()
}

// Query 查询
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BlogTopicRelationService) Query(ctx *gin.Context, ct modBlogTopicRelation.QueryCt) (rt rg.Rs[pagePg.Paginator[modBlogTopicRelation.Vo]]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	var query entityBlog.BlogTopicRelationEntity
	copier.Copy(&query, &ct)
	slice := make([]modBlogTopicRelation.Vo, 0)
	rt.Data.Data = slice
	r := c.relation
	page, err := r.FindAllPage(ctx, query, optionsPg.WithOption(func(arg *optionsPg.OptionParams) {
		if ct.PageSize < 1 {
			ct.PageSize = 20
		}
		arg.Pageable = new(pagePg.PageablePageSize(0, ct.PageNum, ct.PageSize))
		//自定义查询
		arg.Db = arg.Db.Order("sort asc").Order("create_at desc")
		//自定义查询
		if strPg.IsNotBlank(ct.Wd) {
			arg.Db = arg.Db.Where("name like ?", "%"+ct.Wd+"%")
		}
	}), optionsPg.WithCtx(ctx))
	if nil != err {
		return rt.Ok()
	}

	if page.Total > 0 && page.Data != nil && len(page.Data) > 0 {
		pg := pagePg.NewPaginatorByPageable[modBlogTopicRelation.Vo](page.Pageable)
		//字段赋值
		for _, item := range page.Data {
			var vo modBlogTopicRelation.Vo
			copier.Copy(&vo, &item)
			slice = append(slice, vo)
		}
		pg.Data = slice
		pg.Pageable = page.Pageable
		rt.Data = pg
		return rt.Ok()
	}
	return rt.Ok()
}
