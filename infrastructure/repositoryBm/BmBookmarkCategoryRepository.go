package repositoryBm

import (
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBm"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(BmBookmarkCategoryRepository))

	gs.Provide(new(support.BaseService[BmBookmarkCategoryRepository]))
}

type BmBookmarkCategoryRepository struct {
	repositoryPg.BaseRepository[entityBm.BmBookmarkCategoryEntity, int64]
}
