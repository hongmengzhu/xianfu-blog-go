package modApiBmBookmark

import (
	"github.com/hongmengzhu/xianfu-blog-go/app/models/bm/modBmBookmarkCategory/modApiBmBookmarkCategory"
)

type VoAll struct {
	My           []Vo                          `json:"my"`
	MyCategory   []modApiBmBookmarkCategory.Vo `json:"MyCategory"`
	Team         []Vo                          `json:"team"`
	TeamCategory []modApiBmBookmarkCategory.Vo `json:"teamCategory"`
	Org          []Vo                          `json:"org"`
	OrgCategory  []modApiBmBookmarkCategory.Vo `json:"orgCategory"`
}
