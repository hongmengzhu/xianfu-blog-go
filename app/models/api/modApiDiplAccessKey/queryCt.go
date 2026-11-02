package modApiDiplAccessKey

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type QueryCt struct {
	model.BaseQueryCt
	ID            typePg.Uint64String `json:"id" label:"" `
	Name          string              `json:"name" label:"名称" ` // 名称
	No            string              `json:"no" label:"编号代号"`
	Description   string              `json:"description" label:"描述" `   // 描述
	State         typePg.Int8         `json:"state" label:"状态:1启用;2禁用" ` // 状态:1启用;2禁用
	DiplNo        string              `json:"diplNo" label:"接口编号" `
	Ano           string              `json:"ano" label:"编号代号"`
	PrincipalType string              `json:"principalType" label:"主体类型" `
}
