package modBmBookmark

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type CreateUpdateCt struct {
	ID             typePg.Uint64String `json:"id" form:"id" validate:"required" label:"id" `
	Name           string              `json:"name" form:"name" validate:"required,min=1,max=955" label:"名称" `
	NameFl         string              `json:"nameFl" label:"名称外文" `
	NameFull       string              `json:"nameFull" label:"全称" `
	Description    string              `json:"description" label:"描述" `
	CategoryNo     string              `json:"categoryNo" label:"分类" `
	Content        string              `json:"content" label:"内容" `
	Editor         string              `json:"editor" label:"编辑器类型" `
	Author         string              `json:"author" label:"作者" `
	TypeContent    string              `json:"typeContent" label:"内容类型" `
	TypeSource     string              `json:"typeSource" label:"内容来源" `
	TypeDataSource string              `json:"typeDataSource" label:"数据来源" `
	Where          []string            `json:"where" label:"可用范围" `
	Tags           []string            `json:"tags" label:"标签" `
	Jump           typePg.Int8         `json:"jump" label:"跳转类型:1跳转;2不跳转" `
	Source         string              `json:"source" label:"来源" `
	TypeComment    string              `json:"typeComment" label:"评论类型" `
	TypeReading    string              `json:"typeReading" label:"阅读类型" `
	TypeDomain     string              `json:"typeDomain" label:"领域" `
	Remark         string              `json:"remark" label:"备注" `
	SeoKeywords    string              `json:"seoKeywords" label:"seo关键词" `
	SeoDescription string              `json:"seoDescription" label:"seo描述" `
	PageTitle      string              `json:"pageTitle" label:"网页标题" `
	TypeOwner      string              `json:"typeOwner" label:"类型" `
	TeamNo         string              `json:"teamNo" label:"团队编号" `
	UrlSource      string              `json:"urlSource" label:"来源网址" `
	UrlRewriting   string              `json:"urlRewriting" label:"URL重写" `
	OperationTime  *typePg.Time        `json:"operationTime" label:"操作时间" `
	PublishTime    *typePg.Time        `json:"publishTime" label:"发布时间" `
	Attachment     map[string]string   `json:"attachment" label:"图集" `
}
