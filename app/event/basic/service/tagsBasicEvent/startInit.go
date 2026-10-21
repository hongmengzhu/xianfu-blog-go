package tagsBasicEvent

import (
	"context"
	"strings"

	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/cachePg/bigCachePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/pangu-2/go-tools/tools/cryptPg"
)

// StartInit
// @Description: 启动后初始化一些数据
type StartInit struct {
	rdb *bigCachePg.BatchString `autowire:"?"`

	sp *Sp `autowire:"?"`
}

func NewStartInit(sp *Sp) *StartInit {
	return &StartInit{
		sp:  sp,
		rdb: sp.rdt,
	}
}

func (c *StartInit) Processor(ctx context.Context) error {
	t := entityBasic.BasicTagsRelationEntity{
		State: enumStatePg.ENABLE.Index(),
	}
	infos := c.sp.TagRela.FindAll(ctx, t)
	if nil != infos && len(infos) > 0 {
		data := make(map[string]interface{})
		for _, item := range infos {
			//防止过长
			md5 := cryptPg.Md5(strings.TrimSpace(item.Name))
			tmp := item.CategoryRoot + ":" + md5
			data[tmp] = strings.TrimSpace(item.TagNo)
		}
		c.rdb.SetPipeline(ctx, data)
	}
	return nil
}
