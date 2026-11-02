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
	ser      configPg.Server `value:"${server}"`
	addr     string          `value:"${spring.gin.server.addr:=}"`
	addrEcho string          `value:"${spring.echo.server.addr:=}"`
}

func (b *ZzBootstrap) Run(ctx context.Context) error {
	fmt.Println()
	// echo 统一入口（对外服务端口，按域名分发模块）
	if "" != b.addrEcho {
		fmt.Printf("host: %+v\n", "localhost")
		fmt.Printf("port: %+v\n", b.addrEcho)
		fmt.Printf("url: http://localhost%+v\n", b.addrEcho)
	}
	// gin 引擎作为 echo 的兜底处理器，仅在配置了 spring.gin.server.addr 时展示
	if "" != b.addr {
		fmt.Println()
		fmt.Printf("gin host: %+v\n", "localhost")
		fmt.Printf("gin port: %+v\n", b.addr)
		fmt.Printf("gin url: http://localhost%+v\n", b.addr)
	}
	fmt.Println()
	return nil
}
