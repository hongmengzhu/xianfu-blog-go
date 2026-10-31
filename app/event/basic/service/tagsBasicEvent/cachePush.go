package tagsBasicEvent

import (
	"context"
	"strings"

	"github.com/hongmengzhu/xianfu-blog-go/app/models/event/basic/modEventBasicTags"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/cachePg/bigCachePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
)

// CachePush
// @Description: 缓存更新
type CachePush struct {
	rdt *bigCachePg.BatchString `autowire:"?"`
	sp  *Sp                     `autowire:"?"`
	dto modEventBasicTags.TagsCacheDto
}

func NewCachePush(sp *Sp, dto modEventBasicTags.TagsCacheDto) *CachePush {
	return &CachePush{
		sp:  sp,
		rdt: sp.rdt,
		dto: dto,
	}
}

func (c *CachePush) Processor(ctx context.Context) error {
	t := entityBasic.BasicTagsRelationEntity{
		State: enumStatePg.ENABLE.Index(),
	}
	if nil == c.dto.CategoryRoot {
		return nil
	}
	if len(c.dto.CategoryRoot) < 1 {
		return nil
	}
	infos, result := c.sp.TagRela.FindAllByCategoryRootIn(ctx, t, c.dto.CategoryRoot)
	if result {
		data := make(map[string]interface{})
		for _, item := range infos {
			tmp := item.CategoryRoot + ":" + strings.TrimSpace(item.Name)
			data[tmp] = strings.TrimSpace(item.TagNo)
		}
		c.rdt.SetPipeline(ctx, data)
	}
	return nil
}
