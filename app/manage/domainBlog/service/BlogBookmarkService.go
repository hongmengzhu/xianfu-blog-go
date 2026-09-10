package service

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/blog/modBlogBookmark"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBlog"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryBlog"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/holderPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg/optionsPg"
	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/dbPg/pagePg"
	"github.com/pangu-2/go-tools/tools/noPg"
	"github.com/pangu-2/go-tools/tools/numberPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(BlogBookmarkService))
}

// BlogBookmarkService 书签
// @Description:
type BlogBookmarkService struct {
	sv    *repositoryBlog.BlogBookmarkRepository         `autowire:"?"`
	catDb *repositoryBlog.BlogBookmarkCategoryRepository `autowire:"?"`
}

// Create 新增
//
//	@Description:
//	@receiver c
//	@param ct
//	@return rt
func (c *BlogBookmarkService) Create(ctx *gin.Context, ct modBlogBookmark.CreateUpdateCt) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%#v", ct)
	var info entityBlog.BlogBookmarkEntity
	err := copier.Copy(&info, &ct)
	if err != nil {
		log.Infof(ctx, log.TagAppDef, "copier.Copy error: %+v", err)
	}
	if "" == ct.Name {
		return rt.ErrorMessage("名称不能为空")
	}
	holder := holderPg.GetContextAccount(ctx)
	info.No = noPg.No()
	info.TenantNo = holder.GetTenantNo()
	//附件
	if nil != ct.Attachment && len(ct.Attachment) > 0 {
		bytes, _ := json.Marshal(ct.Attachment)
		info.Attachments = string(bytes)
	}
	err, _ = c.sv.Create(ctx, &info)
	if err != nil {
		return rt.ErrorMessage("保存失败 " + err.Error())
	}
	return rg.OkData(numberPg.Int64ToString(info.ID))
}

// Update 更新
//
//	@Description:
//	@receiver c
//	@param ct
//	@return rt
func (c *BlogBookmarkService) Update(ctx *gin.Context, ct modBlogBookmark.CreateUpdateCt) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%#v", ct)
	var info entityBlog.BlogBookmarkEntity
	err := copier.Copy(&info, &ct)
	if err != nil {
		log.Infof(ctx, log.TagAppDef, "copier.Copy error: %+v", err)
	}
	if ct.ID < 1 {
		return rt.ErrorMessage("id错误")
	}
	find, b := c.sv.FindById(ctx, ct.ID.ToInt64())
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	//附件
	if nil != ct.Attachment && len(ct.Attachment) > 0 {
		bytes, _ := json.Marshal(ct.Attachment)
		info.Attachments = string(bytes)
	}
	//清除不可更新字段
	info.ID = 0
	info.No = ""
	err = c.sv.Update(ctx, info, find.ID)
	if err != nil {
		return rt.ErrorMessage(err.Error())
	}
	return rt.Ok()
}

// Detail 详情
//
//	@Description:
//	@receiver c
//	@param id
func (c *BlogBookmarkService) Detail(ctx *gin.Context, id int64) (rt rg.Rs[modBlogBookmark.DetailVo]) {
	if id < 1 {
		return rt.ErrorMessage("id错误")
	}
	find, b := c.sv.FindById(ctx, id)
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	var info modBlogBookmark.DetailVo
	copier.Copy(&info, &find)
	info.Tags = make([]string, 0)
	info.Where = make([]string, 0)
	info.Attachment = make(map[string]string)
	//附件
	if strPg.IsNotBlank(find.Attachments) {
		var imagesMap map[string]string
		err := json.Unmarshal([]byte(find.Attachments), &imagesMap)
		if err == nil {
			info.Attachment = imagesMap
		}
	}
	//标签
	if nil != find.Tags.Data() {
		info.Tags = find.Tags.Data()
	}
	//可用范围
	if nil != find.Where.Data() {
		info.Where = find.Where.Data()
	}
	//分类名称
	if strPg.IsNotBlank(find.CategoryNo) {
		catFind, catOk := c.catDb.FindByNo(ctx, find.CategoryNo, optionsPg.WithCtx(ctx))
		if catOk {
			info.CategoryName = catFind.Name
		}
	}
	return rt.OkData(info)
}

// Enable 启用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BlogBookmarkService) Enable(ctx *gin.Context, ct model.BaseIdsCt[string]) (rt rg.Rs[string]) {
	return c.State(ctx, ct.Ids, enumStatePg.ENABLE)
}

// Disable 禁用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BlogBookmarkService) Disable(ctx *gin.Context, ct model.BaseIdsCt[string]) (rt rg.Rs[string]) {
	return c.State(ctx, ct.Ids, enumStatePg.GetType(enumStatePg.DISABLE))
}

// State 状态 启用/禁用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BlogBookmarkService) State(ctx *gin.Context, ids []string, state enumStatePg.State) (rt rg.Rs[string]) {
	if len(ids) < 1 {
		return rt.ErrorMessage("id错误")
	}
	r := c.sv
	finds, b := r.FindAllByIdStringIn(ctx, ids, optionsPg.WithCtx(ctx))
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	for _, info := range finds {
		if info.State != state.IndexInt8() {
			r.Update(ctx, entityBlog.BlogBookmarkEntity{State: state.IndexInt8()}, info.ID)
		}
	}
	return rt.Ok()
}

// StateEnableDisable 状态 设置 有效 停用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BlogBookmarkService) StateEnableDisable(ctx *gin.Context, ids []string, state enumStatePg.State) (rt rg.Rs[string]) {
	if !state.IsEnableDisable() {
		return rt.ErrorMessage("状态错误")
	}
	return c.State(ctx, ids, state)
}

// LogicalDeletion 逻辑删除
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BlogBookmarkService) LogicalDeletion(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
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
			log.Infof(ctx, log.TagAppDef, "id=%v,TenantNo=%v", info.ID, info.TenantNo)
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

// LogicalRecovery 逻辑删除恢复
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BlogBookmarkService) LogicalRecovery(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
	if len(ids) < 1 {
		return rt.ErrorMessage("id错误")
	}
	repository := c.sv
	finds, b := repository.FindAllByIdStringIn(ctx, ids, optionsPg.WithCtx(ctx))
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	for _, info := range finds {
		enum := enumStatePg.State(info.State)
		if ok, reverse := enum.ReverseCancelLayAside(); ok {
			repository.Update(ctx, entityBlog.BlogBookmarkEntity{State: reverse.IndexInt8()}, info.ID)
		}
	}
	return rt.Ok()
}

// PhysicalDeletion 物理删除
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BlogBookmarkService) PhysicalDeletion(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
	if len(ids) < 1 {
		return rt.ErrorMessage("id错误")
	}
	cn := c.sv
	finds, b := cn.FindAllByIdStringIn(ctx, ids, optionsPg.WithCtx(ctx))
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	idsNew := make([]int64, 0)
	for _, info := range finds {
		log.Infof(ctx, log.TagAppDef, "id=%v,TenantNo=%v", info.ID, info.TenantNo)
		idsNew = append(idsNew, info.ID)
	}
	if len(idsNew) > 0 {
		cn.DeleteByIds(ctx, idsNew)
	}
	return rt.Ok()
}

// Query 查询
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BlogBookmarkService) Query(ctx *gin.Context, ct modBlogBookmark.QueryCt) (rt rg.Rs[pagePg.Paginator[modBlogBookmark.Vo]]) {
	var query entityBlog.BlogBookmarkEntity
	copier.Copy(&query, &ct)
	slice := make([]modBlogBookmark.Vo, 0)
	rt.Data.Data = slice
	r := c.sv
	page, err := r.FindAllPage(ctx, query, optionsPg.WithOption(func(arg *optionsPg.OptionParams) {
		if ct.PageSize < 1 {
			ct.PageSize = 20
		}
		arg.Pageable = new(pagePg.PageablePageSize(0, ct.PageNum, ct.PageSize))
		arg.Db = arg.Db.Order("create_at desc")
		if strPg.IsNotBlank(ct.Wd) {
			arg.Db = arg.Db.Where("name like ?", "%"+ct.Wd+"%")
		}
		//可用范围
		if nil != ct.Where && len(ct.Where) > 0 {
			for _, tag := range ct.Where {
				arg.Db = arg.Db.Where("where @> ?", "[\""+tag+"\"]")
			}
		}
		//标签
		if nil != ct.TagsQuery && len(ct.TagsQuery) > 0 {
			for _, tag := range ct.TagsQuery {
				arg.Db = arg.Db.Where("tags @> ?", "[\""+tag+"\"]")
			}
		}
	}), optionsPg.WithCtx(ctx))
	if nil != err {
		return rt.Ok()
	}

	if page.Total > 0 && page.Data != nil && len(page.Data) > 0 {
		pg := pagePg.NewPaginatorByPageable[modBlogBookmark.Vo](page.Pageable)
		mapCategory := make(map[string]*entityBlog.BlogBookmarkCategoryEntity)
		idsCategory := make([]string, 0)
		for _, item := range page.Data {
			if strPg.IsNotBlank(item.CategoryNo) {
				idsCategory = append(idsCategory, item.CategoryNo)
			}
		}
		//分类
		if len(idsCategory) > 0 {
			tmp, result := c.catDb.FindAllByNoIn(ctx, idsCategory, optionsPg.WithCtx(ctx))
			if result {
				for _, item := range tmp {
					mapCategory[item.No] = item
				}
			}
		}
		//字段赋值
		for _, item := range page.Data {
			var vo modBlogBookmark.Vo
			copier.Copy(&vo, &item)
			vo.Tags = make([]string, 0)
			vo.Where = make([]string, 0)
			vo.Attachments = make(map[string]string)
			//附件
			if strPg.IsNotBlank(item.Attachments) {
				json.Unmarshal([]byte(item.Attachments), &vo.Attachments)
			}
			if nil == vo.Attachments {
				vo.Attachments = make(map[string]string)
			}
			//标签
			if nil != item.Tags.Data() {
				vo.Tags = item.Tags.Data()
			}
			//可用范围
			if nil != item.Where.Data() {
				vo.Where = item.Where.Data()
			}
			//分类
			if obj, ok := mapCategory[item.CategoryNo]; ok {
				vo.CategoryName = obj.Name
				vo.CategoryNo = obj.No
			}
			slice = append(slice, vo)
		}
		pg.Data = slice
		pg.Pageable = page.Pageable
		rt.Data = pg
		return rt.Ok()
	}
	return rt.Ok()
}

// QueryAll 查询
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BlogBookmarkService) QueryAll(ctx *gin.Context, ct modBlogBookmark.QueryCt) (rt rg.Rs[[]modBlogBookmark.Vo]) {
	var query entityBlog.BlogBookmarkEntity
	copier.Copy(&query, &ct)
	rt.Data = []modBlogBookmark.Vo{}
	infos := c.sv.FindAll(ctx, query)
	if len(infos) > 0 {
		slice := make([]modBlogBookmark.Vo, 0)
		for _, item := range infos {
			var vo modBlogBookmark.Vo
			copier.Copy(&vo, &item)
			slice = append(slice, vo)
		}
		rt.Data = slice
	}
	return rt.Ok()
}

// SelectNodePublic 查询
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BlogBookmarkService) SelectNodePublic(ctx *gin.Context, ct modBlogBookmark.QueryCt) (rt rg.Rs[[]model.BaseNode]) {
	var query entityBlog.BlogBookmarkEntity
	copier.Copy(&query, &ct)
	slice := make([]model.BaseNode, 0)
	rt.Data = slice
	infos := c.sv.FindAll(ctx, query)
	if len(infos) > 0 {
		for _, item := range infos {
			slice = append(slice, model.BaseNode{Key: numberPg.Int64ToString(item.ID), Label: item.Name})
		}
		rt.Data = slice
	}
	return rt.Ok()
}

// SelectNodeAllPublic 查询
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BlogBookmarkService) SelectNodeAllPublic(ctx *gin.Context, ct modBlogBookmark.QueryCt) (rt rg.Rs[[]model.BaseNode]) {
	var query entityBlog.BlogBookmarkEntity
	copier.Copy(&query, &ct)
	slice := make([]model.BaseNode, 0)
	rt.Data = slice
	infos := c.sv.FindAll(ctx, query)
	if len(infos) > 0 {
		for _, item := range infos {
			var vo modBlogBookmark.Vo
			copier.Copy(&vo, &item)
			slice = append(slice, model.BaseNode{Key: numberPg.Int64ToString(item.ID), Label: item.Name, Extend: vo})
		}
		rt.Data = slice
	}
	return rt.Ok()
}

// ExistName 查重
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BlogBookmarkService) ExistName(ctx *gin.Context, ct model.BaseExistWdCt[string]) (rt rg.Rs[string]) {
	if "" == ct.Wd {
		return rt.ErrorMessage("查询内容不能为空")
	}
	id := "0"
	if strPg.IsNotBlank(ct.Id) {
		id = ct.Id
	}
	_, result := c.sv.FindByNameAndIdNot(ctx, ct.Wd, id, optionsPg.WithCtx(ctx))
	if result {
		return rt.ErrorMessage("重复，已存在")
	}
	return rt.OkMessage("可以使用")
}

// ExistNo 查重
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BlogBookmarkService) ExistNo(ctx *gin.Context, ct model.BaseExistWdCt[string]) (rt rg.Rs[string]) {
	if "" == ct.Wd {
		return rt.ErrorMessage("查询内容不能为空")
	}
	id := "0"
	if strPg.IsNotBlank(ct.Id) {
		id = ct.Id
	}
	_, result := c.sv.FindByNoAndIdNot(ctx, ct.Wd, id, optionsPg.WithCtx(ctx))
	if result {
		return rt.ErrorMessage("重复，已存在")
	}
	return rt.OkMessage("可以使用")
}
