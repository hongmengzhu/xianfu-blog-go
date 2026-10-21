package configPg

import "github.com/hongmengzhu/xianfu-blog-go/pkg/configPg/pg"

type Pg struct {
	Domain pg.Domain `json:"domain" value:"${domain}" label:"域模块"`
	Jwt    pg.Jwt    `json:"jwt" value:"${jwt}"`
	//Upload        pg.Upload        `json:"upload" value:"${upload}"`
	Template pg.Template `json:"template" value:"${template}"`
	//Elasticsearch pg.Elasticsearch `json:"elasticsearch" value:"${elasticsearch}"`
	MultiTenant pg.MultiItem  `json:"multiTenant" value:"${multiTenant}" toml:"multiTenant"` //租户表
	Multi       pg.Multi      `json:"multi" value:"${multi}" toml:"multi"`                   //多租户/商户表
	Data        pg.Data       `json:"data" value:"${data}" toml:"data"`                      //数据
	Profiles    pg.Profiles   `value:"${profiles}" toml:"profiles"`                          //配置
	Attachment  pg.Attachment `json:"attachment" value:"${attachment}" toml:"attachment"`    //附件
	Mail        pg.Mail       `json:"mail" value:"${mail}" toml:"mail"`                      //邮件
	//Rocketmq    pg.Rocketmq   `json:"rocketmq" value:"${rocketmq}"` toml:"rocketmq"`                          //消息队列
	Directory pg.Directory `json:"directory" value:"${directory}" toml:"directory" label:"目录/文件夹"`
}

// 设置默认配置
//func (c *Pg) SetDefault() {
//	if len(c.Upload.UploadFolder) < 1 {
//		c.Upload.UploadFolder = "/uploads"
//	}
//	if len(c.Upload.UploadFolder) < 1 {
//		c.Upload.StaticAccessPath = "/uploads"
//	}
//}
