package modelBlogBookmarkCategory

import "github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"

type Vo struct {
	ID       typePg.Int64String `json:"id" label:"id" `
	Name     string             `json:"name" label:"名称" `
	NameFl   string             `json:"nameFl" label:"名称外文" `
	No       string             `json:"no" label:"编号" `
	Code     string             `json:"code" label:"码值" `
	NameFull string             `json:"nameFull" label:"全称" `
	ParentNo string             `json:"parentNo" label:"上级" `
}
