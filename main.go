package main

import (
	"flag"
	"fmt"
	"html/template"
	"os"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/hongmengzhu/xianfu-blog-go/app/middleware"
	"github.com/hongmengzhu/xianfu-blog-go/cmd"
	_ "github.com/hongmengzhu/xianfu-blog-go/pkg/cachePg/bigCachePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/dbManagerPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/templatePg"
	"github.com/labstack/echo/v4"
	"github.com/pangu-2/go-tools/tools/datetimePg"
	"github.com/pangu-2/go-tools/tools/ioPg"
	"go-spring.org/spring/gs"
	"go-spring.org/starter-echo"
	StarterGin "go-spring.org/starter-gin"
)

func init() {
	//
	dirPath := make([]string, 0)
	// 日志目录路径
	dirPath = append(dirPath, "data/logs")
	//配置文件目录
	dirPath = append(dirPath, "data/config")
	//附件
	dirPath = append(dirPath, "data/attachment")
	if nil != dirPath {
		for _, path := range dirPath {
			// 创建目录
			err := ioPg.DirectoryCreate(path)
			if err != nil {
				fmt.Printf("创建目录失败: [%v] => %v", path, err)
			}
		}
	}
	//
	//指定环境
	//gs.SetActiveProfiles("dev")
}

func main() {

	// 构建信息，golang版本 commit id 时间
	var version bool
	flag.BoolVar(&version, "v", false, "version")
	flag.Parse()
	// Unset env vars that leak from the developer shell so runs are reproducible
	// and consistent with sibling examples.
	_ = os.Unsetenv("_")
	_ = os.Unsetenv("TERM")
	_ = os.Unsetenv("TERM_SESSION_ID")
	if version {
		fmt.Printf("go version: %s\nBuild version: %s\nBuild commit: %s\nBuild time: %s\n",
			runtime.Version(), cmd.BuildVersion, cmd.BuildGitCommit, cmd.BuildTime)
		return
	}
	gs.Provide(new(dbManagerPg.Manager))
	// 提供 gin 引擎 Bean：承载非域名模块（后台、接口等）的原有路由
	gs.Provide(NewGinEngine, gs.TagArg("?"))
	// 提供 RouterRegister Bean，由官方 starter-echo 创建并持有 *echo.Echo 与 http.Server
	gs.Provide(NewEchoRegister, gs.TagArg("?"))
	//
	gs.Configure(func(app gs.App) {
		// 指定配置文件目录, 如果不设置，默认 conf 目录
		app.Property("spring.app.config.dir", "./data/config")
		//
		app.Property("env", "dev")
		app.Property("spring.http.server.enabled", "false")
	}).Run()
}

// NewGinEngine 创建 gin 引擎，注册静态文件、模板与全部非域名模块路由。
// 该引擎不再独立监听端口，而是由 echo 作为兜底处理器调用，
// 保证原有管理后台与接口路由行为不变。
//
//	@Description: 创建 gin 引擎
//	@param registrars
//	@return *gin.Engine
func NewGinEngine(registrars []routerPg.RouteRegistrar) StarterGin.RouterRegister {
	return func(e *gin.Engine) {
		// 静态文件目录
		e.Static("/assets", "./data/assets")
		e.Static("/attachment", "./data/attachment")

		// 模板函数，需在 LoadHTMLGlob 之前设置
		e.SetFuncMap(template.FuncMap{
			"unescaped":  templatePg.Unescaped,
			"dateformat": datetimePg.Format,
		})

		// 加载模板文件，需在路由注册之前
		// 加载模板文件，需在路由注册之前（目录不存在时跳过）
		if _, err := os.Stat("data/templates"); err == nil {
			e.LoadHTMLGlob("data/templates/**/**/*")
		}

		// 注册所有路由
		for _, r := range registrars {
			r.RegisterRoutes(e)
		}
	}
}

// NewEchoRegister 返回 starter-echo 的路由注册器。
//
// 按域名（Host）把各模块路由挂载到 e.Host 路由组上：访问 blog.local 命中
// 博客模块，访问 bookmark.local 命中书签模块；未命中任何模块域名的请求
// 全部交给 gin 引擎处理，实现「一个端口、多域名、多模块」。
//
//	@Description: 注册 echo 路由
//	@param ginEngine 兜底的 gin 引擎
//	@param hosts 接入 echo 的域名模块
//	@return StarterEcho.RouterRegister
func NewEchoRegister(ginEngine *gin.Engine, hosts []routerPg.HostRouteRegistrar) StarterEcho.RouterRegister {
	return func(e *echo.Echo) {
		// 模板渲染器，命名规则与 gin LoadHTMLGlob 保持一致
		if renderer := templatePg.NewEchoRenderer("data/templates/**/**/*"); nil != renderer {
			e.Renderer = renderer
		}

		// 静态文件目录
		e.Static("/assets", "./data/assets")
		e.Static("/attachment", "./data/attachment")

		// 域名归一化：去掉端口并转小写，使 e.Host 精确匹配在本地开发端口下同样生效
		e.Pre(newHostNormalizeMiddleware(hosts))

		// 同一域名可能被多个模块注册，复用同一个域名路由组
		groups := make(map[string]*echo.Group)
		for _, h := range hosts {
			for _, domain := range h.Hosts() {
				domain = normalizeDomain(domain)
				if "" == domain {
					continue
				}
				g, exists := groups[domain]
				if !exists {
					g = e.Host(domain)
					// 域名路由组是独立路由器，静态目录需在该组上单独注册
					g.Static("/assets", "./data/assets")
					g.Static("/attachment", "./data/attachment")
					groups[domain] = g
				}
				h.RegisterRoutes(g)
			}
		}

		// 兜底：其余所有请求交给 gin 引擎
		e.Any("/", echo.WrapHandler(ginEngine))
		e.Any("/*", echo.WrapHandler(ginEngine))
	}
}

// normalizeDomain 归一化域名，统一小写并去掉首尾空白。
//
//	@Description: 归一化域名
//	@param domain
//	@return string
func normalizeDomain(domain string) string {
	return strings.ToLower(strings.TrimSpace(domain))
}

// newHostNormalizeMiddleware 域名归一化中间件。
//
// echo 的 e.Host 路由组按请求 Host 精确匹配（含端口），而开发环境访问
// 形如 blog.local:9981。该前置中间件在路由匹配前把已知模块域名的 Host
// 去掉端口并转小写，其它域名保持原样，交给 gin 兜底。
//
//	@Description: 域名归一化中间件
//	@param hosts
//	@return echo.MiddlewareFunc
func newHostNormalizeMiddleware(hosts []routerPg.HostRouteRegistrar) echo.MiddlewareFunc {
	domains := make(map[string]bool)
	for _, h := range hosts {
		for _, domain := range h.Hosts() {
			if domain = normalizeDomain(domain); "" != domain {
				domains[domain] = true
			}
		}
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			host := normalizeDomain(c.Request().Host)
			if index := strings.LastIndex(host, ":"); index > 0 {
				host = host[:index]
			}
			if domains[host] {
				c.Request().Host = host
			}
			return next(c)
		}
	}
}
