package postgresqlPg

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/configPg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 数据库工厂 Postgresql
type Factory struct {
	database configPg.Database `value:"${database}"`
}

func (factory *Factory) CreateDB() (*gorm.DB, error) {
	log.Infof(context.Background(), log.TagAppDef, "[init].[Postgresql]===================")
	log.Debugf(context.Background(), log.TagAppDef, "数据库连接地址： enabled:%+v,URL:=>%s", factory.database.Enabled, factory.database.URL)
	//fmt.Printf("数据库连接地址： enabled:%+v,URL:=>%s", factory.database.Enabled, factory.database.URL)
	if !factory.database.Enabled {
		log.Debugf(context.Background(), log.TagAppDef, "未启用数据库： %s", factory.database.Enabled)
		return nil, nil
	}
	db, err := gorm.Open(postgres.Open(factory.database.URL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Errorf(context.Background(), log.TagAppDef, "open gorm postgresql %s error: %v", factory.database.URL, err)
		panic(errors.New(err.Error()))
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Errorf(context.Background(), log.TagAppDef, "failed to get database connection")
	}
	// 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(10)
	// 设置最大打开连接数
	sqlDB.SetMaxOpenConns(100)
	// 设置连接的最大生存时间
	sqlDB.SetConnMaxLifetime(time.Hour)
	//
	dbName := db.Migrator().CurrentDatabase()
	// check if dbPg exists
	stmt := fmt.Sprintf("SELECT * FROM pg_database WHERE datname = '%s';", dbName)
	rs := db.Raw(stmt)
	if rs.Error != nil {
		return nil, rs.Error
	}
	// if not create it
	var rec = make(map[string]interface{})
	if rs.Find(rec); len(rec) == 0 {
		stmt := fmt.Sprintf("CREATE DATABASE %s;", dbName)
		if rs := db.Exec(stmt); rs.Error != nil {
			return nil, rs.Error
		}

		// close dbPg connection
		sql, err := db.DB()
		defer func() {
			_ = sql.Close()
		}()
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}

func init() {
	gs.Provide(newClient, gs.TagArg("${database}")).
		Name("GormDb").
		Condition(
			gs.OnProperty("database.enabled").HavingValue("true").MatchIfMissing(),
			gs.OnMissingBean[*gorm.DB]("GormDB"),
		).
		Name("__default__")
	//gs.Group("${database.instances}", newClient, nil)
}

func newClient(c configPg.Database) (*gorm.DB, error) {
	log.Infof(context.Background(), log.TagAppDef, "[init].[Postgresql]===================")
	log.Debugf(context.Background(), log.TagAppDef, "数据库连接地址： enabled:%+v,URL:=>%s", c.Enabled, c.URL)
	//fmt.Printf("数据库连接地址： enabled:%+v,URL:=>%s", c.Enabled, c.URL)
	if !c.Enabled {
		log.Debugf(context.Background(), log.TagAppDef, "未启用数据库： %s", c.Enabled)
		return nil, nil
	}
	db, err := gorm.Open(postgres.Open(c.URL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Errorf(context.Background(), log.TagAppDef, "open gorm postgresql %s error: %v", c.URL, err)
		panic(errors.New(err.Error()))
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Errorf(context.Background(), log.TagAppDef, "failed to get database connection")
	}
	// 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(10)
	// 设置最大打开连接数
	sqlDB.SetMaxOpenConns(100)
	// 设置连接的最大生存时间
	sqlDB.SetConnMaxLifetime(time.Hour)
	//
	dbName := db.Migrator().CurrentDatabase()
	// check if dbPg exists
	stmt := fmt.Sprintf("SELECT * FROM pg_database WHERE datname = '%s';", dbName)
	rs := db.Raw(stmt)
	if rs.Error != nil {
		return nil, rs.Error
	}
	// if not create it
	var rec = make(map[string]interface{})
	if rs.Find(rec); len(rec) == 0 {
		stmt := fmt.Sprintf("CREATE DATABASE %s;", dbName)
		if rs := db.Exec(stmt); rs.Error != nil {
			return nil, rs.Error
		}

		// close dbPg connection
		sql, err := db.DB()
		defer func() {
			_ = sql.Close()
		}()
		if err != nil {
			return nil, err
		}
	}
	return db, nil
}
