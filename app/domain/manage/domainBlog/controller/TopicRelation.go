package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/manage/domainBlog/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/blog/modBlogTopicRelation"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(TopicRelationController)).Name("ManageTopicRelationController").Export(gs.As[routerPg.RouteRegistrar]())
}

// TopicRelationController 文章话题关系
// @Description:
type TopicRelationController struct {
	routerPg.RouteRegistrar
	Sp *authPg.GroupManageMiddlewareSp   `autowire:"?"`
	sv *service.BlogTopicRelationService `autowire:"?"`
}

// RegisterRoutes 注册路由
//
//	@Description:
//	@receiver c
//	@param e
func (c *TopicRelationController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/manage/blog/topic-relation", authPg.GroupManageMiddleware(c.Sp))
	group.POST("/addByTopic", c.AddByTopic)
	group.POST("/physicalDeletion", c.PhysicalDeletion)
	group.POST("/delete", c.PhysicalDeletion)
	group.POST("/query", c.Query)
}

// AddByTopic 创建
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *TopicRelationController) AddByTopic(ctx *gin.Context) {
	var ct modBlogTopicRelation.AddByTopicCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.AddByTopic(ctx, ct))
}

// PhysicalDeletion 物理删除
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *TopicRelationController) PhysicalDeletion(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.PhysicalDeletion(ctx, ct.Ids))
}

// Query 查询列表
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *TopicRelationController) Query(ctx *gin.Context) {
	var ct modBlogTopicRelation.QueryCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Query(ctx, ct))
}
