package mfa

import (
	"errors"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

const (
	// MfaTotpPeriodInSeconds TOTP 周期（秒）
	MfaTotpPeriodInSeconds = 30
)

// TotpMfaUtil TOTP MFA 工具
type TotpMfaUtil struct {
	*MfaProps
	period     uint
	secretSize uint
	digits     otp.Digits
}

// NewTotpMfaUtil 创建 TOTP MFA 工具实例
func NewTotpMfaUtil(config *MfaProps) *TotpMfaUtil {
	if config == nil {
		config = &MfaProps{
			MfaType: TotpType,
		}
	}
	return &TotpMfaUtil{
		MfaProps:   config,
		period:     MfaTotpPeriodInSeconds,
		secretSize: 20,
		digits:     otp.DigitsSix,
	}
}

// Initiate 初始化 TOTP（生成密钥和 URL）
func (mfa *TotpMfaUtil) Initiate(userId string, issuer string) (*MfaProps, error) {
	if issuer == "" {
		issuer = "XianfuBlog"
	}

	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: userId,
		Period:      mfa.period,
		SecretSize:  mfa.secretSize,
		Digits:      mfa.digits,
	})
	if err != nil {
		return nil, err
	}

	mfaProps := MfaProps{
		MfaType: mfa.MfaType,
		Secret:  key.Secret(),
		URL:     key.URL(),
	}
	return &mfaProps, nil
}

// SetupVerify 设置阶段验证 TOTP 验证码
func (mfa *TotpMfaUtil) SetupVerify(passcode string) error {
	result, err := totp.ValidateCustom(passcode, mfa.Secret, time.Now().UTC(), totp.ValidateOpts{
		Period:    MfaTotpPeriodInSeconds,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	if err != nil {
		return err
	}
	if result {
		return nil
	}
	return errors.New("TOTP 验证码错误")
}

// Verify 验证 TOTP 验证码
func (mfa *TotpMfaUtil) Verify(passcode string) error {
	result, err := totp.ValidateCustom(passcode, mfa.Secret, time.Now().UTC(), totp.ValidateOpts{
		Period:    MfaTotpPeriodInSeconds,
		Skew:      1,
		Digits:    otp.DigitsSix,
		Algorithm: otp.AlgorithmSHA1,
	})
	if err != nil {
		return err
	}
	if result {
		return nil
	}
	return errors.New("TOTP 验证码错误")
}
