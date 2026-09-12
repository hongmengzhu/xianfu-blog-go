package modelBlogBookmark

import "github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"

type CreateUpdate struct {
	ID          typePg.Uint64String `json:"id" form:"id" label:"id" `
	Name        string              `json:"name" form:"name" validate:"required,min=1,max=255" label:"名称" `
	Description string              `json:"description" label:"描述" `
	ParentNo    string              `json:"parentNo" label:"上级" `
}
