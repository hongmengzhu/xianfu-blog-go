package tagsBasicEvent

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/cachePg/bigCachePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/repositoryBasic"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(Sp))
}

type Sp struct {
	dao     *repositoryBasic.BasicAttachmentRepository   `autowire:"?"`
	TagRela *repositoryBasic.BasicTagsRelationRepository `autowire:"?"`
	TagCate *repositoryBasic.BasicTagsCategoryRepository `autowire:"?"`
	TagsDb  *repositoryBasic.BasicTagsRepository         `autowire:"?"`
	rdt     *bigCachePg.BatchString                      `autowire:"?"`
}
