package repositoryBasic

import (
	"context"
	"errors"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
	"gorm.io/gorm"
)

func init() {
	gs.Provide(new(BasicModelRulesRepository))

	gs.Provide(new(support.BaseService[BasicModelRulesRepository]))
}

type BasicModelRulesRepository struct {
	repositoryPg.BaseRepository[entityBasic.BasicModelRulesEntity, int64]
}

func (c *BasicModelRulesRepository) DeleteAllByValueNoAndIds(ctx context.Context, no string, ids []string) (info []*entityBasic.BasicModelRulesEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("id in ?", ids).Where("value_no=?", no).Delete(&c.Entity)
	if tx.Error != nil {
		// record not found 跳过日志
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Error(ctx, log.TagAppDef, tx.Error, "err")
		}
		return
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}

// FindAllByValueNo 查询模型/事件下所有字段规则。
func (c *BasicModelRulesRepository) FindAllByValueNo(ctx context.Context, valueNo string) (info []*entityBasic.BasicModelRulesEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("value_no=?", valueNo).Order("sort asc,create_at").Find(&info)
	if tx.Error != nil {
		if !errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			log.Error(ctx, log.TagAppDef, tx.Error, "err")
		}
		return
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}

// DeleteAllByValueNo 物理删除模型/事件下所有规则，用于设计器整体保存时重建。
func (c *BasicModelRulesRepository) DeleteAllByValueNo(ctx context.Context, valueNo string) error {
	tx := c.DbModel().WithContext(ctx).Where("value_no=?", valueNo).Delete(&c.Entity)
	return tx.Error
}
