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
	gs.Provide(new(BasicConfigModelRulesRepository))

	gs.Provide(new(support.BaseService[BasicConfigModelRulesRepository]))
}

type BasicConfigModelRulesRepository struct {
	repositoryPg.BaseRepository[entityBasic.BasicConfigModelRulesEntity, int64]
}

// FindAllByModelNo 查询模型下所有模型级规则。
func (c *BasicConfigModelRulesRepository) FindAllByModelNo(ctx context.Context, modelNo string) (info []*entityBasic.BasicConfigModelRulesEntity, result bool) {
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
func (c *BasicConfigModelRulesRepository) DeleteAllByModelNo(ctx context.Context, modelNo string) error {
	tx := c.DbModel().WithContext(ctx).Where("model_no=?", modelNo).Delete(&c.Entity)
	return tx.Error
}
