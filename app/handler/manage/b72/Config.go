package b72

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/domainBasic/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/basic/modBasicConfig"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/model"
	routerPg "github.com/hongmengzhu/xianfu-blog-go/pkg/core/routerPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/authPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(ConfigController)).Name("ManageConfigController").Export(gs.As[routerPg.RouteRegistrar]())
}

// ConfigController 配置（BasicConfig）
type ConfigController struct {
	routerPg.RouteRegistrar
	Sp *authPg.GroupManageMiddlewareSp `autowire:"?"`
	sv *service.BasicConfigService     `autowire:"?"`
}

// RegisterRoutes 注册路由
func (c *ConfigController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/manage/basic/config", authPg.GroupManageMiddleware(c.Sp))
	group.POST("/createUpdate", c.CreateUpdate)
	group.GET("/detail/:id", c.Detail)
	group.POST("/enable", c.Enable)
	group.POST("/disable", c.Disable)
	group.POST("/state", c.State)
	group.POST("/delete", c.Delete)
	group.POST("/recovery", c.Recovery)
	group.POST("/physicalDeletion", c.PhysicalDeletion)
	group.POST("/query", c.Query)
	group.POST("/allByModel", c.AllByModel)
	group.POST("/updateValues", c.UpdateValues)
	group.POST("/cacheOverride", c.CacheOverride)
	group.POST("/existName", c.ExistName)
}

func (c *ConfigController) CreateUpdate(ctx *gin.Context) {
	var ct modBasicConfig.CreateUpdateCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.CreateUpdate(ctx, ct))
}

func (c *ConfigController) Detail(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(200, rg.Error[modBasicConfig.Vo]("id错误"))
		return
	}
	ctx.JSON(200, c.sv.Detail(ctx, id))
}

func (c *ConfigController) Enable(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Enable(ctx, ct))
}

func (c *ConfigController) Disable(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Disable(ctx, ct))
}

func (c *ConfigController) State(ctx *gin.Context) {
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

func (c *ConfigController) Delete(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.LogicalDeletion(ctx, ct.Ids))
}

func (c *ConfigController) Recovery(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.LogicalRecovery(ctx, ct.Ids))
}

func (c *ConfigController) PhysicalDeletion(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.PhysicalDeletion(ctx, ct.Ids))
}

func (c *ConfigController) Query(ctx *gin.Context) {
	var ct modBasicConfig.QueryCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Query(ctx, ct))
}

func (c *ConfigController) AllByModel(ctx *gin.Context) {
	var ct modBasicConfig.QueryCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.AllByModel(ctx, ct))
}

func (c *ConfigController) UpdateValues(ctx *gin.Context) {
	var ct modBasicConfig.UpdateValuesCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.UpdateValues(ctx, ct))
}

func (c *ConfigController) CacheOverride(ctx *gin.Context) {
	c.sv.CacheOverride(ctx)
	ctx.JSON(200, rg.Ok[string](""))
}

func (c *ConfigController) ExistName(ctx *gin.Context) {
	var ct model.BaseExistWdCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.ExistName(ctx, ct))
}
