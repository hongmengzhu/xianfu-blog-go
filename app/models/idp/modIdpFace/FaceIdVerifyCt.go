package modIdpFace

// FaceIdVerifyCt Face ID 验证请求
type FaceIdVerifyCt struct {
	SourceNo string `json:"sourceNo" form:"sourceNo" label:"认证源编号"`
	UserId   string `json:"userId" form:"userId" label:"用户ID"`
	ImageA   string `json:"imageA" form:"imageA" label:"人脸图片A(Base64)"`
	ImageB   string `json:"imageB" form:"imageB" label:"人脸图片B(Base64)"`
}
