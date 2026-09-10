package repositoryRam

import (
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(RamAccountSessionRepository))

	gs.Provide(new(support.BaseService[RamAccountSessionRepository]))
}

type RamAccountSessionRepository struct {
	repositoryPg.BaseRepository[entityRam.RamAccountSessionEntity, int64]
}
