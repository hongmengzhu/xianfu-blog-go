package service

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamLogin"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityRam"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/holderPg/multiTenantPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/idp"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constHeaderPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constsRam/typeDomainPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/clientPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/pangu-2/go-tools/tools/noPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
	"golang.org/x/oauth2"
)

func init() {
	gs.Provide(NewAccountLoginIdpService)
}

// idpBaseConfig 认证源基础配置（从 BaseConfig JSON 解析）
type idpBaseConfig struct {
	ClientId        string `json:"clientId"`
	ClientSecret    string `json:"clientSecret"`
	ClientId2       string `json:"clientId2"`
	ClientSecret2   string `json:"clientSecret2"`
	AppId           string `json:"appId"`
	HostUrl         string `json:"hostUrl"`
	RedirectUrl     string `json:"redirectUrl"`
	DisableSsl      bool   `json:"disableSsl"`
	TokenURL        string `json:"tokenURL"`
	AuthURL         string `json:"authURL"`
	UserInfoURL     string `json:"userInfoURL"`
	AppCertificate  string `json:"appCertificate"`
	RootCertificate string `json:"rootCertificate"`
}

// AccountLoginIdpService IDP OAuth 登录
// @Description:
type AccountLoginIdpService struct {
	daoAccount    *repositoryRam.RamAccountRepository           `autowire:"?"`
	daoSource     *repositoryRam.RamIdentitySourceRepository    `autowire:"?"`
	daoProvider   *repositoryRam.RamIdentityProviderRepository  `autowire:"?"`
	daoBinding    *repositoryRam.RamIdpBindingRepository        `autowire:"?"`
	daoSessionLog *repositoryRam.RamAccountSessionLogRepository `autowire:"?"`
	loginService  *AccountLoginService                          `autowire:"?"`
	mfaService    *AccountMfaService                            `autowire:"?"`
}

func NewAccountLoginIdpService() *AccountLoginIdpService {
	return new(AccountLoginIdpService)
}

// Login IDP OAuth 登录
func (c *AccountLoginIdpService) Login(ctx *gin.Context, ct modRamLogin.IdpLoginCt) (rt rg.Rs[modRamLogin.IdpLoginSuccess]) {
	log.Infof(ctx, log.TagAppDef, "IdpLogin sourceNo=%s, method=%s", ct.SourceNo, ct.Method)

	// 1. 参数校验
	if strPg.IsBlank(ct.SourceNo) {
		return rt.ErrorMessage("认证源编号不能为空")
	}
	if strPg.IsBlank(ct.Code) {
		return rt.ErrorMessage("授权码不能为空")
	}

	// 2. 查询认证源
	source, found := c.daoSource.FindByNo(ctx, ct.SourceNo)
	if !found {
		return rt.ErrorMessage("认证源不存在")
	}
	if !enumStatePg.ENABLE.IsExistInt8(source.State) {
		return rt.ErrorMessage("认证源已禁用")
	}

	// 3. 查询关联的身份提供商
	providerEntity, found := c.daoProvider.FindByNo(ctx, source.Idp)
	if !found {
		return rt.ErrorMessage("身份提供商不存在")
	}
	if !enumStatePg.ENABLE.IsExistInt8(providerEntity.State) {
		return rt.ErrorMessage("身份提供商已禁用")
	}

	// 4. 解析 BaseConfig JSON
	var baseCfg idpBaseConfig
	if strPg.IsNotBlank(source.BaseConfig) {
		if err := json.Unmarshal([]byte(source.BaseConfig), &baseCfg); err != nil {
			log.Errorf(ctx, log.TagAppDef, "解析 BaseConfig 失败: %v", err)
			return rt.ErrorMessage("认证源配置解析失败")
		}
	}

	// 5. 构造 idp.ProviderInfo
	redirectUrl := ct.RedirectUri
	if strPg.IsBlank(redirectUrl) {
		redirectUrl = baseCfg.RedirectUrl
	}
	idpInfo := &idp.ProviderInfo{
		Type:            providerEntity.Code,
		SubType:         source.Protocol,
		ClientId:        baseCfg.ClientId,
		ClientSecret:    baseCfg.ClientSecret,
		ClientId2:       baseCfg.ClientId2,
		ClientSecret2:   baseCfg.ClientSecret2,
		AppId:           baseCfg.AppId,
		HostUrl:         baseCfg.HostUrl,
		RedirectUrl:     baseCfg.RedirectUrl,
		DisableSsl:      baseCfg.DisableSsl,
		CodeVerifier:    ct.CodeVerifier,
		TokenURL:        baseCfg.TokenURL,
		AuthURL:         baseCfg.AuthURL,
		UserInfoURL:     baseCfg.UserInfoURL,
		AppCertificate:  baseCfg.AppCertificate,
		RootCertificate: baseCfg.RootCertificate,
	}

	// 6. 创建 IdProvider 实例
	idProvider, err := idp.GetIdProvider(idpInfo, redirectUrl)
	if err != nil {
		log.Errorf(ctx, log.TagAppDef, "创建 IdProvider 失败: %v", err)
		return rt.ErrorMessage("不支持的认证提供商类型: " + providerEntity.Code)
	}

	// 7. 注入 HTTP 客户端
	idProvider.SetHttpClient(&http.Client{})

	// 8. 用授权码换 Token
	oauthToken, err := idProvider.GetToken(ct.Code)
	if err != nil {
		log.Errorf(ctx, log.TagAppDef, "获取 Token 失败: %v", err)
		return rt.ErrorMessage("获取 Token 失败: " + err.Error())
	}
	if oauthToken == nil || !oauthToken.Valid() {
		return rt.ErrorMessage("获取的 Token 无效")
	}

	// 9. 获取第三方用户信息
	userInfo, err := idProvider.GetUserInfo(oauthToken)
	if err != nil {
		log.Errorf(ctx, log.TagAppDef, "获取用户信息失败: %v", err)
		return rt.ErrorMessage("获取用户信息失败: " + err.Error())
	}
	if userInfo == nil {
		return rt.ErrorMessage("获取的用户信息为空")
	}
	log.Infof(ctx, log.TagAppDef, "OAuth userInfo: id=%s, username=%s, email=%s", userInfo.Id, userInfo.Username, userInfo.Email)

	// 10. 查找绑定
	var binding *entityRam.RamIdpBindingEntity
	bindingFound := false

	// 优先按 ExternalSub 查找
	if strPg.IsNotBlank(userInfo.Id) {
		binding, bindingFound = c.daoBinding.FindByIdpAndExternalSub(ctx, source.Idp, userInfo.Id)
	}
	// 若未找到，尝试 OpenId
	if !bindingFound && strPg.IsNotBlank(userInfo.Extra["openid"]) {
		binding, bindingFound = c.daoBinding.FindByIdpAndOpenId(ctx, source.Idp, userInfo.Extra["openid"])
	}
	// 若未找到，尝试 UnionId
	if !bindingFound && strPg.IsNotBlank(userInfo.UnionId) {
		binding, bindingFound = c.daoBinding.FindByIdpAndUnionId(ctx, source.Idp, userInfo.UnionId)
	}

	isSignup := false

	if bindingFound && binding != nil {
		// 11a. 已绑定 → 获取绑定账号
		if strPg.IsBlank(binding.BindAno) {
			return rt.ErrorMessage("绑定记录存在但未关联账号")
		}
		account, accountFound := c.daoAccount.FindByNo(ctx, binding.BindAno)
		if !accountFound {
			return rt.ErrorMessage("绑定的账号不存在")
		}
		if !enumStatePg.ENABLE.IsExistInt8(account.State) {
			return rt.ErrorMessage("账号已被禁用")
		}

		// 12. 更新 binding 的 token 和登录时间
		now := time.Now()
		c.daoBinding.Update(ctx, entityRam.RamIdpBindingEntity{
			AccessToken:   oauthToken.AccessToken,
			RefreshToken:  oauthToken.RefreshToken,
			LastLoginTime: &now,
			NickName:      userInfo.DisplayName,
			Avatar:        userInfo.AvatarUrl,
			Mail:          userInfo.Email,
			Phone:         userInfo.Phone,
		}, binding.ID)

		// 检查是否需要 MFA
		return c.checkMfaOrLoginSuccess(ctx, account, source, userInfo, oauthToken, isSignup)

	} else {
		// 11b/11c. 未绑定
		if ct.Method == "link" {
			return rt.ErrorMessage("当前用户未绑定该 OAuth 提供商")
		}

		// 检查是否允许自动创建
		if source.AutoCreateUser != 1 {
			return rt.ErrorMessage("该认证源不允许自动创建账号，请先绑定已有账号")
		}

		// 自动创建账号
		isSignup = true
		account, bindEntity, errCreate := c.createAccountAndBinding(ctx, source, providerEntity, userInfo, oauthToken)
		if errCreate != nil {
			return rt.ErrorMessage("自动创建账号失败: " + errCreate.Error())
		}

		// 检查是否需要 MFA
		_ = bindEntity // binding 已在 createAccountAndBinding 中创建
		return c.checkMfaOrLoginSuccess(ctx, account, source, userInfo, oauthToken, isSignup)
	}
}

// createAccountAndBinding 自动创建账号和绑定记录
func (c *AccountLoginIdpService) createAccountAndBinding(
	ctx *gin.Context,
	source *entityRam.RamIdentitySourceEntity,
	providerEntity *entityRam.RamIdentityProviderEntity,
	userInfo *idp.UserInfo,
	oauthToken *oauth2.Token,
) (*entityRam.RamAccountEntity, *entityRam.RamIdpBindingEntity, error) {
	now := time.Now()
	accountNo := noPg.No()

	// 创建 RamAccount
	account := &entityRam.RamAccountEntity{
		No:           accountNo,
		TenantNo:     source.TenantNo,
		OrgNo:        source.OrgNo,
		StoreNo:      source.StoreNo,
		TypeDomain:   typeDomainPg.Manage.String(),
		Account:      userInfo.Username,
		Name:         userInfo.DisplayName,
		RealName:     userInfo.DisplayName,
		Mail:         userInfo.Email,
		Phone:        userInfo.Phone,
		Avatar:       userInfo.AvatarUrl,
		State:        int8(enumStatePg.ENABLE),
		RegisterTime: &now,
		RegisterIP:   ctx.ClientIP(),
		LoginTime:    &now,
		CreateAt:     &now,
	}
	// 如果用户名为空，使用 Id 作为用户名
	if strPg.IsBlank(account.Account) {
		account.Account = userInfo.Id
	}
	// 如果名称为空，使用用户名
	if strPg.IsBlank(account.Name) {
		account.Name = account.Account
	}

	err, _ := c.daoAccount.Create(ctx, account)
	if err != nil {
		log.Errorf(ctx, log.TagAppDef, "创建账号失败: %v", err)
		return nil, nil, err
	}
	log.Infof(ctx, log.TagAppDef, "IDP 自动创建账号: no=%s, account=%s", account.No, account.Account)

	// 创建 RamIdpBinding
	bindingNo := noPg.No()
	binding := &entityRam.RamIdpBindingEntity{
		No:            bindingNo,
		TenantNo:      source.TenantNo,
		OrgNo:         source.OrgNo,
		StoreNo:       source.StoreNo,
		TypeDomain:    typeDomainPg.Manage.String(),
		Idp:           source.Idp,
		ExternalSub:   userInfo.Id,
		OpenId:        userInfo.Extra["openid"],
		UnionId:       userInfo.UnionId,
		BindAno:       accountNo,
		State:         int8(enumStatePg.ENABLE),
		StateBind:     2, // 已绑定
		BindTime:      &now,
		LastLoginTime: &now,
		AccessToken:   oauthToken.AccessToken,
		RefreshToken:  oauthToken.RefreshToken,
		Platform:      providerEntity.Platform,
		Protocol:      source.Protocol,
		Mail:          userInfo.Email,
		Phone:         userInfo.Phone,
		NickName:      userInfo.DisplayName,
		Avatar:        userInfo.AvatarUrl,
		CreateAt:      &now,
	}
	errBind, _ := c.daoBinding.Create(ctx, binding)
	if errBind != nil {
		log.Errorf(ctx, log.TagAppDef, "创建绑定记录失败: %v", errBind)
		return nil, nil, errBind
	}

	return account, binding, nil
}

// loginSuccess 登录成功，生成 token 并记录审计日志
func (c *AccountLoginIdpService) loginSuccess(
	ctx *gin.Context,
	account *entityRam.RamAccountEntity,
	source *entityRam.RamIdentitySourceEntity,
	userInfo *idp.UserInfo,
	oauthToken *oauth2.Token,
	isSignup bool,
) (rt rg.Rs[modRamLogin.IdpLoginSuccess]) {
	// 生成 token（复用 AccountLoginService.MakeToken）
	mult := multiTenantPg.MultiTenantPg{
		TenantNo: []string{account.TenantNo},
	}
	tokenResult := c.loginService.MakeToken(ctx, mult, account, typeDomainPg.Manage, clientPg.Browser, true)
	if tokenResult.ErrorIs() {
		return rt.ErrorMessage(tokenResult.Message)
	}
	dataToken := tokenResult.Data

	// 写入审计日志
	c.saveSessionLog(ctx, account, source, userInfo, 1)

	success := modRamLogin.IdpLoginSuccess{
		AccessToken:  dataToken.Access,
		RefreshToken: dataToken.Refresh,
		Info: modRamLogin.LoginSuccessInfo{
			Account: account.Account,
			Name:    account.Name,
			Avatar:  account.Avatar,
			Roles:   []string{},
		},
		IsSignup: isSignup,
	}
	return rt.OkData(success)
}

// saveSessionLog 保存 IDP 登录审计日志
func (c *AccountLoginIdpService) saveSessionLog(
	ctx *gin.Context,
	account *entityRam.RamAccountEntity,
	source *entityRam.RamIdentitySourceEntity,
	userInfo *idp.UserInfo,
	eventResult int8,
) {
	now := time.Now()
	logEntity := &entityRam.RamAccountSessionLogEntity{
		No:            noPg.No(),
		TenantNo:      account.TenantNo,
		OrgNo:         account.OrgNo,
		StoreNo:       account.StoreNo,
		TypeDomain:    typeDomainPg.Manage.String(),
		Idp:           source.Idp,
		SourceNo:      source.No,
		Protocol:      source.Protocol,
		Ano:           account.No,
		BindAno:       account.No,
		ExternalSub:   userInfo.Id,
		EventCategory: "login",
		EventType:     "login_success",
		EventResult:   eventResult,
		IpAddress:     ctx.ClientIP(),
		UserAgent:     ctx.GetHeader(constHeaderPg.HeaderUserAgent),
		OperateAt:     &now,
		LoginAt:       &now,
		Ip:            ctx.ClientIP(),
		LoginSource:   "idp",
	}
	err, _ := c.daoSessionLog.Create(ctx, logEntity)
	if err != nil {
		log.Errorf(ctx, log.TagAppDef, "保存 IDP 登录审计日志失败: %v", err)
	}
}

// checkMfaOrLoginSuccess 检查是否需要 MFA，如果需要则返回 MFA 令牌，否则直接登录成功
func (c *AccountLoginIdpService) checkMfaOrLoginSuccess(
	ctx *gin.Context,
	account *entityRam.RamAccountEntity,
	source *entityRam.RamIdentitySourceEntity,
	userInfo *idp.UserInfo,
	oauthToken *oauth2.Token,
	isSignup bool,
) (rt rg.Rs[modRamLogin.IdpLoginSuccess]) {
	// 新注册用户跳过 MFA
	if isSignup {
		return c.loginSuccess(ctx, account, source, userInfo, oauthToken, isSignup)
	}

	// 检查是否需要 MFA
	mfaRequired, mfaToken, mfaType := c.mfaService.CheckMfaRequired(ctx, account.No)
	if mfaRequired {
		// 记录审计日志
		c.saveSessionLog(ctx, account, source, userInfo, 2) // 2=MFA待验证
		return rt.OkData(modRamLogin.IdpLoginSuccess{
			MfaRequired: true,
			MfaToken:    mfaToken,
			MfaType:     mfaType,
			Info: modRamLogin.LoginSuccessInfo{
				Account: account.Account,
				Name:    account.Name,
				Avatar:  account.Avatar,
			},
		})
	}

	return c.loginSuccess(ctx, account, source, userInfo, oauthToken, isSignup)
}

// RefreshToken 刷新 token（复用 AccountLoginService 的逻辑）
func (c *AccountLoginIdpService) RefreshToken(ctx *gin.Context, ct modRamLogin.TokenRefreshCt) (rt rg.Rs[modRamLogin.IdpLoginSuccess]) {
	// 复用 AccountLoginService.RefreshToken
	result := c.loginService.RefreshToken(ctx, ct, typeDomainPg.Manage, clientPg.Browser)
	if result.ErrorIs() {
		return rt.ErrorMessage(result.Message)
	}
	data := result.Data
	success := modRamLogin.IdpLoginSuccess{
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
	}
	return rt.OkData(success)
}

// LoginByLdap LDAP 认证源登录（通过 AccountLoginIdpService 统一入口）
// 当认证源 protocol 为 "ldap" 时，调用此方法
func (c *AccountLoginIdpService) LoginByLdap(ctx *gin.Context, sourceNo string, username string, password string) (rt rg.Rs[modRamLogin.IdpLoginSuccess]) {
	// 查找认证源
	source, found := c.daoSource.FindByNo(ctx, sourceNo)
	if !found {
		return rt.ErrorMessage("认证源不存在")
	}
	if !enumStatePg.ENABLE.IsExistInt8(source.State) {
		return rt.ErrorMessage("认证源已禁用")
	}
	if source.Protocol != "ldap" {
		return rt.ErrorMessage("该认证源不是 LDAP 协议")
	}

	// 委托给 LdapService 处理
	// 注意：这里通过直接构造请求并调用 LdapService.Login 实现
	// 实际使用时建议直接注入 LdapService
	return rt.ErrorMessage("LDAP 登录请使用 LdapService 直接调用")
}
