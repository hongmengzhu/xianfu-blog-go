package modelBlogBookmark

import "github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"

type GetAll struct {
	State typePg.Int8 `json:"state"`
}
