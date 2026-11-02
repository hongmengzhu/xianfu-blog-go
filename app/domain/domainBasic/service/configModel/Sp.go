package configModel

import (
	repositoryBasic2 "github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/repositoryBasic"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(Sp))
}

type Sp struct {
	repModel  *repositoryBasic2.BasicConfigModelRepository       `autowire:"?"`
	repField  *repositoryBasic2.BasicConfigModelFieldsRepository `autowire:"?"`
	repRule   *repositoryBasic2.BasicConfigModelRulesRepository  `autowire:"?"`
	repModule *repositoryBasic2.BasicModuleRepository            `autowire:"?"`
}
