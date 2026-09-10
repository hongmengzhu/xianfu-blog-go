package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamLogin"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamMfa"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/entityRam"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constsRam/typeDomainPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/clientPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/holderPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/holderPg/multiTenantPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/mfa"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(new(AccountMfaService))
}

// AccountMfaService MFA 管理服务
type AccountMfaService struct {
	daoMfa       *repositoryRam.RamAccountAuthMfaRepository `autowire:"?"`
	daoAccount   *repositoryRam.RamAccountRepository        `autowire:"?"`
	loginService *AccountLoginService                       `autowire:"?"`
}

// Setup 初始化 MFA 设置（生成 secret、恢复码）
func (c *AccountMfaService) Setup(ctx *gin.Context, ct modRamMfa.SetupCt) (rt rg.Rs[modRamMfa.MfaSetupVo]) {
	holder := holderPg.GetContextAccount(ctx)
	ano := holder.GetAccountNo()

	mfaUtil := mfa.GetMfaUtil(ct.MfaType, nil)
	if mfaUtil == nil {
		return rt.ErrorMessage("不支持的 MFA 类型")
	}

	// 初始化 MFA（生成 secret 和 URL）
	props, err := mfaUtil.Initiate(ano, "XianfuBlog")
	if err != nil {
		log.Errorf(ctx, log.TagAppDef, "MFA Initiate 失败: %v", err)
		return rt.ErrorMessage("初始化 MFA 失败: " + err.Error())
	}

	// 生成恢复码
	recoveryCodes := mfa.GenerateRecoveryCodes(8)

	vo := modRamMfa.MfaSetupVo{
		MfaType:       props.MfaType,
		Secret:        props.Secret,
		URL:           props.URL,
		RecoveryCodes: recoveryCodes,
	}
	return rt.OkData(vo)
}

// SetupVerify 验证 MFA 设置（确认用户已正确配置）
func (c *AccountMfaService) SetupVerify(ctx *gin.Context, ct modRamMfa.SetupVerifyCt) (rt rg.Rs[string]) {
	config := &mfa.MfaProps{
		MfaType: ct.MfaType,
		Secret:  ct.Secret,
	}
	mfaUtil := mfa.GetMfaUtil(ct.MfaType, config)
	if mfaUtil == nil {
		return rt.ErrorMessage("不支持的 MFA 类型")
	}

	err := mfaUtil.SetupVerify(ct.Passcode)
	if err != nil {
		return rt.ErrorMessage("验证失败: " + err.Error())
	}
	return rt.OkMessage("验证成功")
}

// Enable 启用 MFA（写入数据库）
func (c *AccountMfaService) Enable(ctx *gin.Context, ct modRamMfa.EnableCt) (rt rg.Rs[string]) {
	holder := holderPg.GetContextAccount(ctx)
	ano := holder.GetAccountNo()
	tenantNo := holder.GetTenantNo()

	if strPg.IsBlank(ct.Secret) {
		return rt.ErrorMessage("密钥不能为空")
	}

	// 序列化恢复码
	recoveryCodesJson, err := json.Marshal(ct.RecoveryCodes)
	if err != nil {
		return rt.ErrorMessage("恢复码序列化失败")
	}

	// 生成 credential ID
	credIdBytes := make([]byte, 16)
	_, _ = rand.Read(credIdBytes)
	credId := hex.EncodeToString(credIdBytes)

	now := time.Now()
	entity := &entityRam.RamAccountAuthMfaEntity{
		TenantNo:      tenantNo,
		Ano:           ano,
		MfaType:       ct.MfaType,
		Secret:        ct.Secret,
		CredentialID:  credId,
		DeviceName:    "Authenticator",
		RecoveryCodes: string(recoveryCodesJson),
		State:         int8(enumStatePg.ENABLE),
		CreateAt:      &now,
	}

	errCreate, _ := c.daoMfa.Create(ctx, entity)
	if errCreate != nil {
		log.Errorf(ctx, log.TagAppDef, "创建 MFA 记录失败: %v", errCreate)
		return rt.ErrorMessage("启用 MFA 失败: " + errCreate.Error())
	}

	return rt.OkMessage("MFA 已启用")
}

// Disable 禁用 MFA
func (c *AccountMfaService) Disable(ctx *gin.Context, ct modRamMfa.DisableCt) (rt rg.Rs[string]) {
	holder := holderPg.GetContextAccount(ctx)
	ano := holder.GetAccountNo()

	// 查找该账号的所有 MFA 记录
	mfaRecords, found := c.daoMfa.FindByAno(ctx, ano)
	if !found || len(mfaRecords) == 0 {
		return rt.ErrorMessage("未找到 MFA 记录")
	}

	// 如果有 TOTP 类型的记录，需要先验证 passcode
	for _, record := range mfaRecords {
		if record.MfaType == mfa.TotpType {
			if strPg.IsBlank(ct.Passcode) {
				return rt.ErrorMessage("请输入当前验证码以确认禁用")
			}
			mfaUtil := mfa.GetMfaUtil(mfa.TotpType, &mfa.MfaProps{
				MfaType: mfa.TotpType,
				Secret:  record.Secret,
			})
			if err := mfaUtil.Verify(ct.Passcode); err != nil {
				return rt.ErrorMessage("验证码错误，无法禁用 MFA")
			}
			break
		}
	}

	// 删除所有 MFA 记录
	err := c.daoMfa.DeleteByAno(ctx, ano)
	if err != nil {
		log.Errorf(ctx, log.TagAppDef, "删除 MFA 记录失败: %v", err)
		return rt.ErrorMessage("禁用 MFA 失败: " + err.Error())
	}

	return rt.OkMessage("MFA 已禁用")
}

// Preferred 获取 MFA 偏好设置
func (c *AccountMfaService) Preferred(ctx *gin.Context) (rt rg.Rs[modRamMfa.MfaStatusVo]) {
	holder := holderPg.GetContextAccount(ctx)
	ano := holder.GetAccountNo()

	mfaRecords, found := c.daoMfa.FindByAno(ctx, ano)
	methods := make([]modRamMfa.MfaPropsVo, 0)
	enabled := false
	preferredType := ""

	if found && len(mfaRecords) > 0 {
		enabled = true
		for _, record := range mfaRecords {
			methods = append(methods, modRamMfa.MfaPropsVo{
				Enabled:     true,
				IsPreferred: true,
				MfaType:     record.MfaType,
			})
			if preferredType == "" {
				preferredType = record.MfaType
			}
		}
	}

	vo := modRamMfa.MfaStatusVo{
		Enabled:       enabled,
		PreferredType: preferredType,
		Methods:       methods,
	}
	return rt.OkData(vo)
}

// VerifyLogin 登录时验证 MFA（验证通过后签发 Token）
func (c *AccountMfaService) VerifyLogin(ctx *gin.Context, ct modRamMfa.VerifyCt) (rt rg.Rs[modRamLogin.IdpLoginSuccess]) {
	if strPg.IsBlank(ct.MfaToken) {
		return rt.ErrorMessage("MFA 令牌不能为空")
	}

	// 通过 mfaToken 查找 MFA 记录
	record, found := c.daoMfa.FindByMfaToken(ctx, ct.MfaToken)
	if !found {
		return rt.ErrorMessage("MFA 令牌无效或已过期")
	}

	// 检查令牌是否过期
	if record.MfaTokenExpireAt != nil && time.Now().After(*record.MfaTokenExpireAt) {
		return rt.ErrorMessage("MFA 令牌已过期，请重新登录")
	}

	// 验证 MFA 验证码
	if strPg.IsNotBlank(ct.Passcode) {
		mfaUtil := mfa.GetMfaUtil(record.MfaType, &mfa.MfaProps{
			MfaType: record.MfaType,
			Secret:  record.Secret,
		})
		if mfaUtil == nil {
			return rt.ErrorMessage("不支持的 MFA 类型")
		}
		if err := mfaUtil.Verify(ct.Passcode); err != nil {
			return rt.ErrorMessage("MFA 验证码错误")
		}
	} else {
		return rt.ErrorMessage("请输入验证码")
	}

	// 验证通过，清除 mfaToken
	now := time.Now()
	c.daoMfa.Update(ctx, entityRam.RamAccountAuthMfaEntity{
		MfaToken:         "",
		MfaTokenExpireAt: nil,
		LastUsedAt:       &now,
	}, record.ID)

	// 查找账号信息
	account, accountFound := c.daoAccount.FindByNo(ctx, record.Ano)
	if !accountFound {
		return rt.ErrorMessage("账号不存在")
	}

	// 生成 Token
	mult := multiTenantPg.MultiTenantPg{
		TenantNo: []string{account.TenantNo},
	}
	tokenResult := c.loginService.MakeToken(ctx, mult, account, typeDomainPg.Manage, clientPg.Browser, true)
	if tokenResult.ErrorIs() {
		return rt.ErrorMessage(tokenResult.Message)
	}
	dataToken := tokenResult.Data

	success := modRamLogin.IdpLoginSuccess{
		AccessToken:  dataToken.Access,
		RefreshToken: dataToken.Refresh,
		Info: modRamLogin.LoginSuccessInfo{
			Account: account.Account,
			Name:    account.Name,
			Avatar:  account.Avatar,
			Roles:   []string{},
		},
	}
	return rt.OkData(success)
}

// Recover 使用恢复码恢复
func (c *AccountMfaService) Recover(ctx *gin.Context, ct modRamMfa.RecoverCt) (rt rg.Rs[modRamLogin.IdpLoginSuccess]) {
	if strPg.IsBlank(ct.MfaToken) {
		return rt.ErrorMessage("MFA 令牌不能为空")
	}

	// 通过 mfaToken 查找 MFA 记录
	record, found := c.daoMfa.FindByMfaToken(ctx, ct.MfaToken)
	if !found {
		return rt.ErrorMessage("MFA 令牌无效或已过期")
	}

	// 检查令牌是否过期
	if record.MfaTokenExpireAt != nil && time.Now().After(*record.MfaTokenExpireAt) {
		return rt.ErrorMessage("MFA 令牌已过期，请重新登录")
	}

	// 解析恢复码
	var recoveryCodes []string
	if strPg.IsNotBlank(record.RecoveryCodes) {
		if err := json.Unmarshal([]byte(record.RecoveryCodes), &recoveryCodes); err != nil {
			return rt.ErrorMessage("恢复码数据异常")
		}
	}

	// 验证恢复码
	remaining, err := mfa.VerifyRecoveryCode(recoveryCodes, ct.RecoveryCode)
	if err != nil {
		return rt.ErrorMessage("恢复码验证失败: " + err.Error())
	}

	// 更新剩余恢复码
	remainingJson, _ := json.Marshal(remaining)
	now := time.Now()
	c.daoMfa.Update(ctx, entityRam.RamAccountAuthMfaEntity{
		RecoveryCodes:    string(remainingJson),
		MfaToken:         "",
		MfaTokenExpireAt: nil,
		LastUsedAt:       &now,
	}, record.ID)

	// 查找账号信息
	account, accountFound := c.daoAccount.FindByNo(ctx, record.Ano)
	if !accountFound {
		return rt.ErrorMessage("账号不存在")
	}

	// 生成 Token
	mult := multiTenantPg.MultiTenantPg{
		TenantNo: []string{account.TenantNo},
	}
	tokenResult := c.loginService.MakeToken(ctx, mult, account, typeDomainPg.Manage, clientPg.Browser, true)
	if tokenResult.ErrorIs() {
		return rt.ErrorMessage(tokenResult.Message)
	}
	dataToken := tokenResult.Data

	success := modRamLogin.IdpLoginSuccess{
		AccessToken:  dataToken.Access,
		RefreshToken: dataToken.Refresh,
		Info: modRamLogin.LoginSuccessInfo{
			Account: account.Account,
			Name:    account.Name,
			Avatar:  account.Avatar,
			Roles:   []string{},
		},
	}
	return rt.OkData(success)
}

// CheckMfaRequired 检查账号是否需要 MFA（供登录流程调用）
// 返回: 是否需要 MFA, mfaToken(如果需要), mfaType(首选类型)
func (c *AccountMfaService) CheckMfaRequired(ctx context.Context, ano string) (bool, string, string) {
	records, found := c.daoMfa.FindByAno(ctx, ano)
	if !found || len(records) == 0 {
		return false, "", ""
	}

	// 生成 mfaToken
	tokenBytes := make([]byte, 32)
	_, _ = rand.Read(tokenBytes)
	mfaToken := hex.EncodeToString(tokenBytes)

	// 设置 5 分钟过期
	expireAt := time.Now().Add(5 * time.Minute)

	// 将 mfaToken 写入第一条 MFA 记录
	c.daoMfa.Update(ctx, entityRam.RamAccountAuthMfaEntity{
		MfaToken:         mfaToken,
		MfaTokenExpireAt: &expireAt,
	}, records[0].ID)

	return true, mfaToken, records[0].MfaType
}
