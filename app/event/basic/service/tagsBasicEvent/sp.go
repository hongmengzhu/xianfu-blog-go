package tagsBasicEvent

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/cachePg/bigCachePg"
	repositoryBasic2 "github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/repositoryBasic"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(Sp))
}

type Sp struct {
	dao     *repositoryBasic2.BasicAttachmentRepository   `autowire:"?"`
	TagRela *repositoryBasic2.BasicTagsRelationRepository `autowire:"?"`
	TagCate *repositoryBasic2.BasicTagsCategoryRepository `autowire:"?"`
	TagsDb  *repositoryBasic2.BasicTagsRepository         `autowire:"?"`
	rdt     *bigCachePg.BatchString                       `autowire:"?"`
}
