package controllerPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/authPg"
)

type SpManageAuth struct {
	Sp *authPg.GroupManageMiddlewareSp `autowire:"?"`
}
