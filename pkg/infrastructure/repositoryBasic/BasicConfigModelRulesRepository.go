package repositoryBasic

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/support"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(BasicConfigModelRulesRepository))

	gs.Provide(new(support.BaseService[BasicConfigModelRulesRepository]))
}

type BasicConfigModelRulesRepository struct {
	repositoryPg.BaseRepository[entityBasic.BasicConfigModelRulesEntity, int64]
}
