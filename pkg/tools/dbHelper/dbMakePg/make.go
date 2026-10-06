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
	log.Infof(context.Background(), log.TagAppDef, "tableList.len=%+v", len(tableList))
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
		log.Errorf(context.Background(), log.TagAppDef, err, "MakeTable err=%+v \n", err)
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
