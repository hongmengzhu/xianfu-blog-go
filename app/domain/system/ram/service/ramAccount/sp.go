package ramAccount

import (
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryRam"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(Sp))
}

type Sp struct {
	accDb   *repositoryRam.RamAccountRepository              `autowire:"?"`
	authDb  *repositoryRam.RamAccountAuthorizationRepository `autowire:"?"`
	depDb   *repositoryRam.RamDepartmentRepository           `autowire:"?"`
	roleDb  *repositoryRam.RamRoleRepository                 `autowire:"?"`
	teamDb  *repositoryRam.RamTeamRepository                 `autowire:"?"`
	groupDb *repositoryRam.RamGroupRepository                `autowire:"?"`
	levelDb *repositoryRam.RamLevelRepository                `autowire:"?"`
}
