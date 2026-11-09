package repositoryIdp

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityIdp"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(IdpIdentitySourceRepository))

	gs.Provide(new(support.BaseService[IdpIdentitySourceRepository]))
}

type IdpIdentitySourceRepository struct {
	repositoryPg.BaseRepository[entityIdp.IdpIdentitySourceEntity, int64]
	//
}
