package modBasicDataDictionary

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type QueryCt struct {
	model.BaseQueryCt
	Name        string      `json:"name" label:"名称" `          // 名称
	NameFl      string      `json:"nameFl" label:"名称外文" `      // 名称外文
	Code        string      `json:"code" label:"编号代号" `        // 编号代号
	NameFull    string      `json:"nameFull" label:"全称" `      // 全称
	Description string      `json:"description" label:"描述" `   // 描述
	State       typePg.Int8 `json:"state" label:"状态:1启用;2禁用" ` // 状态:1启用;2禁用
	Range       string      `json:"range" label:"范围" `         // 范围
}
