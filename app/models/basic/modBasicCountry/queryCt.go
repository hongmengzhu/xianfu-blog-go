package modBasicCountry

import (
	"time"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type QueryCt struct {
	model.BaseQueryCt
	ID          typePg.Int64String `json:"id" label:"" `
	Name        string             `json:"name" label:"名称" `          // 名称
	NameFl      string             `json:"nameFl" label:"名称外文" `      // 名称外文
	Code        string             `json:"code" label:"编号代号" `        // 编号代号
	NameFull    string             `json:"nameFull" label:"全称" `      // 全称
	Description string             `json:"description" label:"描述" `   // 描述
	State       typePg.Int8        `json:"state" label:"状态:1启用;2禁用" ` // 状态:1启用;2禁用
	CreateBy    string             `json:"createBy" label:"创建人" `     // 创建人
	CreateAt    *time.Time         `json:"createAt" label:"创建时间" `    // 创建时间
	ParentId    string             `json:"parentId" label:"上级" `
	ParentNo    string             `json:"parentNo" label:"上级编号" `
	Iso3        string             `json:"iso3" label:"ISO三字代码" `
	CountryCode string             `json:"countryCode" label:"国际区号" `
	PhoneUse    typePg.Int8        `json:"phoneUse" label:"电话使用1是2否" `
}
