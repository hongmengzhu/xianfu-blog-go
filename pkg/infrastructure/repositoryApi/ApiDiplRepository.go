package repositoryApi

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityApi"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(ApiDiplRepository))

	gs.Provide(new(support.BaseService[ApiDiplRepository]))
}

type ApiDiplRepository struct {
	repositoryPg.BaseRepository[entityApi.ApiDiplEntity, int64]
}
