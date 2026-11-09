package modRamIdpBinding

import (
	"time"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type Vo struct {
	ID            typePg.Uint64String `json:"id" label:"id" `
	Description   string              `json:"description" label:"描述" `   // 描述
	State         typePg.Int8         `json:"state" label:"状态:1启用;2禁用" ` // 状态:1启用;2禁用
	StateBind     typePg.Int8         `json:"stateBind" label:"绑定状态" `
	Idp           string              `json:"idp" label:"身份提供商" `
	OpenId        string              `json:"openId" label:"OpenID" `
	UnionId       string              `json:"unionId" label:"UnionID" `
	AppMark       string              `json:"appMark" label:"应用标识" `
	Platform      string              `json:"platform" label:"平台" `
	Protocol      string              `json:"protocol" label:"协议" `
	Mail          string              `json:"mail" label:"邮箱" `
	Phone         string              `json:"phone" label:"手机" `
	NickName      string              `json:"nickName" label:"昵称" `
	Avatar        string              `json:"avatar" label:"头像" `
	BindAno       string              `json:"bindAno" label:"绑定账号" `
	ExternalSub   string              `json:"externalSub" label:"外部唯一身份标识" `
	BindTime      *time.Time          `json:"bindTime" label:"绑定时间" `
	UnBindTime    *time.Time          `json:"unBindTime" label:"解绑时间" `
	LastLoginTime *time.Time          `json:"lastLoginTime" label:"最后登录时间" `
	CreateAt      *time.Time          `json:"createAt" label:"创建时间" ` // 创建时间
	UpdateAt      *time.Time          `json:"updateAt" label:"更新时间" ` // 更新时间
	CreateBy      string              `json:"createBy" label:"创建人" `  // 创建人
	UpdateBy      string              `json:"updateBy" label:"更新人" `  // 更新人
}
