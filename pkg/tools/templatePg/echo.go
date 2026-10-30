package templatePg

import (
	"html/template"
	"io"
	"path/filepath"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/configPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constHeaderPg"
	"github.com/labstack/echo/v5"
	"github.com/pangu-2/go-tools/tools/datetimePg"
)

// DefaultFuncMap 模板函数集合，gin 与 echo 渲染共用，保证两套入口渲染结果一致。
func DefaultFuncMap() template.FuncMap {
	return template.FuncMap{
		"unescaped":  Unescaped,
		"dateformat": datetimePg.Format,
	}
}

// echoRenderer echo 渲染器，命名规则与 gin 的 LoadHTMLGlob 保持一致，
// 即模板以文件内 {{ define "主题/文件.html" }} 声明的名字注册。
type echoRenderer struct {
	templ *template.Template
}

// Render 实现 echo.Renderer
//
//	@Description: 渲染模板
//	@receiver r
//	@param w
//	@param name 模板名，形如 nisarg/blog/index.html
//	@param data 渲染数据
//	@param c
//	@return error
func (r *echoRenderer) Render(c *echo.Context, w io.Writer, name string, data any) error {
	if viewContext, isMap := data.(map[string]any); isMap {
		viewContext["reverse"] = c.RouteInfo().Reverse
	}
	return r.templ.ExecuteTemplate(w, name, data)
}

// NewEchoRenderer 解析模板目录并返回 echo 渲染器。
// 未匹配到任何模板文件时返回 nil，由调用方保持 echo 的默认渲染行为。
//
//	@Description: 创建 echo 渲染器
//	@param pattern 模板 glob，例如 data/templates/**/**/*
//	@return echo.Renderer
func NewEchoRenderer(pattern string) echo.Renderer {
	files, err := filepath.Glob(pattern)
	if nil != err || len(files) == 0 {
		return nil
	}
	return &echoRenderer{templ: template.Must(template.New("").Funcs(DefaultFuncMap()).ParseGlob(pattern))}
}

// HTMLEcho 渲染模板页面（echo 版本），与 HTML 的 gin 版本产出相同的模板名与数据。
//
//	@Description: 模版页面
//	@param c
//	@param file 模板文件，例如 blog/index
//	@param opts
//	@return error
func HTMLEcho(c *echo.Context, file string, opts ...Option) error {
	server := configPg.Server{}
	if v := c.Get(constHeaderPg.WebTemplatePgServer); nil != v {
		if tmp, ok := v.(configPg.Server); ok {
			server = tmp
		}
	}
	pg := configPg.Pg{}
	if v := c.Get(constHeaderPg.WebTemplatePg); nil != v {
		if tmp, ok := v.(configPg.Pg); ok {
			pg = tmp
		}
	} else {
		pg.Template.Theme = "nisarg"
	}
	code, name, obj := buildTemplateParam(pg, server, c.Request(), file, opts...)
	return c.Render(code, name, obj)
}
