package modApiBlogArticleCategory

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/model"
)

type QueryPublicCt struct {
	model.BaseQueryNodeCt
	ParentNo string `json:"parentNo" label:"上级" `
}
