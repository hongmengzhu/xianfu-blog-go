package service

import (
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamLdap"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamLogin"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityRam"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/auth/holderPg/multiTenantPg"
	ldap2 "github.com/hongmengzhu/xianfu-blog-go/pkg/auth/ldap"
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
	gs.Provide(new(LdapService))
}

// LdapService LDAP 认证+同步 Service
type LdapService struct {
	daoSource     *repositoryRam.RamIdentitySourceRepository    `autowire:"?"`
	daoProvider   *repositoryRam.RamIdentityProviderRepository  `autowire:"?"`
	daoAccount    *repositoryRam.RamAccountRepository           `autowire:"?"`
	daoBinding    *repositoryRam.RamIdpBindingRepository        `autowire:"?"`
	daoSessionLog *repositoryRam.RamAccountSessionLogRepository `autowire:"?"`
	loginService  *AccountLoginService                          `autowire:"?"`
}

// buildLdapConfig 从认证源配置构建 LDAP 连接配置
func (s *LdapService) buildLdapConfig(source *entityRam.RamIdentitySourceEntity) (*ldap2.LdapConfig, error) {
	var baseCfg idpBaseConfig
	if strPg.IsNotBlank(source.BaseConfig) {
		if err := json.Unmarshal([]byte(source.BaseConfig), &baseCfg); err != nil {
			return nil, err
		}
	}

	// 从 BaseConfig 解析 LDAP 特有字段
	var ldapExtra struct {
		Filter          string   `json:"filter"`
		FilterFields    []string `json:"filterFields"`
		BaseDn          string   `json:"baseDn"`
		Port            int      `json:"port"`
		EnableSsl       bool     `json:"enableSsl"`
		AllowSelfSigned bool     `json:"allowSelfSignedCert"`
	}
	// 尝试从 SyncStrategyConfig 解析 LDAP 额外配置
	if strPg.IsNotBlank(source.SyncStrategyConfig) {
		_ = json.Unmarshal([]byte(source.SyncStrategyConfig), &ldapExtra)
	}

	if ldapExtra.BaseDn == "" {
		ldapExtra.BaseDn = baseCfg.HostUrl // fallback
	}
	if ldapExtra.Port == 0 {
		if ldapExtra.EnableSsl {
			ldapExtra.Port = 636
		} else {
			ldapExtra.Port = 389
		}
	}

	return &ldap2.LdapConfig{
		Host:                baseCfg.HostUrl,
		Port:                ldapExtra.Port,
		EnableSsl:           ldapExtra.EnableSsl,
		AllowSelfSignedCert: ldapExtra.AllowSelfSigned,
		Username:            baseCfg.ClientId,
		Password:            baseCfg.ClientSecret,
		BaseDn:              ldapExtra.BaseDn,
		Filter:              ldapExtra.Filter,
		FilterFields:        ldapExtra.FilterFields,
	}, nil
}

// TestConnection 测试 LDAP 连接
func (s *LdapService) TestConnection(ctx *gin.Context, ct modRamLdap.LdapTestCt) (rt rg.Rs[modRamLdap.LdapTestVo]) {
	source, found := s.daoSource.FindByNo(ctx, ct.SourceNo)
	if !found {
		return rt.ErrorMessage("认证源不存在")
	}

	cfg, err := s.buildLdapConfig(source)
	if err != nil {
		return rt.ErrorMessage("配置解析失败: " + err.Error())
	}

	conn, err := ldap2.GetLdapConn(cfg)
	if err != nil {
		return rt.OkData(modRamLdap.LdapTestVo{
			Connected: false,
			Message:   "连接失败: " + err.Error(),
		})
	}
	defer conn.Close()

	// 统计用户数量
	users, err := conn.SearchUsers(cfg)
	if err != nil {
		return rt.OkData(modRamLdap.LdapTestVo{
			Connected: true,
			IsAD:      conn.IsAD,
			Message:   "连接成功，但搜索用户失败: " + err.Error(),
		})
	}

	return rt.OkData(modRamLdap.LdapTestVo{
		Connected: true,
		IsAD:      conn.IsAD,
		Message:   "连接成功",
		UserCount: len(users),
	})
}

// SearchUsers 搜索 LDAP 用户
func (s *LdapService) SearchUsers(ctx *gin.Context, ct modRamLdap.LdapSearchCt) (rt rg.Rs[[]modRamLdap.LdapUserVo]) {
	source, found := s.daoSource.FindByNo(ctx, ct.SourceNo)
	if !found {
		return rt.ErrorMessage("认证源不存在")
	}

	cfg, err := s.buildLdapConfig(source)
	if err != nil {
		return rt.ErrorMessage("配置解析失败: " + err.Error())
	}

	if strPg.IsNotBlank(ct.Filter) {
		cfg.Filter = ct.Filter
	}

	conn, err := ldap2.GetLdapConn(cfg)
	if err != nil {
		return rt.ErrorMessage("LDAP 连接失败: " + err.Error())
	}
	defer conn.Close()

	users, err := conn.SearchUsers(cfg)
	if err != nil {
		return rt.ErrorMessage("搜索用户失败: " + err.Error())
	}

	var result []modRamLdap.LdapUserVo
	for _, u := range users {
		result = append(result, modRamLdap.LdapUserVo{
			Uid:         u.Uid,
			Cn:          u.Cn,
			DisplayName: u.DisplayName,
			Email:       u.GetEmail(),
			Mobile:      u.Mobile,
			Uuid:        u.GetLdapUuid(),
			MemberOf:    u.MemberOf,
		})
	}

	return rt.OkData(result)
}

// SyncUsers 同步 LDAP 用户
func (s *LdapService) SyncUsers(ctx *gin.Context, ct modRamLdap.LdapSyncCt) (rt rg.Rs[modRamLdap.LdapSyncVo]) {
	source, found := s.daoSource.FindByNo(ctx, ct.SourceNo)
	if !found {
		return rt.ErrorMessage("认证源不存在")
	}

	cfg, err := s.buildLdapConfig(source)
	if err != nil {
		return rt.ErrorMessage("配置解析失败: " + err.Error())
	}

	conn, err := ldap2.GetLdapConn(cfg)
	if err != nil {
		return rt.ErrorMessage("LDAP 连接失败: " + err.Error())
	}
	defer conn.Close()

	syncResult, err := ldap2.SyncUsers(conn, cfg)
	if err != nil {
		return rt.ErrorMessage("同步失败: " + err.Error())
	}

	return rt.OkData(modRamLdap.LdapSyncVo{
		NewUsers:     syncResult.NewUsers,
		UpdatedUsers: syncResult.UpdatedUsers,
		FailedUsers:  syncResult.FailedUsers,
		Errors:       syncResult.Errors,
	})
}

// Login LDAP 登录认证
func (s *LdapService) Login(ctx *gin.Context, ct modRamLdap.LdapLoginCt) (rt rg.Rs[modRamLogin.IdpLoginSuccess]) {
	source, found := s.daoSource.FindByNo(ctx, ct.SourceNo)
	if !found {
		return rt.ErrorMessage("认证源不存在")
	}
	if !enumStatePg.ENABLE.IsExistInt8(source.State) {
		return rt.ErrorMessage("认证源已禁用")
	}

	cfg, err := s.buildLdapConfig(source)
	if err != nil {
		return rt.ErrorMessage("配置解析失败: " + err.Error())
	}

	conn, err := ldap2.GetLdapConn(cfg)
	if err != nil {
		return rt.ErrorMessage("LDAP 连接失败: " + err.Error())
	}
	defer conn.Close()

	// 验证用户密码
	ldapUser, err := conn.CheckUserPassword(cfg, ct.Username, ct.Password)
	if err != nil {
		log.Errorf(ctx, log.TagAppDef, "LDAP 用户 %s 认证失败: %v", ct.Username, err)
		return rt.ErrorMessage("LDAP 认证失败: " + err.Error())
	}

	log.Infof(ctx, log.TagAppDef, "LDAP 用户认证成功: uid=%s, cn=%s", ldapUser.Uid, ldapUser.Cn)

	// 查找绑定
	binding, bindingFound := s.daoBinding.FindByIdpAndExternalSub(ctx, source.Idp, ldapUser.GetLdapUuid())

	isSignup := false

	if bindingFound && binding != nil && strPg.IsNotBlank(binding.BindAno) {
		// 已绑定 → 获取账号
		account, accountFound := s.daoAccount.FindByNo(ctx, binding.BindAno)
		if !accountFound {
			return rt.ErrorMessage("绑定的账号不存在")
		}
		if !enumStatePg.ENABLE.IsExistInt8(account.State) {
			return rt.ErrorMessage("账号已被禁用")
		}

		// 更新 binding 登录时间
		now := time.Now()
		s.daoBinding.Update(ctx, entityRam.RamIdpBindingEntity{
			LastLoginTime: &now,
		}, binding.ID)

		return s.loginSuccess(ctx, account, source, ldapUser, isSignup)
	}

	// 未绑定 → 检查是否允许自动创建
	if source.AutoCreateUser != 1 {
		return rt.ErrorMessage("该认证源不允许自动创建账号，请先绑定已有账号")
	}

	// 自动创建账号
	isSignup = true
	account, errCreate := s.createAccountFromLdap(ctx, source, ldapUser)
	if errCreate != nil {
		return rt.ErrorMessage("自动创建账号失败: " + errCreate.Error())
	}

	return s.loginSuccess(ctx, account, source, ldapUser, isSignup)
}

// createAccountFromLdap 从 LDAP 用户创建本地账号
func (s *LdapService) createAccountFromLdap(
	ctx *gin.Context,
	source *entityRam.RamIdentitySourceEntity,
	ldapUser *ldap2.LdapUser,
) (*entityRam.RamAccountEntity, error) {
	now := time.Now()
	accountNo := noPg.No()

	username := ldapUser.Uid
	if strPg.IsBlank(username) {
		username = ldapUser.Cn
	}
	if strPg.IsBlank(username) {
		username = ldapUser.GetLdapUuid()
	}

	account := &entityRam.RamAccountEntity{
		No:           accountNo,
		TenantNo:     source.TenantNo,
		OrgNo:        source.OrgNo,
		StoreNo:      source.StoreNo,
		TypeDomain:   typeDomainPg.Manage.String(),
		Account:      username,
		Name:         ldapUser.GetDisplayName(),
		RealName:     ldapUser.GetDisplayName(),
		Mail:         ldapUser.GetEmail(),
		Phone:        ldapUser.Mobile,
		State:        int8(enumStatePg.ENABLE),
		RegisterTime: &now,
		RegisterIP:   ctx.ClientIP(),
		LoginTime:    &now,
		CreateAt:     &now,
	}
	if strPg.IsBlank(account.Name) {
		account.Name = account.Account
	}

	errAcc, _ := s.daoAccount.Create(ctx, account)
	if errAcc != nil {
		return nil, errAcc
	}

	// 创建绑定记录
	bindingNo := noPg.No()
	binding := &entityRam.RamIdpBindingEntity{
		No:            bindingNo,
		TenantNo:      source.TenantNo,
		OrgNo:         source.OrgNo,
		StoreNo:       source.StoreNo,
		TypeDomain:    typeDomainPg.Manage.String(),
		Idp:           source.Idp,
		ExternalSub:   ldapUser.GetLdapUuid(),
		BindAno:       accountNo,
		State:         int8(enumStatePg.ENABLE),
		StateBind:     2,
		BindTime:      &now,
		LastLoginTime: &now,
		Protocol:      "ldap",
		Mail:          ldapUser.GetEmail(),
		Phone:         ldapUser.Mobile,
		NickName:      ldapUser.GetDisplayName(),
		CreateAt:      &now,
	}
	errBind, _ := s.daoBinding.Create(ctx, binding)
	if errBind != nil {
		log.Errorf(ctx, log.TagAppDef, "创建 LDAP 绑定记录失败: %v", errBind)
	}

	log.Infof(ctx, log.TagAppDef, "LDAP 自动创建账号: no=%s, account=%s", account.No, account.Account)
	return account, nil
}

// loginSuccess LDAP 登录成功
func (s *LdapService) loginSuccess(
	ctx *gin.Context,
	account *entityRam.RamAccountEntity,
	source *entityRam.RamIdentitySourceEntity,
	ldapUser *ldap2.LdapUser,
	isSignup bool,
) (rt rg.Rs[modRamLogin.IdpLoginSuccess]) {
	mult := multiTenantPg.MultiTenantPg{
		TenantNo: []string{account.TenantNo},
	}
	tokenResult := s.loginService.MakeToken(ctx, mult, account, typeDomainPg.Manage, clientPg.Browser, true)
	if tokenResult.ErrorIs() {
		return rt.ErrorMessage(tokenResult.Message)
	}
	dataToken := tokenResult.Data

	// 写入审计日志
	now := time.Now()
	logEntity := &entityRam.RamAccountSessionLogEntity{
		No:            noPg.No(),
		TenantNo:      account.TenantNo,
		OrgNo:         account.OrgNo,
		StoreNo:       account.StoreNo,
		TypeDomain:    typeDomainPg.Manage.String(),
		Idp:           source.Idp,
		SourceNo:      source.No,
		Protocol:      "ldap",
		Ano:           account.No,
		BindAno:       account.No,
		ExternalSub:   ldapUser.GetLdapUuid(),
		EventCategory: "login",
		EventType:     "login_success",
		EventResult:   1,
		IpAddress:     ctx.ClientIP(),
		UserAgent:     ctx.GetHeader(constHeaderPg.HeaderUserAgent),
		OperateAt:     &now,
		LoginAt:       &now,
		Ip:            ctx.ClientIP(),
		LoginSource:   "ldap",
	}
	errLog, _ := s.daoSessionLog.Create(ctx, logEntity)
	if errLog != nil {
		log.Errorf(ctx, log.TagAppDef, "保存 LDAP 登录审计日志失败: %v", errLog)
	}

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
