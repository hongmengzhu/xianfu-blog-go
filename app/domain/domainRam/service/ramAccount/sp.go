package ramAccount

import (
	repositoryRam2 "github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/repositoryRam"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(Sp))
}

type Sp struct {
	accDb   *repositoryRam2.RamAccountRepository              `autowire:"?"`
	authDb  *repositoryRam2.RamAccountAuthorizationRepository `autowire:"?"`
	depDb   *repositoryRam2.RamDepartmentRepository           `autowire:"?"`
	roleDb  *repositoryRam2.RamRoleRepository                 `autowire:"?"`
	teamDb  *repositoryRam2.RamTeamRepository                 `autowire:"?"`
	groupDb *repositoryRam2.RamGroupRepository                `autowire:"?"`
	levelDb *repositoryRam2.RamLevelRepository                `autowire:"?"`
}
