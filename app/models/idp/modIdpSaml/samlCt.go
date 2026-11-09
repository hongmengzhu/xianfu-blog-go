package modRamSaml

// SamlLoginCt SAML 登录请求参数
type SamlLoginCt struct {
	SourceNo string `json:"sourceNo" form:"sourceNo" label:"认证源编号"` // 认证源编号
}

// SamlCallbackCt SAML Response 回调请求参数
type SamlCallbackCt struct {
	SAMLResponse string `json:"SAMLResponse" form:"SAMLResponse"` // SAML Response（Base64 编码）
	RelayState   string `json:"RelayState" form:"RelayState"`     // RelayState（包含 sourceNo 等信息）
}
