package repositoryB72

import (
	"context"
	"errors"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityB72"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
	"gorm.io/gorm"
)

func init() {
	gs.Provide(new(B72ConfigModelRulesRepository))

	gs.Provide(new(support.BaseService[B72ConfigModelRulesRepository]))
}

type B72ConfigModelRulesRepository struct {
	repositoryPg.BaseRepository[entityB72.B72ConfigModelRulesEntity, int64]
}

// FindAllByModelNo 查询模型下所有模型级规则。
func (c *B72ConfigModelRulesRepository) FindAllByModelNo(ctx context.Context, modelNo string) (info []*entityB72.B72ConfigModelRulesEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("model_no=?", modelNo).Order("sort asc,create_at").Find(&info)
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

// DeleteAllByModelNo 物理删除模型下所有模型级规则。
func (c *B72ConfigModelRulesRepository) DeleteAllByModelNo(ctx context.Context, modelNo string) error {
	tx := c.DbModel().WithContext(ctx).Where("model_no=?", modelNo).Delete(&c.Entity)
	return tx.Error
}
