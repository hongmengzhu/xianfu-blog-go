package modRamCas

// CasValidateCt CAS 验证请求
type CasValidateCt struct {
	Ticket  string `json:"ticket" form:"ticket" label:"票据"`
	Service string `json:"service" form:"service" label:"服务地址"`
	PgtUrl  string `json:"pgtUrl" form:"pgtUrl" label:"代理回调地址"`
	Format  string `json:"format" form:"format" label:"响应格式(xml|json)"`
}

// CasLoginCt CAS 登录请求（生成 ST）
type CasLoginCt struct {
	SourceNo string `json:"sourceNo" form:"sourceNo" label:"认证源编号"`
	Service  string `json:"service" form:"service" label:"服务地址" binding:"required"`
}
