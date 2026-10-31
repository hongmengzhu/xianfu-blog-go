package modApiBmBookmark

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/model"
)

type QueryCt struct {
	model.BaseQueryCt
	Name        string   `json:"name" label:"名称" `
	TypeReading string   `json:"typeReading" label:"阅读类型" `
	TypeDomain  string   `json:"typeDomain" label:"领域" `
	TypeOwner   string   `json:"typeOwner" label:"类型" `
	TagsQuery   []string `json:"tagsQuery" label:"标签" `
}
