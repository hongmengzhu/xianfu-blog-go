package service

import (
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/core/blog/serviceCore"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/manage/domainBlog/service/blogArticle"
	modBlogArticle2 "github.com/hongmengzhu/xianfu-blog-go/app/domain/web/blog/model/modBlogArticle"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/web/utils/webPg"
	markdownPg2 "github.com/hongmengzhu/xianfu-blog-go/app/middleware/components/markdownPg"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/basic/modBasicTagsRelation"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/blog/modBlogArticleCategory"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBlog"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryBasic"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryBlog"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/cachePg/rdsPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constTags"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/blog/attachmentTypePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumApprovedPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/sdk/blog/key/blogKeyPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg/optionsPg"
	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/dbPg/pagePg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
	"golang.org/x/exp/slices"
)

func init() {
	gs.Provide(new(ArticleService))
}

type ArticleService struct {
	sv        *repositoryBlog.BlogArticleRepository           `autowire:"?"`
	statDb    *repositoryBlog.BlogArticleStatisticsRepository `autowire:"?"`
	catDb     *repositoryBlog.BlogArticleCategoryRepository   `autowire:"?"`
	sp        *blogArticle.Sp                                 `autowire:"?"`
	tagsRelat *repositoryBasic.BasicTagsRelationRepository    `autowire:"?"`
	rdu       *rdsPg.BatchString                              `autowire:"?"`
	catCore   *serviceCore.CoreArticleCategory                `autowire:"?"`
}

// Detail 详情
//
//	@Description:
//	@receiver c
//	@param id
func (c *ArticleService) Detail(ctx *gin.Context, id string) (rt rg.Rs[modBlogArticle2.DetailVo]) {
	if strPg.IsBlank(id) {
		return rt.ErrorMessage("id错误")
	}
	id = strings.TrimSpace(id)
	find, b := c.sv.FindByIdString(ctx, id)
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	no := webPg.GetTenantNo(ctx)
	if strPg.IsNotBlank(no) && find.TenantNo != no {
		return rt.ErrorMessage("数据不存在")
	}
	var stat modBlogArticle2.StatisticsVo
	//统计
	if strPg.IsNotBlank(find.No) {
		{
			tmp, result2 := c.statDb.FindByArticleNo(ctx, find.No)
			if result2 {
				copier.Copy(&stat, tmp)
			}
		}
	}
	log.Infof(ctx, log.TagAppDef, "find=%+v", find)
	var info modBlogArticle2.DetailVo
	copier.Copy(&info, find)
	//
	info.Statistics = stat
	//
	catNo := make([]string, 0)
	tags := make([]string, 0)
	info.CategoryObj = make([]*modBlogArticleCategory.Cache, 0)
	info.Tags = make([]string, 0)
	info.Where = make([]string, 0)
	info.AttachmentsMap = make(map[string]string)
	mapTagsOnly := make(map[string]bool)
	//附件图
	tagsData := make(map[string]modBasicTagsRelation.AllVo)
	if strPg.IsNotBlank(find.Attachments) {
		var imagesMap map[string]string
		err := json.Unmarshal([]byte(find.Attachments), &imagesMap)
		if err == nil {
			info.AttachmentsMap = imagesMap
		}
	}
	if strPg.IsNotBlank(find.CategoryNo) {
		catNo = append(catNo, find.CategoryNo)
	}
	if nil != find.Categorys.Data() {
		tmp := find.Categorys.Data()
		if len(tmp) > 0 {
			for _, tag := range tmp {
				if strPg.IsNotBlank(tag) && !slices.Contains(catNo, tag) {
					catNo = append(catNo, tag)
				}
			}
		}
	}
	if len(catNo) > 0 {
		{
			//cat, result := c.catDb.FindAllByNoIn(catNo)
			//if result {
			//	for _, item := range cat {
			//		info.CategoryObj = append(info.CategoryObj, item)
			//	}
			//}
			tmp, b := c.catCore.GetAllByKeysRetMap(ctx, catNo, find.TenantNo)
			if b {
				for _, item := range tmp {
					info.CategoryObj = append(info.CategoryObj, item)
				}
			}
		}
	}
	//标签
	if nil != find.Tags.Data() {
		tmp := find.Tags.Data()
		if len(tmp) > 0 {
			for _, tag := range tmp {
				//重复过滤掉
				if _, ok := mapTagsOnly[tag]; ok {
					continue
				}
				mapTagsOnly[tag] = false
				tags = append(tags, tag)
				//标签
				//info.Tags = append(info.Tags, tag)
			}
		}
		if len(tags) > 0 {
			{
				infos, result := c.tagsRelat.FindAllByTagNoInAndCategoryRoot(ctx, tags, constTags.ArticleInfo.String())
				if result {
					for _, item := range infos {
						var vo modBasicTagsRelation.AllVo
						copier.Copy(&vo, &item)
						vo.AttributeMap = make(map[string]interface{})
						//
						if strPg.IsNotBlank(item.Attribute) {
							err := json.Unmarshal([]byte(item.Attribute), &vo.AttributeMap)
							if err != nil {
								log.Errorf(ctx, log.TagAppDef, "json解析失败 %+v", err)
							}
							if obj, ok := vo.AttributeMap["color"]; ok {
								color := make(map[string]interface{})
								if strPg.IsNotBlank(obj.(string)) {
									err := json.Unmarshal([]byte(obj.(string)), &color)
									if err != nil {
										log.Errorf(ctx, log.TagAppDef, "json解析失败 %+v", err)
									}
								}
								vo.AttributeMap["color"] = color
							}
						} else {
							vo.AttributeMap["bordered"] = true
							vo.AttributeMap["type"] = "default"
							vo.AttributeMap["color"] = struct {
							}{}
							vo.AttributeMap["strong"] = false
							vo.AttributeMap["round"] = false
						}
						//
						tagsData[item.TagNo] = vo
						mapTagsOnly[item.TagNo] = true
					}
				}
			}
			{
				for tag, b2 := range mapTagsOnly {
					if !b2 {
						vo := modBasicTagsRelation.AllVo{
							Name:   tag,
							Code:   tag,
							NameFl: tag,
							TagNo:  tag,
						}
						vo.AttributeMap = make(map[string]interface{})
						vo.AttributeMap["bordered"] = true
						vo.AttributeMap["type"] = "default"
						vo.AttributeMap["color"] = struct {
						}{}
						vo.AttributeMap["strong"] = false
						vo.AttributeMap["round"] = false
						tagsData[tag] = vo
					}
				}
			}
			//
			for _, tag := range tags {
				obj := tagsData[tag]
				info.TagsStyle = append(info.TagsStyle, obj)
				info.Tags = append(info.Tags, obj.Name)
			}
		}
	}

	if strPg.IsNotBlank(find.Content) {
		raw := markdownPg2.Markdown([]byte(find.Content))
		info.ContentConv = raw.String()
	}
	//log.Infof(context.Background(), log.TagAppDef, "info:%+v", info)
	//log.Infof(context.Background(), log.TagAppDef, "info.create:%+v", datetimePg.Format(info.CreateAt, "2006"))
	return rt.OkData(info)
}

// Query 查询
//
//	@Description:
//	@receiver c
//	@param ct
func (c *ArticleService) Query(ctx *gin.Context, ct modBlogArticle2.QueryCt) (rt rg.Rs[pagePg.Paginator[modBlogArticle2.Vo]]) {
	var query entityBlog.BlogArticleEntity
	copier.Copy(&query, &ct)
	tenantNo := webPg.GetTenantNo(ctx)
	if strPg.IsNotBlank(tenantNo) {
		query.TenantNo = tenantNo
	}
	//启用
	query.State = enumStatePg.ENABLE.Index()
	//审批通过
	query.PlatformApproved = enumApprovedPg.ApprovedStateApproved.Index()
	slice := make([]modBlogArticle2.Vo, 0)
	rt.Data.Data = slice
	r := c.sv
	page, err := r.FindAllPage(ctx, query, optionsPg.WithOption(func(arg *optionsPg.OptionParams) {
		if ct.PageSize < 1 {
			ct.PageSize = 20
		}
		arg.Pageable = new(pagePg.PageablePageSize(0, ct.PageNum, ct.PageSize))
		//自定义查询
		arg.Db = arg.Db.Order("create_at desc")
		//自定义查询
		if "" != ct.Wd {
			arg.Db = arg.Db.Where("name like ?", "%"+ct.Wd+"%")
		}
		// 时间区间查询
		if nil != ct.CreateAtStart && nil != ct.CreateAtEnd {
			arg.Db = arg.Db.Where("create_at between ? and ?", ct.CreateAtStart, ct.CreateAtEnd)
		}
		//标签
		if nil != ct.TagsQuery && len(ct.TagsQuery) > 0 {
			for _, tag := range ct.TagsQuery {
				if strPg.IsNotBlank(tag) {
					arg.Db = arg.Db.Where("tags @> ?", "[\""+tag+"\"]")
				}
			}
		}
		//多分类
		if nil != ct.CategoryQuery && len(ct.CategoryQuery) > 0 {
			for _, tag := range ct.CategoryQuery {
				//获取缓存，得到 编号
				get, b := c.rdu.Get(ctx, blogKeyPg.ArticleCategoryTenantNoAndNoByCode(tenantNo, tag))
				if b {
					arg.Db = arg.Db.Where("categorys @> ?", "[\""+get+"\"]")
				} else {
					arg.Db = arg.Db.Where("categorys @> ?", "[\""+tag+"\"]")
				}
			}
		}
	}))
	if nil != err {
		return rt.Ok()
	}

	if page.Total > 0 && page.Data != nil && len(page.Data) > 0 {
		ImgDefault := "/assets/imgs/shop/product-1-1.jpg"
		pg := pagePg.NewPaginatorByPageable[modBlogArticle2.Vo](page.Pageable)
		mapCategory := make(map[string]*modBlogArticleCategory.Cache)
		mapStat := make(map[string]*entityBlog.BlogArticleStatisticsEntity)
		idsCategory := make([]string, 0)
		idsNo := make([]string, 0)
		tags := make([]string, 0)
		mapTagsOnly := make(map[string]bool)
		//标签
		tagsData := make(map[string]modBasicTagsRelation.AllVo)
		for _, item := range page.Data {
			if strPg.IsNotBlank(item.CategoryNo) {
				idsCategory = append(idsCategory, item.CategoryNo)
			}
			if item.Categorys.Data() != nil && len(item.Categorys.Data()) > 0 {
				for _, obj := range item.Categorys.Data() {
					idsCategory = append(idsCategory, obj)
				}
			}
			if strPg.IsNotBlank(item.No) {
				idsNo = append(idsNo, item.No)
			}
			//标签
			if nil != item.Tags.Data() {
				if len(item.Tags.Data()) > 0 {
					for _, tag := range item.Tags.Data() {
						//重复过滤掉
						if _, ok := mapTagsOnly[tag]; ok {
							continue
						}
						mapTagsOnly[tag] = false
						tags = append(tags, tag)
						//标签
					}
				}
			}
		}
		//分类
		{
			if len(idsCategory) > 0 {
				tmp, b := c.catCore.GetAllByKeysRetMap(ctx, idsCategory, tenantNo)
				if b {
					mapCategory = tmp
				}
				//tmp, result := c.catDb.FindAllByNoIn(ctx, idsCategory)
				//if result {
				//	for _, item := range tmp {
				//		mapCategory[item.No] = item
				//	}
				//}
			}
		}
		//统计
		{
			if len(idsNo) > 0 {
				tmp, result := c.statDb.FindAllByArticleNoIn(ctx, idsNo)
				if result {
					for _, item := range tmp {
						mapStat[item.ArticleNo] = item
					}
				}
			}
		}
		//标签
		{
			if len(tags) > 0 {
				infos, result := c.tagsRelat.FindAllByTagNoInAndCategoryRoot(ctx, tags, constTags.ArticleInfo.String())
				if result {
					for _, item := range infos {
						var vo modBasicTagsRelation.AllVo
						copier.Copy(&vo, &item)
						//
						vo.AttributeMap = make(map[string]interface{})
						//
						if strPg.IsNotBlank(item.Attribute) {
							err := json.Unmarshal([]byte(item.Attribute), &vo.AttributeMap)
							if err != nil {
								log.Errorf(ctx, log.TagAppDef, "json解析失败 %+v", err)
							}
							if obj, ok := vo.AttributeMap["color"]; ok {
								color := make(map[string]interface{})
								if strPg.IsNotBlank(obj.(string)) {
									err := json.Unmarshal([]byte(obj.(string)), &color)
									if err != nil {
										log.Errorf(ctx, log.TagAppDef, "json解析失败 %+v", err)
									}
								}
								vo.AttributeMap["color"] = color
							}
						} else {
							vo.AttributeMap["bordered"] = true
							vo.AttributeMap["type"] = "default"
							vo.AttributeMap["color"] = struct {
							}{}
							vo.AttributeMap["strong"] = false
							vo.AttributeMap["round"] = false
						}
						//
						tagsData[item.TagNo] = vo
						mapTagsOnly[item.TagNo] = true
					}
				}
				for tag, b := range mapTagsOnly {
					if !b {
						vo := modBasicTagsRelation.AllVo{
							Name:   tag,
							Code:   tag,
							NameFl: tag,
							TagNo:  tag,
						}
						vo.AttributeMap = make(map[string]interface{})
						vo.AttributeMap["bordered"] = true
						vo.AttributeMap["type"] = "default"
						vo.AttributeMap["color"] = struct {
						}{}
						vo.AttributeMap["strong"] = false
						vo.AttributeMap["round"] = false
						tagsData[tag] = vo
					}
				}
			}
		}
		//字段赋值
		for _, item := range page.Data {
			var vo modBlogArticle2.Vo
			copier.Copy(&vo, &item)
			vo.Tags = make([]string, 0)
			vo.TagsStyle = make([]modBasicTagsRelation.AllVo, 0)
			vo.Categorys = make([]string, 0)
			vo.CategoryObj = make([]*modBlogArticleCategory.Cache, 0)
			//
			if strPg.IsNotBlank(item.Attachments) {
				json.Unmarshal([]byte(item.Attachments), &vo.AttachmentsMap)
			}
			if nil == vo.AttachmentsMap {
				vo.AttachmentsMap = make(map[string]string, 0)
			}
			if _, ok := vo.AttachmentsMap[attachmentTypePg.AttachmentList.String()]; !ok {
				vo.AttachmentsMap[attachmentTypePg.AttachmentList.String()] = ImgDefault
			}
			if _, ok := vo.AttachmentsMap[attachmentTypePg.AttachmentMain.String()]; !ok {
				vo.AttachmentsMap[attachmentTypePg.AttachmentMain.String()] = ImgDefault
			}
			//统计
			if obj, ok := mapStat[item.No]; ok {
				vo.Statistics.Comment = obj.Comment
				vo.Statistics.Read = obj.Read
				vo.Statistics.SeoKeywords = obj.SeoKeywords
				vo.Statistics.SeoDescription = obj.SeoDescription
				vo.Statistics.PageTitle = obj.PageTitle
			}
			//分类
			if obj, ok := mapCategory[item.CategoryNo]; ok {
				vo.CategoryName = obj.Name
				vo.CategoryNo = obj.No
			}
			//标签
			if nil != item.Tags.Data() {
				if len(item.Tags.Data()) > 0 {
					for _, tag := range item.Tags.Data() {
						obj := tagsData[tag]
						vo.Tags = append(vo.Tags, obj.Name)
						vo.TagsStyle = append(vo.TagsStyle, obj)
					}
				}
			}
			if strPg.IsNotBlank(item.Content) {
				// 截取前1000个有效文字（图片不计，代码段完整）
				truncated := markdownPg2.TruncateMarkdown(item.Content, 1000)
				raw := markdownPg2.Markdown([]byte(truncated))
				vo.ContentConv = raw.String()
			}
			// 栏目
			if nil != item.Categorys.Data() {
				if len(item.Categorys.Data()) > 0 {
					for _, obj := range item.Categorys.Data() {
						vo.Categorys = append(vo.Categorys, obj)
						if obj2, ok := mapCategory[obj]; ok {
							vo.CategoryObj = append(vo.CategoryObj, obj2)
						}
					}
				}
			}
			if strPg.IsNotBlank(item.CategoryNo) && !slices.Contains(vo.Categorys, item.CategoryNo) {
				vo.Categorys = append(vo.Categorys, item.CategoryNo)
				if obj2, ok := mapCategory[item.CategoryNo]; ok {
					vo.CategoryObj = append(vo.CategoryObj, obj2)
				}
			}
			//
			slice = append(slice, vo)
		}
		pg.Data = slice
		pg.Pageable = page.Pageable
		rt.Data = pg
		return rt.Ok()
	}
	return rt.Ok()
}
