package bmToolEvent

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/cachePg/bigCachePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/repositoryBm"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(Sp))
}

type Sp struct {
	rdt    *bigCachePg.BatchString                `autowire:"?"`
	catRep *repositoryBm.BmToolCategoryRepository `autowire:"?"`
}
