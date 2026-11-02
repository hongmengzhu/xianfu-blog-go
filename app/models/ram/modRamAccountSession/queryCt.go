package modRamAccountSession

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type QueryCt struct {
	model.BaseQueryCt
	Ano           string         `json:"ano" label:"编号" `
	LoginSource   string         `json:"loginSource" label:"登录来源" `
	AppNo         string         `json:"appNo" label:"应用编号" `
	Client        string         `json:"client" label:"客户端" `
	TenantNo      string         `json:"tenantNo" label:"租户编号" `
	System        string         `json:"system" label:"系统" `
	Ip            string         `json:"ip" label:"ip" `
	UpdateAtRange []*typePg.Time `json:"updateAtRange" label:"更新时间" `
}
