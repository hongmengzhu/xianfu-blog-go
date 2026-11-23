package opLogPg

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/consts/constContextPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/consts/constModulePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/sdk/auth/holderPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/sdk/auth/holderPg/holderApiPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/trace"
	"github.com/pangu-2/go-tools/tools/noPg"
	"go-spring.org/log"
)

const (
	// maxRequestBody 请求体读取上限，超过则不记录，避免影响业务与大字段入库。
	maxRequestBody = 64 * 1024
	// contextKey 请求级 traceID 存放键，handler 可通过 opLogPg.TraceID(ctx) 获取。
	contextKey = "opLogPg:traceID"
)

// TraceID 返回当前请求的链路 ID。
func TraceID(ctx *gin.Context) string {
	if v, ok := ctx.Get(contextKey); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// Global 全局操作日志中间件：每个请求记录一条操作日志。
// 必须在注册路由前 Use。
func Global() gin.HandlerFunc {
	return (&Middleware{}).Global()
}

// Global 返回使用当前配置的全局操作日志中间件。
// 配置中的 exclude 命中路径时直接放行，不记录日志；
// 默认不记录读请求（见 Config.IncludeRead），只有写操作才入库。
func (m *Middleware) Global() gin.HandlerFunc {
	matcher := newPathMatcher(m.Config.Exclude)
	includeRead := m.Config.IncludeRead
	return func(c *gin.Context) {
		if matcher.Match(c.Request.URL.Path) {
			c.Next()
			return
		}
		if !shouldRecordMethod(c.Request.Method, includeRead) {
			c.Next()
			return
		}
		traceID := sanitizeTraceID(c.GetHeader("x-request-id"))
		if traceID == "" {
			traceID = trace.GenerateTraceID()
		}
		c.Set(contextKey, traceID)
		c.Header("X-Request-Id", traceID)

		body := captureRequestBody(c)

		// 记录放在 defer 中，panic 时也能落一条（状态码按 500 处理）。
		defer func() {
			if r := recover(); r != nil {
				record(c, traceID, body, http.StatusInternalServerError)
				panic(r)
			}
		}()

		c.Next()
		record(c, traceID, body, c.Writer.Status())
	}
}

// record 组装并异步写入一条操作日志。
func record(c *gin.Context, traceID, body string, statusCode int) {
	if shouldSkip(c) {
		return
	}
	w := defaultWriter.Load()
	if w == nil {
		log.Error(c.Request.Context(), log.TagAppDef, nil, "操作日志 Writer 未初始化，日志跳过："+traceID)
		return
	}
	id := resolveIdentity(c)
	path := c.Request.URL.Path
	moduleRoot, module, moduleName, action, actionName := deriveModuleAction(path)
	if id.typeDomain == "" {
		id.typeDomain = deriveTypeDomain(path)
	}
	now := time.Now()
	entity := &entityBasic.BasicOperationLogEntity{
		No:              noPg.MakeNo("AL"),
		TenantNo:        truncate(id.tenantNo, 80),
		OrgNo:           truncate(id.orgNo, 80),
		StoreNo:         truncate(id.storeNo, 80),
		TypeDomain:      truncate(id.typeDomain, 40),
		AppNo:           truncate(id.AppNo, 80),
		TerminalCode:    truncate(id.TerminalCode, 80),
		Client:          truncate(id.Client, 80),
		DeviceNo:        truncate(id.DeviceNo, 80),
		Module:          truncate(module, 80),
		ModuleName:      truncate(moduleName, 120),
		ModuleRoot:      truncate(moduleRoot, 80),
		Action:          truncate(action, 80),
		ActionName:      truncate(actionName, 80),
		OperatorNo:      truncate(id.operatorNo, 80),
		OperatorAccount: truncate(id.operatorAccount, 120),
		OperatorName:    truncate(id.operatorName, 120),
		Roles:           id.roles,
		TraceID:         traceID,
		StatusCode:      statusCode,
		Method:          truncate(c.Request.Method, 16),
		Path:            truncate(path, 512),
		IP:              truncate(c.ClientIP(), 128),
		Header:          headerJSON(c.Request.Header),
		Param:           truncate(c.Request.URL.RawQuery, maxRequestBody),
		RequestBody:     truncate(body, maxRequestBody),
		CreateAt:        &now,
	}
	w.enqueue(entity)
}

// requestIdentity 当前请求可获取到的操作人/归属信息。
type requestIdentity struct {
	tenantNo        string
	orgNo           string
	storeNo         string
	typeDomain      string
	operatorNo      string
	operatorAccount string
	operatorName    string
	roles           []string
	TerminalCode    string
	Client          string
	AppNo           string
	DeviceNo        string
}

// resolveIdentity 从 manage/system（AUTH_LOGIN）或开放 API（AUTH_LOGIN_API）上下文读取身份。
func resolveIdentity(c *gin.Context) requestIdentity {
	var id requestIdentity
	id.roles = make([]string, 0)
	if v, ok := c.Get(constContextPg.AUTH_LOGIN); ok {
		var h *holderPg.HolderPg
		switch x := v.(type) {
		case holderPg.HolderPg:
			h = &x
		case *holderPg.HolderPg:
			h = x
		}
		if h != nil && h.HolderData != nil {
			acc := h.GetAccount()
			id.tenantNo = normalizeCode(h.GetTenantNo())
			id.orgNo = normalizeCode(h.GetOrgNo())
			id.storeNo = normalizeCode(h.GetStoreNo())
			id.typeDomain = h.GetTypeDomain()
			id.operatorNo = acc.No
			id.operatorAccount = acc.Account
			id.operatorName = acc.Name
			if len(acc.Os.Roles) > 0 {
				id.roles = append(id.roles, acc.Os.Roles...)
			}
			{
				if v2, ok2 := c.Get(constContextPg.AUTH_TerminalCode); ok2 {
					id.TerminalCode = v2.(string)
				}
				if v2, ok2 := c.Get(constContextPg.AUTH_Client); ok2 {
					id.Client = v2.(string)
				}
				if v2, ok2 := c.Get(constContextPg.AUTH_AppNo); ok2 {
					id.AppNo = v2.(string)
				}
				if v2, ok2 := c.Get(constContextPg.AUTH_DeviceNo); ok2 {
					id.DeviceNo = v2.(string)
				}
			}
			return id
		}
	}
	if v, ok := c.Get(constContextPg.AUTH_LOGIN_API); ok {
		var h *holderApiPg.HolderPg
		switch x := v.(type) {
		case holderApiPg.HolderPg:
			h = &x
		case *holderApiPg.HolderPg:
			h = x
		}
		if h != nil {
			d := h.GetDipl()
			id.tenantNo = normalizeCode(h.GetTenantNo())
			id.typeDomain = h.TypeDomain
			id.operatorNo = d.Ano
			id.operatorAccount = d.Ano
			id.operatorName = d.Name
		}
	}
	{
		if v2, ok2 := c.Get(constContextPg.AUTH_TerminalCode); ok2 {
			id.TerminalCode = v2.(string)
		}
		if v2, ok2 := c.Get(constContextPg.AUTH_Client); ok2 {
			id.Client = v2.(string)
		}
		if v2, ok2 := c.Get(constContextPg.AUTH_AppNo); ok2 {
			id.AppNo = v2.(string)
		}
		if v2, ok2 := c.Get(constContextPg.AUTH_DeviceNo); ok2 {
			id.DeviceNo = v2.(string)
		}
	}
	return id
}

// normalizeCode 把 holder 的空值占位 "-1" 归一为空字符串。
func normalizeCode(s string) string {
	s = strings.TrimSpace(s)
	if s == "-1" {
		return ""
	}
	return s
}

// deriveRoute 从路径推导模块/动作：/xianfu/{domain}/{root}/{module}/{action}。
func deriveRoute(path string) (root, module, action string) {
	segs := splitPath(path)
	if len(segs) == 0 {
		return
	}
	domain := segs[0]
	switch {
	case domain == "xianfu" && len(segs) >= 5:
		root, module, action = segs[2], segs[3], segs[len(segs)-1]
	case domain == "xianfu" && len(segs) >= 4:
		module, action = segs[2], segs[len(segs)-1]
	case domain == "xianfu" && len(segs) == 3:
		module, action = segs[1], segs[2]
	case strings.HasPrefix(domain, "api-") && len(segs) >= 3:
		root, module, action = domain, segs[1], segs[len(segs)-1]
	case len(segs) >= 3:
		root, module, action = segs[0], segs[len(segs)-2], segs[len(segs)-1]
	case len(segs) == 2:
		module, action = segs[0], segs[1]
	default:
		action = segs[0]
	}
	return
}

func splitPath(path string) []string {
	path = strings.Trim(path, "/")
	if path == "" {
		return nil
	}
	return strings.Split(path, "/")
}

// deriveTypeDomain holder 未提供域类型时按路径兜底。
func deriveTypeDomain(path string) string {
	switch {
	case strings.HasPrefix(path, "/xianfu/manage/"):
		return "manage"
	case strings.HasPrefix(path, "/xianfu/sys/"):
		return "system"
	case strings.HasPrefix(path, "/api-"):
		return "customer"
	default:
		return "general"
	}
}

// deriveTerminalCode 终端类型兜底。
func deriveTerminalCode(path string) string {
	switch {
	case strings.HasPrefix(path, "/xianfu/manage/"):
		return "manage"
	case strings.HasPrefix(path, "/xianfu/sys/"):
		return "system"
	case strings.HasPrefix(path, "/api-app/"):
		return "app"
	default:
		return "web"
	}
}

// headerValue 依次读取若干请求头，返回第一个非空值。
func headerValue(c *gin.Context, names ...string) string {
	for _, name := range names {
		if v := strings.TrimSpace(c.GetHeader(name)); v != "" {
			return v
		}
	}
	return ""
}

// headerJSON 序列化请求头；敏感头做脱敏，避免把令牌写入日志。
func headerJSON(h http.Header) string {
	if len(h) == 0 {
		return ""
	}
	m := make(map[string]string, len(h))
	for k, vs := range h {
		v := strings.Join(vs, ",")
		// if isSensitiveHeader(k) {
		// 	v = "***"
		// }
		if !isSensitiveHeaderUserAgent(k) {
			continue
		}
		m[k] = truncate(v, 512)
	}
	bs, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return truncate(string(bs), maxRequestBody)
}

func isSensitiveHeader(k string) bool {
	switch strings.ToLower(k) {
	case "authorization", "authorization-sys", "authorization-refresh", "authorization-refresh-sys",
		"proxy-authorization", "cookie", "set-cookie", "token", "x-token", "auth", "x-api-key", "api-key":
		return true
	}
	return false
}
func isSensitiveHeaderUserAgent(k string) bool {
	switch strings.ToLower(k) {
	case "user-agent":
		return true
	}
	return false
}

// shouldRecordMethod 判断当前请求是否需要记入操作日志。
// 默认只记录有副作用的写操作；读请求（GET/HEAD）仅有在配置显式开启时才记录，
// 避免列表/详情接口与页面浏览把操作日志表写满。
func shouldRecordMethod(method string, includeRead bool) bool {
	if includeRead {
		return true
	}
	switch method {
	case http.MethodGet, http.MethodHead:
		return false
	default:
		return true
	}
}

// shouldSkip 过滤静态资源与预检请求，避免操作日志被资源请求淹没。
func shouldSkip(c *gin.Context) bool {
	if c.Request.Method == http.MethodOptions {
		return true
	}
	path := c.Request.URL.Path
	return path == "/assets" || strings.HasPrefix(path, "/assets/") ||
		path == "/attachment" || strings.HasPrefix(path, "/attachment/")
}

// captureRequestBody 读取请求体并原样放回；文件流、非文本类型与超大请求体不记录。
func captureRequestBody(c *gin.Context) string {
	if c.Request == nil || c.Request.Body == nil {
		return ""
	}
	if c.Request.ContentLength <= 0 || c.Request.ContentLength > maxRequestBody {
		return ""
	}
	ct := c.ContentType()
	if isFileStream(ct) || !isTextContentType(ct) {
		return ""
	}
	bs, err := io.ReadAll(io.LimitReader(c.Request.Body, maxRequestBody))
	if err != nil {
		c.Request.Body = io.NopCloser(bytes.NewReader(bs))
		return ""
	}
	if len(bs) >= maxRequestBody {
		// 达到读取上限时把未读完的部分接回去，保证 handler 仍能读到完整请求体。
		c.Request.Body = io.NopCloser(io.MultiReader(bytes.NewReader(bs), c.Request.Body))
	} else {
		c.Request.Body = io.NopCloser(bytes.NewReader(bs))
	}
	return string(bs)
}

// isFileStream 判断是否为文件流/二进制内容，这类请求体不记录。
func isFileStream(contentType string) bool {
	ct := strings.ToLower(contentType)
	switch {
	case strings.HasPrefix(ct, "multipart/form-data"),
		strings.HasPrefix(ct, "application/octet-stream"),
		strings.HasPrefix(ct, "image/"),
		strings.HasPrefix(ct, "video/"),
		strings.HasPrefix(ct, "audio/"),
		strings.HasPrefix(ct, "application/pdf"),
		strings.HasPrefix(ct, "application/zip"),
		strings.HasPrefix(ct, "application/gzip"),
		strings.HasPrefix(ct, "application/x-tar"),
		strings.HasPrefix(ct, "application/msword"),
		strings.HasPrefix(ct, "application/vnd.openxmlformats"),
		strings.HasPrefix(ct, "application/vnd.ms-"):
		return true
	}
	return false
}

// isTextContentType 只记录可读的文本类请求体。
func isTextContentType(contentType string) bool {
	ct := strings.ToLower(contentType)
	switch {
	case strings.HasPrefix(ct, "application/json"),
		strings.HasSuffix(ct, "+json"),
		strings.HasPrefix(ct, "application/x-www-form-urlencoded"),
		strings.HasPrefix(ct, "application/xml"),
		strings.HasPrefix(ct, "application/graphql"),
		strings.HasPrefix(ct, "text/"):
		return true
	}
	return false
}

// sanitizeTraceID 清理请求头里的 traceID：去空白、控制字符并限长。
func sanitizeTraceID(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Map(func(r rune) rune {
		if r < 0x20 || r == 0x7f {
			return -1
		}
		return r
	}, s)
	s = strings.TrimSpace(s)
	if len(s) > 80 {
		s = truncate(s, 80)
	}
	return s
}

func truncate(s string, n int) string {
	if n <= 0 || len(s) <= n {
		return s
	}
	s = s[:n]
	// 截断可能落在多字节字符中间，回退到完整字符，避免入库非法 UTF-8。
	for len(s) > 0 && !utf8.ValidString(s) {
		s = s[:len(s)-1]
	}
	return s
}

// moduleIndexOnce 模块索引只构建一次，避免每次请求遍历 ModuleMap。
var (
	moduleIndexOnce sync.Once
	moduleIndex     map[string]constModulePg.Item
)

// deriveModuleAction 优先按 constModulePg 解析模块，动作按 MatchActionMap 推导；
// 模块未登记时退回路径推导，动作仍按 MatchActionMap 匹配。
func deriveModuleAction(path string) (root, module, moduleName, action, actionName string) {
	if item, act, ok := resolveModule(path); ok {
		action, actionName = matchAction(act)
		return item.ModuleRoot, item.Module, item.Name, action, actionName
	}
	root, module, rawAction := deriveRoute(path)
	action, actionName = matchAction(rawAction)
	return root, module, module, action, actionName
}

// resolveModule 按路由中的 <domain>/<resource> 查 constModulePg.ModuleMap，动作取路由最后一段。
func resolveModule(path string) (constModulePg.Item, string, bool) {
	segs := splitPath(path)
	if len(segs) < 3 {
		return constModulePg.Item{}, "", false
	}
	var key, action string
	switch {
	case segs[0] == "xianfu" && len(segs) >= 5:
		key, action = segs[2]+"/"+segs[3], segs[len(segs)-1]
	case strings.HasPrefix(segs[0], "api-") && len(segs) >= 3:
		key, action = "api/"+segs[1], segs[len(segs)-1]
	default:
		key, action = segs[len(segs)-3]+"/"+segs[len(segs)-2], segs[len(segs)-1]
	}
	item, ok := lookupModule(key)
	return item, action, ok
}

// lookupModule 忽略大小写、- 和 _ 查模块；PathMark 用于路由与模块编码不一致的场景。
func lookupModule(key string) (constModulePg.Item, bool) {
	moduleIndexOnce.Do(func() {
		moduleIndex = make(map[string]constModulePg.Item, len(constModulePg.ModuleMap)*2)
		for path, item := range constModulePg.ModuleMap {
			moduleIndex[normalizeModuleKey(path)] = item
			for _, mark := range item.PathMark {
				moduleIndex[normalizeModuleKey(mark)] = item
			}
		}
	})
	item, ok := moduleIndex[normalizeModuleKey(key)]
	return item, ok
}

func normalizeModuleKey(key string) string {
	key = strings.ToLower(strings.TrimSpace(key))
	key = strings.ReplaceAll(key, "-", "")
	key = strings.ReplaceAll(key, "_", "")
	return strings.Trim(key, "/")
}

// matchAction 解析动作名：先精确匹配 MatchActionMap，再按 Prefix 取最长前缀，
// 未登记的动作保留原始动作名。
func matchAction(action string) (string, string) {
	action = strings.TrimSpace(action)
	if action == "" {
		return "", ""
	}
	if m, ok := MatchActionMap[action]; ok {
		return m.Action, m.ActionName
	}
	bestLen := 0
	var best Match
	found := false
	for _, m := range MatchActionMap {
		for _, prefix := range m.Prefix {
			prefix = strings.TrimSpace(prefix)
			if prefix == "" || len(prefix) <= bestLen {
				continue
			}
			if action == prefix || strings.HasPrefix(action, prefix) {
				best, bestLen, found = m, len(prefix), true
			}
		}
	}
	if found {
		return best.Action, best.ActionName
	}
	return action, action
}
