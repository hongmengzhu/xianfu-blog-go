package configBasic

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/cachePg/bigCachePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/repositoryB72"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(Sp))
}

type Sp struct {
	rdt *bigCachePg.BatchString `autowire:"?"`

	repModel      *repositoryB72.B72ConfigModelRepository       `autowire:"?"`
	repEvent      *repositoryB72.B72ConfigEventRepository       `autowire:"?"`
	repEventField *repositoryB72.B72ConfigEventFieldsRepository `autowire:"?"`
	repConfig     *repositoryB72.B72ConfigRepository            `autowire:"?"`
	repConfigList *repositoryB72.B72ConfigListRepository        `autowire:"?"`
}
