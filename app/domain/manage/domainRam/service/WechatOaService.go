package service

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/ram/modRamWechat"
	"github.com/hongmengzhu/xianfu-blog-go/infrastructure/repositoryRam"
	idp2 "github.com/hongmengzhu/xianfu-blog-go/pkg/auth/idp"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/enum/state/enumStatePg"
	"github.com/pangu-2/go-tools/tools/strPg"
	"github.com/pangu-2/go-tools/tools/wrapperPg/rg"
	"go-spring.org/log"
	"go-spring.org/spring/gs"
)

func init() {
	gs.Provide(NewWechatOaService)
}

// wechatOaBaseConfig 微信公众号认证源基础配置（从 BaseConfig JSON 解析）
type wechatOaBaseConfig struct {
	ClientId2     string `json:"clientId2"`
	ClientSecret2 string `json:"clientSecret2"`
	WechatToken   string `json:"wechatToken"`
}

// wechatEventXML 微信推送事件 XML 结构
type wechatEventXML struct {
	MsgType      string `xml:"MsgType"`
	Event        string `xml:"Event"`
	EventKey     string `xml:"EventKey"`
	FromUserName string `xml:"FromUserName"`
	Ticket       string `xml:"Ticket"`
}

// wechatCacheValue 微信扫码缓存值
type wechatCacheValue struct {
	IsScanned     bool
	WechatUnionId string
}

var (
	wechatOaCacheMap sync.Map
)

// WechatOaService 微信公众号扫码登录 Service
type WechatOaService struct {
	daoSource   *repositoryRam.RamIdentitySourceRepository   `autowire:"?"`
	daoProvider *repositoryRam.RamIdentityProviderRepository `autowire:"?"`
}

// NewWechatOaService 构造函数
func NewWechatOaService() *WechatOaService {
	return new(WechatOaService)
}

// GetQRCode 获取扫码登录二维码
func (s *WechatOaService) GetQRCode(ctx *gin.Context, sourceNo string) (rt rg.Rs[modRamWechat.QRCodeVo]) {
	if strPg.IsBlank(sourceNo) {
		return rt.ErrorMessage("认证源编号不能为空")
	}

	// 查询认证源
	source, found := s.daoSource.FindByNo(ctx, sourceNo)
	if !found {
		return rt.ErrorMessage("认证源不存在")
	}
	if !enumStatePg.ENABLE.IsExistInt8(source.State) {
		return rt.ErrorMessage("认证源已禁用")
	}

	// 解析 BaseConfig
	var cfg wechatOaBaseConfig
	if strPg.IsNotBlank(source.BaseConfig) {
		if err := json.Unmarshal([]byte(source.BaseConfig), &cfg); err != nil {
			log.Errorf(ctx, log.TagAppDef, "解析 BaseConfig 失败: %v", err)
			return rt.ErrorMessage("认证源配置解析失败")
		}
	}
	if strPg.IsBlank(cfg.ClientId2) || strPg.IsBlank(cfg.ClientSecret2) {
		return rt.ErrorMessage("公众号 AppID 或 AppSecret 未配置")
	}

	// 生成 scene_str: sourceNo + 随机字符串
	sceneStr := sourceNo + "_" + randomString(16)

	// 调用微信 API 创建二维码
	qrResult, err := idp2.WechatOaCreateQRCode(cfg.ClientId2, cfg.ClientSecret2, sceneStr)
	if err != nil {
		log.Errorf(ctx, log.TagAppDef, "创建公众号二维码失败: %v", err)
		return rt.ErrorMessage("创建二维码失败: " + err.Error())
	}

	vo := modRamWechat.QRCodeVo{
		Ticket:        qrResult.Ticket,
		QRUrl:         idp2.WechatOaQRCodeShowURL(qrResult.Ticket),
		ExpireSeconds: qrResult.ExpireSeconds,
	}
	return rt.OkData(vo)
}

// PollScanStatus 轮询扫码状态
func (s *WechatOaService) PollScanStatus(ctx *gin.Context, ticket string) (rt rg.Rs[modRamWechat.PollVo]) {
	if strPg.IsBlank(ticket) {
		return rt.ErrorMessage("ticket 不能为空")
	}

	val, ok := wechatOaCacheMap.Load(ticket)
	if !ok {
		return rt.OkData(modRamWechat.PollVo{IsScanned: false})
	}

	cacheVal := val.(*wechatCacheValue)
	vo := modRamWechat.PollVo{
		IsScanned:     cacheVal.IsScanned,
		WechatUnionId: cacheVal.WechatUnionId,
	}

	// 已扫码，返回后删除缓存（一次性使用）
	if cacheVal.IsScanned {
		wechatOaCacheMap.Delete(ticket)
	}

	return rt.OkData(vo)
}

// VerifySignature 验证微信服务器签名
func (s *WechatOaService) VerifySignature(token, nonce, timestamp, signature string) bool {
	tmpArr := sort.StringSlice{token, timestamp, nonce}
	sort.Sort(tmpArr)

	tmpStr := strings.Join(tmpArr, "")
	b := sha1.Sum([]byte(tmpStr))
	res := hex.EncodeToString(b[:])
	return res == signature
}

// HandleEvent 处理微信 webhook 事件（SCAN / SUBSCRIBE）
func (s *WechatOaService) HandleEvent(ctx *gin.Context, bodyBytes []byte, sourceNo string) (rt rg.Rs[string]) {
	// 解析 XML
	var eventData wechatEventXML
	if err := xml.Unmarshal(bodyBytes, &eventData); err != nil {
		log.Errorf(ctx, log.TagAppDef, "解析微信事件 XML 失败: %v", err)
		return rt.ErrorMessage("解析事件失败")
	}

	event := strings.ToUpper(eventData.Event)
	if event != "SCAN" && event != "SUBSCRIBE" {
		return rt.OkData("")
	}

	if strPg.IsBlank(eventData.Ticket) {
		return rt.ErrorMessage("empty ticket")
	}

	// 通过 EventKey（即 scene_str）查找认证源
	// scene_str 格式: sourceNo_randomId
	sceneStr := eventData.EventKey
	var cfg wechatOaBaseConfig

	// 优先从 scene_str 中提取 sourceNo 查找
	if strPg.IsNotBlank(sceneStr) {
		parts := strings.SplitN(sceneStr, "_", 2)
		if len(parts) > 0 {
			source, found := s.daoSource.FindByNo(ctx, parts[0])
			if found && strPg.IsNotBlank(source.BaseConfig) {
				_ = json.Unmarshal([]byte(source.BaseConfig), &cfg)
			}
		}
	}

	// 如果 scene_str 方式未找到，尝试用传入的 sourceNo
	if strPg.IsBlank(cfg.ClientId2) && strPg.IsNotBlank(sourceNo) {
		source, found := s.daoSource.FindByNo(ctx, sourceNo)
		if found && strPg.IsNotBlank(source.BaseConfig) {
			_ = json.Unmarshal([]byte(source.BaseConfig), &cfg)
		}
	}

	// 写入缓存
	wechatOaCacheMap.Store(eventData.Ticket, &wechatCacheValue{
		IsScanned:     true,
		WechatUnionId: eventData.FromUserName,
	})
	log.Infof(ctx, log.TagAppDef, "微信扫码事件: ticket=%s, unionId=%s, event=%s", eventData.Ticket, eventData.FromUserName, event)

	// 同时写入 idp 包的全局缓存（兼容已有的 wechat.go GetUserInfo 逻辑）
	idp2.Lock.Lock()
	if idp2.WechatCacheMap == nil {
		idp2.WechatCacheMap = make(map[string]idp2.WechatCacheMapValue)
	}
	idp2.WechatCacheMap[eventData.Ticket] = idp2.WechatCacheMapValue{
		IsScanned:     true,
		WechatUnionId: eventData.FromUserName,
	}
	idp2.Lock.Unlock()

	return rt.OkData("")
}

// GetWechatOaToken 获取认证源配置的微信 Token（用于签名验证）
func (s *WechatOaService) GetWechatOaToken(ctx *gin.Context, sourceNo string) string {
	source, found := s.daoSource.FindByNo(ctx, sourceNo)
	if !found {
		return ""
	}
	var cfg wechatOaBaseConfig
	if strPg.IsNotBlank(source.BaseConfig) {
		_ = json.Unmarshal([]byte(source.BaseConfig), &cfg)
	}
	return cfg.WechatToken
}

// randomString 生成指定长度的随机字符串
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	sb := make([]byte, n)
	for i := range sb {
		sb[i] = letters[r.Intn(len(letters))]
	}
	return string(sb)
}
