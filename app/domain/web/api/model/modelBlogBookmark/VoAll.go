package modelBlogBookmark

import (
	"github.com/hongmengzhu/xianfu-blog-go/app/domain/web/api/model/modelBlogBookmarkCategory"
)

type VoAll struct {
	My           []Vo                           `json:"my"`
	MyCategory   []modelBlogBookmarkCategory.Vo `json:"MyCategory"`
	Team         []Vo                           `json:"team"`
	TeamCategory []modelBlogBookmarkCategory.Vo `json:"teamCategory"`
}
