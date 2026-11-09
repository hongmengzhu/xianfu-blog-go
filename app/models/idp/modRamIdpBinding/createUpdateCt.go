package modRamIdpBinding

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
	"time"
)

type CreateUpdateCt struct {
	ID           typePg.Uint64String `json:"id" form:"id" label:"id" `
	Description  string              `json:"description" label:"描述" ` // 描述
	Idp          string              `json:"idp" label:"身份提供商" `
	OpenId       string              `json:"openId" label:"OpenID" `
	UnionId      string              `json:"unionId" label:"UnionID" `
	AppMark      string              `json:"appMark" label:"应用标识" `
	Platform     string              `json:"platform" label:"平台" `
	Protocol     string              `json:"protocol" label:"协议" `
	Mail         string              `json:"mail" label:"邮箱" `
	Phone        string              `json:"phone" label:"手机" `
	NickName     string              `json:"nickName" label:"昵称" `
	Avatar       string              `json:"avatar" label:"头像" `
	BindAno      string              `json:"bindAno" label:"绑定账号" `
	ExternalSub  string              `json:"externalSub" label:"外部唯一身份标识" `
	BindTime     *time.Time          `json:"bindTime" label:"绑定时间" `
}
