package configModel

import (
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/domainBasic/service/configCache"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/repositoryBasic"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(Sp))
}

type Sp struct {
	repModel      *repositoryBasic.BasicConfigModelRepository       `autowire:"?"`
	repField      *repositoryBasic.BasicConfigModelFieldsRepository `autowire:"?"`
	repRule       *repositoryBasic.BasicConfigModelRulesRepository  `autowire:"?"`
	repFieldRules *repositoryBasic.BasicModelRulesRepository        `autowire:"?"`
	repConfig     *repositoryBasic.BasicConfigRepository            `autowire:"?"`
	repConfigList *repositoryBasic.BasicConfigListRepository        `autowire:"?"`
	repModule     *repositoryBasic.BasicModuleRepository            `autowire:"?"`
	cache         *configCache.Sp                                   `autowire:"?"`
}
