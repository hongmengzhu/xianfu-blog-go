package repositoryBlog

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityBlog"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(BlogArticleRepository))

	gs.Provide(new(support.BaseService[BlogArticleRepository]))
}

type BlogArticleRepository struct {
	repositoryPg.BaseRepository[entityBlog.BlogArticleEntity, int64]
}
