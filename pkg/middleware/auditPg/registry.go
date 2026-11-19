package auditPg

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"go-spring.org/log"
)

// routeMetaKey 显式审计声明的主键：HTTP 方法 + gin 路由模板（如 /xianfu/sys/basic/area/delete）。
type routeMetaKey struct {
	method string
	path   string
}

var routeMetas sync.Map

// RegisterRoute 显式声明「该路由需要审计」，用于覆盖自动推导结果或配置非常规路由。
// 常规路由无需调用：Global() 会按 route -> 模块 -> 动作的约定自动匹配（见 auto.go）。
func RegisterRoute(method string, path string, meta Meta) {
	m := meta.withDefaults()
	if err := m.Validate(); err != nil {
		panic(fmt.Sprintf("auditPg: %s %s: %v", method, path, err))
	}
	key := routeMetaKey{method: strings.ToUpper(strings.TrimSpace(method)), path: normalizePath(path)}
	if _, loaded := routeMetas.LoadOrStore(key, m); loaded {
		panic(fmt.Sprintf("auditPg: 路由重复声明审计: %s %s", key.method, key.path))
	}
}

// LookupRoute 按方法 + 路由模板查找显式声明的审计元数据。
func LookupRoute(method string, path string) (Meta, bool) {
	v, ok := routeMetas.Load(routeMetaKey{method: strings.ToUpper(method), path: normalizePath(path)})
	if !ok {
		return Meta{}, false
	}
	return v.(Meta), true
}

// Global 返回全局审计中间件：显式声明优先，否则按约定自动推导，两者都未命中则直接放行。
// 必须在注册路由之前 e.Use(auditPg.Global())，依赖 gin 在进入处理器前已设置 c.FullPath()。
func Global() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()
		meta, ok := LookupRoute(c.Request.Method, path)
		if !ok {
			meta, ok = autoMetaForPath(path)
		}
		if !ok {
			c.Next()
			return
		}
		run(c, meta)
	}
}

// VerifyRegistered 启动后校验并统计审计覆盖：
//   - 显式声明的路由必须真实注册，路径写错会漏记，只打错误日志，不阻断启动；
//   - 自动推导按真实路由表统计覆盖数量，便于确认约定路由是否被识别。
func VerifyRegistered(routes gin.RoutesInfo) {
	registered := make(map[routeMetaKey]struct{}, len(routes))
	for _, r := range routes {
		registered[routeMetaKey{method: strings.ToUpper(r.Method), path: normalizePath(r.Path)}] = struct{}{}
	}
	explicit := 0
	routeMetas.Range(func(k, _ any) bool {
		explicit++
		key := k.(routeMetaKey)
		if _, ok := registered[key]; !ok {
			log.Error(context.Background(), log.TagAppDef,
				fmt.Errorf("auditPg: 路由未注册但已声明审计: %s %s", key.method, key.path),
				"审计路由声明与路由表不一致")
		}
		return true
	})

	auto := 0
	for _, r := range routes {
		if _, ok := LookupRoute(r.Method, r.Path); ok {
			continue
		}
		if _, ok := autoMetaForPath(normalizePath(r.Path)); ok {
			auto++
		}
	}
	log.Info(context.Background(), log.TagAppDef, "审计路由覆盖统计",
		log.Int("routes", len(routes)), log.Int("explicit", explicit), log.Int("auto", auto))
}

func normalizePath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return "/"
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if len(path) > 1 && strings.HasSuffix(path, "/") {
		path = strings.TrimSuffix(path, "/")
	}
	return path
}
