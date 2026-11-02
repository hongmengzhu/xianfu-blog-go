package repositoryRam

import (
	"context"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(RamResourceMenuRepository))

	gs.Provide(new(support.BaseService[RamResourceMenuRepository]))
}

type RamResourceMenuRepository struct {
	repositoryPg.BaseRepository[entityRam.RamResourceMenuEntity, int64]
	//
}

func (c *RamResourceMenuRepository) DeleteByMenuId(ctx context.Context, code int64) {
	c.DbModel().WithContext(ctx).Where("menu_id=?", code).Delete(&entityRam.RamResourceMenuEntity{})
}
