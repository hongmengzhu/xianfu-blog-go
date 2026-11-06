package repositoryBm

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityBm"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(BmToolCategoryRepository))

	gs.Provide(new(support.BaseService[BmToolCategoryRepository]))
}

type BmToolCategoryRepository struct {
	repositoryPg.BaseRepository[entityBm.BmToolCategoryEntity, int64]
}
