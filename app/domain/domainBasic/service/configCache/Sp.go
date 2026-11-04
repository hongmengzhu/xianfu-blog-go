package configCache

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/cachePg/bigCachePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/repositoryBasic"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(Sp))
}

// Sp 配置缓存依赖。配置缓存统一使用 tenant_no + model 作为 key，
// 任何配置或模型字段变更后都必须调用 Refresh 使缓存失效/重建。
type Sp struct {
	Rdt       *bigCachePg.BatchString                `autowire:"?"`
	RepConfig *repositoryBasic.BasicConfigRepository `autowire:"?"`
}
