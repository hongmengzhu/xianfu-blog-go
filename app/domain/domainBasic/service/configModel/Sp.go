package configModel

import (
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/domainBasic/service/configCache"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/repositoryB72"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/repositoryBasic"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(Sp))
}

type Sp struct {
	repModel      *repositoryB72.B72ConfigModelRepository       `autowire:"?"`
	repField      *repositoryB72.B72ConfigModelFieldsRepository `autowire:"?"`
	repRule       *repositoryB72.B72ConfigModelRulesRepository  `autowire:"?"`
	repFieldRules *repositoryBasic.BasicModelRulesRepository    `autowire:"?"`
	repConfig     *repositoryB72.B72ConfigRepository            `autowire:"?"`
	repConfigList *repositoryB72.B72ConfigListRepository        `autowire:"?"`
	repModule     *repositoryBasic.BasicModuleRepository        `autowire:"?"`
	cache         *configCache.Sp                               `autowire:"?"`
}
