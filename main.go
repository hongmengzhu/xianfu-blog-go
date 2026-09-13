package main

import (
	"flag"
	"fmt"
	"html/template"
	"os"
	"runtime"

	"github.com/farseer-go/eventBus"
	fsE "github.com/farseer-go/fs"
	"github.com/gin-gonic/gin"
	_ "github.com/hongmengzhu/xianfu-blog-go/app/middleware"
	"github.com/hongmengzhu/xianfu-blog-go/cmd"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/routerPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/templatePg"
	"github.com/pangu-2/go-tools/tools/datetimePg"
	"github.com/pangu-2/go-tools/tools/ioPg"
	"go-spring.org/spring/gs"
	"go-spring.org/starter-gin"
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
	if version {
		fmt.Printf("go version: %s\nBuild version: %s\nBuild commit: %s\nBuild time: %s\n",
			runtime.Version(), cmd.BuildVersion, cmd.BuildGitCommit, cmd.BuildTime)
		return
	}

	// 提供 RouterRegister Bean，官方 starter-gin 自动发现并创建 SimpleGinServer (gs.Server)
	gs.Provide(NewRouterRegister, gs.TagArg("?"))
	//事件监听
	fsE.Initialize[eventBus.Module]("panGu")
	gs.Configure(func(app gs.App) {
		// 指定配置文件目录, 如果不设置，默认 conf 目录
		app.Property("spring.app.config.dir", "./data/config")
		//
		app.Property("env", "dev")
		app.Property("spring.http.server.enabled", "false")
	}).Run()
}

// NewRouterRegister 收集所有 RouteRegistrar，返回官方 starter-gin 要求的
// RouterRegister Bean。
//
// 官方 starter-gin 自己创建 *gin.Engine 和 http.Server，并通过
// NewSimpleGinServer(register RouterRegister, cfg Config) 注入此 Bean。
// 本函数仅负责在 starter 创建的 engine 上注册静态文件、模板和业务路由。
func NewRouterRegister(registrars []routerPg.RouteRegistrar) StarterGin.RouterRegister {
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
