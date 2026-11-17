package service

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modPublic"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamLogin"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/configPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/configPg/pg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/cacheService/cacheRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/consts/constHeaderPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/consts/constsPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/consts/constsRam/typeDomainPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/state/clientPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityRam"
	repositoryRam "github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/repositoryRam"
	authTokenPg "github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/components/authTokenPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/components/cachePg/cacheAuthPubPrivPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/components/loginLimitPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/middleware/event"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/sdk/auth/holderPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/sdk/ram/model/modRamAccount"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/multiTenantPg"
	"github.com/jinzhu/copier"
	"github.com/pangu-2/go-tools/tools/convPg"
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
	Bus                  event.Bus                                           `autowire:"?"`
	pg                   configPg.Pg                                         `value:"${pg}"`
	authLogin            pg.Auth                                             `value:"${pg.auth}"`
}

func NewAccountLoginService() *AccountLoginService {
	return new(AccountLoginService)
}

// accountRoles 返回账号真实角色；仅创始人作为引导管理员保留 administrator。
func accountRoles(info *entityRam.RamAccountEntity) []string {
	roles := make([]string, 0)
	if nil != info {
		for _, role := range info.Os.Data().Roles {
			if strPg.IsNotBlank(role) {
				roles = append(roles, role)
			}
		}
		if len(roles) == 0 && info.Founder == 1 {
			roles = append(roles, "administrator")
		}
	}
	return roles
}

func hasAdministrator(roles []string) bool {
	for _, role := range roles {
		if role == "administrator" {
			return true
		}
	}
	return false
}

func accountAuthCodes(roles []string) []string {
	if hasAdministrator(roles) {
		return []string{"AC_100100", "AC_100110", "AC_100120", "AC_100010"}
	}
	return make([]string, 0)
}

// Login
//
//	@Description:
//	@receiver c
//	@param ct
//	@return rt
func (c *AccountLoginService) Login(ctx *gin.Context, ct modRamLogin.LoginManageCt, tp typeDomainPg.TypeDomain) (rt rg.Rs[modRamLogin.LoginSuccess]) {
	log.Info(ctx, log.TagAppDef, "ct", log.Any("ct", ct), log.Any("tp", tp))
	//
	client := clientPg.Browser
	//
	if strPg.IsBlank(ct.Account) {
		return rt.ErrorMessage("账号不能为空")
	}
	if strPg.IsBlank(ct.Password) {
		return rt.ErrorMessage("密码不能为空")
	}
	if strPg.IsBlank(ct.OrgCode) {
		return rt.ErrorMessage("租户编号不能为空")
	}
	// 登录失败限流：按 域+账号+IP+租户 计数，超限短期锁定
	limitKey := loginLimitPg.Key("manage", tp.Code(), strings.ToLower(ct.Account), ctx.ClientIP(), ct.OrgCode)
	if allowed, wait := loginLimitPg.Allowed(limitKey); !allowed {
		return rt.ErrorMessage("登录尝试过于频繁，请 " + strconv.Itoa(wait) + " 秒后重试")
	}
	pwd := ct.Password
	log.Info(context.Background(), log.TagAppDef, "authLogin", log.Any("authLogin", c.authLogin))
	//解密密码
	if logingEn, ok := c.authLogin.LoginEncrypt["default"]; ok && logingEn {
		login := c.cacheSessionPubPrive.DecodeByLoginManage(ctx, pwd, clientPg.Browser, ct.OrgCode)
		if login.ErrorIs() {
			return rt.ErrorMessage(login.Message)
		}
		pwd = login.Data
	}

	md5 := cryptPg.Md5(ct.Account)
	info, b, err := c.dao.FindByTenantNoAccountMd5AndTypeDomain(ctx, ct.OrgCode, md5, tp.Code())
	if nil != err || !b {
		// 账号不存在与密码错误返回一致文案，避免账号枚举；同样计入限流
		loginLimitPg.Failure(limitKey)
		return rt.ErrorMessage("账号或密码错误")
	}
	if !enumStatePg.ENABLE.IsExistInt8(info.State) {
		return rt.ErrorMessage("账户已被禁用，不能登陆")
	}
	pwdInfo, result := c.daoAuth.FindByTypePasswordANo(ctx, info.No)
	if !result {
		loginLimitPg.Failure(limitKey)
		return rt.ErrorMessage("账号或密码错误")
	}
	if !userPg.PasswordVerify(pwdInfo.Value, pwd, pwdInfo.ExtraData) {
		// 只记录失败事实与账号，禁止输出明文口令/盐/哈希
		log.Debug(ctx, log.TagAppDef, "password verify failed", func() []log.Field {
			return []log.Field{log.Any("accountNo", info.No)}
		})
		loginLimitPg.Failure(limitKey)
		return rt.ErrorMessage("账号或密码错误")
	}
	loginLimitPg.Reset(limitKey)
	//
	//租户默认
	mult := multiTenantPg.MultiTenantPg{
		TenantNo: make([]string, 0),
	}
	mult.TenantNo = append(mult.TenantNo, info.TenantNo)
	//生成 token
	token := c.MakeToken(ctx, mult, info, tp, client, true)
	if token.ErrorIs() {
		return rt.ErrorMessage(token.Message)
	}
	dataToken := token.Data
	//记录登录日志
	c.loginLogSave(ctx, info)
	//
	successInfo := modRamLogin.LoginSuccessInfo{
		Account: info.Account,
		Name:    info.Name,
	}
	roles := accountRoles(info)
	successInfo.Roles = roles
	success := modRamLogin.LoginSuccess{
		Info:         successInfo,
		AccessToken:  dataToken.Access,
		RefreshToken: dataToken.Refresh,
		AuthCode:     accountAuthCodes(roles)}
	rt.Data = success
	return rt.Ok()
}

// LoginBySystem
//
//	@Description:
//	@receiver c
//	@param ct
//	@return rt
func (c *AccountLoginService) LoginBySystem(ctx *gin.Context, ct modRamLogin.LoginCt, tp typeDomainPg.TypeDomain, tenantNo string) (rt rg.Rs[modRamLogin.LoginSuccess]) {
	log.Info(ctx, log.TagAppDef, "ct", log.Any("ct", ct), log.Any("tp", tp))
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
	// 登录失败限流：按 域+账号+IP+租户 计数，超限短期锁定
	limitKey := loginLimitPg.Key("system", tp.Code(), strings.ToLower(ct.Account), ctx.ClientIP(), tenantNo)
	if allowed, wait := loginLimitPg.Allowed(limitKey); !allowed {
		return rt.ErrorMessage("登录尝试过于频繁，请 " + strconv.Itoa(wait) + " 秒后重试")
	}
	pwd := ct.Password
	log.Info(context.Background(), log.TagAppDef, "authLogin", log.Any("authLogin", c.authLogin))
	//解密密码
	if logingEn, ok := c.authLogin.LoginEncrypt["default"]; ok && logingEn {
		login := c.cacheSessionPubPrive.DecodeByLoginSystem(ctx, pwd)
		if login.ErrorIs() {
			//兼容历史前端：旧版 /xianfu/pub/config 默认返回管理域公钥，客户端仍用它加密
			legacy := c.cacheSessionPubPrive.DecodeByLoginManage(ctx, pwd, clientPg.Browser, constsPg.ACCOUNT_MANAGE_No)
			if legacy.ErrorIs() {
				return rt.ErrorMessage(login.Message)
			}
			log.Debug(ctx, log.TagBizDef, "系统登录使用管理域密钥解密成功，建议前端改用 typeDomain=system 获取公钥", nil)
			pwd = legacy.Data
		} else {
			pwd = login.Data
		}
	}

	md5 := cryptPg.Md5(ct.Account)
	info, b, err := c.dao.FindByAccountMd5AndTypeDomainAndTenantNo(ctx, md5, tp.Code(), tenantNo)
	if nil != err || !b {
		// 账号不存在与密码错误返回一致文案，避免账号枚举；同样计入限流
		loginLimitPg.Failure(limitKey)
		return rt.ErrorMessage("账号或密码错误")
	}
	if !enumStatePg.ENABLE.IsExistInt8(info.State) {
		return rt.ErrorMessage("账户已被禁用，不能登陆")
	}
	pwdInfo, result := c.daoAuth.FindByTypePasswordANo(ctx, info.No)
	if !result {
		loginLimitPg.Failure(limitKey)
		return rt.ErrorMessage("账号或密码错误")
	}
	if !userPg.PasswordVerify(pwdInfo.Value, pwd, pwdInfo.ExtraData) {
		// 只记录失败事实与账号，禁止输出明文口令/盐/哈希
		log.Debug(ctx, log.TagAppDef, "password verify failed", func() []log.Field {
			return []log.Field{log.Any("accountNo", info.No)}
		})
		loginLimitPg.Failure(limitKey)
		return rt.ErrorMessage("账号或密码错误")
	}
	loginLimitPg.Reset(limitKey)
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
	successInfo := modRamLogin.LoginSuccessInfo{
		Account: info.Account,
		Name:    info.Name,
	}
	roles := accountRoles(info)
	successInfo.Roles = roles
	success := modRamLogin.LoginSuccess{
		Info:         successInfo,
		AccessToken:  dataToken.Access,
		RefreshToken: dataToken.Refresh,
		AuthCode:     accountAuthCodes(roles)}
	rt.Data = success
	return rt.Ok()
}

func (c *AccountLoginService) loginLogSave(ctx *gin.Context, account *entityRam.RamAccountEntity) {
	now := time.Now()
	var tmp entityRam.RamAccountEntity
	copier.Copy(&tmp, account)
	tmp.LoginTime = &now
	tmp.TenantNo = account.TenantNo
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
	err := c.Bus.Publish(context.Background(), obj)
	if err != nil {
		log.Error(ctx, log.TagAppDef, err, "Bus.Publish err")
	}
}

// MakeToken 生成 token
func (c *AccountLoginService) MakeToken(ctx *gin.Context,
	mult multiTenantPg.MultiTenantPg,
	account *entityRam.RamAccountEntity,
	tp typeDomainPg.TypeDomain,
	client clientPg.Client,
	makeRefresh bool,
) (rt rg.Rs[authTokenPg.ResultAccessRefresh]) {
	token := authTokenPg.ResultAccessRefresh{}
	//
	privatePubKey := authTokenPg.Result{}
	// 获取 密钥对
	key := c.cacheSessionPubPrive.PaseKeyAccessToken(ctx, tp, client, account.TenantNo, nil)
	if strPg.IsNotBlank(key.Public) {
		privatePubKey.PrivateKey = key.Private
		privatePubKey.PublicKey = key.Public
	}
	//生成 令牌
	param := authTokenPg.Param{
		UniqueId:    strPg.GenerateNumberId22(),
		MultiTenant: mult,
		LoginNo:     account.No,
		No:          account.No,
		Name:        account.Name,
		Type:        string(tp),
		Result:      privatePubKey,
		TenantNo:    account.TenantNo,
	}
	ret := authTokenPg.MakePaseToken(param, c.pg.Jwt.System)
	if ret.ErrorIs() {
		return rt.ErrorMessage(ret.Message)
	}
	tokenResult := ret.Data
	token.Access = tokenResult.Token
	//刷新 token
	if makeRefresh {
		{
			//
			privatePubKey2 := authTokenPg.Result{}
			// 获取 密钥对
			key2 := c.cacheSessionPubPrive.PaseKeyRefreshToken(ctx, tp, client, account.TenantNo, nil)
			if strPg.IsNotBlank(key2.Public) {
				privatePubKey2.PrivateKey = key2.Private
				privatePubKey2.PublicKey = key2.Public
			}
			//生成 令牌
			param = authTokenPg.Param{
				UniqueId:    strPg.GenerateNumberId22(),
				MultiTenant: mult,
				LoginNo:     account.No,
				No:          account.No,
				Name:        account.Name,
				Type:        string(tp),
				Result:      privatePubKey2,
				TenantNo:    account.TenantNo,
			}
			ret2 := authTokenPg.MakePaseToken(param, c.pg.Jwt.System)
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
func (c *AccountLoginService) RefreshToken(ctx *gin.Context, ct modRamLogin.TokenRefreshCt, tp typeDomainPg.TypeDomain, client clientPg.Client) (rt rg.Rs[modPublic.LoginToken]) {
	accessToken := strings.TrimSpace(strings.Replace(ctx.GetHeader("Authorization"), authTokenPg.AuthScheme+" ", "", 1))
	refreshToken := strings.TrimSpace(strings.Replace(ctx.GetHeader("Authorization-Refresh"), authTokenPg.AuthScheme+" ", "", 1))
	if strPg.IsBlank(accessToken) || strPg.IsBlank(refreshToken) {
		return rt.ErrorMessage("刷新token失败")
	}

	// 仅用未验签载荷中的 tno 选择对应租户公钥；随后必须完成签名验证。
	accessPayload, ok := authTokenPg.ParseUnverified(accessToken)
	if !ok {
		return rt.ErrorMessage("刷新token失败")
	}
	refreshPayload, ok := authTokenPg.ParseUnverified(refreshToken)
	if !ok {
		return rt.ErrorMessage("刷新token失败")
	}
	var accessClaim, refreshClaim map[string]interface{}
	if json.Unmarshal(accessPayload, &accessClaim) != nil || json.Unmarshal(refreshPayload, &refreshClaim) != nil {
		return rt.ErrorMessage("刷新token失败")
	}
	tenantNoAcc := convPg.ObjToStr(accessClaim[authTokenPg.TenantNo])
	tenantNoRef := convPg.ObjToStr(refreshClaim[authTokenPg.TenantNo])
	if strPg.IsBlank(tenantNoAcc) || tenantNoAcc != tenantNoRef {
		return rt.ErrorMessage("刷新token失败")
	}

	accessKey, found := cacheAuthPubPrivPg.Get(cacheAuthPubPrivPg.AccessTokenKey(tenantNoAcc, tp, client))
	if !found || strPg.IsBlank(accessKey.Public) {
		return rt.ErrorMessage("刷新token失败")
	}
	accessVerified, vr := authTokenPg.VerifyByPublicKey(accessKey.Public, accessToken)
	if vr.ErrorIs() || nil == accessVerified {
		return rt.ErrorMessage("刷新token失败")
	}
	refreshKey, found := cacheAuthPubPrivPg.Get(cacheAuthPubPrivPg.RefreshTokenKey(tenantNoAcc, tp, client))
	if !found || strPg.IsBlank(refreshKey.Public) {
		return rt.ErrorMessage("刷新token失败")
	}
	refreshVerified, vr2 := authTokenPg.VerifyByPublicKey(refreshKey.Public, refreshToken)
	if vr2.ErrorIs() || nil == refreshVerified {
		return rt.ErrorMessage("刷新token失败")
	}

	loginNoAcc, err := accessVerified.GetSubject()
	if nil != err || strPg.IsBlank(loginNoAcc) {
		return rt.ErrorMessage("刷新token失败")
	}
	loginNoRef, err := refreshVerified.GetSubject()
	if nil != err || loginNoRef != loginNoAcc {
		return rt.ErrorMessage("刷新token失败")
	}
	tnoAcc, err := accessVerified.GetString(authTokenPg.TenantNo)
	if nil != err || tnoAcc != tenantNoAcc {
		return rt.ErrorMessage("刷新token失败")
	}
	tnoRef, err := refreshVerified.GetString(authTokenPg.TenantNo)
	if nil != err || tnoRef != tnoAcc {
		return rt.ErrorMessage("刷新token失败")
	}
	typAcc, _ := accessVerified.GetString(authTokenPg.Type)
	typRef, _ := refreshVerified.GetString(authTokenPg.Type)
	if typAcc != tp.Code() || typRef != tp.Code() {
		return rt.ErrorMessage("刷新token失败")
	}

	info, result := c.dao.FindByNo(ctx, loginNoAcc)
	if !result || info == nil ||
		info.TenantNo != tnoAcc || info.TypeDomain != tp.Code() ||
		!enumStatePg.ENABLE.IsExistInt8(info.State) {
		return rt.ErrorMessage("刷新token失败")
	}
	mult := multiTenantPg.MultiTenantPg{TenantNo: []string{info.TenantNo}}
	tokenRet := c.MakeToken(ctx, mult, info, tp, client, true)
	if tokenRet.ErrorIs() {
		return rt.ErrorMessage(tokenRet.Message)
	}
	return rt.OkData(modPublic.LoginToken{
		AccessToken:      tokenRet.Data.Access,
		AccessTokenCamel: tokenRet.Data.Access,
		RefreshToken:     tokenRet.Data.Refresh,
		ExpiresIn:        int64(c.pg.Jwt.System.Expire),
	})
}
