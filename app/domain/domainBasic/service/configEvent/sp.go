package configEvent

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/repositoryB72"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/repositoryBasic"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(Sp))
}

type Sp struct {
	repModel      *repositoryB72.B72ConfigModelRepository       `autowire:"?"`
	repField      *repositoryB72.B72ConfigModelFieldsRepository `autowire:"?"`
	repRule       *repositoryB72.B72ConfigModelRulesRepository  `autowire:"?"`
	repModule     *repositoryBasic.BasicModuleRepository        `autowire:"?"`
	repEvent      *repositoryB72.B72ConfigEventRepository       `autowire:"?"`
	repEventField *repositoryB72.B72ConfigEventFieldsRepository `autowire:"?"`
}
