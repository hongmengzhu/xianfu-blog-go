package mfa

import (
	"errors"
	"testing"
)

// 短信/邮件 MFA 在没有接入发送通道之前必须 fail closed：
// 任意非空验证码都不能被当成校验通过。
func TestChannelMfaFailsClosed(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		util MfaInterface
	}{
		{name: "sms", util: NewSmsMfaUtil(nil)},
		{name: "email", util: NewEmailMfaUtil(nil)},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if _, err := tc.util.Initiate("account-1", "XianfuBlog"); !errors.Is(err, ErrChannelMfaUnavailable) {
				t.Fatalf("Initiate err = %v, want ErrChannelMfaUnavailable", err)
			}
			if err := tc.util.SetupVerify("123456"); !errors.Is(err, ErrChannelMfaUnavailable) {
				t.Fatalf("SetupVerify err = %v, want ErrChannelMfaUnavailable", err)
			}
			if err := tc.util.Verify("123456"); !errors.Is(err, ErrChannelMfaUnavailable) {
				t.Fatalf("Verify err = %v, want ErrChannelMfaUnavailable", err)
			}
			// 空验证码同样必须被拒绝
			if err := tc.util.Verify(""); !errors.Is(err, ErrChannelMfaUnavailable) {
				t.Fatalf("Verify(empty) err = %v, want ErrChannelMfaUnavailable", err)
			}
		})
	}
}
