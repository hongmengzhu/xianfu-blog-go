package blogTopic

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/event/basic/model/modEventBasicTags"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/blog/modBlogTopic"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityBlog"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryBasic"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryBlog"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/holderPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/automatedPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constTags"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/blog/typeContentPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/blog/typeDataSourcePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/blog/typeSourcePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/blog/wherePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/content/enumEditorPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/yesNoPg/yesNoIntPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/yesNoPg/yesNoString"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/event"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/repositoryPg/optionsPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/versionPg"
	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/noPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
	"gorm.io/datatypes"
)

func init() {
	gs.Provide(new(Sp))
}

type Sp struct {
	sv            *repositoryBlog.BlogTopicRepository           `autowire:"?"`
	statisticsDb  *repositoryBlog.BlogTopicStatisticsRepository `autowire:"?"`
	catDb         *repositoryBlog.BlogTopicCategoryRepository   `autowire:"?"`
	AttachmentDao *repositoryBasic.BasicAttachmentRepository    `autowire:"?"`
	Bus           event.Bus                                     `autowire:"?"`
}

func New(sp *Sp, holder holderPg.HolderPg, ct modBlogTopic.CreateUpdateCt, isUpdate bool) *CreateUpdate {
	return &CreateUpdate{
		sp:             sp,
		holder:         holder,
		ct:             ct,
		isUpdate:       isUpdate,
		module:         constTags.ArticleInfo.Index(),
		images:         make(map[string]string),
		record:         &entityBlog.BlogTopicEntity{},
		recordCategory: &entityBlog.BlogTopicCategoryEntity{},
		entitySave:     &entityBlog.BlogTopicEntity{},
		tags:           make([]string, 0),
	}
}

type CreateUpdate struct {
	sp *Sp `autowire:"?"`

	isUpdate bool //是否更新
	ct       modBlogTopic.CreateUpdateCt
	holder   holderPg.HolderPg
	//
	record         *entityBlog.BlogTopicEntity //查询记录专用
	recordCategory *entityBlog.BlogTopicCategoryEntity
	entitySave     *entityBlog.BlogTopicEntity //保存专用
	images         map[string]string           //图片集
	module         string
	tags           []string
}

// Process
//
//	@Description: 处理
//	@receiver c
//	@param ctx
//	@return rt
func (c *CreateUpdate) Process(ctx *gin.Context) (rt rg.Rs[string]) {
	log.Infof(ctx, log.TagAppDef, "ct=%+v", c.ct)
	c.entitySave = &entityBlog.BlogTopicEntity{}
	ret := c.verify(ctx)
	if ret.ErrorIs() {
		return ret
	}
	return c.save(ctx)
}

// verify 验证
//
//	@Description:
//	@receiver c
//	@param ctx
//	@return rt
func (c *CreateUpdate) verify(ctx *gin.Context) (rt rg.Rs[string]) {
	ct := c.ct
	copier.Copy(c.entitySave, &c.ct)
	if "" == ct.Name {
		return rt.ErrorMessage("名称不能为空")
	}
	//if "" == ct.CategoryNo {
	//	return rt.ErrorMessage("分类不能为空")
	//}
	if strPg.IsBlank(c.entitySave.Code) {
		c.entitySave.Code = automatedPg.CREATE_CODE
	}
	if strPg.IsBlank(ct.TypeContent) {
		c.entitySave.TypeContent = typeContentPg.ORIGINAL.String()
	}
	if !typeContentPg.IsExistTypeContent(c.entitySave.TypeContent) {
		return rt.ErrorMessage("内容类型错误")
	}
	if strPg.IsBlank(ct.TypeSource) {
		c.entitySave.TypeSource = typeSourcePg.HANDWRITTEN.String()
	}
	if !typeSourcePg.IsExistTypeSource(c.entitySave.TypeSource) {
		return rt.ErrorMessage("类型源错误")
	}

	if strPg.IsBlank(ct.TypeDataSource) {
		c.entitySave.TypeDataSource = typeDataSourcePg.PLATFORM.String()
	}
	if !typeDataSourcePg.IsExistTypeDataSource(c.entitySave.TypeDataSource) {
		return rt.ErrorMessage("数据源错误")
	}

	if strPg.IsBlank(ct.Editor) {
		c.entitySave.Editor = enumEditorPg.Markdown.String()
	}
	if !enumEditorPg.IsExistEditor(c.entitySave.Editor) {
		return rt.ErrorMessage("编辑器类型错误")
	}
	if ct.Jump < 0 {
		c.entitySave.Jump = yesNoIntPg.No.IndexInt8()
	}
	if _, ok := yesNoIntPg.IsExistInt8(c.entitySave.Jump); !ok {
		return rt.ErrorMessage("跳转类型错误")
	}

	if strPg.IsBlank(ct.TypeComment) {
		c.entitySave.TypeComment = yesNoString.No.String()
	}
	if _, ok := yesNoString.IsExistString(c.entitySave.TypeComment); !ok {
		return rt.ErrorMessage("评论类型错误")
	}
	//判断是否是自动,不是自动
	if !automatedPg.IsCreateCode(c.entitySave.Code) {
		//判断格式是否满足要求
		if !automatedPg.FormatVerify(c.entitySave.Code) {
			return rt.ErrorMessage("标志格式不能为空")
		}
		//不是自动
		_, result := c.sp.sv.FindByCode(ctx, c.entitySave.Code, optionsPg.WithCtx(ctx))
		if result {
			return rt.ErrorMessage("标志已存在")
		}
	}
	c.entitySave.Name = strings.TrimSpace(c.entitySave.Name)
	//更新
	if c.isUpdate {
		if ct.ID.ToInt64() <= 0 {
			return rt.ErrorMessage("ID不能为空")
		}
		result := false
		c.record, result = c.sp.sv.FindById(ctx, ct.ID.ToInt64())
		if !result {
			return rt.ErrorMessage("数据不存在")
		}
		if strPg.IsNotBlank(c.ct.CategoryNo) {
			c.recordCategory, result = c.sp.catDb.FindByNo(ctx, c.ct.CategoryNo, optionsPg.WithCtx(ctx))
			if !result {
				return rt.ErrorMessage("数据不存在")
			}
		}

	} else {
		//判断是否是自动,不是自动
		//if !automatedPg.IsCreateCode(ct.No) {
		//	//判断格式是否满足要求
		//	if !automatedPg.FormatVerify(ct.No) {
		//		return rt.ErrorMessage("编号格式不能为空")
		//	}
		//	//不是自动
		//	_, result := r.FindByNo(ctx, ct.No)
		//	if result {
		//		return rt.ErrorMessage("编号已存在")
		//	}
		//}
	}

	return rt.Ok()
}

// save
//
//	@Description: 保存
//	@receiver c
//	@return rt
func (c CreateUpdate) save(ctx *gin.Context) (rt rg.Rs[string]) {
	where := make([]string, 0)
	tags := make([]string, 0)
	ct := c.ct
	if nil != ct.Where && len(ct.Where) > 0 {
		isFind := false
		for _, v := range ct.Where {
			if wherePg.IsExistWhere(v) {
				where = append(where, v)
				isFind = true
			}
		}
		if !isFind {
			return rt.ErrorMessage("编辑器类型错误")
		}
	} else {
		where = append(where, wherePg.ALL.String())
	}
	if nil != ct.Tags && len(ct.Tags) > 0 {
		for _, v := range ct.Tags {
			if strPg.IsNotBlank(v) {
				tags = append(tags, strings.TrimSpace(v))
			}
		}
		c.tags = tags
	}
	if len(where) < 1 {
		return rt.ErrorMessage("编辑器错误")
	}
	c.entitySave.Where = datatypes.NewJSONType(where)
	c.entitySave.Tags = datatypes.NewJSONType(tags)
	catDb := c.sp.catDb
	if strPg.IsNotBlank(ct.CategoryNo) {
		_, result := catDb.FindByNo(ctx, ct.CategoryNo, optionsPg.WithCtx(ctx))
		if !result {
			return rt.ErrorMessage("分类不存在")
		}
	}
	statistics := entityBlog.BlogTopicStatisticsEntity{}
	statistics.SeoKeywords = strings.TrimSpace(c.ct.SeoKeywords)
	statistics.SeoDescription = strings.TrimSpace(c.ct.SeoDescription)
	statistics.PageTitle = strings.TrimSpace(c.ct.PageTitle)
	//生成 版本时间戳
	c.entitySave.Version = versionPg.Make()
	r := c.sp.sv
	if c.isUpdate {
		var info entityBlog.BlogTopicEntity
		copier.Copy(&info, c.entitySave)
		info.ID = 0
		log.Infof(ctx, log.TagAppDef, "info=%+v", info)
		err := r.Update(ctx, info, c.ct.ID.ToInt64())
		if err != nil {
			log.Debugf(ctx, log.TagAppDef, "save err=%+v", err.Error())
			return rt.ErrorMessage("保存失败：" + err.Error())
		}
		//
		err = c.sp.statisticsDb.Update(ctx, statistics, c.ct.ID.ToInt64())
		if err != nil {
			log.Errorf(ctx, log.TagAppDef, "save statistics err=%+v", err)
		}
	} else {
		log.Infof(ctx, log.TagAppDef, "info=%+v", c.entitySave)
		c.entitySave.TenantNo = c.holder.GetTenantNo()
		c.entitySave.No = noPg.No()
		if automatedPg.IsCreateCode(c.entitySave.Code) {
			c.entitySave.Code = c.entitySave.No
		}
		err, _ := r.Create(ctx, c.entitySave)
		if err != nil {
			log.Debugf(ctx, log.TagAppDef, "save err=%+v", err)
			return rt.ErrorMessage("保存失败：" + err.Error())
		}
		log.Infof(ctx, log.TagAppDef, "save=%+v", c.entitySave)
		c.record = c.entitySave
		//
		statistics.ID = c.entitySave.ID
		statistics.TopicNo = c.entitySave.No
		err, _ = c.sp.statisticsDb.Create(ctx, &statistics)
		if err != nil {
			log.Errorf(ctx, log.TagAppDef, "save statistics err=%+v", err)
		}
	}
	//
	c.attachment(ctx)
	// 更新 附件图
	var imageEnt entityBlog.BlogTopicEntity
	//附件图
	images, err := json.Marshal(c.images)
	if err == nil {
		imageEnt.Attachments = string(images)
	} else {
		imageEnt.Attachments = "{}"
	}
	c.sp.sv.Update(ctx, imageEnt, c.record.ID)
	//更新标签
	c.tagsListener(ctx)

	return rt.Ok()
}

// attachment
//
//	@Description: 附件处理
//	@receiver c
//	@param ctx
func (c *CreateUpdate) attachment(ctx *gin.Context) {
	//更新附件
	if nil == c.ct.Attachment {
		c.images = make(map[string]string)
		return
	}
	c.images = c.ct.Attachment
}

func (c *CreateUpdate) tagsListener(ctx *gin.Context) {
	if len(c.tags) < 1 {
		return
	}
	//保存到数据库

	err := c.sp.Bus.Publish(context.Background(), modEventBasicTags.TagsRelation{
		Category: constTags.ArticleInfo.Index(),
		Tags:     c.tags,
		Holder:   c.holder,
	})
	if err != nil {
		log.Errorf(ctx, log.TagAppDef, "error:%+v", err)
	}
}
