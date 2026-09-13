package listennerBlog

import (
	"context"

	"github.com/hongmengzhu/xianfu-blog-go/app/event/blog/model/modEventBlogArticleCategory"
	"github.com/hongmengzhu/xianfu-blog-go/app/event/blog/service/articleBlogEvent"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/event"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
	_ "go-spring.org/spring/gs"
)

func init() {
	// 注册订阅者 Bean，导出为 event.Listener 供容器自动收集。
	gs.Provide(new(ArticleCategoryCacheListener)).Export(gs.As[event.Listener]())
}

// ArticleCategoryCacheListener 文章分类处理
type ArticleCategoryCacheListener struct {
	sp *articleBlogEvent.Sp `autowire:"?"`
}

func (c *ArticleCategoryCacheListener) Register(bus event.Bus) {
	// 创建适配器
	event.Subscribe(bus, func(ctx context.Context, e modEventBlogArticleCategory.CacheDto) error {
		err := articleBlogEvent.NewCategoryCache(c.sp, e).Processor(context.Background())
		if nil != err {
			log.Errorf(ctx, log.TagAppDef, "博客.分类.缓存:%+v", err)
		}
		return nil
	})
}

// Run 启动加载
//
//	@Description:
//	@receiver c
//	@param ctx
//func (c *ArticleCategoryCacheListener) Run(ctx context.Context) error {
//	log.Infof(ctx, log.TagAppDef, "[init].listener.[博客.分类.缓存]===================")
//	//博客文章 分类
//	eventBus.RegisterEvent(constEventBusPg.BlogArticleCategoryCache).RegisterSubscribe(constEventBusPg.BlogArticleCategoryCache, func(message any, _ core.EventArgs) {
//		//log.Infof(ctx, log.TagAppDef,"[init].listener.[博客.分类.缓存]22===================")
//		dto := message.(modEventBlogArticleCategory.CacheDto)
//		//log.Infof(ctx, log.TagAppDef,"dto=%+v", dto)
//		err := articleBlogEvent.NewCategoryCache(c.sp, dto).Processor(context.Background())
//		if nil != err {
//			log.Errorf(ctx, log.TagAppDef, "博客.分类.缓存:%+v", err)
//		}
//		message = nil
//	})
//	return nil
//}
