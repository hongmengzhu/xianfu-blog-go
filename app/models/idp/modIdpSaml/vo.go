package modRamSaml

// SamlRedirectVo SAML 登录重定向响应
type SamlRedirectVo struct {
	RedirectUrl string `json:"redirectUrl"` // 重定向 URL
	Method      string `json:"method"`      // GET 或 POST
	PostBody    string `json:"postBody"`    // POST 方式时的 HTML body
}
