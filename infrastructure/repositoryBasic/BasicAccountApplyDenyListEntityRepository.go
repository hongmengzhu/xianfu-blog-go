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
	gs.Provide(new(BasicAccountApplyDenyListEntityRepository))

	gs.Provide(new(support.BaseService[BasicAccountApplyDenyListEntityRepository]))
}

type BasicAccountApplyDenyListEntityRepository struct {
	repositoryPg.BaseRepository[entityBasic.BasicAccountApplyDenyListEntity, int64]
}

func (c *BasicAccountApplyDenyListEntityRepository) FindByExprAndIdNot(ctx context.Context, name string, id string) (info *entityBasic.BasicTagsRelationEntity, result bool) {
	tx := c.Db().WithContext(ctx).Where("expr=?", name).Where("id <> ?", id).First(&info)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}
