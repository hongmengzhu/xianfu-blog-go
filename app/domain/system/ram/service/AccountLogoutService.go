package service

import (
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/holderPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/configPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(NewAccountLogoutService)
}

// AccountLogoutService 退出
// @Description:
type AccountLogoutService struct {
	dao *repositoryRam.RamAccountRepository `autowire:"?"`
	pg  configPg.Pg                         `value:"${pg}"`
}

func NewAccountLogoutService() *AccountLogoutService {
	return new(AccountLogoutService)
}

func (c *AccountLogoutService) Logout(holder holderPg.HolderPg) (rt rg.Rs[string]) {
	return rt.Ok()
}
