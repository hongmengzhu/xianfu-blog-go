package modRamAccount

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type PasswordCt struct {
	ID          typePg.Uint64String `json:"id" form:"id" binding:"required" validate:"required" label:"id" `
	PasswordNew string              `json:"passwordNew" binding:"required" validate:"required" label:"新密码" ` // 新密码
	Code        string              `json:"code" label:"验证码" `                                               // 验证码
	Captcha     string              `json:"captcha" label:"验证码" `                                            // 验证码
	CaptchaNo   string              `json:"captchaNo" label:"验证码编号" `                                        // 验证码编号
	AuthCode    string              `json:"authCode" label:"授权码" `                                           // 授权码
	MfaCode     string              `json:"mfaCode" label:"多因素身份验证码" `                                       // 多因素身份验证码
}
