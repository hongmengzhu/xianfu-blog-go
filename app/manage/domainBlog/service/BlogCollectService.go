package service

import (
	"encoding/json"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/manage/domainBlog/service/blogCollect"
	"github.com/hongmengzhu/xianfu-blog-go/app/manage/domainBlog/utilsBlog"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/basic/modBasicTagsRelation"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/blog/modBlogCollect"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBlog"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityTc"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryBasic"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryBlog"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryTc"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/holderPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/cachePg/rdsPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constTags"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg/optionsPg"
	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/dbPg/pagePg"
	"github.com/pangu-2/go-tools/tools/numberPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(BlogCollectService))
}

// BlogCollectService 说明分类,协议，服务
// @Description:
type BlogCollectService struct {
	sv        *repositoryBlog.BlogCollectRepository         `autowire:"?"`
	catDb     *repositoryBlog.BlogCollectCategoryRepository `autowire:"?"`
	topic     *repositoryBlog.BlogTopicRepository           `autowire:"?"`
	topicRel  *repositoryBlog.BlogTopicRelationRepository   `autowire:"?"`
	tagsRelat *repositoryBasic.BasicTagsRelationRepository  `autowire:"?"`
	ten       *repositoryTc.TcTenantRepository              `autowire:"?"`
	sp        *blogCollect.Sp                               `autowire:"?"`

	rdu *rdsPg.BatchString `autowire:"?"`
}

// Create 新增
//
//	@Description:
//	@receiver c
//	@param ct
//	@return rt
func (c *BlogCollectService) Create(ctx *gin.Context, ct modBlogCollect.CreateUpdateCt) (rt rg.Rs[string]) {
	return blogCollect.New(c.sp, holderPg.GetContextAccount(ctx), ct, false).Process(ctx)
}

// Update 更新
//
//	@Description:
//	@receiver c
//	@param ct
//	@return rt
func (c *BlogCollectService) Update(ctx *gin.Context, ct modBlogCollect.CreateUpdateCt) (rt rg.Rs[string]) {
	return blogCollect.New(c.sp, holderPg.GetContextAccount(ctx), ct, true).Process(ctx)
}

// Detail 详情
//
//	@Description:
//	@receiver c
//	@param id
func (c *BlogCollectService) Detail(ctx *gin.Context, id int64) (rt rg.Rs[modBlogCollect.DetailVo]) {
	if id < 1 {
		return rt.ErrorMessage("id错误")
	}
	find, b := c.sv.FindById(ctx, id)
	if !b {
		return rt.ErrorMessage("数据不存在")
	}
	var info modBlogCollect.DetailVo
	copier.Copy(&info, &find)
	tags := make([]string, 0)
	info.Tags = make([]string, 0)
	info.TagsStyle = make([]modBasicTagsRelation.AllVo, 0)
	info.Attachment = make(map[string]string)
	mapTagsOnly := make(map[string]bool)
	tagsData := make(map[string]modBasicTagsRelation.AllVo)
	//附件图
	if strPg.IsNotBlank(find.Attachments) {
		var imagesMap map[string]string
		err := json.Unmarshal([]byte(find.Attachments), &imagesMap)
		if err == nil {
			info.Attachment = imagesMap
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
				infos, result := c.tagsRelat.FindAllByTagNoInAndCategoryRoot(ctx, tags, constTags.CollectInfo.String())
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

	//统计

	//话题

	return rt.OkData(info)
}

// Enable 启用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BlogCollectService) Enable(ctx *gin.Context, ct model.BaseIdsCt[string]) (rt rg.Rs[string]) {
	return c.State(ctx, ct.Ids, enumStatePg.ENABLE)
}

// Disable 禁用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BlogCollectService) Disable(ctx *gin.Context, ct model.BaseIdsCt[string]) (rt rg.Rs[string]) {
	return c.State(ctx, ct.Ids, enumStatePg.GetType(enumStatePg.DISABLE))
}

// State 状态 启用/禁用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BlogCollectService) State(ctx *gin.Context, ids []string, state enumStatePg.State) (rt rg.Rs[string]) {
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
			r.Update(ctx, entityBlog.BlogCollectEntity{State: state.IndexInt8()}, info.ID)
		}
	}
	return rt.Ok()
}

// StateEnableDisable 状态 设置 有效 停用
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BlogCollectService) StateEnableDisable(ctx *gin.Context, ids []string, state enumStatePg.State) (rt rg.Rs[string]) {
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
func (c *BlogCollectService) LogicalDeletion(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
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
			// 有效 停用，反转 为对应的 取消 弃置
			if ok, reverse := enum.ReverseEnableDisable(); ok {
				repository.Update(ctx, entityBlog.BlogCollectEntity{State: reverse.IndexInt8()}, info.ID)
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
func (c *BlogCollectService) LogicalRecovery(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
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
		//  取消 弃置 批量删除，反转 为对应的 有效 停用 停用
		if ok, reverse := enum.ReverseCancelLayAside(); ok {
			repository.Update(ctx, entityBlog.BlogCollectEntity{State: reverse.IndexInt8()}, info.ID)
		}
	}
	return rt.Ok()
}

// PhysicalDeletion 物理删除
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BlogCollectService) PhysicalDeletion(ctx *gin.Context, ids []string) (rt rg.Rs[string]) {
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
func (c *BlogCollectService) Query(ctx *gin.Context, ct modBlogCollect.QueryCt) (rt rg.Rs[pagePg.Paginator[modBlogCollect.Vo]]) {
	var query entityBlog.BlogCollectEntity
	copier.Copy(&query, &ct)
	slice := make([]modBlogCollect.Vo, 0)
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
		if strPg.IsNotBlank(ct.Wd) {
			arg.Db = arg.Db.Where("name like ?", "%"+ct.Wd+"%")
		}
		//发布范围
		if nil != ct.Where && len(ct.Where) > 0 {
			for _, tag := range ct.Where {
				arg.Db = arg.Db.Where("where @> ?", "[\""+tag+"\"]")
			}
		}
		//标签
		if nil != ct.TagsQuery && len(ct.TagsQuery) > 0 {
			for _, tag := range ct.TagsQuery {
				//获取缓存，得到 标签编号
				get, b := c.rdu.Get(ctx, utilsBlog.TagCacheKey(tag))
				if b {
					arg.Db = arg.Db.Where("tags @> ?", get)
				} else {
					arg.Db = arg.Db.Where("tags @> ?", "[\""+tag+"\"]")
				}
			}
		}
	}), optionsPg.WithCtx(ctx))
	if nil != err {
		return rt.Ok()
	}

	if page.Total > 0 && page.Data != nil && len(page.Data) > 0 {
		pg := pagePg.NewPaginatorByPageable[modBlogCollect.Vo](page.Pageable)
		mapCategory := make(map[string]*entityBlog.BlogCollectCategoryEntity)
		mapTenant := make(map[string]*entityTc.TcTenantEntity)
		idsCategory := make([]string, 0)
		idsTenant := make([]string, 0)
		idsNo := make([]string, 0)
		tags := make([]string, 0)
		mapTagsOnly := make(map[string]bool)
		//标签
		tagsData := make(map[string]modBasicTagsRelation.AllVo)
		for _, item := range page.Data {
			if strPg.IsNotBlank(item.CategoryNo) {
				idsCategory = append(idsCategory, item.CategoryNo)
			}
			if strPg.IsNotBlank(item.TenantNo) && !slices.Contains(idsTenant, item.TenantNo) {
				idsTenant = append(idsTenant, item.TenantNo)
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
		//租户
		{
			if len(idsTenant) > 0 {
				tmp, result := c.ten.FindAllByNoIn(ctx, idsTenant, optionsPg.WithCtx(ctx))
				if result {
					for _, item := range tmp {
						mapTenant[item.No] = item
					}
				}
			}
		}
		//分类
		{
			if len(idsCategory) > 0 {
				tmp, result := c.catDb.FindAllByNoIn(ctx, idsCategory, optionsPg.WithCtx(ctx))
				if result {
					for _, item := range tmp {
						mapCategory[item.No] = item
					}
				}
			}
		}
		//统计
		{

		}
		//标签
		{
			if len(tags) > 0 {
				infos, result := c.tagsRelat.FindAllByTagNoInAndCategoryRoot(ctx, tags, constTags.CollectInfo.String())
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
			var vo modBlogCollect.Vo
			copier.Copy(&vo, &item)
			vo.Tags = make([]string, 0)
			vo.TagsStyle = make([]modBasicTagsRelation.AllVo, 0)
			//
			if strPg.IsNotBlank(item.Attachments) {
				json.Unmarshal([]byte(item.Attachments), &vo.Attachments)
			}
			if nil == vo.Attachments {
				vo.Attachments = make(map[string]string)
			}
			//统计

			//分类
			if obj, ok := mapCategory[item.CategoryNo]; ok {
				vo.CategoryName = obj.Name
				vo.CategoryNo = obj.No
			}
			//租户
			if obj, ok := mapTenant[item.TenantNo]; ok {
				vo.TenantNoName = obj.Name
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

// SelectNodePublic 查询
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BlogCollectService) SelectNodePublic(ctx *gin.Context, ct modBlogCollect.QueryCt) (rt rg.Rs[[]model.BaseNode]) {
	var query entityBlog.BlogCollectEntity
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
func (c *BlogCollectService) SelectNodeAllPublic(ctx *gin.Context, ct modBlogCollect.QueryCt) (rt rg.Rs[[]model.BaseNode]) {
	var query entityBlog.BlogCollectEntity
	copier.Copy(&query, &ct)
	slice := make([]model.BaseNode, 0)
	rt.Data = slice
	infos := c.sv.FindAll(ctx, query)
	if len(infos) > 0 {

		for _, item := range infos {
			var vo modBlogCollect.Vo
			copier.Copy(&vo, &item)
			slice = append(slice, model.BaseNode{Key: numberPg.Int64ToString(item.ID), Label: item.Name, Extend: vo})
		}
		rt.Data = slice
	}
	return rt.Ok()
}

// SelectPublic 查询
//
//	@Description:
//	@receiver c
//	@param ct
func (c *BlogCollectService) SelectPublic(ctx *gin.Context, ct modBlogCollect.QueryCt) (rt rg.Rs[[]modBlogCollect.Vo]) {
	var query entityBlog.BlogCollectEntity
	copier.Copy(&query, &ct)
	rt.Data = []modBlogCollect.Vo{}
	infos := c.sv.FindAll(ctx, query)
	if len(infos) > 0 {
		slice := make([]modBlogCollect.Vo, 0)
		for _, item := range infos {
			var vo modBlogCollect.Vo
			copier.Copy(&vo, &item)
			slice = append(slice, vo)
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
func (c *BlogCollectService) ExistName(ctx *gin.Context, ct model.BaseExistWdCt[string]) (rt rg.Rs[string]) {
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
func (c *BlogCollectService) ExistNo(ctx *gin.Context, ct model.BaseExistWdCt[string]) (rt rg.Rs[string]) {
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
