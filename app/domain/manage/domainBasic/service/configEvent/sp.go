package configEvent

import (
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryBasic"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(Sp))
}

type Sp struct {
	repModel      *repositoryBasic.BasicConfigModelRepository       `autowire:"?"`
	repField      *repositoryBasic.BasicConfigModelFieldsRepository `autowire:"?"`
	repRule       *repositoryBasic.BasicConfigModelRulesRepository  `autowire:"?"`
	repModule     *repositoryBasic.BasicModuleRepository            `autowire:"?"`
	repEvent      *repositoryBasic.BasicConfigEventRepository       `autowire:"?"`
	repEventField *repositoryBasic.BasicConfigEventFieldsRepository `autowire:"?"`
}
