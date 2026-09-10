package repositoryRam

import (
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(RamAccountDeviceRepository))

	gs.Provide(new(support.BaseService[RamAccountDeviceRepository]))
}

type RamAccountDeviceRepository struct {
	repositoryPg.BaseRepository[entityRam.RamAccountDeviceEntity, int64]
}
