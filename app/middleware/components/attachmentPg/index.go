package attachmentPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/components/attachmentPg/drive"
	"github.com/hongmengzhu/xianfu-blog-go/app/middleware/components/attachmentPg/types"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(drive.Local)).
		Export(gs.As[types.FileProvider]())
	// 本地存储，当 minio 不存在时才注册
	// 可以添加其它判断条件，例如 aliyun 等
	//On(cond.Group(cond.And, cond.OnMissingBean((*minio.Client)(nil))))

	//gs.Provide(new(Minio)).
	//	Export((*FileProvider)(nil)).
	//	On(cond.OnBean((*minio.Client)(nil)))
}
