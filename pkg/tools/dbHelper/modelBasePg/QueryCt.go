package modelBasePg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type QueryCt struct {
	model.BaseQueryCt
	State typePg.Int8 `json:"state" label:"状态:1启用;2禁用" `
}
