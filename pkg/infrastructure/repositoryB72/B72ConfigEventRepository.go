package repositoryB72

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityB72"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(B72ConfigEventRepository))

	gs.Provide(new(support.BaseService[B72ConfigEventRepository]))
}

type B72ConfigEventRepository struct {
	repositoryPg.BaseRepository[entityB72.B72ConfigEventEntity, int64]
}
