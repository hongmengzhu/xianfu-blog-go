package repositoryRam

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(RamAccountLoginLogRepository))

	gs.Provide(new(support.BaseService[RamAccountLoginLogRepository]))
}

type RamAccountLoginLogRepository struct {
	repositoryPg.BaseRepository[entityRam.RamAccountLoginLogEntity, int64]
}
