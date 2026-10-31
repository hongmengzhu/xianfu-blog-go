package modBmBookmark

import (
	"time"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type QueryCt struct {
	model.BaseQueryCt
	ID             typePg.Uint64String `json:"id" label:"" `
	Name           string              `json:"name" label:"名称" `
	NameFl         string              `json:"nameFl" label:"名称外文" `
	No             string              `json:"no" label:"编号代号" `
	Description    string              `json:"description" label:"描述" `
	State          typePg.Int8         `json:"state" label:"状态:1启用;2禁用" `
	CreateBy       string              `json:"createBy" label:"创建人" `
	CreateAt       *time.Time          `json:"createAt" label:"创建时间" `
	CategoryNo     string              `json:"categoryNo" label:"分类" `
	Content        string              `json:"content" label:"内容" `
	Editor         string              `json:"editor" label:"编辑器类型" `
	Author         string              `json:"author" label:"作者" `
	TypeContent    string              `json:"typeContent" label:"内容类型" `
	TypeSource     string              `json:"typeSource" label:"内容来源" `
	TypeDataSource string              `json:"typeDataSource" label:"数据来源" `
	Where          []string            `json:"where" label:"可用范围" `
	Jump           typePg.Int8         `json:"jump" label:"跳转类型:1跳转;2不跳转" `
	Source         string              `json:"source" label:"来源" `
	TypeComment    string              `json:"typeComment" label:"评论类型" `
	TypeReading    string              `json:"typeReading" label:"阅读类型" `
	TypeDomain     string              `json:"typeDomain" label:"领域" `
	TypeOwner      string              `json:"typeOwner" label:"类型" `
	TeamNo         string              `json:"teamNo" label:"团队编号" `
	Tags           []string            `json:"tags" label:"标签" `
	TagsQuery      []string            `json:"tagsQuery" label:"标签" `
}
