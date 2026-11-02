package repositoryRam

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(RamLevelRepository))

	gs.Provide(new(support.BaseService[RamLevelRepository]))
}

type RamLevelRepository struct {
	repositoryPg.BaseRepository[entityRam.RamLevelEntity, int64]
}
