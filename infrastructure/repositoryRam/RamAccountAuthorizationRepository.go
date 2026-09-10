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

	gs.Provide(new(RamAccountAuthorizationRepository))

	gs.Provide(new(support.BaseService[RamAccountAuthorizationRepository]))
}

type RamAccountAuthorizationRepository struct {
	repositoryPg.BaseRepository[entityRam.RamAccountAuthorizationEntity, int64]
	//
}

func (c *RamAccountAuthorizationRepository) FindByTypePasswordANo(ctx context.Context, code string) (info *entityRam.RamAccountAuthorizationEntity, result bool) {
	tx := c.DbModel().WithContext(ctx).Where("type=?", "password").Where("ano=?", code).Find(&info)
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return nil, false
	}
	if 0 == tx.RowsAffected {
		return nil, false
	}
	return info, true
}

func (c *RamAccountAuthorizationRepository) DeleteByAno(ctx context.Context, code string) (result bool) {
	tx := c.DbModel().WithContext(ctx).Where("ano=?", code).Delete(&entityRam.RamAccountAuthorizationEntity{})
	if tx.Error != nil {
		log.Errorf(ctx, log.TagAppDef, "", tx.Error)
		return false
	}
	if 0 == tx.RowsAffected {
		return false
	}
	return true
}
