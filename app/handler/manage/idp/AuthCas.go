package idp

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/domainIdp/service"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/idp/modIdpCas"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/routerPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(AuthCasController)).Name("ManageIdpCasController").Export(gs.As[routerPg.RouteRegistrar]())
}

// AuthCasController CAS 协议端点
type AuthCasController struct {
	routerPg.RouteRegistrar
	sv *service.IdpCasService `autowire:"?"`
}

// RegisterRoutes 注册路由
func (c *AuthCasController) RegisterRoutes(e *gin.Engine) {
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
func (c *AuthCasController) Login(ctx *gin.Context) {
	var ct modIdpCas.CasLoginCt
	if err := ctx.ShouldBind(&ct); err != nil {
		ctx.JSON(200, rg.ErrorDefault[string]())
		return
	}
	ctx.JSON(200, c.sv.Login(ctx, ct))
}

// Validate CAS 1.0 验证端点（返回纯文本 yes/no）
func (c *AuthCasController) Validate(ctx *gin.Context) {
	ct := modIdpCas.CasValidateCt{
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
func (c *AuthCasController) ServiceValidate(ctx *gin.Context) {
	ct := modIdpCas.CasValidateCt{
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
func (c *AuthCasController) P3ServiceValidate(ctx *gin.Context) {
	ct := modIdpCas.CasValidateCt{
		Ticket:  ctx.Query("ticket"),
		Service: ctx.Query("service"),
		PgtUrl:  ctx.Query("pgtUrl"),
		Format:  ctx.Query("format"),
	}

	xmlResponse := c.sv.ServiceValidate(ctx, ct)

	ctx.Header("Content-Type", "application/xml; charset=utf-8")
	ctx.String(http.StatusOK, xmlResponse)
}
