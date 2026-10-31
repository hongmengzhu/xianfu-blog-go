package utilsBlog

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constTags"
	"github.com/pangu-2/go-tools/tools/cryptPg"
	"strings"
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
