package repositoryRam

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(RamIdentityProviderRepository))

	gs.Provide(new(support.BaseService[RamIdentityProviderRepository]))
}

type RamIdentityProviderRepository struct {
	repositoryPg.BaseRepository[entityRam.RamIdentityProviderEntity, int64]
	//
}
