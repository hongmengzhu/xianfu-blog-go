package mfa

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// MfaProps MFA 属性
type MfaProps struct {
	Enabled       bool     `json:"enabled"`
	IsPreferred   bool     `json:"isPreferred"`
	MfaType       string   `json:"mfaType"`
	Secret        string   `json:"secret,omitempty"`
	CountryCode   string   `json:"countryCode,omitempty"`
	URL           string   `json:"url,omitempty"`
	RecoveryCodes []string `json:"recoveryCodes,omitempty"`
}

// MfaInterface MFA 接口
type MfaInterface interface {
	// Initiate 初始化 MFA（生成 secret、URL 等）
	Initiate(userId string, issuer string) (*MfaProps, error)
	// SetupVerify 设置阶段验证（确认用户已正确配置）
	SetupVerify(passcode string) error
	// Verify 验证 MFA 验证码
	Verify(passcode string) error
}

// MFA 类型常量
const (
	SmsType  = "sms"
	EmailType = "email"
	TotpType = "app"
)

// GetMfaUtil 根据 MFA 类型获取对应的工具实例
func GetMfaUtil(mfaType string, config *MfaProps) MfaInterface {
	switch mfaType {
	case SmsType:
		return NewSmsMfaUtil(config)
	case EmailType:
		return NewEmailMfaUtil(config)
	case TotpType:
		return NewTotpMfaUtil(config)
	}
	return nil
}

// GenerateRecoveryCodes 生成恢复码
func GenerateRecoveryCodes(count int) []string {
	codes := make([]string, 0, count)
	for i := 0; i < count; i++ {
		b := make([]byte, 8)
		_, _ = rand.Read(b)
		codes = append(codes, hex.EncodeToString(b))
	}
	return codes
}

// VerifyRecoveryCode 验证恢复码是否匹配，匹配成功则返回删除该码后的剩余码列表
func VerifyRecoveryCode(recoveryCodes []string, code string) ([]string, error) {
	if len(recoveryCodes) == 0 {
		return nil, fmt.Errorf("没有可用的恢复码")
	}
	for i, c := range recoveryCodes {
		if c == code {
			// 删除已使用的恢复码
			remaining := make([]string, 0, len(recoveryCodes)-1)
			remaining = append(remaining, recoveryCodes[:i]...)
			remaining = append(remaining, recoveryCodes[i+1:]...)
			return remaining, nil
		}
	}
	return nil, fmt.Errorf("恢复码不正确")
}
