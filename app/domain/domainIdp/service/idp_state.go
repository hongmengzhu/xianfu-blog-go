package service

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/idp/modIdpLogin"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamLogin"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/state/enumStatePg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"golang.org/x/oauth2"
)

const idpStateTTL = 10 * time.Minute

type idpStateEntry struct {
	SourceNo    string
	RedirectUri string
	ExpiresAt   time.Time
}

var idpStates sync.Map

// Authorize 由服务端生成 OAuth state 与授权地址，供前端跳转到 IdP。
// 只有服务端签发的 state 才能在 /idp/login 中通过校验。
func (c *AccountLoginIdpService) Authorize(ctx *gin.Context, ct modRamLogin.IdpLoginCt) (rt rg.Rs[map[string]string]) {
	if strPg.IsBlank(ct.SourceNo) {
		return rt.ErrorMessage("认证源编号不能为空")
	}
	source, found := c.daoSource.FindByNo(ctx, ct.SourceNo)
	if !found {
		return rt.ErrorMessage("认证源不存在")
	}
	if !enumStatePg.ENABLE.IsExistInt8(source.State) {
		return rt.ErrorMessage("认证源已禁用")
	}
	providerEntity, found := c.daoProvider.FindByNo(ctx, source.Idp)
	if !found || !enumStatePg.ENABLE.IsExistInt8(providerEntity.State) {
		return rt.ErrorMessage("身份提供商不可用")
	}
	var baseCfg modIdpLogin.IdpBaseConfig
	if strPg.IsNotBlank(source.BaseConfig) {
		if err := json.Unmarshal([]byte(source.BaseConfig), &baseCfg); nil != err {
			return rt.ErrorMessage("认证源配置解析失败")
		}
	}
	redirectUrl, err := c.resolveRedirectUri(ctx, source, baseCfg.RedirectUrl, ct.RedirectUri)
	if nil != err {
		return rt.ErrorMessage("回调地址不在白名单")
	}
	providerType := providerEntity.Code
	if providerType != "Custom" && providerType != "Custom Flexible" && !isCustomProviderType(providerType) {
		return rt.ErrorMessage("该提供商不支持服务端 state 发起，请使用其官方 SDK")
	}
	buf := make([]byte, 24)
	if _, err := rand.Read(buf); err != nil {
		return rt.ErrorMessage("生成 state 失败")
	}
	state := hex.EncodeToString(buf)
	idpStates.Store(state, idpStateEntry{SourceNo: ct.SourceNo, RedirectUri: redirectUrl, ExpiresAt: time.Now().Add(idpStateTTL)})
	idpStates.Range(func(key, value any) bool {
		if entry, ok := value.(idpStateEntry); ok && time.Now().After(entry.ExpiresAt) {
			idpStates.Delete(key)
		}
		return true
	})
	cfg := &oauth2.Config{
		ClientID:     baseCfg.ClientId,
		ClientSecret: baseCfg.ClientSecret,
		RedirectURL:  redirectUrl,
		Endpoint:     oauth2.Endpoint{AuthURL: baseCfg.AuthURL, TokenURL: baseCfg.TokenURL},
	}
	authURL := cfg.AuthCodeURL(state, oauth2.AccessTypeOffline)
	return rt.OkData(map[string]string{"state": state, "authUrl": authURL, "redirectUri": redirectUrl})
}

func isCustomProviderType(t string) bool {
	return len(t) >= 6 && t[:6] == "Custom"
}

// consumeState 校验并一次性消费 state。
func (c *AccountLoginIdpService) consumeState(state, sourceNo, redirectUri string) error {
	if strPg.IsBlank(state) {
		return fmt.Errorf("state 不能为空")
	}
	value, ok := idpStates.LoadAndDelete(state)
	if !ok {
		return fmt.Errorf("state 无效或已使用")
	}
	entry, ok := value.(idpStateEntry)
	if !ok || time.Now().After(entry.ExpiresAt) {
		return fmt.Errorf("state 已过期")
	}
	if entry.SourceNo != sourceNo || entry.RedirectUri != redirectUri {
		return fmt.Errorf("state 与请求不匹配")
	}
	return nil
}
