package configBasic

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/cachePg/bigCachePg"
	repositoryBasic2 "github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/repositoryBasic"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(Sp))
}

type Sp struct {
	rdt *bigCachePg.BatchString `autowire:"?"`

	repModel      *repositoryBasic2.BasicConfigModelRepository       `autowire:"?"`
	repEvent      *repositoryBasic2.BasicConfigEventRepository       `autowire:"?"`
	repEventField *repositoryBasic2.BasicConfigEventFieldsRepository `autowire:"?"`
	repConfig     *repositoryBasic2.BasicConfigRepository            `autowire:"?"`
	repConfigList *repositoryBasic2.BasicConfigListRepository        `autowire:"?"`
}
