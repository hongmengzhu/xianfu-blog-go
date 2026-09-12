package configBasic

import (
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/cachePg/rdsPg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(Sp))
}

type Sp struct {
	rdt *rdsPg.BatchString `autowire:"?"`

	repModel      *repositoryBasic.BasicConfigModelRepository       `autowire:"?"`
	repEvent      *repositoryBasic.BasicConfigEventRepository       `autowire:"?"`
	repEventField *repositoryBasic.BasicConfigEventFieldsRepository `autowire:"?"`
	repConfig     *repositoryBasic.BasicConfigRepository            `autowire:"?"`
	repConfigList *repositoryBasic.BasicConfigListRepository        `autowire:"?"`
}
