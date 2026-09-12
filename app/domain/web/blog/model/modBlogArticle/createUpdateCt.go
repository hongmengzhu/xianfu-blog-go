package modBlogArticle

import (
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/components/attachmentPg/modAttachment"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type CreateUpdateCt struct {
	ID             typePg.Uint64String                   `json:"id" form:"id" validate:"required" label:"id" `
	Name           string                                `json:"name" form:"name" validate:"required,min=1,max=255" label:"名称" ` // 名称
	NameFl         string                                `json:"nameFl" label:"名称外文" `                                           // 名称外文
	No             string                                `json:"no" label:"编号代号" `                                               // 编号代号
	NameFull       string                                `json:"nameFull" label:"全称" `                                           // 全称
	Description    string                                `json:"description" label:"描述" `                                        // 描述
	CategoryNo     string                                `json:"categoryNo" label:"分类" `
	Content        string                                `json:"content" label:"内容" `
	Editor         string                                `json:"editor" label:"编辑器类型" `
	TypeContent    string                                `json:"typeContent" label:"内容类型" `
	TypeSource     string                                `json:"typeSource" label:"内容来源" `
	TypeDataSource string                                `json:"typeDataSource" label:"数据来源" `
	Where          []string                              `json:"where" label:"发布范围" `
	Jump           typePg.Int8                           `json:"jump" label:"跳转类型:1跳转;2不跳转" `
	Source         string                                `json:"source" label:"来源" `
	TypeComment    string                                `json:"typeComment" label:"评论类型" `
	TypeReading    string                                `json:"typeReading" label:"阅读类型" `
	Tags           []string                              `json:"tags" label:"标签" `
	SeoKeywords    string                                `json:"seoKeywords" label:"seo关键词" `
	SeoDescription string                                `json:"seoDescription" label:"seo描述" `
	PageTitle      string                                `json:"pageTitle" label:"网页标题" `
	OperationTime  *typePg.Time                          `json:"operationTime" label:"操作时间" `
	Attachment     map[string][]modAttachment.Attachment `json:"attachment" label:"图集" `
}
