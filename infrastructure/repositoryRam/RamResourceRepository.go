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
	gs.Provide(new(RamResourceRepository))

	gs.Provide(new(support.BaseService[RamResourceRepository]))
}

type RamResourceRepository struct {
	repositoryPg.BaseRepository[entityRam.RamResourceEntity, int64]
	//
}

func (c *RamResourceRepository) FindByParentNoRoot(ctx context.Context) (info []*entityRam.RamResourceEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("parent_no='' or parent_no is null ").Find(&info)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}

func (c *RamResourceRepository) FindByParentIdRoot(ctx context.Context) (info []*entityRam.RamResourceEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("parent_id='' ").Find(&info)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}

// CountByParentIdString
//
//	@Description: 统计
//	@receiver c
//	@param pid
//	@return info
//	@return result
func (c *RamResourceRepository) CountByParentIdString(ctx context.Context, pid string) (total int64, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("parent_id= ? ", pid).Count(&total)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return 0, false
	}
	if 0 == tx.RowsAffected {
		return 0, false
	}
	return total, true
}
