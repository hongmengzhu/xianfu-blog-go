package attachmentPg

import (
	"context"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/urlGuardPg"
)

// SafeFetchClient 用于服务端抓取远程附件，限制超时、跳转次数并校验每一跳目标。
var SafeFetchClient = urlGuardPg.NewSafeClient(0, 3)

// ValidateFetchURL 校验目标 URL 是否允许服务端发起请求
func ValidateFetchURL(ctx context.Context, raw string) error {
	return urlGuardPg.ValidateFetchURL(ctx, raw)
}
