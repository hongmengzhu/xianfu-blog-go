package repositoryBasic

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(BasicTagsRepository))

	gs.Provide(new(support.BaseService[BasicTagsRepository]))
}

type BasicTagsRepository struct {
	repositoryPg.BaseRepository[entityBasic.BasicTagsEntity, int64]
}
