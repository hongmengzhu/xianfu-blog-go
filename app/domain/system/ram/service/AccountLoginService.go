package service

import (
	"context"
	"time"

	"github.com/farseer-go/eventBus"
	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/core/cache/cacheRam"
	authTokenPg2 "github.com/hongmengzhu/xianfu-blog-go/app/middleware/components/authTokenPg"
	modRamLogin2 "github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamLogin"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityRam"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/holderPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/holderPg/multiTenantPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/configPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/configPg/pg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constEventBusPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constHeaderPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constsRam/typeDomainPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/clientPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/sdk/ram/model/modRamAccount"
	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/cryptPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/userPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(NewAccountLoginService)
}

// AccountLoginService 登录
// @Description:
type AccountLoginService struct {
	dao                  *repositoryRam.RamAccountRepository                 `autowire:"?"`
	daoAuth              *repositoryRam.RamAccountAuthorizationRepository    `autowire:"?"`
	sessionAk            *repositoryRam.RamAccountSessionAccessKeyRepository `autowire:"?"`
	cacheSessionPubPrive *cacheRam.CacheSessionPubPrive                      `autowire:"?" `
	pg                   configPg.Pg                                         `value:"${pg}"`
	authLogin            pg.Auth                                             `value:"${pg.auth}"`
}

func NewAccountLoginService() *AccountLoginService {
	return new(AccountLoginService)
}

// Login
//
//	@Description:
//	@receiver c
//	@param ct
//	@return rt
func (c *AccountLoginService) Login(ctx *gin.Context, ct modRamLogin2.LoginCt, tp typeDomainPg.TypeDomain, tenantNo string) (rt rg.Rs[modRamLogin2.LoginSuccess]) {
	log.Infof(ctx, log.TagAppDef, "tp=%+v,ct=%+v", tp, ct)
	//
	client := clientPg.Browser
	//
	if !c.dao.Config().Domain.System {
		return rt.ErrorMessage("系统管理模块已禁用，不允许操作")
	}
	if strPg.IsBlank(ct.Account) {
		return rt.ErrorMessage("账号不能为空")
	}
	if strPg.IsBlank(ct.Password) {
		return rt.ErrorMessage("密码不能为空")
	}
	pwd := ct.Password
	log.Infof(context.Background(), log.TagAppDef, "authLogin=%+v", c.authLogin)
	//解密密码
	if logingEn, ok := c.authLogin.LoginEncrypt["default"]; ok && logingEn {
		login := c.cacheSessionPubPrive.DecodeByLoginSystem(ctx, pwd)
		if login.ErrorIs() {
			return rt.ErrorMessage(login.Message)
		}
		pwd = login.Data
	}
	log.Infof(context.Background(), log.TagAppDef, "pwd=%+v", pwd)

	md5 := cryptPg.Md5(ct.Account)
	info, b, err := c.dao.FindByAccountMd5AndTypeDomainAndTenantNo(ctx, md5, tp.Code(), tenantNo)
	if nil != err {
		return rt.Error()
	}
	if !b {
		return rt.ErrorMessage("账号不存在")
	}
	if !enumStatePg.ENABLE.IsExistInt8(info.State) {
		return rt.ErrorMessage("账户已被禁用，不能登陆")
	}
	pwdInfo, result := c.daoAuth.FindByTypePasswordANo(ctx, info.No)
	if !result {
		return rt.ErrorMessage("用户密码未设置")
	}
	if !userPg.PasswordVerify(pwdInfo.Value, pwd, pwdInfo.ExtraData) {
		log.Debugf(ctx, log.TagAppDef, "pwd=%+v,[%+v],[value]=%+v,[extra]=%+v", pwd, userPg.PasswordSalt(pwd, pwdInfo.ExtraData), pwdInfo.Value, pwdInfo.ExtraData)
		return rt.ErrorMessage("账号密码错误")
	}
	//
	//租户默认
	mult := multiTenantPg.MultiTenantPg{
		TenantNo: make([]string, 0),
	}
	mult.TenantNo = append(mult.TenantNo, tp.Code())
	//生成 token
	token := c.MakeToken(ctx, mult, info, tp, client, true)
	if token.ErrorIs() {
		return rt.ErrorMessage(token.Message)
	}
	dataToken := token.Data
	//记录登录日志
	c.loginLogSave(ctx, info)
	//
	successInfo := modRamLogin2.LoginSuccessInfo{
		Account: info.Account,
		Name:    info.Name,
	}
	success := modRamLogin2.LoginSuccess{
		Info:         successInfo,
		AccessToken:  dataToken.Access,
		RefreshToken: dataToken.Refresh,
		AuthCode:     []string{"AC_100100", "AC_100110", "AC_100120", "AC_100010"}}
	rt.Data = success
	return rt.Ok()
}

func (c *AccountLoginService) loginLogSave(ctx *gin.Context, account *entityRam.RamAccountEntity) {
	now := time.Now()
	var tmp entityRam.RamAccountEntity
	copier.Copy(&tmp, account)
	tmp.LoginTime = &now
	tmp.TenantNo = typeDomainPg.System.Code()
	obj := modRamAccount.LoginLogDto{
		Account:     tmp,
		Ano:         tmp.No,
		AppNo:       "",
		Client:      "",
		LoginSource: "",
		Ip:          ctx.ClientIP(),
		ExtraData:   make(map[string]any),
	}
	ua := ctx.GetHeader(constHeaderPg.HeaderUserAgent)
	obj.ExtraData[constHeaderPg.HeaderUserAgent] = ua
	//保存到数据库
	eventBus.PublishEventAsync(constEventBusPg.RamAccountLoginLog, obj)
}

// MakeToken 生成 token
func (c *AccountLoginService) MakeToken(ctx *gin.Context,
	mult multiTenantPg.MultiTenantPg,
	account *entityRam.RamAccountEntity,
	tp typeDomainPg.TypeDomain,
	client clientPg.Client,
	makeRefresh bool,
) (rt rg.Rs[authTokenPg2.ResultAccessRefresh]) {
	token := authTokenPg2.ResultAccessRefresh{}
	//
	privatePubKey := authTokenPg2.Result{}
	// 获取 密钥对
	key := c.cacheSessionPubPrive.PaseKeyAccessToken(ctx, tp, client, account.TenantNo, nil)
	if strPg.IsNotBlank(key.Public) {
		privatePubKey.PrivateKey = key.Private
		privatePubKey.PublicKey = key.Public
	}
	//生成 令牌
	param := authTokenPg2.Param{
		UniqueId:    strPg.GenerateNumberId22(),
		MultiTenant: mult,
		LoginNo:     account.No,
		No:          account.No,
		Name:        account.Name,
		Type:        string(tp),
		Result:      privatePubKey,
		TenantNo:    account.TenantNo,
	}
	ret := authTokenPg2.MakePaseToken(param, c.pg.Jwt.System)
	if ret.ErrorIs() {
		return rt.ErrorMessage(ret.Message)
	}
	tokenResult := ret.Data
	token.Access = tokenResult.Token
	//刷新 token
	if makeRefresh {
		{
			//
			privatePubKey2 := authTokenPg2.Result{}
			// 获取 密钥对
			key2 := c.cacheSessionPubPrive.PaseKeyRefreshToken(ctx, tp, client, account.TenantNo, nil)
			if strPg.IsNotBlank(key2.Public) {
				privatePubKey2.PrivateKey = key2.Private
				privatePubKey2.PublicKey = key2.Public
			}
			//生成 令牌
			param = authTokenPg2.Param{
				UniqueId:    strPg.GenerateNumberId22(),
				MultiTenant: mult,
				LoginNo:     account.No,
				No:          account.No,
				Name:        account.Name,
				Type:        string(tp),
				Result:      privatePubKey2,
				TenantNo:    account.TenantNo,
			}
			ret2 := authTokenPg2.MakePaseToken(param, c.pg.Jwt.System)
			if ret2.ErrorIs() {
				return rt.ErrorMessage(ret.Message)
			}
			tokenResult2 := ret2.Data
			token.Refresh = tokenResult2.Token
		}
	}
	//
	return rt.OkData(token)
}

func (c *AccountLoginService) Logout(holder holderPg.HolderPg) (rt rg.Rs[string]) {
	return rt.Ok()
}

// RefreshToken
//
//	@Description:  刷新
//	@receiver c
func (c *AccountLoginService) RefreshToken(ctx *gin.Context, ct modRamLogin2.TokenRefreshCt) (rt rg.Rs[modRamLogin2.LoginSuccess]) {
	token := ctx.GetHeader("Authorization")
	rt.Data = modRamLogin2.LoginSuccess{Token: token, AccessToken: token}
	return rt.Ok()
}
