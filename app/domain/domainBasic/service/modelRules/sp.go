package modelRules

import (
	repositoryBasic2 "github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/repositoryBasic"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(Sp))
}

type Sp struct {
	repModel       *repositoryBasic2.BasicConfigModelRepository       `autowire:"?"`
	repModelFields *repositoryBasic2.BasicConfigModelFieldsRepository `autowire:"?"`
	repEvent       *repositoryBasic2.BasicConfigEventRepository       `autowire:"?"`
	repEventFields *repositoryBasic2.BasicConfigEventFieldsRepository `autowire:"?"`
	repRules       *repositoryBasic2.BasicModelRulesRepository        `autowire:"?"`
}
