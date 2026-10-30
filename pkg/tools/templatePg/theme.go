package templatePg

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/configPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constHeaderPg"
	"github.com/pangu-2/go-tools/tools/strPg"
)

// Theme
//
//	@Description: 获取主题文件名称
//	@param themeName
//	@param file
//	@return string
func themeByPg(pg configPg.Pg, file string) string {
	if strPg.IsBlank(pg.Template.Theme) {
		return "nisarg/" + file + ".html"
	}
	return pg.Template.Theme + "/" + file + pg.Template.Suffix
}

// Theme
//
//	@Description: 获取主题文件名称
//	@param themeName
//	@param file
//	@return string
func themeBy(pg configPg.Pg, theme, file string) string {
	if strPg.IsBlank(theme) {
		return "nisarg/" + file + ".html"
	}
	return theme + "/" + file + pg.Template.Suffix
}

// buildTemplateParam 组装模板参数（与 Web 框架无关），返回状态码、模板名与渲染数据。
//
//	@Description: 组装模板参数
//	@param pg 主题配置
//	@param server 站点配置
//	@param req 当前请求
//	@param file 模板文件
//	@return int 状态码
//	@return string 模板名
//	@return map[string]any 渲染数据
func buildTemplateParam(pg configPg.Pg, server configPg.Server, req *http.Request, file string, opts ...Option) (int, string, map[string]any) {
	param := &TemplateParameter{
		Code:    http.StatusOK,
		HtmlObj: MakeHtmlObjDefaultMap(),
	}
	for _, opt := range opts {
		opt(param)
		if strPg.IsNotBlank(param.SitePage.SiteUrl) {
			param.SitePage.SiteUrl = server.Domain
		}
		if strPg.IsNotBlank(param.SitePage.Title) && strPg.IsNotBlank(param.SitePage.SiteName) {
			param.SitePage.Title = param.SitePage.Title + " - " + param.SitePage.SiteName
		}
	}
	param.HtmlObj["ctxRequest"] = NewHttpPgRequest(req)
	param.HtmlObj["sitePage"] = param.SitePage
	if param.DataIs {
		param.HtmlObj["data"] = param.Data
		param.HtmlObj["dataIs"] = param.DataIs
	} else if nil != param.Data {
		param.HtmlObj["data"] = param.Data
	}
	param.HtmlObj["pg"] = pg
	param.HtmlObj["server"] = server
	return param.Code, themeByPg(pg, file), param.HtmlObj
}

// HTML
//
//	@Description: 模版页面
//	@param themeName
//	@param file
//	@return string
func HTML(ctx *gin.Context, file string, opts ...Option) {
	server := configPg.Server{}
	serverTmp := ctx.MustGet(constHeaderPg.WebTemplatePgServer)
	server = serverTmp.(configPg.Server)
	//
	pg := configPg.Pg{}
	pgTmp, exists := ctx.Get(constHeaderPg.WebTemplatePg)
	if !exists {
		pg.Template.Theme = "nisarg"
	} else {
		pg = pgTmp.(configPg.Pg)
	}
	//
	code, name, obj := buildTemplateParam(pg, server, ctx.Request, file, opts...)
	ctx.HTML(code, name, obj)
}
