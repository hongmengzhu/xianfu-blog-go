package repositoryRam

import (
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(RamCasbinRepository))

	gs.Provide(new(support.BaseService[RamCasbinRepository]))
}

type RamCasbinRepository struct {
	repositoryPg.BaseRepository[gormadapter.CasbinRule, uint]
	//
}
