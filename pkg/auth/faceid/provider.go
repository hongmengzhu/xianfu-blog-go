package faceid

// FaceIdProvider Face ID 提供商接口
type FaceIdProvider interface {
	// Check 对比两张人脸图片（Base64 编码），返回是否匹配
	Check(base64ImageA string, base64ImageB string) (bool, error)
}

// GetFaceIdProvider 根据类型获取 Face ID 提供商
// typ: "Local UniFace" | "Aliyun" 等
func GetFaceIdProvider(typ string, clientId string, clientSecret string, endPoint string) FaceIdProvider {
	if typ == "Local UniFace" {
		return NewLocalUniFaceProvider(endPoint, clientSecret)
	}
	return NewAliyunFaceIdProvider(clientId, clientSecret, endPoint)
}
