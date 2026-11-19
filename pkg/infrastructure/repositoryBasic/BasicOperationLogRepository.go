package repositoryBasic

import (
	"context"
	"errors"
	"time"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/spring/gs"
	"gorm.io/gorm"
)

func init() {
	gs.Provide(new(BasicOperationLogRepository))

	gs.Provide(new(support.BaseService[entityBasic.BasicOperationLogEntity]))
}

type BasicOperationLogRepository struct {
	repositoryPg.BaseRepository[entityBasic.BasicOperationLogEntity, int64]
}

// CreateLogsBatch 批量写入操作日志；batchSize 建议 100。
// 审计日志是追加写模型，此处只负责落库，不做任何更新/回写。
func (r *BasicOperationLogRepository) CreateLogsBatch(ctx context.Context, logs []*entityBasic.BasicOperationLogEntity, batchSize int) error {
	if len(logs) == 0 {
		return nil
	}
	if batchSize <= 0 {
		batchSize = 100
	}
	return r.DbModel().WithContext(ctx).CreateInBatches(logs, batchSize).Error
}

// DeleteBefore 按批删除 cutoff 之前的日志，返回删除条数。
// 单批上限 batchSize，循环执行直到没有更多可删，避免长事务与大锁。
func (r *BasicOperationLogRepository) DeleteBefore(ctx context.Context, cutoff time.Time, batchSize int) (int64, error) {
	if batchSize <= 0 {
		batchSize = 5000
	}
	var total int64
	for {
		sub := r.DbModel().WithContext(ctx).
			Model(&entityBasic.BasicOperationLogEntity{}).
			Select("id").
			Where("operate_at < ?", cutoff).
			Limit(batchSize)
		tx := r.DbModel().WithContext(ctx).
			Where("id IN (?)", sub).
			Delete(&entityBasic.BasicOperationLogEntity{})
		if tx.Error != nil {
			return total, tx.Error
		}
		total += tx.RowsAffected
		if tx.RowsAffected < int64(batchSize) {
			break
		}
	}
	return total, nil
}

// OperationLogFilter 统计/聚合过滤条件。
type OperationLogFilter struct {
	TenantNo   string
	TypeDomain string
	Module     string
	Action     string
	OperatorNo string
	TargetType string
	TargetNo   string
	Keyword    string
	Result     int8
	StartAt    *time.Time
	EndAt      *time.Time
}

func (f OperationLogFilter) Apply(db *gorm.DB) *gorm.DB {
	if f.TenantNo != "" {
		db = db.Where("tenant_no = ?", f.TenantNo)
	}
	if f.TypeDomain != "" {
		db = db.Where("type_domain = ?", f.TypeDomain)
	}
	if f.Module != "" {
		db = db.Where("module = ?", f.Module)
	}
	if f.Action != "" {
		db = db.Where("action = ?", f.Action)
	}
	if f.OperatorNo != "" {
		db = db.Where("operator_no = ?", f.OperatorNo)
	}
	if f.TargetType != "" {
		db = db.Where("target_type = ?", f.TargetType)
	}
	if f.TargetNo != "" {
		db = db.Where("target_no = ?", f.TargetNo)
	}
	if f.Keyword != "" {
		like := "%" + f.Keyword + "%"
		db = db.Where("no like ? OR description like ? OR operator_name like ? OR target_name like ?", like, like, like, like)
	}
	if f.Result != 0 {
		db = db.Where("result = ?", f.Result)
	}
	if f.StartAt != nil {
		db = db.Where("operate_at >= ?", *f.StartAt)
	}
	if f.EndAt != nil {
		db = db.Where("operate_at < ?", *f.EndAt)
	}
	return db
}

// OperationLogDayStat 按天统计。
type OperationLogDayStat struct {
	Day     string
	Total   int64
	Success int64
	Fail    int64
	Partial int64
	Unknown int64
}

// OperationLogGroupRow 维度聚合行。
type OperationLogGroupRow struct {
	Key   string
	Name  string
	Count int64
}

// OperationLogOperatorRow 操作人聚合行。
type OperationLogOperatorRow struct {
	OperatorNo      string
	OperatorAccount string
	OperatorName    string
	Count           int64
}

// StatByDay 按天统计趋势与总数。
func (r *BasicOperationLogRepository) StatByDay(ctx context.Context, f OperationLogFilter) ([]OperationLogDayStat, error) {
	rows := make([]OperationLogDayStat, 0)
	err := f.Apply(r.DbModel().WithContext(ctx).Model(&entityBasic.BasicOperationLogEntity{})).
		Select(`to_char(operate_at, 'YYYY-MM-DD') AS day,
			count(*) AS total,
			count(*) FILTER (WHERE result = 1) AS success,
			count(*) FILTER (WHERE result = 2) AS fail,
			count(*) FILTER (WHERE result = 3) AS partial,
			count(*) FILTER (WHERE result = 0) AS unknown`).
		Group("day").
		Order("day asc").
		Scan(&rows).Error
	return rows, err
}

// GroupBy 按维度聚合；field 仅接受白名单值。
func (r *BasicOperationLogRepository) GroupBy(ctx context.Context, f OperationLogFilter, field string, limit int) ([]OperationLogGroupRow, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	columns := map[string][2]string{
		"module":   {"module", "module_name"},
		"action":   {"action", "action_name"},
		"risk":     {"risk_level::text", ""},
		"result":   {"result::text", ""},
		"operator": {"operator_no", "operator_name"},
	}
	col, ok := columns[field]
	if !ok {
		return nil, errors.New("audit: unsupported group field " + field)
	}
	nameExpr, nameCol := "''", ""
	if col[1] != "" {
		nameCol = col[1]
		nameExpr = "max(" + nameCol + ")"
	}
	rows := make([]OperationLogGroupRow, 0)
	err := f.Apply(r.DbModel().WithContext(ctx).Model(&entityBasic.BasicOperationLogEntity{})).
		Select(col[0] + " AS key, " + nameExpr + " AS name, count(*) AS count").
		Group(col[0]).
		Order("count desc").
		Limit(limit).
		Scan(&rows).Error
	return rows, err
}

// TopOperators 操作人聚合，用于下拉。
func (r *BasicOperationLogRepository) TopOperators(ctx context.Context, f OperationLogFilter, keyword string, limit int) ([]OperationLogOperatorRow, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	db := f.Apply(r.DbModel().WithContext(ctx).Model(&entityBasic.BasicOperationLogEntity{}))
	if keyword != "" {
		like := "%" + keyword + "%"
		db = db.Where("operator_no like ? OR operator_account like ? OR operator_name like ?", like, like, like)
	}
	rows := make([]OperationLogOperatorRow, 0)
	err := db.
		Select("operator_no, max(operator_account) AS operator_account, max(operator_name) AS operator_name, count(*) AS count").
		Group("operator_no").
		Order("count desc").
		Limit(limit).
		Scan(&rows).Error
	return rows, err
}

// FindByPageForExport 导出用分页查询（按 operate_at 倒序，limit 由调用方控制）。
func (r *BasicOperationLogRepository) FindByPageForExport(ctx context.Context, f OperationLogFilter, limit int) ([]entityBasic.BasicOperationLogEntity, error) {
	if limit <= 0 || limit > 50000 {
		limit = 10000
	}
	rows := make([]entityBasic.BasicOperationLogEntity, 0)
	err := f.Apply(r.DbModel().WithContext(ctx)).
		Order("operate_at desc").
		Limit(limit).
		Find(&rows).Error
	return rows, err
}
