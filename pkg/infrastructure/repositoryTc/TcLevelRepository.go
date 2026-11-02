package repositoryTc

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityTc"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(TcLevelRepository))

	gs.Provide(new(support.BaseService[TcLevelRepository]))
}

type TcLevelRepository struct {
	repositoryPg.BaseRepository[entityTc.TcLevelEntity, int64]
}
