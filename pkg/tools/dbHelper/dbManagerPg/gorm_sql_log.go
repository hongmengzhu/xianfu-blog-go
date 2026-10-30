package dbManagerPg

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go-spring.org/log"
	gormcore "go-spring.org/starter-gorm"
	"gorm.io/gorm/logger"
)

// gormSlowThreshold 与 gorm 内置默认值保持一致：耗时超过该值的 SQL 记为慢查询。
const gormSlowThreshold = 200 * time.Millisecond

// DBs 返回已装配的 gorm 实例，未配置的方言为 nil，不返回。
func (c *Manager) DBs() []*gormcore.DB {
	out := make([]*gormcore.DB, 0, 3)
	for _, db := range []*gormcore.DB{c.dbPostgres, c.dbMysql, c.dbSqlite} {
		if db != nil {
			out = append(out, db)
		}
	}
	return out
}

// GormSQLLog 按配置项 database.log-level 设置所有 gorm 实例的 SQL 日志级别。
// 实现 gs.Runner 接口：在 Bean 装配完成之后、服务开始监听之前执行，
// 此时还没有并发查询，替换 Logger 是安全的。
//
// level 取值：silent / error / warn / info（info 与 debug 等价，输出全部 SQL）。
// 留空表示不覆盖，沿用框架默认行为（只输出慢查询与错误）。
type GormSQLLog struct {
	Manager *Manager `autowire:"?"`
	Level   string   `value:"${database.log-level:=}"`
}

// Run 设置各实例的 gorm Logger。
func (g *GormSQLLog) Run(ctx context.Context) error {
	if g.Manager == nil {
		return nil
	}
	if strings.TrimSpace(g.Level) == "" {
		return nil
	}
	l, err := NewGormSQLLogger(g.Level)
	if err != nil {
		return err
	}
	for _, db := range g.Manager.DBs() {
		db.Logger = l
		log.Debug(ctx, log.TagAppDef, "gorm sql logger installed", func() []log.Field {
			return []log.Field{log.String("level", g.Level)}
		})
	}
	return nil
}

// NewGormSQLLogger 构造把 GORM 日志桥接到 go-spring 日志系统的 logger。
// GORM 没有 Debug 级别，一次 SQL 一条日志对应的是 logger.Info；
// 每行输出走 log.Debug，因此是否显示取决于全局日志级别（根 logger 需为 DEBUG）。
func NewGormSQLLogger(level string) (logger.Interface, error) {
	lv, err := parseGormLogLevel(level)
	if err != nil {
		return nil, err
	}
	return logger.New(gormLogWriter{}, logger.Config{
		SlowThreshold:             gormSlowThreshold,
		LogLevel:                  lv,
		IgnoreRecordNotFoundError: true,
		// false 表示 SQL 里带真实参数值（调试用）；带敏感字段时建议置 true 走占位符。
		ParameterizedQueries: false,
		Colorful:             false,
	}), nil
}

// parseGormLogLevel 把配置里的级别字符串转换为 gorm 的日志级别。
func parseGormLogLevel(s string) (logger.LogLevel, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug", "info":
		return logger.Info, nil
	case "warn", "warning":
		return logger.Warn, nil
	case "error":
		return logger.Error, nil
	case "silent", "off", "none":
		return logger.Silent, nil
	}
	return logger.Silent, fmt.Errorf("数据库 gorm 日志级别 %q 非法，可选值: silent/error/warn/info/debug", s)
}

// gormLogWriter 实现 gorm logger.Writer，把 GORM 的每条消息转到 go-spring 日志，
// 使其落到已配置的 appender（控制台 / 文件），而不是写到标准输出。
type gormLogWriter struct{}

// Printf 实现 gorm logger.Writer 接口。
// GORM 每条消息是“调用位置 + 耗时/影响行数 + SQL”的多行文本，
// 这里压缩成单行（用 " | " 连接），避免日志布局把换行转义成字面 \n。
func (gormLogWriter) Printf(format string, args ...any) {
	lines := strings.Split(fmt.Sprintf(format, args...), "\n")
	parts := make([]string, 0, len(lines))
	for _, line := range lines {
		if s := strings.TrimSpace(line); s != "" {
			parts = append(parts, s)
		}
	}
	msg := strings.Join(parts, " | ")
	if msg == "" {
		return
	}
	log.Debug(context.Background(), log.TagAppDef, "gorm sql", func() []log.Field {
		return []log.Field{log.String("msg", msg)}
	})
}
