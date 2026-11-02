package cacheDomainPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/app/models/tc/dtoTenantDomain"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/configPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/mapPg"
	"go-spring.org/spring/gs"

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

// IsLocalHostExist
//
//	@Description: 本地域名
//	@receiver c
//	@param host
//	@return bool
func (c *TenantDomainCache) IsLocalHostExist(host string) bool {
	if strings.Contains(host, "localhost") {
		return true
	}
	if strings.Contains(host, "192.168.") {
		return true
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
	str := strings.Replace(server.Domain, "http://", "", -1)
	str = strings.Replace(str, "https://", "", -1)
	// 查找最后一个冒号的位置
	lastColonIndex := strings.LastIndex(str, ":")
	// 如果找到了冒号，并且冒号后面还有字符，则截取冒号前面的部分
	if lastColonIndex != -1 && lastColonIndex < len(str)-1 {
		str = str[:lastColonIndex]
	}
	if strings.Contains(host, str) {
		return true
	}
	return false
}

// GetMap
//
//	@Description: 指定域名
func (c *TenantDomainCache) GetMap(host string) (dtoTenantDomain.Dto, bool) {
	dto, exist := c.DomainData.Load(host)
	return dto, exist
}
