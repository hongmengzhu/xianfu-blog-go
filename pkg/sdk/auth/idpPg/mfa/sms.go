package mfa

// SmsMfaUtil SMS MFA 工具
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

// Initiate 初始化 SMS MFA。
//
// 短信验证码需要在 Initiate 阶段生成并调用网关发送，当前项目没有接入任何短信通道，
// 因此这里直接返回 ErrChannelMfaUnavailable：既不能凭空生成一个用户永远收不到的 Secret，
// 也不能让 SetupVerify/Verify 变成"验证码非空即通过"的摆设。
func (mfa *SmsMfaUtil) Initiate(userId string, issuer string) (*MfaProps, error) {
	return nil, ErrChannelMfaUnavailable
}

// SetupVerify 设置阶段验证 SMS 验证码
func (mfa *SmsMfaUtil) SetupVerify(passcode string) error {
	return ErrChannelMfaUnavailable
}

// Verify 验证 SMS 验证码
func (mfa *SmsMfaUtil) Verify(passcode string) error {
	return ErrChannelMfaUnavailable
}

// Initiate 初始化 Email MFA。
//
// 同 SMS：没有可用的邮件发送通道（configPg.Pg.Mail 目前无人使用），
// 无法把一个真实的验证码送达用户，因此拒绝启用而不是放行任意验证码。
func (mfa *EmailMfaUtil) Initiate(userId string, issuer string) (*MfaProps, error) {
	return nil, ErrChannelMfaUnavailable
}

// SetupVerify 设置阶段验证 Email 验证码
func (mfa *EmailMfaUtil) SetupVerify(passcode string) error {
	return ErrChannelMfaUnavailable
}

// Verify 验证 Email 验证码
func (mfa *EmailMfaUtil) Verify(passcode string) error {
	return ErrChannelMfaUnavailable
}
