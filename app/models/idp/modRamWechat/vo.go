package modRamWechat

// QRCodeVo 公众号扫码登录二维码响应
type QRCodeVo struct {
	Ticket        string `json:"ticket"`
	QRUrl         string `json:"qrUrl"`
	ExpireSeconds int    `json:"expireSeconds"`
}

// PollVo 轮询扫码状态响应
type PollVo struct {
	IsScanned     bool   `json:"isScanned"`
	WechatUnionId string `json:"wechatUnionId,omitempty"`
}
