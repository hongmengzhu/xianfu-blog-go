package utilsBlog

import (
	"strings"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/consts/constTags"
	"github.com/pangu-2/go-tools/tools/cryptPg"
)

// TagCacheKey
//
//	@Description: 博客标签key
//	@param tag
//	@return string
func TagCacheKey(tag string) string {
	md5 := cryptPg.Md5(strings.TrimSpace(tag))
	return constTags.ArticleInfo.Index() + ":" + md5
}
