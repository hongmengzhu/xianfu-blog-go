package data

import (
	"context"
	"fmt"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/configPg"
	_ "go-spring.org/spring/gs"
)

// ZzBootstrap
// @Description: 显示服务启动信息
type ZzBootstrap struct {
	ser  configPg.Server `value:"${server}"`
	addr string          `value:"${spring.gin.server.addr}"`
}

func (b *ZzBootstrap) Run(ctx context.Context) error {
	fmt.Println()
	fmt.Printf("host: %+v\n", "localhost")
	fmt.Printf("port: %+v\n", b.addr)
	fmt.Printf("url: http://localhost%+v\n", b.addr)
	fmt.Println()
	return nil
}
