package main

import (
	"flag"
	"fmt"
	"html/template"
	"os"
	"runtime"

	"github.com/gin-gonic/gin"
	_ "github.com/hongmengzhu/xianfu-blog-go/app/middleware"
	"github.com/hongmengzhu/xianfu-blog-go/cmd"
	_ "github.com/hongmengzhu/xianfu-blog-go/pkg/cachePg/bigCachePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/dbManagerPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/templatePg"
	"github.com/pangu-2/go-tools/tools/datetimePg"
	"github.com/pangu-2/go-tools/tools/ioPg"
	"go-spring.org/spring/gs"
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
