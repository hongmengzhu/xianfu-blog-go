package controllerPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/authPg"
)

type SpSystemAuth struct {
	Sp *authPg.GroupSystemMiddlewareSp `autowire:""`
}
