package service

import (
	"strings"
	"time"

	"github.com/duke-git/lancet/v2/strutil"
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/manage/domainBlog/service/blogCollect"
	modBlogCollect2 "github.com/hongmengzhu/xianfu-blog-go/app/domain/web/api/model/modBlogCollect"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBlog"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryBlog"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/holderPg/holderApiPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/blog/typeContentPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/blog/typeDataSourcePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/blog/typeReadingPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/blog/typeSourcePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model/modelBasePg"
	"github.com/pangu-2/go-tools/tools/cryptPg"
	"github.com/pangu-2/go-tools/tools/noPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
	"gorm.io/datatypes"
)

func init() {
	gs.Provide(new(CollectService))
}

type CollectService struct {
	sv    *repositoryBlog.BlogCollectRepository         `autowire:"?"`
	catDb *repositoryBlog.BlogCollectCategoryRepository `autowire:"?"`
	sp    *blogCollect.Sp                               `autowire:"?"`
}

// Push
//
//	@Description: 推送文章连接
//	@receiver c
func (c *CollectService) Push(ctx *gin.Context, ct modBlogCollect2.PushCt) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)

	if strPg.IsBlank(ct.Title) {
		return rt.ErrorMessage("标题不能为空")
	}
	if strPg.IsBlank(ct.Url) {
		return rt.ErrorMessage("url地址不能为空")
	}
	urlSource := strings.TrimSpace(ct.Url)
	urlSourceMd5 := cryptPg.Md5(urlSource)
	//规则
	//重复检查
	duplicate := true
	if ct.Rule != nil && len(ct.Rule) > 0 {
		for _, v := range ct.Rule {
			if strPg.IsBlank(v) {
				continue
			}
			//重复:允许重复数据
			if "duplicatePass" == strings.TrimSpace(v) {
				//不检查重复数据
				duplicate = false
			}
		}
	}
	//重复
	if duplicate {
		info, result := c.sv.FindAllByUrlSourceMd5(ctx, urlSourceMd5)
		if result && nil != info && len(info) > 0 {
			return rt.ErrorMessage("该链接已存在")
		}
	}

	holder := holderApiPg.GetContextAccount(ctx)
	save := entityBlog.BlogCollectEntity{
		Name:        strings.TrimSpace(ct.Title),
		UrlSource:   urlSource,
		Description: strings.TrimSpace(ct.Description),
	}

	if strPg.IsBlank(ct.CategoryNo) {
		save.CategoryNo = "default"
	} else {
		info, result := c.catDb.FindByNo(ctx, ct.CategoryNo)
		if !result {
			return rt.ErrorMessage("分类不存在")
		}
		save.CategoryNo = info.No
	}
	save.Content = strings.TrimSpace(ct.Content)
	save.Editor = strings.TrimSpace(ct.Editor)
	save.Author = strings.TrimSpace(ct.Author)
	save.Source = strings.TrimSpace(ct.Source)
	tags := make([]string, 0)
	if nil != ct.Tags && len(ct.Tags) > 0 {
		for _, v := range ct.Tags {
			if strPg.IsNotBlank(v) {
				tags = append(tags, strings.TrimSpace(v))
			}
		}
	}
	save.Tags = datatypes.NewJSONType(tags)

	save.No = noPg.No()
	save.Code = save.No
	save.TenantNo = holder.GetTenantNo()
	save.UrlSourceMd5 = urlSourceMd5
	save.TypeReading = typeReadingPg.UNREAD.Index()
	save.TypeDataSource = typeDataSourcePg.EXTERNAL.Index()
	save.TypeSource = typeSourcePg.COLLECTION.Index()
	save.TypeContent = typeContentPg.REPOST.Index()
	now := time.Now()
	save.OperationTime = &now
	if nil == ct.PublishTime {
		save.PublishTime = &now
	} else {
		toTime := ct.PublishTime.ToTime()
		save.PublishTime = &toTime
	}
	//
	//
	err, _ := c.sv.Create(ctx, &save)
	if err != nil {
		log.Debugf(ctx, log.TagAppDef, "save err=%+v", err)
		return rt.ErrorMessage("保存失败：" + err.Error())
	}
	return rt.Ok()
}

// PushAll
//
//	@Description: 推送文章连接
//	@receiver c
func (c *CollectService) PushAll(ctx *gin.Context, ct modBlogCollect2.PushAll) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", ct)
	if nil == ct.Data || len(ct.Data) <= 0 {
		return rt.ErrorMessage("数据不能为空")
	}
	errResult := make([]modelBasePg.ItemResult, 0)
	//规则
	//重复检查
	duplicate := true
	duplicateCount := 0
	if ct.Rule != nil && len(ct.Rule) > 0 {
		for _, v := range ct.Rule {
			if strPg.IsBlank(v) {
				continue
			}
			//重复:允许重复数据
			if "duplicatePass" == strings.TrimSpace(v) {
				//不检查重复数据
				duplicate = false
			}
		}
	}
	holder := holderApiPg.GetContextAccount(ctx)
	tmpIndex := make(map[string]int64)
	mapUrl := make(map[string]*entityBlog.BlogCollectEntity)
	mapCat := make(map[string]*entityBlog.BlogCollectCategoryEntity)
	saveTmp := make([]*entityBlog.BlogCollectEntity, 0)
	save := make([]*entityBlog.BlogCollectEntity, 0)
	catNo := make([]string, 0)
	urlMd5 := make([]string, 0)
	now := time.Now()
	for i, item := range ct.Data {
		if strPg.IsBlank(item.Title) {
			errResult = append(errResult, modelBasePg.ItemResult{
				Row: int64(i),
				Msg: "标题不能为空",
			})
			continue
		}
		if strPg.IsBlank(item.Url) {
			errResult = append(errResult, modelBasePg.ItemResult{
				Row: int64(i),
				Msg: "链接不能为空",
			})
			continue
		}
		if strPg.IsNotBlank(item.CategoryNo) {
			catNo = append(catNo, strings.TrimSpace(item.CategoryNo))
		}
		obj := entityBlog.BlogCollectEntity{
			UrlSource: strings.TrimSpace(item.Url),
		}
		obj.UrlSourceMd5 = cryptPg.Md5(obj.UrlSource)
		//
		urlMd5 = append(urlMd5, obj.UrlSourceMd5)
		//
		//
		obj.Description = strings.TrimSpace(item.Description)
		obj.Name = strings.TrimSpace(item.Title)
		if len(obj.Name) > 300 {
			obj.Name = strutil.Substring(item.Title, 0, 300)
			obj.Description = strings.TrimSpace(item.Title)
		}

		tmpIndex[obj.UrlSourceMd5] = int64(i)
		//
		obj.Content = strings.TrimSpace(item.Content)
		obj.Editor = strings.TrimSpace(item.Editor)
		obj.Author = strings.TrimSpace(item.Author)
		obj.Source = strings.TrimSpace(item.Source)
		tags := make([]string, 0)
		if nil != item.Tags && len(item.Tags) > 0 {
			for _, v := range item.Tags {
				if strPg.IsNotBlank(v) {
					tags = append(tags, strings.TrimSpace(v))
				}
			}
		}
		obj.Tags = datatypes.NewJSONType(tags)
		//
		obj.No = noPg.No()
		obj.Code = obj.No
		obj.TenantNo = holder.GetTenantNo()
		obj.TypeReading = typeReadingPg.UNREAD.Index()
		obj.TypeDataSource = typeDataSourcePg.EXTERNAL.Index()
		obj.TypeSource = typeSourcePg.COLLECTION.Index()
		obj.TypeContent = typeContentPg.REPOST.Index()
		obj.OperationTime = &now
		if nil == item.PublishTime {
			obj.PublishTime = &now
		} else {
			obj.PublishTime = new(item.PublishTime.ToTime())
		}
		//
		saveTmp = append(saveTmp, &obj)
	}
	// 连接
	//重复
	if duplicate && len(urlMd5) > 0 {
		info, result := c.sv.FindAllByUrlSourceMd5In(ctx, urlMd5)
		if result && nil != info && len(info) > 0 {
			for _, entity := range info {
				duplicateCount++
				mapUrl[entity.UrlSourceMd5] = entity
			}
		}
	}
	if len(catNo) > 0 {
		{
			info, result := c.catDb.FindAllByNoIn(ctx, catNo)
			if result {
				for _, item := range info {
					mapCat[item.No] = item
				}
			}
		}
	}
	if len(saveTmp) > 0 {
		for _, item := range saveTmp {
			if _, ok := mapUrl[item.UrlSourceMd5]; ok {
				errResult = append(errResult, modelBasePg.ItemResult{
					Row: tmpIndex[item.UrlSourceMd5],
					Msg: "该链接已存在",
				})
				continue
			}
			//
			if strPg.IsBlank(item.CategoryNo) {
				item.CategoryNo = "default"
			} else {
				if obj, ok := mapCat[item.CategoryNo]; ok {
					item.CategoryNo = obj.No
				} else {
					item.CategoryNo = ""
				}
			}
			log.Infof(ctx, log.TagAppDef, "item=%+v", item)
			//
			save = append(save, item)
		}
		//
		//
		{
			tx := c.sv.DbModel().CreateInBatches(save, 1000000)
			if tx.Error != nil {
				log.Errorf(ctx, log.TagAppDef, "save err=%+v", tx.Error)
				return rt.ErrorMessage("保存失败：")
			}
			if 0 == tx.RowsAffected {
				return rt.ErrorMessage("保存失败，没有更新任何数据")
			}
		}
		clear(saveTmp)
		clear(save)
		clear(mapUrl)
		clear(mapCat)
		rt.Extend = make(map[string]interface{})
		rt.Extend["duplicateCount"] = duplicateCount
		rt.Extend["errors"] = errResult
		return rt.Ok()
	}
	clear(saveTmp)
	clear(save)
	clear(mapUrl)
	clear(mapCat)
	rt.Extend = make(map[string]interface{})
	rt.Extend["duplicateCount"] = duplicateCount
	rt.Extend["errors"] = errResult
	return rt.ErrorMessage("没有数据需要保存")
}
