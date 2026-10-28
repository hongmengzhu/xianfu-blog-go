package service

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/blog/modBlogBookmark/modApiBlogBookmark"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/blog/modBlogBookmarkCategory/modApiBlogBookmarkCategory"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBlog"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryBlog"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/holderPg/holderApiPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/blog/bookmarkTypeOwnerPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/blog/typeReadingPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg/optionsPg"
	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/numberPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
	"gorm.io/gorm"
)

func init() {
	gs.Provide(new(BookmarkService))
}

type BookmarkService struct {
	sv     *repositoryBlog.BlogBookmarkRepository         `autowire:"?"`
	catDb  *repositoryBlog.BlogBookmarkCategoryRepository `autowire:"?"`
	catSer *BookmarkCategoryService                       `autowire:"?"`
}

// GetAll 获取所有
func (c *BookmarkService) GetAll(ctx *gin.Context) (rt rg.Rs[modApiBlogBookmark.VoAll]) {
	//
	holder := holderApiPg.GetContextAccount(ctx)
	//
	ret := modApiBlogBookmark.VoAll{}
	ret.My = make([]modApiBlogBookmark.Vo, 0)
	ret.Team = make([]modApiBlogBookmark.Vo, 0)
	ret.Org = make([]modApiBlogBookmark.Vo, 0)
	ret.MyCategory = make([]modApiBlogBookmarkCategory.Vo, 0)
	ret.TeamCategory = make([]modApiBlogBookmarkCategory.Vo, 0)
	ret.OrgCategory = make([]modApiBlogBookmarkCategory.Vo, 0)
	//
	all := c.catSer.GetAll(ctx)
	if all.SuccessIs() {
		ret.MyCategory = all.Data.My
		ret.TeamCategory = all.Data.Team
		ret.OrgCategory = all.Data.Org
	}
	//
	{
		q := entityBlog.BlogBookmarkEntity{}
		q.State = enumStatePg.ENABLE.Index()
		q.TenantNo = holder.GetTenantNo()
		q.Ano = holder.GetAno()
		q.TypeOwner = bookmarkTypeOwnerPg.MY.Code()
		infos := c.sv.FindAll(ctx, q)
		if nil != infos {
			for _, item := range infos {
				var vo modApiBlogBookmark.Vo
				copier.Copy(&vo, &item)
				//
				ret.My = append(ret.My, vo)
			}
		}
	}
	//
	{
		q := entityBlog.BlogBookmarkEntity{}
		q.State = enumStatePg.ENABLE.Index()
		q.TenantNo = holder.GetTenantNo()
		q.TypeOwner = bookmarkTypeOwnerPg.TEAM.Code()
		infos := c.sv.FindAll(ctx, q)
		if nil != infos {
			for _, item := range infos {
				var vo modApiBlogBookmark.Vo
				copier.Copy(&vo, &item)
				//
				ret.Team = append(ret.Team, vo)
			}
		}
	}
	//
	{
		q := entityBlog.BlogBookmarkEntity{}
		q.State = enumStatePg.ENABLE.Index()
		q.TenantNo = holder.GetTenantNo()
		q.TypeOwner = bookmarkTypeOwnerPg.ORG.Code()
		infos := c.sv.FindAll(ctx, q)
		if nil != infos {
			for _, item := range infos {
				var vo modApiBlogBookmark.Vo
				copier.Copy(&vo, &item)
				//
				ret.Org = append(ret.Org, vo)
			}
		}
	}
	//
	return rt.OkData(ret)
}

// GetMy 获取my
func (c *BookmarkService) GetMy(ctx *gin.Context, ct modApiBlogBookmark.QueryCt) (rt rg.Rs[[]modApiBlogBookmark.Vo]) {
	//
	log.Info(ctx, log.TagAppDef, "ct", log.Any("ct", ct))
	//
	holder := holderApiPg.GetContextAccount(ctx)
	//
	data := make([]modApiBlogBookmark.Vo, 0)
	//
	{
		q := entityBlog.BlogBookmarkEntity{}
		q.State = enumStatePg.ENABLE.Index()
		q.TenantNo = holder.GetTenantNo()
		q.Ano = holder.GetAno()
		q.TypeOwner = bookmarkTypeOwnerPg.MY.Code()
		infos := c.sv.FindAll(ctx, q, optionsPg.WithCondition(func(db *gorm.DB) *gorm.DB {
			// 未读完
			if ("un" + typeReadingPg.COMPLETED.Index()) == ct.TypeReading {
				db = db.Where("type_reading in ?", []string{typeReadingPg.UNREAD.Index(), typeReadingPg.READING.Index()})
				//排序
				db = db.Order("CASE typeReading WHEN 'unread' THEN 1 WHEN 'reading' THEN 2 WHEN 'completed' THEN 3 ELSE 4 END")
			}
			return db
		}))
		if nil != infos {
			for _, item := range infos {
				var vo modApiBlogBookmark.Vo
				copier.Copy(&vo, &item)
				//
				data = append(data, vo)
			}
		}
	}
	//
	return rt.OkData(data)
}

// GetOrg 获取
func (c *BookmarkService) GetOrg(ctx *gin.Context) (rt rg.Rs[[]modApiBlogBookmark.Vo]) {
	//
	holder := holderApiPg.GetContextAccount(ctx)
	//
	data := make([]modApiBlogBookmark.Vo, 0)
	//
	{
		q := entityBlog.BlogBookmarkEntity{}
		q.State = enumStatePg.ENABLE.Index()
		q.TenantNo = holder.GetTenantNo()
		q.Ano = holder.GetAno()
		q.TypeOwner = bookmarkTypeOwnerPg.ORG.Code()
		infos := c.sv.FindAll(ctx, q)
		if nil != infos {
			for _, item := range infos {
				var vo modApiBlogBookmark.Vo
				copier.Copy(&vo, &item)
				//
				data = append(data, vo)
			}
		}
	}
	//
	return rt.OkData(data)
}

// Create 创建
func (c *BookmarkService) Create(ctx *gin.Context, ct modApiBlogBookmark.CreateUpdateCt) (rt rg.Rs[string]) {
	log.Info(ctx, log.TagAppDef, "ct", log.Any("ct", ct))
	var info entityBlog.BlogBookmarkEntity
	err := copier.Copy(&info, &ct)
	if err != nil {
		log.Info(ctx, log.TagAppDef, "copier.Copy error:", log.Any("err", err.Error()))
	}
	if ct.Name == "" {
		return rt.ErrorMessage("名称不能为空")
	}
	if ct.UrlSource == "" {
		return rt.ErrorMessage("来源地址不能为空")
	}
	now := time.Now()
	if nil == ct.PublishTime {
		info.PublishTime = &now
	}
	err, _ = c.sv.Create(ctx, &info)
	if err != nil {
		return rt.ErrorMessage("保存失败 " + err.Error())
	}
	return rg.OkData(numberPg.Int64ToString(info.ID))
}

// Update 创建
func (c *BookmarkService) Update(ctx *gin.Context, ct modApiBlogBookmark.CreateUpdateCt) (rt rg.Rs[string]) {
	log.Info(ctx, log.TagAppDef, "ct", log.Any("ct", ct))
	if ct.ID.ToInt64() <= 0 {
		return rt.ErrorMessage("ID不能为空")
	}
	find, b := c.sv.FindById(ctx, ct.ID.ToInt64())
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	var info entityBlog.BlogBookmarkEntity
	err := copier.Copy(&info, &ct)
	if err != nil {
		log.Info(ctx, log.TagAppDef, "copier.Copy error:", log.Any("err", err.Error()))
	}
	if ct.Name == "" {
		return rt.ErrorMessage("名称不能为空")
	}
	if ct.UrlSource == "" {
		return rt.ErrorMessage("来源地址不能为空")
	}
	now := time.Now()
	if nil == ct.PublishTime {
		info.PublishTime = &now
	}
	err = c.sv.Update(ctx, info, find.ID)
	if err != nil {
		return rt.ErrorMessage("保存失败 " + err.Error())
	}
	return rg.OkData(numberPg.Int64ToString(info.ID))
}

// LogicalDeletion 逻辑删除
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BookmarkService) LogicalDeletion(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
	if len(ids) < 1 {
		return rt.ErrorMessage("id错误")
	}
	repository := c.sv
	finds, b := repository.FindAllByIdStringIn(ctx, ids, optionsPg.WithCtx(ctx))
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	if c.sv.Config().Data.Delete {
		for _, info := range finds {
			log.Info(ctx, log.TagAppDef, "LogicalDeletion", log.Any("id", info.ID), log.Any("TenantNo", info.TenantNo))
		}
		repository.DeleteByIdsString(ctx, ids)
	} else {
		for _, info := range finds {
			enum := enumStatePg.State(info.State)
			if ok, reverse := enum.ReverseEnableDisable(); ok {
				repository.Update(ctx, entityBlog.BlogBookmarkEntity{State: reverse.IndexInt8()}, info.ID)
			}
		}
	}
	return rt.Ok()
}
