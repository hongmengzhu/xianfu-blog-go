package repositoryRam

import (
	"context"

	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(RamResourceRelationRepository))

	gs.Provide(new(support.BaseService[RamResourceRelationRepository]))
}

type RamResourceRelationRepository struct {
	repositoryPg.BaseRepository[entityRam.RamResourceRelationEntity, int64]
	//
}

func (c *RamResourceRelationRepository) DeleteByAuthorityId(ctx context.Context, code int64) {
	c.DbModel().WithContext(ctx).Where("authority_id=?", code).Delete(&entityRam.RamResourceRelationEntity{})
}

func (c *RamResourceRelationRepository) FindByMark(ctx context.Context, code string) (info *entityRam.RamResourceRelationEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("mark=?", code).First(&info)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}
