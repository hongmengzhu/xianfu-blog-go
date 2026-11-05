package modelRules

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/repositoryB72"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/repositoryBasic"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(Sp))
}

type Sp struct {
	repModel       *repositoryB72.B72ConfigModelRepository       `autowire:"?"`
	repModelFields *repositoryB72.B72ConfigModelFieldsRepository `autowire:"?"`
	repEvent       *repositoryB72.B72ConfigEventRepository       `autowire:"?"`
	repEventFields *repositoryB72.B72ConfigEventFieldsRepository `autowire:"?"`
	repRules       *repositoryBasic.BasicModelRulesRepository    `autowire:"?"`
}
