package deviceauth

import "time"

// 状态常量
const (
	StatusPending     = "pending"
	StatusApproved    = "approved"
	StatusDenied      = "denied"
	StatusTokenIssued = "token_issued"

	DefaultExpiresIn = 120 // 默认过期时间（秒）
	DefaultInterval  = 5   // 默认轮询间隔（秒）
)

// DeviceAuthCache 设备授权缓存条目
type DeviceAuthCache struct {
	DeviceCode  string    // 设备码
	UserCode    string    // 用户码
	ClientId    string    // 客户端 ID
	UserName    string    // 授权用户账号（审批后填充）
	UserNo      string    // 授权用户 No（审批后填充）
	TenantNo    string    // 租户编号
	Scope       string    // 请求范围
	RequestAt   time.Time // 请求时间
	Status      string    // 状态：pending / approved / denied / token_issued
	CancelToken string    // 取消令牌
	ExpiresIn   int       // 过期时间（秒）
}

// IsExpired 判断是否已过期
func (d *DeviceAuthCache) IsExpired() bool {
	expiresIn := d.ExpiresIn
	if expiresIn <= 0 {
		expiresIn = DefaultExpiresIn
	}
	return d.RequestAt.Add(time.Duration(expiresIn) * time.Second).Before(time.Now())
}

// DeviceAuthResponse 设备授权响应（RFC 8628）
type DeviceAuthResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationUri string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
	CancelToken     string `json:"cancel_token,omitempty"`
}
