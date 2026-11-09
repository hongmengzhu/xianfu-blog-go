package repositoryIdp

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityIdp"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(RamIdentitySourceRepository))

	gs.Provide(new(support.BaseService[RamIdentitySourceRepository]))
}

type RamIdentitySourceRepository struct {
	repositoryPg.BaseRepository[entityIdp.IdpIdentitySourceEntity, int64]
	//
}
