package modelBlogBookmark

import (
	"time"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type Vo struct {
	ID            typePg.Uint64String `json:"id" form:"id" label:"id" `
	No            string              `json:"no"`
	Name          string              `json:"name"`
	State         typePg.Int8         `json:"state"`
	CreateAt      *time.Time          `json:"createAt" label:"创建时间" `
	CategoryNo    string              `json:"categoryNo" label:"分类" `
	CategoryName  string              `json:"categoryName" label:"分类" `
	TeamNo        string              `json:"teamNo" label:"团队" `
	UrlSource     string              `json:"urlSource" label:"来源地址" `
	OperationTime *time.Time          `json:"operationTime" label:"操作时间" `
	PageTitle     string              `json:"pageTitle" label:"网页标题" `
}
