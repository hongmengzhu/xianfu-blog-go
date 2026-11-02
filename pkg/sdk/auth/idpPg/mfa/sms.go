package mfa

import (
	"errors"
)

// SmsMfaUtil SMS MFA 工具（简化实现）
type SmsMfaUtil struct {
	*MfaProps
}

// EmailMfaUtil Email MFA 工具（复用 SMS 逻辑）
type EmailMfaUtil struct {
	*MfaProps
}

// NewSmsMfaUtil 创建 SMS MFA 工具实例
func NewSmsMfaUtil(config *MfaProps) *SmsMfaUtil {
	if config == nil {
		config = &MfaProps{
			MfaType: SmsType,
		}
	}
	return &SmsMfaUtil{config}
}

// NewEmailMfaUtil 创建 Email MFA 工具实例
func NewEmailMfaUtil(config *MfaProps) *EmailMfaUtil {
	if config == nil {
		config = &MfaProps{
			MfaType: EmailType,
		}
	}
	return &EmailMfaUtil{config}
}

// Initiate 初始化 SMS MFA
func (mfa *SmsMfaUtil) Initiate(userId string, issuer string) (*MfaProps, error) {
	return &MfaProps{
		MfaType: mfa.MfaType,
	}, nil
}

// SetupVerify 设置阶段验证 SMS 验证码
func (mfa *SmsMfaUtil) SetupVerify(passcode string) error {
	// TODO: 接入实际的短信验证码校验服务
	if passcode == "" {
		return errors.New("验证码不能为空")
	}
	return nil
}

// Verify 验证 SMS 验证码
func (mfa *SmsMfaUtil) Verify(passcode string) error {
	// TODO: 接入实际的短信验证码校验服务
	if passcode == "" {
		return errors.New("验证码不能为空")
	}
	return nil
}

// Initiate 初始化 Email MFA
func (mfa *EmailMfaUtil) Initiate(userId string, issuer string) (*MfaProps, error) {
	return &MfaProps{
		MfaType: mfa.MfaType,
	}, nil
}

// SetupVerify 设置阶段验证 Email 验证码
func (mfa *EmailMfaUtil) SetupVerify(passcode string) error {
	// TODO: 接入实际的邮件验证码校验服务
	if passcode == "" {
		return errors.New("验证码不能为空")
	}
	return nil
}

// Verify 验证 Email 验证码
func (mfa *EmailMfaUtil) Verify(passcode string) error {
	// TODO: 接入实际的邮件验证码校验服务
	if passcode == "" {
		return errors.New("验证码不能为空")
	}
	return nil
}
