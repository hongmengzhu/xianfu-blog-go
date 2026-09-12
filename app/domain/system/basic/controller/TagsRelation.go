package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/system/basic/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/basic/modBasicTagsRelation"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(TagsRelationController)).Name("SystemTagsRelationController").Export(gs.As[routerPg.RouteRegistrar]())
}

type TagsRelationController struct {
	routerPg.RouteRegistrar
	Sp *authPg.GroupSystemMiddlewareSp   `autowire:""`
	sv *service.BasicTagsRelationService `autowire:"?"`
}

func (c *TagsRelationController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/sys/basic/tags-relation", authPg.GroupSystemMiddleware(c.Sp))
	group.GET("/detail/:id", c.Detail)
	group.POST("/enable", c.Enable)
	group.POST("/disable", c.Disable)
	group.POST("/delete", c.Delete)
	group.POST("/recovery", c.Recovery)
	group.POST("/physicalDeletion", c.PhysicalDeletion)
	group.POST("/query", c.Query)
	group.POST("/all", c.All)
	group.POST("/selectNodeAll", c.SelectNodeAll)
	group.POST("/selectNodeAllPublic", c.SelectNodeAllPublic)
	group.POST("/existName", c.ExistName)
	group.POST("/existCode", c.ExistCode)
	group.GET("/getCategoryRoot/:category", c.GetCategory)
	group.POST("/getCategoryTagsAll/:category", c.GetCategoryTagsAll)
	group.POST("/getCategoryTags/:category", c.GetCategoryTags)
}

func (c *TagsRelationController) Create(ctx *gin.Context) {
	var ct modBasicTagsRelation.CreateCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Create(ctx, ct))
}

func (c *TagsRelationController) Update(ctx *gin.Context) {
	var ct modBasicTagsRelation.UpdateCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Update(ctx, ct))
}

func (c *TagsRelationController) Delete(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.LogicalDeletion(ctx, ct.Ids))
}

func (c *TagsRelationController) Recovery(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.LogicalRecovery(ctx, ct.Ids))
}

func (c *TagsRelationController) PhysicalDeletion(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.PhysicalDeletion(ctx, ct.Ids))
}

func (c *TagsRelationController) Detail(ctx *gin.Context) {
	param := ctx.Param("id")
	if "" == param {
		ctx.JSON(200, rg.Error[string]("id不能为空"))
		return
	}
	ctx.JSON(200, c.sv.Detail(ctx, strPg.ToInt64(param)))
}

func (c *TagsRelationController) Enable(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Enable(ctx, ct))
}

func (c *TagsRelationController) Disable(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Disable(ctx, ct))
}

func (c *TagsRelationController) State(ctx *gin.Context) {
	var ct model.BaseStateIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	state, ok := enumStatePg.IsExistInt64(ct.State)
	if !ok {
		ctx.JSON(200, rg.Error[string]("类型不正确"))
		return
	}
	ctx.JSON(200, c.sv.StateEnableDisable(ctx, ct.Ids, state))
}

func (c *TagsRelationController) Query(ctx *gin.Context) {
	var ct modBasicTagsRelation.QueryCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Query(ctx, ct))
}

func (c *TagsRelationController) SelectPublic(ctx *gin.Context) {
	ct := modBasicTagsRelation.QueryCt{State: enumStatePg.ENABLE.IndexPg()}
	ctx.JSON(200, c.sv.SelectPublic(ctx, ct))
}

func (c *TagsRelationController) All(ctx *gin.Context) {
	var ct modBasicTagsRelation.AllCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.AllByLink(ctx, ct))
}

func (c *TagsRelationController) SelectNodePublic(ctx *gin.Context) {
	ct := modBasicTagsRelation.QueryCt{State: enumStatePg.ENABLE.IndexPg()}
	ctx.JSON(200, c.sv.SelectNodePublic(ctx, ct))
}

func (c *TagsRelationController) SelectNodeAll(ctx *gin.Context) {
	ct := modBasicTagsRelation.QueryCt{}
	ctx.JSON(200, c.sv.SelectNodeAllPublic(ctx, ct))
}

func (c *TagsRelationController) SelectNodeAllPublic(ctx *gin.Context) {
	ct := modBasicTagsRelation.QueryCt{State: enumStatePg.ENABLE.IndexPg()}
	ctx.JSON(200, c.sv.SelectNodeAllPublic(ctx, ct))
}

func (c *TagsRelationController) ExistName(ctx *gin.Context) {
	var ct modBasicTagsRelation.ExistWdCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.ExistName(ctx, ct))
}

func (c *TagsRelationController) ExistCode(ctx *gin.Context) {
	var ct modBasicTagsRelation.ExistWdCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.ExistCode(ctx, ct))
}

func (c *TagsRelationController) GetCategory(ctx *gin.Context) {
	param := ctx.Param("category")
	if "" == param {
		ctx.JSON(200, rg.Error[string]("id不能为空"))
		return
	}
	ctx.JSON(200, c.sv.GetCategory(ctx, param))
}

func (c *TagsRelationController) GetCategoryTagsAll(ctx *gin.Context) {
	var ct modBasicTagsRelation.QueryCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	param := ctx.Param("category")
	if "" == param {
		ctx.JSON(200, rg.Error[string]("category 不能为空"))
		return
	}
	ctx.JSON(200, c.sv.GetCategoryTagsAll(ctx, param, ct))
}

func (c *TagsRelationController) GetCategoryTags(ctx *gin.Context) {
	var ct modBasicTagsRelation.QueryCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	param := ctx.Param("category")
	if "" == param {
		ctx.JSON(200, rg.Error[string]("category 不能为空"))
		return
	}
	ctx.JSON(200, c.sv.GetCategoryTags(ctx, param, ct))
}
