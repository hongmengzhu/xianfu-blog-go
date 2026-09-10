package saml

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net/url"

	saml2 "github.com/russellhaering/gosaml2"
	dsig "github.com/russellhaering/goxmldsig"
)

// SAMLServiceProvider 导出 gosaml2 的 SAMLServiceProvider 类型
type SAMLServiceProvider = saml2.SAMLServiceProvider

// SpConfig SAML Service Provider 配置
type SpConfig struct {
	// SP 信息
	Issuer      string // SP EntityID / Issuer
	AcsURL      string // Assertion Consumer Service URL
	IdpSsoURL   string // IdP SSO 登录地址
	IdpIssuer   string // IdP EntityID / Issuer
	EnableSign  bool   // 是否签名 AuthnRequest

	// IdP 证书（PEM 或 Base64 编码的 DER）
	IdpCert string

	// SP 证书和私钥（PEM 格式），EnableSign=true 时需要
	SpCert string // SP 证书 PEM
	SpKey  string // SP 私钥 PEM
}

// SpUserInfo SAML Response 解析后的用户信息
type SpUserInfo struct {
	NameID      string            // SAML NameID
	Attributes  map[string]string // SAML 属性映射
}

// BuildServiceProvider 构建 SAML Service Provider
func BuildServiceProvider(cfg *SpConfig) (*saml2.SAMLServiceProvider, error) {
	certStore, err := buildSpCertificateStore(cfg)
	if err != nil {
		return nil, fmt.Errorf("构建 IdP 证书存储失败: %w", err)
	}

	sp := &saml2.SAMLServiceProvider{
		ServiceProviderIssuer:       cfg.Issuer,
		AssertionConsumerServiceURL: cfg.AcsURL,
		SignAuthnRequests:           false,
		IDPCertificateStore:         &certStore,
		SPKeyStore:                  dsig.RandomKeyStoreForTest(),
	}

	if cfg.IdpSsoURL != "" {
		sp.IdentityProviderSSOURL = cfg.IdpSsoURL
	}
	if cfg.IdpIssuer != "" {
		sp.IdentityProviderIssuer = cfg.IdpIssuer
	}

	if cfg.EnableSign {
		sp.SignAuthnRequests = true
		keyStore, err := buildSpKeyStore(cfg)
		if err != nil {
			return nil, fmt.Errorf("构建 SP 密钥存储失败: %w", err)
		}
		sp.SPKeyStore = keyStore
	}

	return sp, nil
}

// GenerateSamlRequest 生成 SAML AuthnRequest URL
// 返回 (authURL, method, error)
// 如果 EnableSign=true 返回 POST body 和 method="POST"
// 否则返回 redirect URL 和 method="GET"
func GenerateSamlRequest(sp *saml2.SAMLServiceProvider, relayState string) (auth string, method string, err error) {
	if sp.SignAuthnRequests {
		post, err := sp.BuildAuthBodyPost(relayState)
		if err != nil {
			return "", "", fmt.Errorf("构建 SAML POST body 失败: %w", err)
		}
		return string(post), "POST", nil
	}

	authURL, err := sp.BuildAuthURL(relayState)
	if err != nil {
		return "", "", fmt.Errorf("构建 SAML redirect URL 失败: %w", err)
	}
	return authURL, "GET", nil
}

// ParseSamlResponse 解析 SAML Response，提取用户信息
func ParseSamlResponse(samlResponse string, sp *saml2.SAMLServiceProvider) (*SpUserInfo, error) {
	samlResponse, _ = url.QueryUnescape(samlResponse)

	assertionInfo, err := sp.RetrieveAssertionInfo(samlResponse)
	if err != nil {
		return nil, fmt.Errorf("解析 SAML Assertion 失败: %w", err)
	}

	userInfo := &SpUserInfo{
		NameID:     assertionInfo.NameID,
		Attributes: make(map[string]string),
	}

	for _, attr := range assertionInfo.Values {
		if len(attr.Values) > 0 {
			userInfo.Attributes[attr.Name] = attr.Values[0].Value
		}
	}

	return userInfo, nil
}

// buildSpCertificateStore 构建 IdP 证书存储
func buildSpCertificateStore(cfg *SpConfig) (dsig.MemoryX509CertificateStore, error) {
	certEncodedData := cfg.IdpCert
	if certEncodedData == "" {
		return dsig.MemoryX509CertificateStore{}, fmt.Errorf("IdP 证书为空")
	}

	var certData []byte
	block, _ := pem.Decode([]byte(certEncodedData))
	if block != nil {
		// PEM 格式
		certData = block.Bytes
	} else {
		// 尝试 Base64 解码
		var err error
		certData, err = base64.StdEncoding.DecodeString(certEncodedData)
		if err != nil {
			return dsig.MemoryX509CertificateStore{}, fmt.Errorf("解析 IdP 证书失败: %w", err)
		}
	}

	idpCert, err := x509.ParseCertificate(certData)
	if err != nil {
		return dsig.MemoryX509CertificateStore{}, fmt.Errorf("解析 IdP X509 证书失败: %w", err)
	}

	return dsig.MemoryX509CertificateStore{
		Roots: []*x509.Certificate{idpCert},
	}, nil
}

// buildSpKeyStore 构建 SP 密钥存储（用于签名 AuthnRequest）
func buildSpKeyStore(cfg *SpConfig) (dsig.X509KeyStore, error) {
	if cfg.SpCert == "" || cfg.SpKey == "" {
		return nil, fmt.Errorf("SP 证书或私钥为空")
	}

	keyPair, err := tls.X509KeyPair([]byte(cfg.SpCert), []byte(cfg.SpKey))
	if err != nil {
		return nil, fmt.Errorf("加载 SP 证书/私钥对失败: %w", err)
	}

	return &dsig.TLSCertKeyStore{
		PrivateKey:  keyPair.PrivateKey,
		Certificate: keyPair.Certificate,
	}, nil
}
