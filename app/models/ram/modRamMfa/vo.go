package modRamMfa

// MfaPropsVo MFA 属性视图
type MfaPropsVo struct {
	Enabled       bool     `json:"enabled"`
	IsPreferred   bool     `json:"isPreferred"`
	MfaType       string   `json:"mfaType"`
	Secret        string   `json:"secret,omitempty"`
	URL           string   `json:"url,omitempty"`
	RecoveryCodes []string `json:"recoveryCodes,omitempty"`
}

// MfaSetupVo MFA 设置响应
type MfaSetupVo struct {
	MfaType       string   `json:"mfaType"`
	Secret        string   `json:"secret"`
	URL           string   `json:"url"`
	RecoveryCodes []string `json:"recoveryCodes"`
}

// MfaStatusVo MFA 状态响应
type MfaStatusVo struct {
	Enabled       bool         `json:"enabled"`
	PreferredType string       `json:"preferredType"`
	Methods       []MfaPropsVo `json:"methods"`
}
