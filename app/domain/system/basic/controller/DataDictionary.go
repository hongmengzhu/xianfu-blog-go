package controller

import (
	"github.com/gin-gonic/gin"
	service2 "github.com/hongmengzhu/xianfu-blog-go/app/domain/system/basic/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/basic/modBasicDataDictionary"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/r"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(DataDictionaryController)).Name("SystemDataDictionaryController").Export(gs.As[routerPg.RouteRegistrar]())
}

type DataDictionaryController struct {
	routerPg.RouteRegistrar
	Sp         *authPg.GroupSystemMiddlewareSp         `autowire:"?"`
	sv         *service2.BasicDataDictionaryService    `autowire:"?"`
	dictSubRep *service2.BasicDataDictionarySubService `autowire:"?"`
}

func (c *DataDictionaryController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/sys/basic/data-dictionary", authPg.GroupSystemMiddleware(c.Sp))
	group.POST("/createUpdate", c.CreateUpdate)
	group.GET("/detail/:id", c.Detail)
	group.GET("/typeCodePublic/:id", c.TypeCodePublicGet)
	group.POST("/typeCodePublic", c.TypeCodeList)
	group.POST("/typeCodeAllPublic", c.TypeCodeAllPublic)
	group.POST("/enable", c.Enable)
	group.POST("/disable", c.Disable)
	group.POST("/delete", c.Delete)
	group.POST("/recovery", c.Recovery)
	group.POST("/physicalDeletion", c.PhysicalDeletion)
	group.POST("/query", c.Query)
	group.POST("/selectNodeAll", c.SelectNodeAll)
	group.POST("/selectNodeAllPublic", c.SelectNodeAllPublic)
	group.POST("/existName", c.ExistName)
	group.POST("/existCode", c.ExistCode)
}

func (c *DataDictionaryController) CreateUpdate(ctx *gin.Context) {
	var ct modBasicDataDictionary.CreateUpdateCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.CreateUpdate(ctx, ct))
}

func (c *DataDictionaryController) Delete(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.LogicalDeletion(ctx, ct.Ids))
}

func (c *DataDictionaryController) Recovery(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.LogicalRecovery(ctx, ct.Ids))
}

func (c *DataDictionaryController) PhysicalDeletion(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.PhysicalDeletion(ctx, ct.Ids))
}

func (c *DataDictionaryController) Detail(ctx *gin.Context) {
	param := ctx.Param("id")
	if "" == param {
		ctx.JSON(200, rg.Error[string]("id不能为空"))
		return
	}
	ctx.JSON(200, c.sv.Detail(ctx, param))
}

func (c *DataDictionaryController) Enable(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Enable(ctx, ct))
}

func (c *DataDictionaryController) Disable(ctx *gin.Context) {
	var ct model.BaseIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Disable(ctx, ct))
}

func (c *DataDictionaryController) State(ctx *gin.Context) {
	var ct model.BaseStateIdsCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	state, ok := enumStatePg.IsExistInt64(ct.State)
	if !ok {
		ctx.JSON(200, r.Error("类型不正确"))
		return
	}
	ctx.JSON(200, c.sv.StateEnableDisable(ctx, ct.Ids, state))
}

func (c *DataDictionaryController) Query(ctx *gin.Context) {
	var ct modBasicDataDictionary.QueryCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.Query(ctx, ct))
}

func (c *DataDictionaryController) SelectNodeAll(ctx *gin.Context) {
	var ct modBasicDataDictionary.SelectNodeCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.SelectNodeAllPublic(ctx, ct))
}
func (c *DataDictionaryController) SelectNodeAllPublic(ctx *gin.Context) {
	var ct modBasicDataDictionary.SelectNodeCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.SelectNodeAllPublic(ctx, ct))
}

func (c *DataDictionaryController) ExistName(ctx *gin.Context) {
	var ct model.BaseExistWdCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.ExistName(ctx, ct))
}

func (c *DataDictionaryController) ExistCode(ctx *gin.Context) {
	var ct model.BaseExistWdCt[string]
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.sv.ExistCode(ctx, ct))
}

func (c *DataDictionaryController) TypeCodeAllPublic(ctx *gin.Context) {
	var ct modBasicDataDictionary.SelectNodeAllCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.dictSubRep.TypeCodeAllPublic(ctx, ct))
}

func (c *DataDictionaryController) TypeCodeList(ctx *gin.Context) {
	var ct modBasicDataDictionary.SelectNodeCt
	if !routerPg.BindJson(ctx, &ct) {
		return
	}
	ctx.JSON(200, c.dictSubRep.TypeCodeList(ctx, ct))
}
func (c *DataDictionaryController) TypeCodePublicGet(ctx *gin.Context) {
	param := ctx.Param("id")
	if "" == param {
		ctx.JSON(200, rg.Error[string]("码值类型不能为空"))
		return
	}
	var ct modBasicDataDictionary.SelectNodeCt
	ct.TypeCode = param
	ctx.JSON(200, c.dictSubRep.TypeCodeList(ctx, ct))
}
