package repositoryBlog

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityBlog"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(BlogCollectCategoryRepository))

	gs.Provide(new(support.BaseService[BlogCollectCategoryRepository]))
}

type BlogCollectCategoryRepository struct {
	repositoryPg.BaseRepository[entityBlog.BlogCollectCategoryEntity, int64]
}
