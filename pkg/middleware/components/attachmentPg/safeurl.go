package attachmentPg

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// SafeFetchClient 用于服务端抓取远程附件，限制超时、跳转次数并校验每一跳目标。
var SafeFetchClient = &http.Client{
	Timeout: 10 * time.Second,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		if len(via) >= 3 {
			return fmt.Errorf("跳转次数过多")
		}
		if err := ValidateFetchURL(req.Context(), req.URL.String()); nil != err {
			return err
		}
		return nil
	},
}

// ValidateFetchURL 仅允许 http/https 且目标解析后的所有 IP 均不是私网/环回/链路本地地址。
func ValidateFetchURL(ctx context.Context, raw string) error {
	u, err := url.Parse(strings.TrimSpace(raw))
	if nil != err {
		return fmt.Errorf("URL 格式错误")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("仅允许 http/https 协议")
	}
	host := u.Hostname()
	if strings.TrimSpace(host) == "" {
		return fmt.Errorf("URL 缺少主机名")
	}
	if ip := net.ParseIP(host); nil != ip {
		if !isPublicIP(ip) {
			return fmt.Errorf("不允许访问内网地址")
		}
		return nil
	}
	lookupCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	addrs, err := net.DefaultResolver.LookupIPAddr(lookupCtx, host)
	if nil != err || len(addrs) == 0 {
		return fmt.Errorf("域名解析失败")
	}
	for _, addr := range addrs {
		if !isPublicIP(addr.IP) {
			return fmt.Errorf("不允许访问内网地址")
		}
	}
	return nil
}

func isPublicIP(ip net.IP) bool {
	if ip == nil || ip.IsLoopback() || ip.IsUnspecified() || ip.IsMulticast() ||
		ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsPrivate() {
		return false
	}
	// 云元数据/特殊网段
	if ip4 := ip.To4(); nil != ip4 {
		if ip4[0] == 169 && ip4[1] == 254 {
			return false
		}
		if ip4[0] == 100 && ip4[1] >= 64 && ip4[1] <= 127 {
			return false // CGNAT
		}
	}
	return true
}
