package idp_e2e

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net/http"
	"strings"
	"testing"
	"time"
)

// makeSelfSignedCertPEM 生成用于 SAML SP 构建的 IdP 证书（PEM）。
func makeSelfSignedCertPEM(t *testing.T) string {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("生成 RSA 密钥失败: %v", err)
	}
	tmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "E2E IdP"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatalf("创建证书失败: %v", err)
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}))
}

func runSAML(t *testing.T, e *env) {
	e.t = t
	ctrl := authIdpURL + "/saml"

	pid, pno := e.makeIdpProvider(e.unique("provider-saml"), "saml2")
	t.Cleanup(func() { _ = pid })
	idpURL := "http://mock-idp.example.com/sso"

	// 带证书的 SAML 认证源
	sid := e.createID(baseURL+"/identity-source", js(map[string]any{
		"name": e.unique("source-saml"), "protocol": "saml2", "idp": pno,
		"baseConfig": js(map[string]any{
			"idpSsoUrl": idpURL, "idpIssuer": "urn:e2e:idp",
			"spIssuer": "https://sp.example.com/saml/acs",
		}),
	}), "ram_identity_source")
	sno := e.sourceNo(sid)
	cert := makeSelfSignedCertPEM(t)
	e.createID(baseURL+"/idp-credential", js(map[string]any{
		"idp": pno, "sourceNo": sno, "credType": "sp_cert", "value": cert,
		"description": "e2e saml cert",
	}), "ram_idp_credential")

	t.Run("login-url-redirect", func(t *testing.T) {
		e.t = t
		r := e.get(ctrl + "/login/" + sno)
		if r.Status != http.StatusFound {
			t.Fatalf("SAML 登录期望 302，实际 status=%d body=%s", r.Status, r.Body)
		}
		loc := r.Header.Get("Location")
		if !strings.HasPrefix(loc, idpURL) || !strings.Contains(loc, "SAMLRequest=") {
			t.Fatalf("重定向地址异常: %s", loc)
		}
	})
	t.Run("login-unknown-source", func(t *testing.T) {
		e.t = t
		r := e.get(ctrl + "/login/not-exist-" + e.runID)
		if r.Code == "200" || !strings.Contains(r.Message, "认证源不存在") {
			t.Fatalf("未知认证源期望失败，实际 code=%s msg=%s", r.Code, r.Message)
		}
	})
	t.Run("login-wrong-protocol", func(t *testing.T) {
		e.t = t
		_, pno2 := e.makeIdpProvider(e.unique("provider-oauth-for-saml"), "oauth2")
		oauthSid := e.createID(baseURL+"/identity-source", js(map[string]any{
			"name": e.unique("source-oauth-for-saml"), "protocol": "oauth2", "idp": pno2,
		}), "ram_identity_source")
		r := e.get(ctrl + "/login/" + e.sourceNo(oauthSid))
		if r.Code == "200" || !strings.Contains(r.Message, "不是 SAML2 协议") {
			t.Fatalf("非 SAML2 协议期望失败，实际 code=%s msg=%s", r.Code, r.Message)
		}
	})
	t.Run("login-source-disabled", func(t *testing.T) {
		e.t = t
		e.ok(baseURL+"/identity-source/disable", js(map[string]any{"ids": []string{sid}}))
		r := e.get(ctrl + "/login/" + sno)
		if r.Code == "200" || !strings.Contains(r.Message, "认证源已禁用") {
			t.Fatalf("停用认证源期望失败，实际 code=%s msg=%s", r.Code, r.Message)
		}
		c := e.postForm(ctrl+"/callback", map[string]string{
			"SAMLResponse": "<x/>", "RelayState": js(map[string]any{"sourceNo": sno}),
		})
		if c.Code == "200" || !strings.Contains(c.Message, "认证源已禁用") {
			t.Fatalf("停用认证源回调期望失败，实际 code=%s msg=%s", c.Code, c.Message)
		}
		e.ok(baseURL+"/identity-source/enable", js(map[string]any{"ids": []string{sid}}))
	})
	t.Run("login-without-cert", func(t *testing.T) {
		e.t = t
		noCertSid := e.createID(baseURL+"/identity-source", js(map[string]any{
			"name": e.unique("source-saml-nocert"), "protocol": "saml2", "idp": pno,
			"baseConfig": js(map[string]any{"idpSsoUrl": idpURL, "idpIssuer": "urn:e2e:idp"}),
		}), "ram_identity_source")
		r := e.get(ctrl + "/login/" + e.sourceNo(noCertSid))
		if r.Code == "200" || !strings.Contains(r.Message, "构建 SAML SP 失败") {
			t.Fatalf("无证书期望构建失败，实际 code=%s msg=%s", r.Code, r.Message)
		}
	})
	t.Run("callback-missing-response", func(t *testing.T) {
		e.t = t
		r := e.postForm(ctrl+"/callback", map[string]string{})
		if r.Code == "200" || !strings.Contains(r.Message, "SAML Response 不能为空") {
			t.Fatalf("缺 SAMLResponse 期望失败，实际 code=%s msg=%s", r.Code, r.Message)
		}
	})
	t.Run("callback-bad-relay-state", func(t *testing.T) {
		e.t = t
		r := e.postForm(ctrl+"/callback", map[string]string{"SAMLResponse": "<x/>", "RelayState": "not-json"})
		if r.Code == "200" || !strings.Contains(r.Message, "RelayState 解析失败") {
			t.Fatalf("坏 RelayState 期望失败，实际 code=%s msg=%s", r.Code, r.Message)
		}
	})
	t.Run("callback-garbage-assertion", func(t *testing.T) {
		e.t = t
		r := e.postForm(ctrl+"/callback", map[string]string{
			"SAMLResponse": "not-a-saml-response", "RelayState": js(map[string]any{"sourceNo": sno}),
		})
		if r.Code == "200" || !strings.Contains(r.Message, "解析 SAML 响应失败") {
			t.Fatalf("非法断言期望失败，实际 code=%s msg=%s", r.Code, r.Message)
		}
	})
}
