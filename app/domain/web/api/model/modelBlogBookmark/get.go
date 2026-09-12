package modelBlogBookmark

import "github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"

type Get struct {
	State typePg.Int8 `json:"state"`
}
