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
	param.HtmlObj["ctxRequest"] = NewHttpPg(ctx)
	param.HtmlObj["sitePage"] = param.SitePage
	if param.DataIs {
		param.HtmlObj["data"] = param.Data
		param.HtmlObj["dataIs"] = param.DataIs
	} else if nil != param.Data {
		param.HtmlObj["data"] = param.Data
	}
	param.HtmlObj["pg"] = pg
	param.HtmlObj["server"] = server
	//
	//fmt.Printf("模版HtmlObj: %+v\n", param.HtmlObj)
	ctx.HTML(param.Code, themeByPg(pg, file), param.HtmlObj)
}

// Html
//
//	@Description: 模版页面
//	@param themeName
//	@param file
//	@return string
//func Html(ctx *gin.Context, file string, opts ...Option) {
//	param := &TemplateParameter{
//		Code: http.StatusOK,
//		Data: make(map[string]any),
//	}
//	for _, opt := range opts {
//		opt(param)
//	}
//	ctx.HTML(param.Code, Theme(ctx, file), param.Data)
//}
