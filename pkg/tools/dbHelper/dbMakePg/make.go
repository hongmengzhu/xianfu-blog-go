package dbMakePg

import (
	"context"
	"fmt"
	"time"

	"github.com/duke-git/lancet/v2/datetime"
	"github.com/pangu-2/go-tools/tools/dbPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"go-spring.org/log"
	"gorm.io/gorm"
)

var lastTime time.Time
var tableList = make([]string, 0)

// GetTables 获取表名称
//
//	@Description:
//	@param db
//	@return []string
func GetTables(db *gorm.DB) []string {
	//获取表名称
	tableList2, err := db.Migrator().GetTables()
	if nil != err {
		panic(err)
	}
	return tableList2
}

// GetTablesByTime 获取表名称
//
//	@Description:
//	@param db
//	@param minute
//	@return []string
func GetTablesByTime(db *gorm.DB, minute int64) []string {
	if len(tableList) == 0 {
		tableList = GetTables(db)
	}
	newTime := datetime.AddMinute(lastTime, minute)
	if lastTime.IsZero() {
		//获取表名称
		tableList = GetTables(db)
		lastTime = time.Now()
	} else if newTime.Before(time.Now()) {
		//获取表名称
		tableList = GetTables(db)
		lastTime = time.Now()
	}
	log.Info(context.Background(), log.TagAppDef, "tableList", log.Any("len", len(tableList)))
	return tableList
}

// GetTablesByTime3Minute 获取表名称
//
//	@Description:
//	@param db
//	@return []string
func GetTablesByTime3Minute(db *gorm.DB) []string {
	return GetTablesByTime(db, 3)
}

// MakeTable 生成表，创建表注释
//
//	@Description:
//	@param db
//	@param tmp
//	@param tableName
//	@param tableComment
func MakeTable(db *gorm.DB, tmp interface{}, tableName, tableComment string) {
	//如果没有创建，那么创建
	err := db.AutoMigrate(tmp)
	if err != nil {
		//panic(err)
		log.Error(context.Background(), log.TagAppDef, err, "MakeTable err")
		return
	}
	if strPg.IsNotBlank(tableComment) {
		// 创建备注
		db.Raw(fmt.Sprintf("COMMENT ON TABLE %s IS '%s';", tableName, tableComment)).Row()
	}
}

// MakeSequenceSql
//
//	@Description: 生成 更新序号 开始值,带判断的 sql
//	@receiver b
//	@param table
//	@param seq
//	@return string
func MakeSequenceSql(table string, seq int64) string {
	// 仅当序列当前值小于目标值时才 setval，避免回退已用序列；
	// is_called=true 表示 nextval 返回 seq+1，保留 seq 作为起始号段
	return fmt.Sprintf(
		"SELECT setval('%s_id_seq', GREATEST((SELECT last_value FROM %s_id_seq), %d), true);",
		table, table, seq,
	)
}

// MakeMysqlAutoIncSql MySQL生成修改自增起始SQL
func MakeMysqlAutoIncSql(table string, start int64) string {
	return fmt.Sprintf("ALTER TABLE %s AUTO_INCREMENT = %d;", table, start)
}

func MakeDataBaseAutoInc(dialect, table string, start int64) (sql string) {
	switch dialect {
	case "postgres":
		sql = MakeSequenceSql(table, start)
	case "mysql":
		sql = MakeMysqlAutoIncSql(table, start)
	}
	return
}

// SyncSequences 按表内最大主键校正自增序列（PostgreSQL）
//
//	@Description: 初始化数据时存在显式指定主键（ID=1、2...）的插入，PostgreSQL
//	不会因此推进序列，导致序列停在起始值、落后于 max(id)，后续 INSERT 取到的
//	id 与已有行冲突，报 duplicate key value violates unique constraint "..._pkey"
//	（SQLSTATE 23505）。这里在启动初始化完成后统一把序列推进到 max(id)+1，
//	保证之后正常分配的 id 不会撞车。只前进不后退，避免破坏已保留的号段。
//	MySQL 的自增值会随显式插入自动前进，无需处理。
//	@param db gorm 连接
//	@return error
func SyncSequences(db *gorm.DB) error {
	if nil == db {
		return nil
	}
	if "postgres" != db.Dialector.Name() {
		return nil
	}
	// 所有带 nextval 默认值的自增列
	type seqColumn struct {
		TableName  string `gorm:"column:table_name"`
		ColumnName string `gorm:"column:column_name"`
	}
	columns := make([]seqColumn, 0)
	err := db.Raw(`SELECT c.relname AS table_name, a.attname AS column_name
		FROM pg_class c
		JOIN pg_attribute a ON a.attrelid = c.oid AND a.attnum > 0 AND NOT a.attisdropped
		JOIN pg_attrdef d ON d.adrelid = c.oid AND d.adnum = a.attnum
		WHERE c.relkind = 'r'
		  AND c.relnamespace = 'public'::regnamespace
		  AND pg_get_expr(d.adbin, d.adrelid) LIKE 'nextval%'`).Scan(&columns).Error
	if nil != err {
		log.Error(context.Background(), log.TagAppDef, err, "查询自增序列失败")
		return err
	}
	fixed := make([]string, 0)
	for _, column := range columns {
		var seqName string
		err = db.Raw(`SELECT pg_get_serial_sequence(?, ?)`, column.TableName, column.ColumnName).Scan(&seqName).Error
		if nil != err || strPg.IsBlank(seqName) {
			continue
		}
		var lastValue, maxId int64
		var isCalled bool
		row := db.Raw(`SELECT last_value, is_called FROM ?::regclass`, seqName).Row()
		if nil == row {
			continue
		}
		if err = row.Scan(&lastValue, &isCalled); nil != err {
			continue
		}
		err = db.Raw(fmt.Sprintf(`SELECT COALESCE(MAX("%s"), 0) FROM "%s"`, column.ColumnName, column.TableName)).Scan(&maxId).Error
		if nil != err {
			continue
		}
		// 序列下一次将要分配的值
		next := lastValue
		if isCalled {
			next = lastValue + 1
		}
		if maxId < next {
			// 序列已领先最大主键，无需调整
			continue
		}
		// is_called=false 表示 setval 之后 nextval 直接返回该值
		err = db.Exec("SELECT setval(?::regclass, ?, false)", seqName, maxId+1).Error
		if nil != err {
			log.Error(context.Background(), log.TagAppDef, err, "校正自增序列失败",
				log.Any("table", column.TableName), log.Any("column", column.ColumnName))
			continue
		}
		fixed = append(fixed, fmt.Sprintf("%s.%s: %d -> %d", column.TableName, column.ColumnName, next, maxId+1))
	}
	if len(fixed) > 0 {
		log.Info(context.Background(), log.TagAppDef, "[init].[自增序列修正]", log.Any("fixed", fixed))
	} else {
		log.Info(context.Background(), log.TagAppDef, "[init].[自增序列修正] 无需调整")
	}
	return nil
}

// MakeJsonOsContainsCond 按数据库方言生成 os JSON 字段内数组包含指定元素的查询条件
//
//	@Description:
//	@param db gorm 数据库连接（用于 Dialector.Name() 判断方言）
//	@param key os 内的 JSON 字段名，如 departments/roles/levels/groups/teams
//	@param val 要匹配的元素值
//	@return 条件 SQL 与绑定参数
func MakeJsonOsContainsCond(db *gorm.DB, key string, val string) (string, []any) {
	switch db.Dialector.Name() {
	case "mysql":
		// MySQL 5.7+：JSON_CONTAINS(目标, 候选值[, 路径])，JSON_QUOTE 保证值按 JSON 字符串正确转义
		return "JSON_CONTAINS(os, JSON_QUOTE(?), '$." + key + "')", []any{val}
	case "sqlite":
		// SQLite：无 JSON_CONTAINS，用 json_each 展开数组后匹配元素（路径不存在时 json_extract 返回 NULL，EXISTS 为 false）
		return "EXISTS (SELECT 1 FROM json_each(json_extract(os, '$." + key + "')) WHERE json_each.value = ?)", []any{val}
	default:
		// 默认 postgres：os->'key' @> '["val"]'
		return "os->'" + key + "' @> ? ", []any{dbPg.StrToArrayJsonExpr(val)}
	}
}
