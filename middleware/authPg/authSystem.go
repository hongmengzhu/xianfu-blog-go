package authPg

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/system/ram/service"
	"github.com/hongmengzhu/xianfu-blog-go/middleware/components/authTokenPg"
	"github.com/hongmengzhu/xianfu-blog-go/middleware/components/cachePg/cacheAuthPubPrivPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/holderPg/multiTenantPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/configPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constContextPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constHeaderPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

// 验证 ownerCode 必填
var systemUrlPaths = make([]string, 0)

func init() {
	gs.Provide(&GroupSystemMiddlewareSp{})
	//商品模块
	systemUrlPaths = append(systemUrlPaths, "manage/gc")
}

// 中间件 服务
type GroupSystemMiddlewareSp struct {
	sv *service.RamAccountMiddlewareService `autowire:"?"`
	pg configPg.Pg                          `value:"${pg}"`
}

// 权限验证 中间件
func GroupSystemMiddleware(m *GroupSystemMiddlewareSp) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader(constHeaderPg.HeaderAuthorization_system)
		if strPg.IsNotBlank(token) {
			get, b := cacheAuthPubPrivPg.Get(cacheAuthPubPrivPg.KeySystem())
			if !b {
				c.JSON(200, rg.Error[string]("密钥不存在"))
				c.Abort()
				return
			}
			token = strings.Replace(token, authTokenPg.AuthScheme+" ", "", -1)
			t, r := authTokenPg.VerifyByPublicKey(get.Public, token)
			if r.ErrorIs() {
				log.Debugf(context.Background(), log.TagAppDef, "JWT= %+v", r)
				c.JSON(200, r)
				c.Abort()
				return
			}
			//获取登录 编号
			loginNo, err2 := t.GetSubject()
			if nil != err2 {
				c.JSON(200, rg.Error[string]("解析失败"))
				c.Abort()
				return
			}
			if strPg.IsBlank(loginNo) {
				c.JSON(200, rg.Error[string]("解析失败"))
				c.Abort()
				return
			}
			//获取登录信息
			rt := m.sv.FindByLoginNo(c, loginNo)
			if rt.ErrorIs() {
				c.JSON(200, rt)
				c.Abort()
				return
			}
			log.Debugf(context.Background(), log.TagAppDef, "rt= %+v", rt)
			log.Debugf(context.Background(), log.TagAppDef, "rt.Rule= %+v", rt.Data.Rule)
			//
			var data = rt.Data
			//如果配置了systemUrlPaths，则必须携带ownerCode
			if len(systemUrlPaths) > 0 && data.GetOwnerNo() == "" {
				path := c.Request.URL.Path
				for _, item := range systemUrlPaths {
					//必须登录 验证，url 包含的路径，那么提示错误
					if "" != item && strings.Index(path, item) != -1 {
						log.Tracef(context.Background(), log.TagAppDef, "pass %+v", c.Request.URL.Path)
						c.JSON(200, rg.Error[string]("所有者错误"))
						c.Abort()
						return
					}
				}
			}
			c.Set(constContextPg.CTX_MULITI_TENANT, multiTenantPg.GetMultiTableSys())

			//
			c.Set(constContextPg.AUTH_LOGIN, rt.Data)
			c.Set(constContextPg.CTX_RULE, rt.Data.Rule)
			//value, exists := c.Get(constContextPg.CTX_RULE)
			//m.log.Infof("Set.constContextPg.CTX_RULE.get= %+v,", value, exists)
			//m.log.WithContext(c.Context()).Tracef("::AdminManageFilter")

			//详细权限
			//err = baseAuth.CasbinHandler(ctx, rt.Data, m.pg, m.casbin)
			//if nil != err {
			//	c.JSON(r.Error(err.Error()))
			//	return
			//}
			c.Next()
			return
		}
		c.JSON(200, rg.Error[string](constHeaderPg.HeaderAuthorization+" 参数不能为空"))
		c.Abort()
		return
	}
}
