package repositoryTc

import (
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityTc"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(TcOrganizationRepository))

	gs.Provide(new(support.BaseService[TcOrganizationRepository]))
}

type TcOrganizationRepository struct {
	repositoryPg.BaseRepository[entityTc.TcOrganizationEntity, int64]
}
