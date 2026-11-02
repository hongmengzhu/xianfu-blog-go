package modAttachment

type AttachmentDto struct {
	Name          string `json:"name" label:"名称" `
	SourceName    string `json:"sourceName" label:"名称" `
	Url           string `json:"url" label:"url" `   //全路径
	File          string `json:"file" label:"相对路径" ` //开头 以 / 为开头
	Size          int64  `json:"size" label:"大小" `
	Module        string `json:"module" label:"模块" `
	Value         string `json:"value" label:"值id" `
	Tag           string `json:"tag" label:"标签" `
	Label         string `json:"label" label:"标记" `
	Domain        string `json:"domain" label:"域名" `
	Ext           string `json:"ext" label:"文件扩展名" `
	No            string `json:"no" comment:"流水号" `
	Method        string `json:"method" comment:"方式" `
	Category      string `json:"category" comment:"上传分类" `
	Client        string `json:"client" comment:"客户端" `
	ProtocolSpace string `json:"protocolSpace" comment:"协议空间" `
}
