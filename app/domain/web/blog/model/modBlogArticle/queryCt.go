package modBlogArticle

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type QueryCt struct {
	model.BaseQueryCt
	ID            typePg.Uint64String `json:"id" label:"" `
	Name          string              `json:"name" label:"名称" `
	No            string              `json:"no" label:"编号代号" `
	Q             string              `json:"q" label:"关键词"  form:"q"`
	TagsQuery     []string            `json:"tagsQuery" label:"标签" form:"tags"`
	CategoryQuery []string            `json:"categoryQuery" label:"标签" form:"categoryQuery"`
	CreateAtStart *typePg.Time        `json:"createAtStart" label:"创建时间开始" `
	CreateAtEnd   *typePg.Time        `json:"createAtEnd" label:"创建时间结束" `
}
