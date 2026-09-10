package controllerPg

import "github.com/hongmengzhu/xianfu-blog-go/middleware/authPg"

type SpManageAuth struct {
	Sp *authPg.GroupManageMiddlewareSp `autowire:""`
}
