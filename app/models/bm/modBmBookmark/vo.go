package modBmBookmark

import (
	"time"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type Vo struct {
	ID               typePg.Uint64String `json:"id" label:"id" `
	Name             string              `json:"name" label:"名称" `
	NameFl           string              `json:"nameFl" label:"名称外文" `
	No               string              `json:"no" label:"编号代号" `
	Code             string              `json:"code" label:"标志" `
	NameFull         string              `json:"nameFull" label:"全称" `
	State            typePg.Int8         `json:"state" label:"状态:1启用;2禁用" `
	Description      string              `json:"description" label:"描述" `
	CreateAt         *time.Time          `json:"createAt" label:"创建时间" `
	UpdateAt         *time.Time          `json:"updateAt" label:"更新时间" `
	CreateBy         string              `json:"createBy" label:"创建人" `
	UpdateBy         string              `json:"updateBy" label:"更新人" `
	Sort             int64               `json:"sort" label:"排序" `
	CategoryNo       string              `json:"categoryNo" label:"分类" `
	CategoryName     string              `json:"categoryName" label:"分类" `
	Version          string              `json:"version" label:"版本" `
	Content          string              `json:"content" label:"内容" `
	Editor           string              `json:"editor" label:"编辑器类型" `
	Tags             []string            `json:"tags" label:"标签" `
	Author           string              `json:"author" label:"作者" `
	TypeContent      string              `json:"typeContent" label:"内容类型" `
	TypeSource       string              `json:"typeSource" label:"内容来源" `
	TypeDataSource   string              `json:"typeDataSource" label:"数据来源" `
	Where            []string            `json:"where" label:"可用范围" `
	Attachments      map[string]string   `json:"attachment" label:"图集" `
	UrlSource        string              `json:"urlSource" label:"来源网址" `
	UrlRewriting     string              `json:"urlRewriting" label:"URL重写" `
	Jump             typePg.Int8         `json:"jump" label:"跳转类型:1跳转;2不跳转" `
	Source           string              `json:"source" label:"来源" `
	TypeComment      string              `json:"typeComment" label:"评论类型" `
	OperationTime    *typePg.Time        `json:"operationTime" label:"操作时间" `
	PublishTime      *typePg.Time        `json:"publishTime" label:"发布时间" `
	TypeReading      string              `json:"typeReading" label:"阅读类型" `
	TypeDomain       string              `json:"typeDomain" label:"领域" `
	Remark           string              `json:"remark" label:"备注" `
	Comment          int64               `json:"comment" label:"评论" `
	Read             int64               `json:"read" label:"阅读" `
	SeoKeywords      string              `json:"seoKeywords" label:"seo关键词" `
	SeoDescription   string              `json:"seoDescription" label:"seo描述" `
	PageTitle        string              `json:"pageTitle" label:"网页标题" `
	TypeOwner        string              `json:"typeOwner" label:"类型" `
	TeamNo           string              `json:"teamNo" label:"团队编号" `
	TenantNo         string              `json:"tenantNo" label:"租户编号" `
	TenantNoName     string              `json:"tenantNoName" label:"租户" `
	PlatformApproved string              `json:"platformApproved" label:"平台审批状态" `
}
