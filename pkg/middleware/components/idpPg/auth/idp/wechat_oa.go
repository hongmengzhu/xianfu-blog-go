package idp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// ==================== 公众号 access_token 缓存 ====================

type wechatOaTokenCache struct {
	AccessToken string
	ExpireAt    time.Time
}

var (
	oaTokenCache     sync.Map // key: clientId → *wechatOaTokenCache
	oaTokenCacheLock sync.Mutex
)

// WechatOaGetAccessToken 获取公众号 access_token（带缓存，有效期2小时，提前5分钟刷新）
func WechatOaGetAccessToken(clientId string, clientSecret string) (string, error) {
	// 尝试从缓存读取
	if val, ok := oaTokenCache.Load(clientId); ok {
		cache := val.(*wechatOaTokenCache)
		if time.Now().Before(cache.ExpireAt) {
			return cache.AccessToken, nil
		}
	}

	// 缓存不存在或已过期，重新获取
	oaTokenCacheLock.Lock()
	defer oaTokenCacheLock.Unlock()

	// double-check
	if val, ok := oaTokenCache.Load(clientId); ok {
		cache := val.(*wechatOaTokenCache)
		if time.Now().Before(cache.ExpireAt) {
			return cache.AccessToken, nil
		}
	}

	accessToken, errMsg, err := GetWechatOfficialAccountAccessToken(clientId, clientSecret)
	if err != nil {
		return "", fmt.Errorf("获取公众号 access_token 失败: %w", err)
	}
	if errMsg != "" {
		return "", fmt.Errorf("获取公众号 access_token 失败: %s", errMsg)
	}

	// 缓存 access_token，有效期设为 1小时55分（提前5分钟过期）
	oaTokenCache.Store(clientId, &wechatOaTokenCache{
		AccessToken: accessToken,
		ExpireAt:    time.Now().Add(115 * time.Minute),
	})

	return accessToken, nil
}

// ==================== 公众号二维码 ====================

// WechatOaQRCodeResult 公众号二维码创建结果
type WechatOaQRCodeResult struct {
	Ticket        string `json:"ticket"`
	ExpireSeconds int    `json:"expire_seconds"`
	URL           string `json:"url"`
}

// WechatOaCreateQRCode 创建公众号带参数临时二维码
// sceneStr: 场景值字符串（如 sourceNo + randomId）
func WechatOaCreateQRCode(clientId string, clientSecret string, sceneStr string) (*WechatOaQRCodeResult, error) {
	accessToken, err := WechatOaGetAccessToken(clientId, clientSecret)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	createURL := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token=%s", accessToken)
	params := fmt.Sprintf(`{"expire_seconds": 3600, "action_name": "QR_STR_SCENE", "action_info": {"scene": {"scene_str": "%s"}}}`, sceneStr)

	req, err := http.NewRequest("POST", createURL, bytes.NewReader([]byte(params)))
	if err != nil {
		return nil, fmt.Errorf("创建二维码请求失败: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求微信二维码接口失败: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result WechatOaQRCodeResult
	if err = json.Unmarshal(respBytes, &result); err != nil {
		return nil, fmt.Errorf("解析二维码响应失败: %w", err)
	}

	if result.Ticket == "" {
		return nil, fmt.Errorf("创建二维码失败，响应: %s", string(respBytes))
	}

	return &result, nil
}

// WechatOaQRCodeShowURL 通过 ticket 换取二维码展示 URL
func WechatOaQRCodeShowURL(ticket string) string {
	return fmt.Sprintf("https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=%s", ticket)
}
