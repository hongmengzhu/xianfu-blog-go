package modBlogArticleCategory

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
)

type QueryPublicCt struct {
	model.BaseQueryNodeCt
	ParentNo string `json:"parentNo" label:"上级" `
}
