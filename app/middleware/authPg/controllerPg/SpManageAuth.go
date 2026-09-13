package controllerPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
)

type SpManageAuth struct {
	Sp *authPg.GroupManageMiddlewareSp `autowire:"?"`
}
