package authCasbinPg

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/manage/domainRam/service"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/holderPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/configPg"
)

func CasbinHandler(ctx *gin.Context, hpg holderPg.HolderPg, pg configPg.Pg, casbin *service.RamResourceCasbinService) error {
	if pg.Profiles.Active != "develop" {
		//获取请求的PATH
		path := ctx.Request.RequestURI
		obj := strings.TrimPrefix(path, "")
		// 获取请求方法
		act := ctx.Request.Method
		// 获取用户的角色
		sub := hpg.GetAccount().RoleNo
		e := casbin.Casbin() // 判断策略中是否存在
		success, _ := e.Enforce(sub, obj, act)
		if !success {
			return errors.New("权限不足")
		}
	}
	return nil
}
