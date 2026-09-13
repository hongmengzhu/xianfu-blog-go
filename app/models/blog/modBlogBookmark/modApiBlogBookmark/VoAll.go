package modApiBlogBookmark

import (
	"github.com/hongmengzhu/xianfu-blog-go/app/models/blog/modBlogBookmarkCategory/modApiBlogBookmarkCategory"
)

type VoAll struct {
	My           []Vo                            `json:"my"`
	MyCategory   []modApiBlogBookmarkCategory.Vo `json:"MyCategory"`
	Team         []Vo                            `json:"team"`
	TeamCategory []modApiBlogBookmarkCategory.Vo `json:"teamCategory"`
}
