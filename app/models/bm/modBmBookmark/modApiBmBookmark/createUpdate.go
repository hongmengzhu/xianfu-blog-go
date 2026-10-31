package modApiBmBookmark

import (
	"time"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type CreateUpdateCt struct {
	ID          typePg.Uint64String `json:"id" form:"id" label:"id" `
	Name        string              `json:"name" form:"name" validate:"required,min=1,max=255" label:"名称" `
	Description string              `json:"description" label:"描述" `
	Content     string              `json:"content" label:"内容" `
	UrlSource   string              `json:"urlSource" label:"来源地址" `
	PublishTime *time.Time          `json:"publishTime" label:"发布时间" `
	TypeReading string              `json:"typeReading" label:"阅读类型|未看|在看|已看;" `
	Remark      string              `json:"remark" label:"备注" `
	Tags        []string            `json:"tags" label:"标签" `
	CategoryNo  string              `json:"categoryNo" label:"分类" `
	OrgNo       string              `json:"orgNo" label:"组织" `
}
