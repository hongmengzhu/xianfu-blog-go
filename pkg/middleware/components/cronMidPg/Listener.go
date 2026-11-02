package cronMidPg

import (
	"context"
	"fmt"
)

// Listener 启动
// @Description:
type Listener struct {
}

// OnAppStart 启动加载
//
//	@Description:
//	@receiver c
//	@param ctx
func (c *Listener) OnAppStart() {
	Start()
}

// OnAppStop 停止
//
//	@Description:
//	@receiver starter
//	@param ctx
func (c *Listener) OnAppStop(ctx context.Context) {
	err := Shutdown()
	if err != nil {
		fmt.Println(err)
	}
}
