package modRamAccount

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

// CreateUpdateAccountCt
// @Description: 更新账号
type CreateUpdateAccountCt struct {
	ID                   typePg.Uint64String `json:"id" form:"id" validate:"required" label:"id" `
	Account              string              `json:"account"  validate:"required,min=1,max=255" label:"账户" ` // 账户
	Cc                   string              `json:"cc" label:"国际区号"`                                        // 国际区号
	Phone                string              `json:"phone" label:"手机号" `                                     // 手机号
	Mail                 string              `json:"mail" label:"邮箱" `                                       // 邮箱
	Code                 string              `json:"code" label:"编码" `                                       // 编码
	RegisterTime         *typePg.Time        `json:"registerTime" label:"注册时间" `                             // 注册时间
	Description          string              `json:"description" label:"描述" `                                // 描述
	PrincipalTypeService typePg.Int8         `json:"principalTypeService" label:"主账号类型" `                    // 主账号类型
}
