package service

import (
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/web/api/model/modelBlogBookmark"
	"github.com/hongmengzhu/xianfu-blog-go/app/web/api/model/modelBlogBookmarkCategory"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBlog"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryBlog"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/holderPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/holderPg/holderApiPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/automatedPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constNodePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/blog/bookmarkTypeOwnerPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg/optionsPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/noPg"
	"github.com/pangu-2/go-tools/tools/numberPg"
	"github.com/pangu-2/go-tools/tools/slicePg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(BookmarkCategoryService))
}

type BookmarkCategoryService struct {
	sv *repositoryBlog.BlogBookmarkCategoryRepository `autowire:"?"`
}

// GetAll 获取所有
func (c *BookmarkCategoryService) GetAll(ctx *gin.Context) (rt rg.Rs[modelBlogBookmarkCategory.VoAll]) {
	//
	holder := holderApiPg.GetContextAccount(ctx)
	//
	vo := modelBlogBookmarkCategory.VoAll{}
	vo.My = make([]modelBlogBookmarkCategory.Vo, 0)
	vo.Team = make([]modelBlogBookmarkCategory.Vo, 0)
	//
	{
		q := entityBlog.BlogBookmarkCategoryEntity{}
		q.State = enumStatePg.ENABLE.Index()
		q.TenantNo = holder.GetTenantNo()
		q.Ano = holder.GetAno()
		q.TypeOwner = bookmarkTypeOwnerPg.MY.Code()
		infos := c.sv.FindAll(ctx, q)
		if nil != infos {
			for _, info := range infos {
				vo.My = append(vo.My, modelBlogBookmarkCategory.Vo{
					ID:       typePg.Int64String(info.ID),
					Name:     info.Name,
					NameFl:   info.NameFl,
					NameFull: info.NameFull,
					No:       info.No,
					Code:     info.Code,
					ParentNo: info.ParentNo,
				})
			}
		}
	}

	{
		q := entityBlog.BlogBookmarkCategoryEntity{}
		q.State = enumStatePg.ENABLE.Index()
		q.TenantNo = holder.GetTenantNo()
		q.TypeOwner = bookmarkTypeOwnerPg.TEAM.Code()
		infos := c.sv.FindAll(ctx, q)
		if nil != infos {
			for _, info := range infos {
				vo.Team = append(vo.Team, modelBlogBookmarkCategory.Vo{
					ID:       typePg.Int64String(info.ID),
					Name:     info.Name,
					NameFl:   info.NameFl,
					NameFull: info.NameFull,
					No:       info.No,
					Code:     info.Code,
					ParentNo: info.ParentNo,
				})
			}
		}
	}
	//
	return rt.OkData(vo)
}

// GetMy 获取所有
func (c *BookmarkCategoryService) GetMy(ctx *gin.Context) (rt rg.Rs[[]modelBlogBookmarkCategory.Vo]) {
	//
	holder := holderApiPg.GetContextAccount(ctx)
	//
	data := make([]modelBlogBookmarkCategory.Vo, 0)
	//
	{
		q := entityBlog.BlogBookmarkCategoryEntity{}
		q.State = enumStatePg.ENABLE.Index()
		q.TenantNo = holder.GetTenantNo()
		q.Ano = holder.GetAno()
		q.TypeOwner = bookmarkTypeOwnerPg.MY.Code()
		infos := c.sv.FindAll(ctx, q)
		if nil != infos {
			for _, info := range infos {
				data = append(data, modelBlogBookmarkCategory.Vo{
					ID:       typePg.Int64String(info.ID),
					Name:     info.Name,
					NameFl:   info.NameFl,
					NameFull: info.NameFull,
					No:       info.No,
					Code:     info.Code,
					ParentNo: info.ParentNo,
				})
			}
		}
	}
	//
	return rt.OkData(data)
}

// Save 创建或更新
//
//	@Description: ct.ID > 0 表示更新，否则表示创建
//	@receiver c
//	@param ct
//	@return rt
func (c *BookmarkCategoryService) Save(ctx *gin.Context, ct modelBlogBookmark.CreateUpdate) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%#v", ct)
	if "" == ct.Name {
		return rt.ErrorMessage("名称不能为空")
	}
	isUpdate := ct.ID > 0
	var info entityBlog.BlogBookmarkCategoryEntity
	err := copier.Copy(&info, &ct)
	if err != nil {
		log.Infof(ctx, log.TagAppDef, "copier.Copy error: %+v", err)
	}
	r := c.sv
	parent := &entityBlog.BlogBookmarkCategoryEntity{}
	if isUpdate {
		//更新
		if strPg.IsBlank(info.Code) {
			info.Code = ""
		} else {
			_, result := r.FindByCodeAndIdNot(ctx, info.Code, ct.ID.ToString(), optionsPg.WithCtx(ctx))
			if result {
				return rt.ErrorMessage("标志已存在")
			}
		}
		find, b := r.FindById(ctx, ct.ID.ToInt64())
		if !b {
			return rt.ErrorMessage("数据不存在")
		}
		//上级
		var childData []*entityBlog.BlogBookmarkCategoryEntity
		if strPg.IsNotBlank(ct.ParentNo) {
			result := false
			parent, result = r.FindByNo(ctx, ct.ParentNo, optionsPg.WithCtx(ctx))
			if !result {
				return rt.ErrorMessage("上级不存在")
			}
			if parent.ID == ct.ID.ToInt64() {
				return rt.ErrorMessage("上级不能等于自己")
			}
			//新的上级 不等于 旧的上级时,检测是否已经在新的子集已存在
			if parent.No != find.ParentNo {
				result2 := false
				childData, result2 = r.FindAllByNoLink(ctx, find.No)
				if result2 {
					for _, item := range childData {
						if item.No == parent.No {
							return rt.ErrorMessage("无法保存，不能设置为自己的子集")
						}
					}
				}
			}
		}
		//设置上级 link
		if strPg.IsNotBlank(ct.ParentNo) {
			info.IdLink = constNodePg.NoLinkAssemble(parent.IdLink, numberPg.Int64ToString(find.ID))
			info.NoLink = constNodePg.NoLinkAssemble(parent.NoLink, find.No)
			info.ParentNo = parent.No
			info.ParentId = numberPg.Int64ToString(parent.ID)
		} else {
			info.IdLink = constNodePg.NoLinkDefault(numberPg.Int64ToString(find.ID))
			info.NoLink = constNodePg.NoLinkDefault(find.No)
			info.ParentNo = ""
			info.ParentId = ""
		}
		info.No = ""
		log.Infof(ctx, log.TagAppDef, "info.IdLink=%+v", info.IdLink)
		err = r.Update(ctx, info, info.ID)
		if err != nil {
			log.Infof(ctx, log.TagAppDef, "update error=%+v", err)
			return rt.ErrorMessage(err.Error())
		}
		log.Infof(ctx, log.TagAppDef, "save.info=%+v", info)
		//更改上级后，相关子集修改
		if strPg.IsNotBlank(ct.ParentNo) && nil != childData {
			maps := slicePg.ToMapArray(childData, func(t *entityBlog.BlogBookmarkCategoryEntity) (string, *entityBlog.BlogBookmarkCategoryEntity) {
				if strPg.IsBlank(t.ParentNo) {
					return constNodePg.ROOT, t
				}
				return t.ParentNo, t
			})
			if strPg.IsBlank(info.ParentNo) {
				info.ParentNo = constNodePg.ROOT
			}
			for _, item := range maps[info.ParentNo] {
				item.IdLink = constNodePg.NoLinkAssemble(info.IdLink, numberPg.Int64ToString(find.ID))
				item.NoLink = constNodePg.NoLinkAssemble(info.NoLink, item.No)
				c.childParentIdLink(maps, item)
			}
			log.Infof(ctx, log.TagAppDef, "maps=%+v", maps)
			for _, val := range maps {
				for _, item := range val {
					if item.ID == find.ID {
						continue
					}
					r.Update(ctx, entityBlog.BlogBookmarkCategoryEntity{IdLink: item.IdLink,
						NoLink: item.NoLink},
						item.ID)
				}
			}
		}
		return rt.Ok()
	}
	//创建
	if strPg.IsBlank(info.Code) {
		info.Code = automatedPg.CREATE_CODE
	}
	//判断是否是自动,不是自动
	if !automatedPg.IsCreateCode(info.Code) {
		//判断格式是否满足要求
		if !automatedPg.FormatVerify(info.Code) {
			return rt.ErrorMessage("标志格式不能为空")
		}
		//不是自动
		_, result := r.FindByCode(ctx, info.Code, optionsPg.WithCtx(ctx))
		if result {
			return rt.ErrorMessage("标志已存在")
		}
	}
	if strPg.IsNotBlank(ct.ParentNo) {
		result := false
		parent, result = r.FindByNo(ctx, ct.ParentNo, optionsPg.WithCtx(ctx))
		if !result {
			return rt.ErrorMessage("上级不存在")
		}
	}
	info.No = noPg.No()
	//自动设置编号
	if automatedPg.IsCreateCode(info.Code) {
		info.Code = strPg.GenerateNumberId22()
	}
	holder := holderPg.GetContextAccount(ctx)
	info.TenantNo = holder.GetTenantNo()
	log.Infof(ctx, log.TagAppDef, "info=%+v", info)
	err, _ = r.Create(ctx, &info)
	if err != nil {
		return rt.ErrorMessage("保存失败 " + err.Error())
	}
	//设置上级 link
	if strPg.IsNotBlank(ct.ParentNo) {
		info.IdLink = constNodePg.NoLinkAssemble(parent.IdLink, numberPg.Int64ToString(info.ID))
		info.NoLink = constNodePg.NoLinkAssemble(parent.NoLink, info.No)
		info.ParentNo = parent.No
		info.ParentId = numberPg.Int64ToString(parent.ID)
	} else {
		info.IdLink = constNodePg.NoLinkDefault(numberPg.Int64ToString(info.ID))
		info.NoLink = constNodePg.NoLinkDefault(info.No)
		info.ParentId = ""
		info.ParentNo = ""
	}
	err = r.Update(ctx, info, info.ID)
	if err != nil {
		return rt.ErrorMessage(err.Error())
	}
	return rg.OkData(numberPg.Int64ToString(info.ID))
}

// ChildParentIdLink 子集 上级 link更新
//
//	@Description:
//	@receiver c
//	@param id
func (c *BookmarkCategoryService) childParentIdLink(maps map[string][]*entityBlog.BlogBookmarkCategoryEntity, parent *entityBlog.BlogBookmarkCategoryEntity) {
	key := parent.ParentNo
	if strPg.IsBlank(parent.ParentNo) {
		key = constNodePg.ROOT
	}
	entities := maps[key]
	for _, item := range entities {
		item.IdLink = constNodePg.NoLinkAssemble(parent.IdLink, numberPg.Int64ToString(item.ID))
		item.NoLink = constNodePg.NoLinkAssemble(parent.NoLink, item.No)
	}
}
