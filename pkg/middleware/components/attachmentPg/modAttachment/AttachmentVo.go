package modAttachment

type AttachmentVo struct {
	Name       string `json:"name" form:"name" label:"名称" `
	SourceName string `json:"sourceName" label:"原名称" `
	Url        string `json:"url" label:"全部网址" `
	Domain     string `json:"domain" label:"域名" `
}
