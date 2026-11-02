package eventBasicEvent

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/cachePg/bigCachePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/repositoryBasic"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(Sp))
}

type Sp struct {
	rdt            *bigCachePg.BatchString                           `autowire:"?"`
	repModel       *repositoryBasic.BasicConfigModelRepository       `autowire:"?"`
	repModelFields *repositoryBasic.BasicConfigModelFieldsRepository `autowire:"?"`
	repEvent       *repositoryBasic.BasicConfigEventRepository       `autowire:"?"`
	repEventFields *repositoryBasic.BasicConfigEventFieldsRepository `autowire:"?"`
	repRules       *repositoryBasic.BasicModelRulesRepository        `autowire:"?"`
}
