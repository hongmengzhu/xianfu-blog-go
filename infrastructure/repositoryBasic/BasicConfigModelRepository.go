package repositoryBasic

import (
	"context"

	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(BasicConfigModelRepository))

	gs.Provide(new(support.BaseService[BasicConfigModelRepository]))
}

type BasicConfigModelRepository struct {
	repositoryPg.BaseRepository[entityBasic.BasicConfigModelEntity, int64]
}

func (c *BasicConfigModelRepository) FindByModel(ctx context.Context, code string) (info *entityBasic.BasicConfigModelEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("model=?", code).First(&info)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}

func (c *BasicConfigModelRepository) FindByModelAndIdNot(ctx context.Context, code string, id string) (info *entityBasic.BasicConfigModelEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("model=?", code).Where("id <>?", id).First(&info)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}
