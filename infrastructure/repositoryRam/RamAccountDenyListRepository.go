package repositoryRam

import (
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(RamAccountDenyListRepository))

	gs.Provide(new(support.BaseService[RamAccountDenyListRepository]))
}

type RamAccountDenyListRepository struct {
	repositoryPg.BaseRepository[entityRam.RamAccountDenyListEntity, int64]
}
