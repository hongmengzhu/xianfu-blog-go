package modBasicModelRules

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type QueryPublicCt struct {
	model.BaseQueryNodeCt
	State typePg.Int8 `json:"state" label:"状态:1启用;2禁用" `
	//
	Description string `json:"description" label:"描述" ` // 描述
	Name        string `json:"name" label:"名称" `
	ValueNo     string `json:"valueNo" label:"值编号/模块编号" `
}
