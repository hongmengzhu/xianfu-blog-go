package b72

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/domainBasic/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/basic/modBasicConfigCategory"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/model"
	routerPg "github.com/hongmengzhu/xianfu-blog-go/pkg/core/routerPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/authPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(ConfigCategoryController)).Name("ManageConfigCategoryController").Export(gs.As[routerPg.RouteRegistrar]())
}

// ConfigCategoryController 配置分类（BasicConfigCategory）
type ConfigCategoryController struct {
	routerPg.RouteRegistrar
	Sp *authPg.GroupManageMiddlewareSp     `autowire:"?"`
	sv *service.BasicConfigCategoryService `autowire:"?"`
}

func (c *ConfigCategoryController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/manage/b72/config-category", authPg.GroupManageMiddleware(c.Sp))
	group.POST("/createUpdate", c.CreateUpdate)
	group.GET("/detail/:id", c.Detail)
	group.POST("/enable", c.Enable)
	group.POST("/disable", c.Disable)
	group.POST("/state", c.State)
	group.POST("/delete", c.Delete)
	group.POST("/recovery", c.Recovery)
	group.POST("/physicalDeletion", c.PhysicalDeletion)
	group.POST("/query", c.Query)
	group.POST("/existName", c.ExistName)
	group.POST("/existCode", c.ExistCode)
}

func (c *ConfigCategoryController) CreateUpdate(ctx *gin.Context) {
	var ct modBasicConfigCategory.CreateUpdateCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	if ct.ID.ToInt64() > 0 {
		ctx.JSON(200, c.sv.Update(ctx, ct))
		return
	}
	ctx.JSON(200, c.sv.Create(ctx, ct))
}

func (c *ConfigCategoryController) Detail(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(200, rg.Error[modBasicConfigCategory.Vo]("id错误"))
		return
	}
	ctx.JSON(200, c.sv.Detail(ctx, id))
}

func (c *ConfigCategoryController) Enable(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Enable(ctx, ct))
}

func (c *ConfigCategoryController) Disable(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Disable(ctx, ct))
}

func (c *ConfigCategoryController) State(ctx *gin.Context) {
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

func (c *ConfigCategoryController) Delete(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.LogicalDeletion(ctx, ct.Ids))
}

func (c *ConfigCategoryController) Recovery(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.LogicalRecovery(ctx, ct.Ids))
}

func (c *ConfigCategoryController) PhysicalDeletion(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.PhysicalDeletion(ctx, ct.Ids))
}

func (c *ConfigCategoryController) Query(ctx *gin.Context) {
	var ct modBasicConfigCategory.QueryCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Query(ctx, ct))
}

func (c *ConfigCategoryController) ExistName(ctx *gin.Context) {
	var ct model.BaseExistWdCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.ExistName(ctx, ct))
}

func (c *ConfigCategoryController) ExistCode(ctx *gin.Context) {
	var ct model.BaseExistWdCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.ExistCode(ctx, ct))
}
