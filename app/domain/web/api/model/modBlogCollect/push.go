package modBlogCollect

import "github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"

type PushCt struct {
	CategoryNo  string       `json:"categoryNo"`
	Url         string       `json:"url"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Editor      string       `json:"editor" label:"编辑器类型"`
	Content     string       `json:"content" label:"内容"`
	Author      string       `json:"author" label:"作者"`
	Source      string       `json:"source" label:"来源平台"`
	Tags        []string     `json:"tags" label:"标签"`
	Rule        []string     `json:"rule" label:"规则"`
	PublishTime *typePg.Time `json:"publishTime" label:"源文章发布时间"`
}
