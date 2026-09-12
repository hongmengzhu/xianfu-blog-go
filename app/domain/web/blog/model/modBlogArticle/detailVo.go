package modBlogArticle

import (
	"html/template"
	"time"

	"github.com/hongmengzhu/xianfu-blog-go/app/models/basic/modBasicTagsRelation"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/blog/modBlogArticleCategory"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type DetailVo struct {
	ID              typePg.Uint64String             `json:"id" label:"id" `
	Name            string                          `json:"name" label:"名称" `     // 名称
	NameFl          string                          `json:"nameFl" label:"名称外文" ` // 名称外文
	No              string                          `json:"no" label:"编号代号" `
	NameFull        string                          `json:"nameFull" label:"全称" `      // 全称
	State           typePg.Int8                     `json:"state" label:"状态:1启用;2禁用" ` // 状态:1启用;2禁用
	Description     string                          `json:"description" label:"描述" `   // 描述
	CreateAt        time.Time                       `json:"createAt" label:"创建时间" `    // 创建时间
	UpdateAt        time.Time                       `json:"updateAt" label:"更新时间" `    // 更新时间
	CreateBy        string                          `json:"createBy" label:"创建人" `     // 创建人
	UpdateBy        string                          `json:"updateBy" label:"更新人" `     // 更新人
	CategoryNo      string                          `json:"categoryNo" label:"分类" `
	CategoryName    string                          `json:"categoryName" label:"分类" `
	Version         string                          `json:"version" label:"版本" `
	Content         string                          `json:"content" label:"内容" `
	ContentConv     string                          `json:"contentConv" label:"内容" `
	ContentConvHtml template.HTML                   `json:"contentConvHtml" label:"内容" `
	Editor          string                          `json:"editor" label:"编辑器类型" `
	Tags            []string                        `json:"tags" label:"标签" `
	Where           []string                        `json:"where" label:"发布范围" `
	TypeContent     string                          `json:"typeContent" label:"内容类型" `
	TypeSource      string                          `json:"typeSource" label:"内容来源" `
	TypeDataSource  string                          `json:"typeDataSource" label:"数据来源" `
	Jump            typePg.Int8                     `json:"jump" label:"跳转类型:1跳转;2不跳转" `
	Source          string                          `json:"source" label:"来源" `
	TypeComment     string                          `json:"typeComment" label:"评论类型" `
	TypeReading     string                          `json:"typeReading" label:"阅读类型" `
	OperationTime   *typePg.Time                    `json:"operationTime" label:"操作时间" `
	AttachmentsMap  map[string]string               `json:"attachment" label:"图集" `
	Author          string                          `json:"author" label:"作者" `
	UrlSource       string                          `json:"urlSource" label:"来源网址" `
	Statistics      StatisticsVo                    `json:"statistics" label:"统计" `
	TagsStyle       []modBasicTagsRelation.AllVo    `json:"tagsStyle" label:"标签" `
	CategoryObj     []*modBlogArticleCategory.Cache `json:"categoryObj" label:"分类"`
}
