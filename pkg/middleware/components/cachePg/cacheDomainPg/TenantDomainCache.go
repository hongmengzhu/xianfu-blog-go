package cacheDomainPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/app/models/tc/dtoTenantDomain"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/configPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/mapPg"
	"go-spring.org/spring/gs"
	"net"

	"strings"
)

func init() {
	gs.Provide(&TenantDomainCache{Domain: &mapPg.SafeMap[string, string]{}, DomainData: &mapPg.SafeMap[string, dtoTenantDomain.Dto]{}})
}

// TenantDomainCache
// @Description: 域名缓存
type TenantDomainCache struct {
	// 域名：域名数据
	DomainData *mapPg.SafeMap[string, dtoTenantDomain.Dto]
	// 域名：租户编号
	Domain *mapPg.SafeMap[string, string]
}

// NormalizeHost 规范化 Host：去协议、去端口、小写、去尾部点
func NormalizeHost(host string) string {
	host = strings.TrimSpace(strings.ToLower(host))
	if i := strings.Index(host, "://"); i >= 0 {
		host = host[i+3:]
	}
	if i := strings.Index(host, "/"); i >= 0 {
		host = host[:i]
	}
	// IPv6 形如 [::1]:8080
	if strings.HasPrefix(host, "[") {
		if i := strings.Index(host, "]"); i >= 0 {
			host = host[1:i]
		}
	} else if i := strings.LastIndex(host, ":"); i >= 0 {
		// 仅当冒号后全是数字时视为端口
		port := host[i+1:]
		if port != "" && isAllDigits(port) {
			host = host[:i]
		}
	}
	host = strings.TrimSuffix(host, ".")
	return host
}

func isAllDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return len(s) > 0
}

// IsLocalHostExist 精确判定的本地调试域名（不再使用子串匹配，避免 Host 伪造）
func (c *TenantDomainCache) IsLocalHostExist(host string) bool {
	h := NormalizeHost(host)
	if "localhost" == h || "127.0.0.1" == h || "::1" == h {
		return true
	}
	if ip := net.ParseIP(h); nil != ip {
		// 仅局域网调试 IP（192.168/10/172.16-31 与回环）
		return ip.IsLoopback() || ip.IsPrivate()
	}
	return false
}

// IsServerHostExist
//
//	@Description: 指定域名
//	@receiver c
//	@param host
//	@param server
//	@return bool
func (c *TenantDomainCache) IsServerHostExist(host string, server configPg.Server) bool {
	str := NormalizeHost(server.Domain)
	if "" == str {
		return false
	}
	return NormalizeHost(host) == str
}

// GetMap
//
//	@Description: 指定域名
func (c *TenantDomainCache) GetMap(host string) (dtoTenantDomain.Dto, bool) {
	dto, exist := c.DomainData.Load(host)
	return dto, exist
}
