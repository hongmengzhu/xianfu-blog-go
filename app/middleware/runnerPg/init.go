package runnerPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/runnerPg/data"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/runnerPg/table"
	"go-spring.org/spring/gs"
)

func init() {
	//初始化 表
	gs.Provide(new(table.AInitTable)).Export(gs.As[gs.Runner]())
	// 发布者 监听
	gs.Provide(new(Publisher)).Export(gs.As[gs.Runner]())

	//初始化租户域名
	gs.Provide(new(data.InitTenantDomain)).Export(gs.As[gs.Runner]())
	//初始化 基础数据
	gs.Provide(new(data.IBasicData)).Export(gs.As[gs.Runner]())
	//

	//
	//附件 上传
	//gs.Provide(new(listenerBasic.AttachmentListener)).Export(gs.As[gs.Runner]())
	//gs.Provide(new(listenerBasic.TagsListener)).Export(gs.As[gs.Runner]())
	//gs.Provide(new(listenerBasic.EventCacheListener)).Export(gs.As[gs.Runner]())
	//gs.Provide(new(listenerBasic.RulesCacheListener)).Export(gs.As[gs.Runner]())
	//文章 分类 换成
	//gs.Provide(new(listennerBlog.ArticleCategoryCacheListener)).Export(gs.As[gs.Runner]())

	// ram相关
	//gs.Provide(new(listenerRam.RamListener)).Export(gs.As[gs.Runner]())
	//
	// 初始化会话密钥
	gs.Provide(new(data.InitSessionPubPrive)).Export(gs.As[gs.Runner]())
	// 基础
	gs.Provide(new(data.ZInitCacheBasic)).Export(gs.As[gs.Runner]())
	// 初始化博客缓存
	gs.Provide(new(data.ZInitCacheBlog)).Export(gs.As[gs.Runner]())
	// 初始化标签缓存
	gs.Provide(new(data.ZInitTagsCache)).Export(gs.As[gs.Runner]())
	//超管账号初始化
	gs.Provide(new(data.ZInitAccountAdmin)).Export(gs.As[gs.Runner]())
	//初始化Dipl缓存
	gs.Provide(new(data.ZInitDiplCache)).Export(gs.As[gs.Runner]())
	//显示服务启动信息
	gs.Provide(new(data.ZzBootstrap)).Export(gs.As[gs.Runner]())
}
