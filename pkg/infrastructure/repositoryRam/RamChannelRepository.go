package repositoryRam

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(RamChannelRepository))

	gs.Provide(new(support.BaseService[RamChannelRepository]))
}

type RamChannelRepository struct {
	repositoryPg.BaseRepository[entityRam.RamChannelEntity, int64]
}
