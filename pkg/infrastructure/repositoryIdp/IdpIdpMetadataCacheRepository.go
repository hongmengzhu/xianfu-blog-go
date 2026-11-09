package repositoryIdp

import (
	"context"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityIdp"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(RamIdpMetadataCacheRepository))

	gs.Provide(new(support.BaseService[RamIdpMetadataCacheRepository]))
}

type RamIdpMetadataCacheRepository struct {
	repositoryPg.BaseRepository[entityIdp.IdpIdpMetadataCacheEntity, int64]
	//
}

// FindBySourceNo 按认证源编号查找元数据缓存
func (c *RamIdpMetadataCacheRepository) FindBySourceNo(ctx context.Context, sourceNo string) (info *entityIdp.IdpIdpMetadataCacheEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("source_no=?", sourceNo).First(&info)
	if tx.Error != nil || tx.RowsAffected == 0 {
		return nil, false
	}
	return info, true
}
