package opLogPg

import (
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(Middleware))
}

// Config 操作日志配置。
type Config struct {
	// Exclude 不记录日志的 URL 前缀，支持 *（匹配单段）与 **（匹配任意段）。
	Exclude []string `json:"exclude" toml:"exclude" value:"${exclude:=}"`
	// IncludeRead 是否记录读请求（GET/HEAD）。
	// 默认 false：只记录有副作用的写操作（POST/PUT/PATCH/DELETE 等），
	// 避免读接口与页面浏览造成操作日志全量入库；需要审计读请求时置为 true。
	IncludeRead bool `json:"includeRead" toml:"includeRead" value:"${includeRead:=false}"`
}

// Middleware 操作日志中间件配置载体。
type Middleware struct {
	Config Config `value:"${pg.operationLog}"`
}
