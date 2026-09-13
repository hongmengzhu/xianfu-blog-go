package repositoryBasic

import (
	"context"
	"errors"

	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBasic"
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
			log.Errorf(ctx, log.TagAppDef, "err=%+v", tx.Error)
		}
		return
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}
