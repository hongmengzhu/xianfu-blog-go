package dbMakePg

import (
	"context"
	"time"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/configPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/tablePg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type CreateTable struct {
	Database  configPg.Database `value:"database"`
	DB        *gorm.DB
	tableList []string
}

// DbOpen
//
//	@Description: 打开数据库
//	@receiver c
//	@return rt
func (c *CreateTable) DbOpen() (rt rg.Rs[string]) {
	db, err := gorm.Open(postgres.Open(c.Database.URL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
	})
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.SetMaxIdleConns(5)
		sqlDB.SetMaxOpenConns(200)
		sqlDB.SetConnMaxLifetime(0)
	}
	c.DB = db
	if err != nil {
		log.Error(context.Background(), log.TagAppDef, err, "open gorm postgresql", log.Any("url", c.Database.URL))
		return rt.ErrorMessage(err.Error())
	}
	return rt.Ok()
}

func (c *CreateTable) Db() *gorm.DB {
	return c.DB
}

func (c *CreateTable) TableCreateOne(entity interface{}) (rt rg.Rs[string]) {
	isExists := c.Db().Migrator().HasTable(entity)
	if isExists {
		log.Info(context.Background(), log.TagAppDef, "表存在")
	} else {
		err2 := c.Db().AutoMigrate(entity)
		if err2 != nil {
			log.Error(context.Background(), log.TagAppDef, err2, "创建表异常")
			return
		}
	}

	return rt.Ok()
}

// 生成表，创建表注释
func (c *CreateTable) dbRun(tmp interface{}, tableName, tableComment string) {
	// 判断表是否已创建
	if !slice.Contain(c.tableList, tableName) {
		log.Info(context.Background(), log.TagAppDef, "创建表", log.Any("tableName", tableName))
		MakeTable(c.Db(), tmp, tableName, tableComment)

		log.Info(context.Background(), log.TagAppDef, "创建表[完成]", log.Any("tableName", tableName))
	} else {
		log.Info(context.Background(), log.TagAppDef, "表已存在", log.Any("tableName", tableName))
	}
}

// 生成表，创建表注释
func (c *CreateTable) dbRunByDb(db *gorm.DB, tmp interface{}, tableName, tableComment string) {
	// 判断表是否已创建
	if !slice.Contain(c.tableList, tableName) {
		log.Info(context.Background(), log.TagAppDef, "创建表", log.Any("tableName", tableName))
		MakeTable(db, tmp, tableName, tableComment)

		log.Info(context.Background(), log.TagAppDef, "创建表[完成]", log.Any("tableName", tableName))
	} else {
		log.Info(context.Background(), log.TagAppDef, "表已存在", log.Any("tableName", tableName))
	}
}

// TableCreateAll
//
//	@Description: 批量创建表
//	@receiver c
//	@param dst
//	@return rt
func (c *CreateTable) TableCreateAll(dst []interface{}) (rt rg.Rs[string]) {
	c.tableList = GetTablesByTime3Minute(c.Db())
	log.Debug(context.Background(), log.TagAppDef, "已存在的表", func() []log.Field {
		return []log.Field{log.Any("tableList", c.tableList)}
	})
	i := 0
	for _, item := range dst {
		name := tablePg.GetTableName(item)
		log.Debug(context.Background(), log.TagAppDef, "创建表: %s", func() []log.Field {
			return []log.Field{log.Any("tableName", name)}
		})
		comment := tablePg.GetTableComment(item)
		if strPg.IsNotBlank(name) {
			c.dbRun(item, name, comment)
		}
		if i%3 == 0 {
			time.Sleep(10000)
		}
		i++
	}

	return rt.Ok()
}

// TableCreateAllByTransaction
//
//	@Description: 批量创建表
//	@receiver c
//	@param dst
//	@return rt
func (c *CreateTable) TableCreateAllByTransaction(dst []interface{}) (rt rg.Rs[string]) {
	c.tableList = GetTablesByTime3Minute(c.Db())
	log.Debug(context.Background(), log.TagAppDef, "已存在的表: %+v", func() []log.Field {
		return []log.Field{log.Any("tableList", c.tableList)}
	})
	//
	err := c.Db().Transaction(func(tx *gorm.DB) error {
		i := 0
		for _, item := range dst {
			name := tablePg.GetTableName(item)
			log.Debug(context.Background(), log.TagAppDef, "创建表: %s", func() []log.Field {
				return []log.Field{log.Any("tableName", name)}
			})
			comment := tablePg.GetTableComment(item)
			if strPg.IsNotBlank(name) {
				c.dbRunByDb(tx, item, name, comment)
			}
			if i%3 == 0 {
				time.Sleep(10000)
			}
			i++
		}

		return nil
	})
	if err != nil {
		log.Error(context.Background(), log.TagAppDef, err, "创建表异常")
	}

	return rt.Ok()
}
