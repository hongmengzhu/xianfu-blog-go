package modBasicConfig

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type QueryCt struct {
	model.BaseQueryCt
	ID    typePg.Int64String `json:"id" label:"id"`
	State typePg.Int8        `json:"state" label:"状态:1启用;2禁用"`

	ModelNo     string `json:"modelNo" label:"模型编号"`
	Model       string `json:"model" label:"模型标识"`
	Module      string `json:"module" label:"模块"`
	ModuleSub   string `json:"moduleSub" label:"子模块"`
	Field       string `json:"field" label:"字段标识"`
	Name        string `json:"name" label:"名称"`
	Description string `json:"description" label:"描述"`
	CategoryNo  string `json:"categoryNo" label:"配置分类编号"`
}
