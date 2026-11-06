package modTcTenantDomainModule

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type CreateUpdateCt struct {
	ID          typePg.Int64String `json:"id" form:"id" label:"id" `
	Name        string             `json:"name" form:"name" validate:"required,min=1,max=255" label:"名称" ` // 名称
	NameFl      string             `json:"nameFl" label:"名称外文" `                                           // 名称外文
	Code        string             `json:"code" form:"code"  label:"编号代号" `                                // 编号代号
	NameFull    string             `json:"nameFull" label:"全称" `                                           // 全称
	Description string             `json:"description" label:"描述" `                                        // 描述
}
