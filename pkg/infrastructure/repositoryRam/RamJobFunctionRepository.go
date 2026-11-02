package repositoryRam

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(RamJobFunctionRepository))

	gs.Provide(new(support.BaseService[RamJobFunctionRepository]))
}

type RamJobFunctionRepository struct {
	repositoryPg.BaseRepository[entityRam.RamJobFunctionEntity, int64]
}
