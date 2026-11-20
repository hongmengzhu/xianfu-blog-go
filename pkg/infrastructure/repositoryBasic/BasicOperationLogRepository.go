package repositoryBasic

import (
	"context"
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
func (r *BasicOperationLogRepository) DeleteBefore(ctx context.Context, cutoff time.Time, batchSize int) (int64, error) {
	if batchSize <= 0 {
		batchSize = 5000
	}
	var total int64
	for {
		sub := r.DbModel().WithContext(ctx).
			Model(&entityBasic.BasicOperationLogEntity{}).
			Select("id").
			Where("create_at < ?", cutoff).
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

// OperationLogFilter 操作日志过滤条件。
type OperationLogFilter struct {
	TraceID    string
	Method     string
	Path       string
	TypeDomain string
	Module     string
	Action     string
	OperatorNo string
	IP         string
	Keyword    string
	StatusCode int
	StartAt    *time.Time
	EndAt      *time.Time
}

func (f OperationLogFilter) Apply(db *gorm.DB) *gorm.DB {
	if f.TraceID != "" {
		db = db.Where("trace_id = ?", f.TraceID)
	}
	if f.Method != "" {
		db = db.Where("method = ?", f.Method)
	}
	if f.Path != "" {
		db = db.Where("path = ?", f.Path)
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
	if f.IP != "" {
		db = db.Where("ip = ?", f.IP)
	}
	if f.StatusCode > 0 {
		db = db.Where("status_code = ?", f.StatusCode)
	}
	if f.Keyword != "" {
		like := "%" + f.Keyword + "%"
		db = db.Where("(trace_id LIKE ? OR path LIKE ? OR module LIKE ? OR action LIKE ? OR operator_no LIKE ? OR operator_account LIKE ? OR operator_name LIKE ? OR ip LIKE ? OR param LIKE ? OR request_body LIKE ?)",
			like, like, like, like, like, like, like, like, like, like)
	}
	if f.StartAt != nil {
		db = db.Where("create_at >= ?", *f.StartAt)
	}
	if f.EndAt != nil {
		db = db.Where("create_at < ?", *f.EndAt)
	}
	return db
}
