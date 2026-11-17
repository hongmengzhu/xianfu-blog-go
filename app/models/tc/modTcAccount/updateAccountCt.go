package modTcAccount

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

// UpdateAccountCt
// @Description: 更新账号
type UpdateAccountCt struct {
	ID           typePg.Uint64String `json:"id" form:"id" binding:"required" validate:"required" label:"id" `
	Account      string              `json:"account"  binding:"required,min=1,max=255" validate:"required,min=1,max=255" label:"账户" ` // 账户
	Cc           string              `json:"cc" label:"国际区号"`                                                                         // 国际区号
	Phone        string              `json:"phone" label:"手机号" `                                                                      // 手机号
	Mail         string              `json:"mail" label:"邮箱" `                                                                        // 邮箱
	Code         string              `json:"code" label:"编码" `                                                                        // 编码
	RegisterTime *typePg.Time        `json:"registerTime" label:"注册时间" `                                                              // 注册时间
	Description  string              `json:"description" label:"描述" `                                                                 // 描述
}
