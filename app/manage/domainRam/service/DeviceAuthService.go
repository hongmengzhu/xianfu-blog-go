package service

import (
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	modRamDeviceAuth2 "github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamDeviceAuth"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryRam"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/consts/constsRam/typeDomainPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/deviceauth"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/clientPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/holderPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/holderPg/multiTenantPg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(NewDeviceAuthService)
}

// DeviceAuthService 设备授权（RFC 8628）
type DeviceAuthService struct {
	dao     *repositoryRam.RamAccountRepository `autowire:"?"`
	loginSv *AccountLoginService                `autowire:"?"`
	store   deviceauth.Store
}

func NewDeviceAuthService() *DeviceAuthService {
	return &DeviceAuthService{
		store: deviceauth.NewMemoryStore(),
	}
}

// generateDeviceCode 生成 32 位随机设备码
func generateDeviceCode() string {
	return strPg.GenerateNumberId22() + strPg.GenerateNumberId22()[:10]
}

// generateUserCode 生成 6 位随机用户码（大写字母）
func generateUserCode() string {
	const letters = "ABCDEFGHJKLMNPQRSTUVWXYZ"
	b := make([]byte, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// generateCancelToken 生成取消令牌
func generateCancelToken() string {
	return strPg.GenerateNumberId22()
}

// StartDeviceAuth 发起设备授权
func (c *DeviceAuthService) StartDeviceAuth(ctx *gin.Context, ct modRamDeviceAuth2.DeviceAuthRequestCt) (rt rg.Rs[modRamDeviceAuth2.DeviceAuthVo]) {
	log.Infof(ctx, log.TagAppDef, "StartDeviceAuth ct=%+v", ct)

	// 生成 deviceCode，确保唯一
	deviceCode := generateDeviceCode()
	userCode := generateUserCode()
	retryCount := 0
	for retryCount < 5 {
		if _, ok := c.store.LoadByUserCode(userCode); !ok {
			break
		}
		userCode = generateUserCode()
		retryCount++
	}
	if retryCount >= 5 {
		return rt.ErrorMessage("生成用户码失败，请重试")
	}

	cancelToken := generateCancelToken()
	expiresIn := deviceauth.DefaultExpiresIn
	now := time.Now()

	// 按 deviceCode 存储
	deviceCache := deviceauth.DeviceAuthCache{
		DeviceCode: deviceCode,
		UserCode:   userCode,
		ClientId:   ct.ClientId,
		Scope:      ct.Scope,
		RequestAt:  now,
		Status:     deviceauth.StatusPending,
		ExpiresIn:  expiresIn,
	}
	// 按 userCode 存储（含 cancelToken）
	userCache := deviceauth.DeviceAuthCache{
		DeviceCode:  deviceCode,
		UserCode:    userCode,
		ClientId:    ct.ClientId,
		Scope:       ct.Scope,
		RequestAt:   now,
		Status:      deviceauth.StatusPending,
		CancelToken: cancelToken,
		ExpiresIn:   expiresIn,
	}

	c.store.Save(deviceCode, deviceCache)
	c.store.Save(userCode, userCache)

	vo := modRamDeviceAuth2.DeviceAuthVo{
		DeviceCode:      deviceCode,
		UserCode:        userCode,
		VerificationUri: "/device/verify",
		ExpiresIn:       expiresIn,
		Interval:        deviceauth.DefaultInterval,
		CancelToken:     cancelToken,
	}
	return rt.OkData(vo)
}

// PollDeviceStatus 设备轮询状态
func (c *DeviceAuthService) PollDeviceStatus(ctx *gin.Context, ct modRamDeviceAuth2.DevicePollCt) (rt rg.Rs[modRamDeviceAuth2.DeviceStatusVo]) {
	cache, ok := c.store.LoadByDeviceCode(ct.DeviceCode)
	if !ok {
		return rt.OkData(modRamDeviceAuth2.DeviceStatusVo{Status: "expired"})
	}
	if cache.IsExpired() {
		c.store.Delete(ct.DeviceCode)
		return rt.OkData(modRamDeviceAuth2.DeviceStatusVo{Status: "expired"})
	}

	switch cache.Status {
	case deviceauth.StatusPending:
		return rt.OkData(modRamDeviceAuth2.DeviceStatusVo{Status: "pending"})
	case deviceauth.StatusDenied:
		c.store.Delete(ct.DeviceCode)
		return rt.OkData(modRamDeviceAuth2.DeviceStatusVo{Status: "denied"})
	case deviceauth.StatusApproved:
		// 签发 Token
		return c.issueTokenForDevice(ctx, cache)
	case deviceauth.StatusTokenIssued:
		// 已签发，直接返回 token
		return c.issueTokenForDevice(ctx, cache)
	default:
		return rt.OkData(modRamDeviceAuth2.DeviceStatusVo{Status: cache.Status})
	}
}

// issueTokenForDevice 为已授权的设备签发 Token
func (c *DeviceAuthService) issueTokenForDevice(ctx *gin.Context, cache deviceauth.DeviceAuthCache) (rt rg.Rs[modRamDeviceAuth2.DeviceStatusVo]) {
	if cache.UserNo == "" {
		return rt.OkData(modRamDeviceAuth2.DeviceStatusVo{Status: "pending"})
	}

	account, found := c.dao.FindByNo(ctx, cache.UserNo)
	if !found || account == nil {
		return rt.ErrorMessage("用户不存在")
	}
	if !enumStatePg.ENABLE.IsExistInt8(account.State) {
		return rt.ErrorMessage("账户已被禁用")
	}

	mult := multiTenantPg.MultiTenantPg{
		TenantNo: []string{account.TenantNo},
	}
	tokenResult := c.loginSv.MakeToken(ctx, mult, account, typeDomainPg.System, clientPg.Browser, true)
	if tokenResult.ErrorIs() {
		return rt.ErrorMessage(tokenResult.Message)
	}
	data := tokenResult.Data

	// 更新状态为 token_issued
	cache.Status = deviceauth.StatusTokenIssued
	c.store.Update(cache)

	return rt.OkData(modRamDeviceAuth2.DeviceStatusVo{
		Status:       "approved",
		AccessToken:  data.Access,
		RefreshToken: data.Refresh,
	})
}

// ApproveDeviceAuth 用户在浏览器中授权
func (c *DeviceAuthService) ApproveDeviceAuth(ctx *gin.Context, ct modRamDeviceAuth2.DeviceApproveCt) (rt rg.Rs[string]) {
	cache, ok := c.store.LoadByUserCode(ct.UserCode)
	if !ok {
		return rt.ErrorMessage("用户码无效或已过期")
	}
	if cache.IsExpired() {
		c.store.Delete(cache.DeviceCode)
		return rt.ErrorMessage("设备授权已过期")
	}
	if cache.Status != deviceauth.StatusPending {
		return rt.ErrorMessage("该设备授权已处理")
	}

	// 从上下文中获取当前登录用户
	holder := holderPg.GetContextAccount(ctx)
	if holder.GetAccountNo() == "" {
		return rt.ErrorMessage("请先登录")
	}
	accountHolder := holder.GetAccount()

	cache.Status = deviceauth.StatusApproved
	cache.UserName = accountHolder.Account
	cache.UserNo = accountHolder.No
	cache.TenantNo = accountHolder.TenantNo
	c.store.Update(cache)

	log.Infof(ctx, log.TagAppDef, "ApproveDeviceAuth userCode=%s, user=%s", ct.UserCode, accountHolder.Account)
	return rt.OkMessage("授权成功")
}

// CancelDeviceAuth 取消设备授权
func (c *DeviceAuthService) CancelDeviceAuth(ctx *gin.Context, ct modRamDeviceAuth2.DeviceCancelCt) (rt rg.Rs[string]) {
	cache, ok := c.store.LoadByUserCode(ct.UserCode)
	if !ok {
		return rt.ErrorMessage("用户码无效")
	}
	if cache.CancelToken == "" || cache.CancelToken != ct.CancelToken {
		return rt.ErrorMessage("取消令牌不匹配")
	}

	c.store.Delete(cache.DeviceCode)
	log.Infof(ctx, log.TagAppDef, "CancelDeviceAuth userCode=%s", ct.UserCode)
	return rt.OkMessage("已取消")
}

// CompleteDeviceAuth 完成设备授权（建立会话）
func (c *DeviceAuthService) CompleteDeviceAuth(ctx *gin.Context, ct modRamDeviceAuth2.DeviceCompleteCt) (rt rg.Rs[modRamDeviceAuth2.DeviceStatusVo]) {
	cache, ok := c.store.LoadByDeviceCode(ct.DeviceCode)
	if !ok {
		return rt.ErrorMessage("设备码无效")
	}
	if cache.IsExpired() {
		c.store.Delete(ct.DeviceCode)
		return rt.ErrorMessage("设备授权已过期")
	}
	if cache.Status != deviceauth.StatusTokenIssued {
		return rt.ErrorMessage("设备授权尚未完成 Token 签发")
	}

	// 清理缓存
	c.store.Delete(ct.DeviceCode)

	return rt.OkData(modRamDeviceAuth2.DeviceStatusVo{
		Status: "completed",
	})
}
