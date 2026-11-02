package repositoryBasic

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(BasicConfigEventRepository))

	gs.Provide(new(support.BaseService[BasicConfigEventRepository]))
}

type BasicConfigEventRepository struct {
	repositoryPg.BaseRepository[entityBasic.BasicConfigEventEntity, int64]
}
