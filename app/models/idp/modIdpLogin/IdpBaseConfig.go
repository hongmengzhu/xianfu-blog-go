package modIdpLogin

// IdpBaseConfig 认证源基础配置（从 BaseConfig JSON 解析）
type IdpBaseConfig struct {
	ClientId        string `json:"clientId"`
	ClientSecret    string `json:"clientSecret"`
	ClientId2       string `json:"clientId2"`
	ClientSecret2   string `json:"clientSecret2"`
	AppId           string `json:"appId"`
	HostUrl         string `json:"hostUrl"`
	RedirectUrl     string `json:"redirectUrl"`
	DisableSsl      bool   `json:"disableSsl"`
	TokenURL        string `json:"tokenURL"`
	AuthURL         string `json:"authURL"`
	UserInfoURL     string `json:"userInfoURL"`
	AppCertificate  string `json:"appCertificate"`
	RootCertificate string `json:"rootCertificate"`
}
