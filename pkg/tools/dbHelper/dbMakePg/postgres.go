package dbMakePg

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

// migratePGTimeColumn PG方言下，将指定表时间字段转为 timestamptz
func migratePGTimeColumn(db *gorm.DB) {
	if db.Dialector.Name() != "postgres" {
		return
	}

	tableCols := map[string][]string{
		"ram_group":           {"create_at", "update_at"},
		"api_dipl_category":   {"create_at", "update_at"},
		"api_dipl_access_key": {"create_at", "update_at", "expiry_date"},
	}

	for tableName, cols := range tableCols {
		for _, col := range cols {
			sql := fmt.Sprintf(`ALTER TABLE %s ALTER COLUMN %s TYPE timestamptz USING %s AT TIME ZONE 'UTC';`, tableName, col, col)
			err := db.Exec(sql).Error
			if err != nil {
				log.Printf("warn: alter table=%s column=%s to timestamptz, err=%v", tableName, col, err)
			}
		}
	}
}
