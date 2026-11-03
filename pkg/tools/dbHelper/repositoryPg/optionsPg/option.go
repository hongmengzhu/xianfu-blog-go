package optionsPg

import (
	"github.com/gin-gonic/gin"
	"github.com/pangu-2/go-tools/tools/dbPg/pagePg"
	"gorm.io/gorm"
)

type OptionParams struct {
	Ctx       *gin.Context
	Db        *gorm.DB
	Condition Condition
	Pageable  *pagePg.Pageable
}

type Option func(arg *OptionParams)

func WithOption(opt Option) Option {
	return opt
}
func WithDb(db *gorm.DB) Option {
	return func(arg *OptionParams) {
		arg.Db = db
	}
}
func WithCtx(ctx *gin.Context) Option {
	return func(arg *OptionParams) {
		arg.Ctx = ctx
	}
}

// WithPageable 设置分页参数
func WithPageable(pageable pagePg.Pageable) Option {
	return func(arg *OptionParams) {
		arg.Pageable = &pageable
	}
}

func WithScopes(page *pagePg.Pageable) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		size := int64(0)
		offset := int(0)
		if nil != page {
			if page.PageSize > 0 {
				size = page.PageSize
			}
			if page.PageNum > 0 {
				offset = int((page.PageNum - 1) * size)
			}
		}
		return db.Offset(offset).Limit(int(size))
	}
}
