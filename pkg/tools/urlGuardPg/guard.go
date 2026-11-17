package urlGuardPg

// 出站地址安全校验：仅允许 http/https 且目标必须是公网地址，
// 用于所有由管理员配置或用户传入的服务端 HTTP 请求，防止 SSRF。

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const dnsTimeout = 3 * time.Second

// ValidateFetchURL 校验目标 URL 是否允许服务端发起请求
func ValidateFetchURL(ctx context.Context, raw string) error {
	raw = strings.TrimSpace(raw)
	if "" == raw {
		return fmt.Errorf("地址不能为空")
	}
	u, err := url.Parse(raw)
	if nil != err {
		return fmt.Errorf("地址格式错误")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("仅允许 http/https 地址")
	}
	host := u.Hostname()
	if "" == host {
		return fmt.Errorf("地址缺少域名")
	}
	if strings.EqualFold(host, "localhost") {
		return fmt.Errorf("不允许访问本地地址")
	}
	if ip := net.ParseIP(host); nil != ip {
		if !IsPublicIP(ip) {
			return fmt.Errorf("不允许访问内网地址")
		}
		return nil
	}
	if nil == ctx {
		ctx = context.Background()
	}
	dialCtx, cancel := context.WithTimeout(ctx, dnsTimeout)
	defer cancel()
	addrs, err := net.DefaultResolver.LookupIPAddr(dialCtx, host)
	if nil != err || len(addrs) < 1 {
		return fmt.Errorf("域名解析失败")
	}
	for _, addr := range addrs {
		if !IsPublicIP(addr.IP) {
			return fmt.Errorf("不允许访问内网地址")
		}
	}
	return nil
}

// IsPublicIP 判断是否为公网可访问地址
func IsPublicIP(ip net.IP) bool {
	if nil == ip {
		return false
	}
	if ip.IsLoopback() || ip.IsUnspecified() || ip.IsMulticast() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsInterfaceLocalMulticast() {
		return false
	}
	if ip.IsPrivate() {
		return false
	}
	// 云元数据与 CGNAT 段
	if ip4 := ip.To4(); nil != ip4 {
		if ip4[0] == 169 && ip4[1] == 254 {
			return false
		}
		if ip4[0] == 100 && ip4[1] >= 64 && ip4[1] <= 127 {
			return false
		}
	}
	return true
}

// NewSafeClient 创建带跳转校验的 HTTP 客户端
func NewSafeClient(timeout time.Duration, maxRedirects int) *http.Client {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	if maxRedirects <= 0 {
		maxRedirects = 3
	}
	return &http.Client{
		Timeout: timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= maxRedirects {
				return fmt.Errorf("跳转次数过多")
			}
			if err := ValidateFetchURL(req.Context(), req.URL.String()); nil != err {
				return err
			}
			return nil
		},
	}
}
