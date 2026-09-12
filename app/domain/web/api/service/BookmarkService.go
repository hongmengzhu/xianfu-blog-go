package service

import (
	"github.com/gin-gonic/gin"
	modelBlogBookmark2 "github.com/hongmengzhu/xianfu-blog-go/app/domain/web/api/model/modelBlogBookmark"
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/web/api/model/modelBlogBookmarkCategory"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBlog"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryBlog"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/holderPg/holderApiPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/blog/bookmarkTypeOwnerPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
)

type BookmarkService struct {
	sv     *repositoryBlog.BlogBookmarkRepository         `autowire:"?"`
	catDb  *repositoryBlog.BlogBookmarkCategoryRepository `autowire:"?"`
	catSer *BookmarkCategoryService                       `autowire:"?"`
}

// GetAll 获取所有
func (c *BookmarkService) GetAll(ctx *gin.Context) (rt rg.Rs[modelBlogBookmark2.VoAll]) {
	//
	holder := holderApiPg.GetContextAccount(ctx)
	//
	ret := modelBlogBookmark2.VoAll{}
	ret.My = make([]modelBlogBookmark2.Vo, 0)
	ret.Team = make([]modelBlogBookmark2.Vo, 0)
	ret.MyCategory = make([]modelBlogBookmarkCategory.Vo, 0)
	ret.TeamCategory = make([]modelBlogBookmarkCategory.Vo, 0)
	//
	all := c.catSer.GetAll(ctx)
	if all.SuccessIs() {
		ret.MyCategory = all.Data.My
		ret.TeamCategory = all.Data.Team
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
				var vo modelBlogBookmark2.Vo
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
				var vo modelBlogBookmark2.Vo
				copier.Copy(&vo, &item)
				//
				ret.Team = append(ret.Team, vo)
			}
		}
	}
	//
	return rt.OkData(ret)
}

// GetMy 获取所有
func (c *BookmarkService) GetMy(ctx *gin.Context) (rt rg.Rs[[]modelBlogBookmark2.Vo]) {
	//
	holder := holderApiPg.GetContextAccount(ctx)
	//
	data := make([]modelBlogBookmark2.Vo, 0)
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
				var vo modelBlogBookmark2.Vo
				copier.Copy(&vo, &item)
				//
				data = append(data, vo)
			}
		}
	}
	//
	return rt.OkData(data)
}
