package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/manage/domainRam/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamCas"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(CasController)).Name("ManageCasController").Export(gs.As[routerPg.RouteRegistrar]())
}

// CasController CAS 协议端点
type CasController struct {
	routerPg.RouteRegistrar
	sv *service.CasService `autowire:"?"`
}

// RegisterRoutes 注册路由
func (c *CasController) RegisterRoutes(e *gin.Engine) {
	group := e.Group("/xianfu/auth/manage/cas")
	// CAS 登录（生成 ST）
	group.POST("/login", c.Login)
	// CAS 1.0 验证
	group.GET("/validate", c.Validate)
	// CAS 2.0 验证
	group.GET("/serviceValidate", c.ServiceValidate)
	// CAS 3.0 验证
	group.GET("/p3/serviceValidate", c.P3ServiceValidate)
}

// Login CAS 登录（生成 Service Ticket）
func (c *CasController) Login(ctx *gin.Context) {
	var ct modRamCas.CasLoginCt
	if err := ctx.ShouldBind(&ct); err != nil {
		ctx.JSON(200, rg.ErrorDefault[string]())
		return
	}
	ctx.JSON(200, c.sv.Login(ctx, ct))
}

// Validate CAS 1.0 验证端点（返回纯文本 yes/no）
func (c *CasController) Validate(ctx *gin.Context) {
	ct := modRamCas.CasValidateCt{
		Ticket:  ctx.Query("ticket"),
		Service: ctx.Query("service"),
	}

	ctx.Header("Content-Type", "text/html; charset=utf-8")

	if ct.Service == "" || ct.Ticket == "" {
		ctx.String(http.StatusOK, "no\n")
		return
	}

	result := c.sv.Validate(ctx, ct)
	if result.ErrorIs() {
		ctx.String(http.StatusOK, "no\n")
		return
	}

	data := result.Data
	if data.Valid {
		ctx.String(http.StatusOK, "yes\n%s\n", data.User)
	} else {
		ctx.String(http.StatusOK, "no\n")
	}
}

// ServiceValidate CAS 2.0 验证端点（返回 XML）
func (c *CasController) ServiceValidate(ctx *gin.Context) {
	ct := modRamCas.CasValidateCt{
		Ticket:  ctx.Query("ticket"),
		Service: ctx.Query("service"),
		PgtUrl:  ctx.Query("pgtUrl"),
		Format:  ctx.Query("format"),
	}

	xmlResponse := c.sv.ServiceValidate(ctx, ct)

	ctx.Header("Content-Type", "application/xml; charset=utf-8")
	ctx.String(http.StatusOK, xmlResponse)
}

// P3ServiceValidate CAS 3.0 验证端点（返回 XML，包含更多属性）
func (c *CasController) P3ServiceValidate(ctx *gin.Context) {
	ct := modRamCas.CasValidateCt{
		Ticket:  ctx.Query("ticket"),
		Service: ctx.Query("service"),
		PgtUrl:  ctx.Query("pgtUrl"),
		Format:  ctx.Query("format"),
	}

	xmlResponse := c.sv.ServiceValidate(ctx, ct)

	ctx.Header("Content-Type", "application/xml; charset=utf-8")
	ctx.String(http.StatusOK, xmlResponse)
}
