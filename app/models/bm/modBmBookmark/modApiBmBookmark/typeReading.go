package modApiBmBookmark

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type TypeReadingCt struct {
	ID          typePg.Uint64String `json:"id" form:"id" label:"id" `
	TypeReading string              `json:"typeReading" label:"阅读类型|未看|在看|已看;" `
	Remark      string              `json:"remark" label:"备注" `
	OrgNo       string              `json:"orgNo" label:"组织" `
}
