package modBlogArticleCategory

import "github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"

type Vo struct {
	ID          typePg.Uint64String `json:"id" label:"id" `
	Name        string              `json:"name" label:"名称" `
	NameFl      string              `json:"nameFl" label:"名称外文" `
	No          string              `json:"no" label:"编号代号" `
	State       typePg.Int8         `json:"state" label:"状态:1启用;2禁用" `
	Description string              `json:"description" label:"描述" `
	CreateAt    *typePg.Time        `json:"createAt" label:"创建时间" `
	ParentNo    string              `json:"parentNo" label:"上级" `
	ParentId    string              `json:"parentId" label:"上级" `
}
