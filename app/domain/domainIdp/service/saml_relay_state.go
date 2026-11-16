package service

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/pangu-2/go-tools/tools/strPg"
)

// relayStatePayload SAML RelayState 的签名载荷（不直接暴露给客户端）。
type relayStatePayload struct {
	SourceNo  string `json:"sourceNo"`
	Nonce     string `json:"nonce"`
	ExpiresAt int64  `json:"expiresAt"`
}

const relayStateTTL = 10 * time.Minute

// 已使用 nonce，保证 RelayState 只能消费一次；过期项在写入时顺带清理。
var relayStateNonces sync.Map

func (c *SamlLoginService) relayStateKey() []byte {
	secret := c.pg.Jwt.System.Secret
	if strPg.IsBlank(secret) {
		// 配置缺失时使用固定域标识，避免空密钥；生产应通过配置覆盖。
		secret = "xianfu-saml-relay-state"
	}
	return []byte("saml-relay-state::" + secret)
}

func (c *SamlLoginService) makeRelayState(sourceNo string) (string, error) {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	payload := relayStatePayload{
		SourceNo:  sourceNo,
		Nonce:     hex.EncodeToString(buf),
		ExpiresAt: time.Now().Add(relayStateTTL).Unix(),
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	body := base64.RawURLEncoding.EncodeToString(raw)
	sig := hmac.New(sha256.New, c.relayStateKey())
	sig.Write([]byte(body))
	return body + "." + hex.EncodeToString(sig.Sum(nil)), nil
}

func (c *SamlLoginService) verifyRelayState(relayState string) (string, error) {
	if strPg.IsBlank(relayState) {
		return "", fmt.Errorf("RelayState 为空")
	}
	var body, signature string
	{
		idx := -1
		for i := 0; i < len(relayState); i++ {
			if relayState[i] == '.' {
				idx = i
				break
			}
		}
		if idx <= 0 || idx >= len(relayState)-1 {
			return "", fmt.Errorf("RelayState 格式错误")
		}
		body, signature = relayState[:idx], relayState[idx+1:]
	}
	expect := hmac.New(sha256.New, c.relayStateKey())
	expect.Write([]byte(body))
	if !hmac.Equal([]byte(signature), []byte(hex.EncodeToString(expect.Sum(nil)))) {
		return "", fmt.Errorf("RelayState 签名无效")
	}
	raw, err := base64.RawURLEncoding.DecodeString(body)
	if err != nil {
		return "", fmt.Errorf("RelayState 解码失败")
	}
	var payload relayStatePayload
	if err := json.Unmarshal(raw, &payload); err != nil {
		return "", fmt.Errorf("RelayState 解析失败")
	}
	if strPg.IsBlank(payload.SourceNo) || strPg.IsBlank(payload.Nonce) {
		return "", fmt.Errorf("RelayState 内容无效")
	}
	if time.Now().Unix() > payload.ExpiresAt {
		return "", fmt.Errorf("RelayState 已过期")
	}
	if _, loaded := relayStateNonces.LoadOrStore(payload.Nonce, payload.ExpiresAt); loaded {
		return "", fmt.Errorf("RelayState 已使用")
	}
	// 顺带清理过期 nonce，避免无限增长。
	relayStateNonces.Range(func(key, value any) bool {
		if ts, ok := value.(int64); ok && time.Now().Unix() > ts {
			relayStateNonces.Delete(key)
		}
		return true
	})
	return payload.SourceNo, nil
}
