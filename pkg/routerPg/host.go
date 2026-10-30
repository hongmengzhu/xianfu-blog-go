package routerPg

import (
	"github.com/labstack/echo/v5"
)

// HostRouteRegistrar 按域名（Host）注册路由的模块接口。
//
// 与 RouteRegistrar（gin）不同，实现该接口的模块由 echo 承载：启动时按
// Hosts() 声明的域名建立 e.Host 路由组，把模块路由挂到对应域名下，
// 从而同一个端口可以按访问域名访问不同的模块。
type HostRouteRegistrar interface {
	// Hosts 返回模块绑定的域名列表，例如 []string{"blog.local"}。
	Hosts() []string

	// RegisterRoutes 在域名对应的路由组上注册路由。
	RegisterRoutes(g *echo.Group)
}
