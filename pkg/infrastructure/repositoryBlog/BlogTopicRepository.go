package repositoryBlog

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityBlog"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(BlogTopicRepository))

	gs.Provide(new(support.BaseService[BlogTopicRepository]))
}

type BlogTopicRepository struct {
	repositoryPg.BaseRepository[entityBlog.BlogTopicEntity, int64]
}
