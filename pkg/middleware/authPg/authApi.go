package authPg

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/configPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/consts/constContextPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/consts/constHeaderPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/consts/constsRam/typeDomainPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/components/cachePg/cacheDiplPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/components/cachePg/cacheDomainPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/sdk/auth/holderPg/holderApiPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/multiTenantPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(&GroupApiMiddlewareSp{})
}

// 中间件 服务
type GroupApiMiddlewareSp struct {
	domain *cacheDomainPg.TenantDomainCache `autowire:"?"`
	pg     configPg.Pg                      `value:"${pg}"`
	server configPg.Server                  `value:"${server}"`
}

// 权限验证 中间件
func GroupApiMiddleware(m *GroupApiMiddlewareSp) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader(constHeaderPg.HeaderAuthorization)
		if strPg.IsBlank(header) {
			c.JSON(http.StatusOK, rg.Error[string](constHeaderPg.HeaderAuthorization+" 参数不能为空"))
			c.Abort()
			return
		}
		split := strings.SplitN(header, ":", 2)
		if len(split) != 2 || strPg.IsBlank(split[0]) || strPg.IsBlank(split[1]) {
			c.JSON(http.StatusOK, rg.Error[string]("Authorization 格式错误"))
			c.Abort()
			return
		}
		get, ok := cacheDiplPg.Get(split[0])
		if !ok || !cacheDiplPg.HashShaVerify(get.Key, get.Secret, split[1]) {
			c.JSON(http.StatusOK, rg.Error[string]("Authorization 密钥无效"))
			c.Abort()
			return
		}
		if enumStatePg.ENABLE.IndexInt8() != get.State {
			c.JSON(http.StatusOK, rg.Error[string]("Authorization 密钥已禁用"))
			c.Abort()
			return
		}
		if nil != get.ExpiryDate && time.Now().After(*get.ExpiryDate) {
			c.JSON(http.StatusOK, rg.Error[string]("Authorization 密钥已过期"))
			c.Abort()
			return
		}
		// 开放 API 必须绑定租户，避免空租户导致查询退化为全平台。
		if strPg.IsBlank(get.TenantNo) {
			c.JSON(http.StatusOK, rg.Error[string]("Authorization 密钥未绑定租户"))
			c.Abort()
			return
		}
		pg := holderApiPg.HolderPg{
			HolderData: holderApiPg.DiplHolder{
				Name:     get.Name,
				No:       get.No,
				TenantNo: get.TenantNo,
				Ano:      get.Ano,
			}}
		pg.TypeDomain = typeDomainPg.Manage.Code()
		pg.MultiTenant = multiTenantPg.MultiTenantPg{TenantNo: []string{get.TenantNo}}
		c.Set(constContextPg.AUTH_LOGIN_API, pg)
		c.Set(constContextPg.CTX_MULITI_TENANT, multiTenantPg.GetMultiTableManage())
		c.Set(constHeaderPg.WebTenantNo, get.TenantNo)
		c.Next()
	}
}
