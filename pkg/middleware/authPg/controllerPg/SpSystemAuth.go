package controllerPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/authPg"
)

type SpSystemAuth struct {
	Sp *authPg.GroupSystemMiddlewareSp `autowire:"?"`
}
