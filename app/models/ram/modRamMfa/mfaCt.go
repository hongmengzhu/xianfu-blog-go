package modRamMfa

// SetupCt MFA 设置请求
type SetupCt struct {
	MfaType string `json:"mfaType" form:"mfaType" label:"MFA类型" binding:"required"` // app/sms/email
}

// SetupVerifyCt MFA 设置验证请求
type SetupVerifyCt struct {
	MfaType     string `json:"mfaType" form:"mfaType" label:"MFA类型" binding:"required"`
	Secret      string `json:"secret" form:"secret" label:"密钥"`
	Passcode    string `json:"passcode" form:"passcode" label:"验证码" binding:"required"`
	Dest        string `json:"dest" form:"dest" label:"目标地址"`
	CountryCode string `json:"countryCode" form:"countryCode" label:"国家码"`
}

// EnableCt MFA 启用请求
type EnableCt struct {
	MfaType       string   `json:"mfaType" form:"mfaType" label:"MFA类型" binding:"required"`
	Secret        string   `json:"secret" form:"secret" label:"密钥"`
	RecoveryCodes []string `json:"recoveryCodes" form:"recoveryCodes" label:"恢复码"`
	Dest          string   `json:"dest" form:"dest" label:"目标地址"`
	CountryCode   string   `json:"countryCode" form:"countryCode" label:"国家码"`
}

// VerifyCt MFA 验证请求（登录时）
type VerifyCt struct {
	MfaType  string `json:"mfaType" form:"mfaType" label:"MFA类型" binding:"required"`
	Passcode string `json:"passcode" form:"passcode" label:"验证码"`
	MfaToken string `json:"mfaToken" form:"mfaToken" label:"MFA临时令牌"`
}

// RecoverCt MFA 恢复请求
type RecoverCt struct {
	MfaToken     string `json:"mfaToken" form:"mfaToken" label:"MFA临时令牌" binding:"required"`
	RecoveryCode string `json:"recoveryCode" form:"recoveryCode" label:"恢复码" binding:"required"`
}

// DisableCt MFA 禁用请求
type DisableCt struct {
	Passcode string `json:"passcode" form:"passcode" label:"验证码"`
}

// MfaRequiredResponse MFA 验证_required_响应
type MfaRequiredResponse struct {
	MfaRequired bool   `json:"mfaRequired"`
	MfaToken    string `json:"mfaToken,omitempty"` // 临时令牌，用于后续 MFA 验证
	MfaType     string `json:"mfaType,omitempty"`  // 首选 MFA 类型
}
