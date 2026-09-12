package modelRules

import (
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryBasic"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(Sp))
}

type Sp struct {
	repModel       *repositoryBasic.BasicConfigModelRepository       `autowire:"?"`
	repModelFields *repositoryBasic.BasicConfigModelFieldsRepository `autowire:"?"`
	repEvent       *repositoryBasic.BasicConfigEventRepository       `autowire:"?"`
	repEventFields *repositoryBasic.BasicConfigEventFieldsRepository `autowire:"?"`
	repRules       *repositoryBasic.BasicModelRulesRepository        `autowire:"?"`
}
