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
}

// Middleware 操作日志中间件配置载体。
type Middleware struct {
	Config Config `value:"${pg.operationLog}"`
}
