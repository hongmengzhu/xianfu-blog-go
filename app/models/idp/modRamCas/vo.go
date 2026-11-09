package modRamCas

// CasLoginVo CAS 登录响应（返回 ST）
type CasLoginVo struct {
	ServiceTicket string `json:"serviceTicket"` // CAS Service Ticket
	Service       string `json:"service"`       // 目标服务地址
	RedirectUrl   string `json:"redirectUrl"`   // 完整重定向 URL
}

// CasValidateVo CAS 验证响应
type CasValidateVo struct {
	Valid   bool   `json:"valid"`
	User    string `json:"user,omitempty"`
	Message string `json:"message,omitempty"`
}
