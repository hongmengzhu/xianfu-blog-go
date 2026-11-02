package modBasicCountry

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type QueryPublicCt struct {
	model.BaseQueryNodeCt
	ID          typePg.Int64String `json:"id" label:"" `
	Name        string             `json:"name" label:"名称" `          // 名称
	Code        string             `json:"code" label:"编号代号" `        // 编号代号
	Description string             `json:"description" label:"描述" `   // 描述
	State       typePg.Int8        `json:"state" label:"状态:1启用;2禁用" ` // 状态:1启用;2禁用
	ParentId    string             `json:"parentId" label:"上级" `
	ParentNo    string             `json:"parentNo" label:"上级编号" `
	Iso3        string             `json:"iso3" label:"ISO三字代码" `
	CountryCode string             `json:"countryCode" label:"国际区号" `
	PhoneUse    typePg.Int8        `json:"phoneUse" label:"电话使用1是2否" `
}
