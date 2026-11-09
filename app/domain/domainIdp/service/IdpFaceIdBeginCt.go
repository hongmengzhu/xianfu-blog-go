package service

// IdpFaceIdBeginCt Face ID 开始验证请求
type IdpFaceIdBeginCt struct {
	SourceNo string `json:"sourceNo" form:"sourceNo" label:"认证源编号"`
	UserId   string `json:"userId" form:"userId" label:"用户ID"`
}
