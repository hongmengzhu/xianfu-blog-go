package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamLogin"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamSaml"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityRam"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/holderPg/multiTenantPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/idp"
	samlPkg "github.com/hongmengzhu/xianfu-blog-go/pkg/auth/saml"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constHeaderPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constsRam/typeDomainPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/clientPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/pangu-2/go-tools/tools/noPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(NewSamlLoginService)
}

// samlSourceConfig SAML 认证源 BaseConfig JSON 结构
type samlSourceConfig struct {
	IdpSsoURL  string `json:"idpSsoUrl"`  // IdP SSO 地址
	IdpIssuer  string `json:"idpIssuer"`  // IdP EntityID
	SpIssuer   string `json:"spIssuer"`   // SP EntityID（可选，默认使用 ACS URL）
	EnableSign bool   `json:"enableSign"` // 是否签名 AuthnRequest
}

// SamlLoginService SAML 登录 Service
type SamlLoginService struct {
	daoAccount    *repositoryRam.RamAccountRepository           `autowire:"?"`
	daoSource     *repositoryRam.RamIdentitySourceRepository    `autowire:"?"`
	daoCredential *repositoryRam.RamIdpCredentialRepository     `autowire:"?"`
	daoMetadata   *repositoryRam.RamIdpMetadataCacheRepository  `autowire:"?"`
	daoBinding    *repositoryRam.RamIdpBindingRepository        `autowire:"?"`
	daoSessionLog *repositoryRam.RamAccountSessionLogRepository `autowire:"?"`
	loginService  *AccountLoginService                          `autowire:"?"`
	mfaService    *AccountMfaService                            `autowire:"?"`
}

func NewSamlLoginService() *SamlLoginService {
	return new(SamlLoginService)
}

// GetSamlLoginUrl 生成 SAML 认证请求 URL
func (c *SamlLoginService) GetSamlLoginUrl(ctx *gin.Context, sourceNo string) (rt rg.Rs[modRamSaml.SamlRedirectVo]) {
	if strPg.IsBlank(sourceNo) {
		return rt.ErrorMessage("认证源编号不能为空")
	}

	// 1. 查询认证源
	source, found := c.daoSource.FindByNo(ctx, sourceNo)
	if !found {
		return rt.ErrorMessage("认证源不存在")
	}
	if !enumStatePg.ENABLE.IsExistInt8(source.State) {
		return rt.ErrorMessage("认证源已禁用")
	}
	if source.Protocol != "saml2" {
		return rt.ErrorMessage("该认证源不是 SAML2 协议")
	}

	// 2. 构建 SP
	sp, cfg, err := c.buildSpFromSource(ctx, source)
	if err != nil {
		log.Errorf(ctx, log.TagAppDef, "构建 SAML SP 失败: %v", err)
		return rt.ErrorMessage("构建 SAML SP 失败: " + err.Error())
	}
	_ = cfg

	// 3. 生成 RelayState（JSON 编码 sourceNo）
	relayState, err := json.Marshal(map[string]string{
		"sourceNo": sourceNo,
	})
	if err != nil {
		return rt.ErrorMessage("生成 RelayState 失败")
	}

	// 4. 生成 SAML Request
	authURL, method, err := samlPkg.GenerateSamlRequest(sp, string(relayState))
	if err != nil {
		log.Errorf(ctx, log.TagAppDef, "生成 SAML Request 失败: %v", err)
		return rt.ErrorMessage("生成 SAML 请求失败: " + err.Error())
	}

	vo := modRamSaml.SamlRedirectVo{
		RedirectUrl: authURL,
		Method:      method,
	}
	if method == "POST" {
		vo.PostBody = authURL
		vo.RedirectUrl = sp.IdentityProviderSSOURL
	}

	return rt.OkData(vo)
}

// HandleSamlCallback 处理 SAML Response 回调
func (c *SamlLoginService) HandleSamlCallback(ctx *gin.Context, ct modRamSaml.SamlCallbackCt) (rt rg.Rs[modRamLogin.IdpLoginSuccess]) {
	if strPg.IsBlank(ct.SAMLResponse) {
		return rt.ErrorMessage("SAML Response 不能为空")
	}

	// 1. 解析 RelayState 获取 sourceNo
	var relayData map[string]string
	if err := json.Unmarshal([]byte(ct.RelayState), &relayData); err != nil {
		log.Errorf(ctx, log.TagAppDef, "解析 RelayState 失败: %v", err)
		return rt.ErrorMessage("RelayState 解析失败")
	}
	sourceNo := relayData["sourceNo"]
	if strPg.IsBlank(sourceNo) {
		return rt.ErrorMessage("RelayState 中缺少 sourceNo")
	}

	// 2. 查询认证源
	source, found := c.daoSource.FindByNo(ctx, sourceNo)
	if !found {
		return rt.ErrorMessage("认证源不存在")
	}
	if !enumStatePg.ENABLE.IsExistInt8(source.State) {
		return rt.ErrorMessage("认证源已禁用")
	}

	// 3. 构建 SP
	sp, _, err := c.buildSpFromSource(ctx, source)
	if err != nil {
		log.Errorf(ctx, log.TagAppDef, "构建 SAML SP 失败: %v", err)
		return rt.ErrorMessage("构建 SAML SP 失败: " + err.Error())
	}

	// 4. 解析 SAML Response
	samlUserInfo, err := samlPkg.ParseSamlResponse(ct.SAMLResponse, sp)
	if err != nil {
		log.Errorf(ctx, log.TagAppDef, "解析 SAML Response 失败: %v", err)
		return rt.ErrorMessage("解析 SAML 响应失败: " + err.Error())
	}
	log.Infof(ctx, log.TagAppDef, "SAML userInfo: NameID=%s, attrs=%v", samlUserInfo.NameID, samlUserInfo.Attributes)

	// 5. 解析属性映射
	userInfo := c.mapSamlAttributes(samlUserInfo, source.AttributeMapping)

	// 6. 查找或创建绑定
	return c.parseAndBindUser(ctx, source, samlUserInfo.NameID, userInfo)
}

// buildSpFromSource 从认证源配置构建 SAML SP
func (c *SamlLoginService) buildSpFromSource(ctx *gin.Context, source *entityRam.RamIdentitySourceEntity) (*samlPkg.SAMLServiceProvider, *samlPkg.SpConfig, error) {
	// 解析 BaseConfig
	var srcCfg samlSourceConfig
	if strPg.IsNotBlank(source.BaseConfig) {
		if err := json.Unmarshal([]byte(source.BaseConfig), &srcCfg); err != nil {
			return nil, nil, fmt.Errorf("解析 BaseConfig 失败: %w", err)
		}
	}

	// 读取 IdP 证书（从 MetadataCache 或 Credential）
	idpCert := ""
	metadataCache, found := c.daoMetadata.FindBySourceNo(ctx, source.No)
	if found && metadataCache != nil {
		idpCert = metadataCache.MetadataRaw
	}

	// 如果元数据缓存中没有证书，尝试从 Credential 读取
	if idpCert == "" {
		credCert, foundCred := c.daoCredential.FindBySourceNoAndCredType(ctx, source.No, "sp_cert")
		if foundCred && credCert != nil {
			idpCert = credCert.Value
		}
	}

	// 读取 SP 证书和私钥
	spCert := ""
	spKey := ""
	credSpCert, foundSpCert := c.daoCredential.FindBySourceNoAndCredType(ctx, source.No, "sp_cert")
	if foundSpCert && credSpCert != nil {
		spCert = credSpCert.Value
	}
	credSpKey, foundSpKey := c.daoCredential.FindBySourceNoAndCredType(ctx, source.No, "sp_private_key")
	if foundSpKey && credSpKey != nil {
		spKey = credSpKey.Value
	}

	// 构建 ACS URL
	acsURL := srcCfg.SpIssuer
	if acsURL == "" {
		acsURL = source.No
	}

	cfg := &samlPkg.SpConfig{
		Issuer:     acsURL,
		AcsURL:     acsURL,
		IdpSsoURL:  srcCfg.IdpSsoURL,
		IdpIssuer:  srcCfg.IdpIssuer,
		EnableSign: srcCfg.EnableSign,
		IdpCert:    idpCert,
		SpCert:     spCert,
		SpKey:      spKey,
	}

	sp, err := samlPkg.BuildServiceProvider(cfg)
	if err != nil {
		return nil, nil, err
	}

	return sp, cfg, nil
}

// mapSamlAttributes 将 SAML 属性映射为 UserInfo
func (c *SamlLoginService) mapSamlAttributes(samlInfo *samlPkg.SpUserInfo, mappingJson string) *idp.UserInfo {
	userInfo := &idp.UserInfo{
		Id:    samlInfo.NameID,
		Extra: make(map[string]string),
	}

	// 默认属性映射
	attrMap := map[string]string{
		"username":    "username",
		"displayName": "displayName",
		"email":       "email",
		"phone":       "phone",
		"avatar":      "avatar",
	}

	// 如果配置了自定义映射，覆盖默认映射
	if strPg.IsNotBlank(mappingJson) {
		var customMap map[string]string
		if err := json.Unmarshal([]byte(mappingJson), &customMap); err == nil {
			for k, v := range customMap {
				attrMap[k] = v
			}
		}
	}

	// 应用映射
	for localField, samlAttr := range attrMap {
		if val, ok := samlInfo.Attributes[samlAttr]; ok {
			switch localField {
			case "username":
				userInfo.Username = val
			case "displayName":
				userInfo.DisplayName = val
			case "email":
				userInfo.Email = val
			case "phone":
				userInfo.Phone = val
			case "avatar":
				userInfo.AvatarUrl = val
			}
		}
	}

	// Fallback: 如果 Username 为空，使用 Email 或 NameID
	if userInfo.Username == "" {
		if userInfo.Email != "" {
			userInfo.Username = userInfo.Email
		} else if userInfo.Id != "" {
			userInfo.Username = userInfo.Id
		}
	}

	return userInfo
}

// parseAndBindUser 解析用户信息并查找/创建绑定
func (c *SamlLoginService) parseAndBindUser(
	ctx *gin.Context,
	source *entityRam.RamIdentitySourceEntity,
	nameID string,
	userInfo *idp.UserInfo,
) (rt rg.Rs[modRamLogin.IdpLoginSuccess]) {
	// 查找绑定（按 ExternalSub = NameID）
	binding, bindingFound := c.daoBinding.FindByIdpAndExternalSub(ctx, source.Idp, nameID)

	isSignup := false

	if bindingFound && binding != nil {
		// 已绑定 → 获取绑定账号
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

		// 更新 binding 的登录时间
		now := time.Now()
		c.daoBinding.Update(ctx, entityRam.RamIdpBindingEntity{
			LastLoginTime: &now,
			NickName:      userInfo.DisplayName,
			Avatar:        userInfo.AvatarUrl,
			Mail:          userInfo.Email,
			Phone:         userInfo.Phone,
		}, binding.ID)

		return c.checkMfaOrLoginSuccess(ctx, account, source, userInfo, isSignup)
	}

	// 未绑定
	// 检查是否允许自动创建
	if source.AutoCreateUser != 1 {
		return rt.ErrorMessage("该认证源不允许自动创建账号，请先绑定已有账号")
	}

	// 自动创建账号
	isSignup = true
	account, err := c.createAccountAndBinding(ctx, source, userInfo, nameID)
	if err != nil {
		return rt.ErrorMessage("自动创建账号失败: " + err.Error())
	}

	return c.checkMfaOrLoginSuccess(ctx, account, source, userInfo, isSignup)
}

// createAccountAndBinding 自动创建账号和绑定记录
func (c *SamlLoginService) createAccountAndBinding(
	ctx *gin.Context,
	source *entityRam.RamIdentitySourceEntity,
	userInfo *idp.UserInfo,
	nameID string,
) (*entityRam.RamAccountEntity, error) {
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
	if strPg.IsBlank(account.Account) {
		account.Account = nameID
	}
	if strPg.IsBlank(account.Name) {
		account.Name = account.Account
	}

	err, _ := c.daoAccount.Create(ctx, account)
	if err != nil {
		log.Errorf(ctx, log.TagAppDef, "SAML 创建账号失败: %v", err)
		return nil, err
	}
	log.Infof(ctx, log.TagAppDef, "SAML 自动创建账号: no=%s, account=%s", account.No, account.Account)

	// 创建 RamIdpBinding
	bindingNo := noPg.No()
	binding := &entityRam.RamIdpBindingEntity{
		No:            bindingNo,
		TenantNo:      source.TenantNo,
		OrgNo:         source.OrgNo,
		StoreNo:       source.StoreNo,
		TypeDomain:    typeDomainPg.Manage.String(),
		Idp:           source.Idp,
		ExternalSub:   nameID,
		BindAno:       accountNo,
		State:         int8(enumStatePg.ENABLE),
		StateBind:     2, // 已绑定
		BindTime:      &now,
		LastLoginTime: &now,
		Protocol:      source.Protocol,
		Mail:          userInfo.Email,
		Phone:         userInfo.Phone,
		NickName:      userInfo.DisplayName,
		Avatar:        userInfo.AvatarUrl,
		CreateAt:      &now,
	}
	errBind, _ := c.daoBinding.Create(ctx, binding)
	if errBind != nil {
		log.Errorf(ctx, log.TagAppDef, "SAML 创建绑定记录失败: %v", errBind)
		return nil, errBind
	}

	return account, nil
}

// checkMfaOrLoginSuccess 检查是否需要 MFA，如果需要则返回 MFA 令牌，否则直接登录成功
func (c *SamlLoginService) checkMfaOrLoginSuccess(
	ctx *gin.Context,
	account *entityRam.RamAccountEntity,
	source *entityRam.RamIdentitySourceEntity,
	userInfo *idp.UserInfo,
	isSignup bool,
) (rt rg.Rs[modRamLogin.IdpLoginSuccess]) {
	// 新注册用户跳过 MFA
	if isSignup {
		return c.loginSuccess(ctx, account, source, userInfo, isSignup)
	}

	// 检查是否需要 MFA
	mfaRequired, mfaToken, mfaType := c.mfaService.CheckMfaRequired(ctx, account.No)
	if mfaRequired {
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

	return c.loginSuccess(ctx, account, source, userInfo, isSignup)
}

// loginSuccess 登录成功，生成 token 并记录审计日志
func (c *SamlLoginService) loginSuccess(
	ctx *gin.Context,
	account *entityRam.RamAccountEntity,
	source *entityRam.RamIdentitySourceEntity,
	userInfo *idp.UserInfo,
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

// saveSessionLog 保存 SAML 登录审计日志
func (c *SamlLoginService) saveSessionLog(
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
		LoginSource:   "saml",
	}
	err, _ := c.daoSessionLog.Create(ctx, logEntity)
	if err != nil {
		log.Errorf(ctx, log.TagAppDef, "保存 SAML 登录审计日志失败: %v", err)
	}
}
