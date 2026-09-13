package dbManagerPg

import (
	"context"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/configPg"
	"go-spring.org/log"
	gormcore "go-spring.org/starter-gorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Manager struct {
	database   configPg.Database `value:"${database}"`
	dbPostgres *gormcore.DB      `autowire:"postgres.primary?"`
	dbMysql    *gormcore.DB      `autowire:"mysql.primary?"`
	dbSqlite   *gormcore.DB      `autowire:"sqlite.primary?"`
}

func (c *Manager) DbPostgres() *gorm.DB {
	return c.dbPostgres.DB
}

func (c *Manager) DbMysql() *gorm.DB {
	return c.dbMysql.DB
}

func (c *Manager) DbSqlite() *gorm.DB {
	return c.dbSqlite.DB
}

func (c *Manager) Db() *gorm.DB {
	log.Debugf(context.Background(), log.TagAppDef, "Database Driver=%+v", c.database.Driver)
	if c.database.Driver == mysql.DefaultDriverName {
		return c.dbMysql.DB
	} else if c.database.Driver == "postgres" {
		return c.dbPostgres.DB
	}
	return c.dbSqlite.DB
}
