package auditPg

import (
	"bytes"
	"encoding/json"
	"io"
	"mime"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityBasic"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/sdk/auth/holderPg"
	"go-spring.org/log"
)

const (
	// maxRequestBody 读取请求体的上限；超过则不捕获，避免影响业务。
	maxRequestBody = 64 * 1024
	// maxResponseBody 响应体捕获上限。
	maxResponseBody = 64 * 1024
	// maxUserAgent User-Agent 截断长度。
	maxUserAgent = 1024
	// maxErrorMessage 失败原因截断长度。
	maxErrorMessage = 8 * 1024
)

// Handle 为单条路由挂审计，仅用于元数据必须在运行时决定的极少数场景。
// 常规用法是在 XxxAudit.go 中 auditPg.RegisterRoute(...)，由 Global() 统一处理；
// 两种方式不要对同一路由混用，否则一条请求两条日志。
func Handle(meta Meta, handlers ...gin.HandlerFunc) gin.HandlersChain {
	m := meta.withDefaults()
	if err := m.Validate(); err != nil {
		panic("auditPg: " + err.Error())
	}
	chain := make(gin.HandlersChain, 0, len(handlers)+1)
	chain = append(chain, func(c *gin.Context) { run(c, m) })
	chain = append(chain, handlers...)
	return chain
}

func run(c *gin.Context, meta Meta) {
	start := time.Now()
	ac := &AuditContext{
		Meta:      meta,
		StartAt:   start,
		OperateAt: start,
		Source:    sourceApi,
	}
	c.Set(contextKey, ac)

	prepareContext(c, ac)
	captureRequestBody(c, ac)
	originalWriter := c.Writer
	cw := newCaptureWriter(originalWriter, maxResponseBody)
	c.Writer = cw

	defer func() {
		if r := recover(); r != nil {
			ac.PanicValue = r
			ac.ErrorMessage = truncate(resultFromPanic(r), maxErrorMessage)
			record(c, ac)
			c.Writer = originalWriter
			panic(r)
		}
		record(c, ac)
		c.Writer = originalWriter
	}()

	c.Next()
}

// prepareContext 填充请求级基础字段与操作人信息。
func prepareContext(c *gin.Context, ac *AuditContext) {
	ac.RequestID = sanitizeRequestID(c.GetHeader("X-Request-Id"))
	if ac.RequestID == "" {
		ac.RequestID = newRequestID()
	}
	if ac.RequestID != "" {
		c.Header("X-Request-Id", ac.RequestID)
	}
	ac.TraceID = traceIDFromHeader(c.GetHeader("X-Trace-Id"), c.GetHeader("traceparent"))
	ac.TypeDomain = domainFromPath(c.Request.URL.Path)

	holder := holderPg.GetContextAccount(c)
	account := holder.GetAccount()
	ac.OperatorNo = holder.GetAccountNo()
	ac.OperatorAccount = account.Account
	ac.OperatorName = account.Name
	ac.RoleNo = account.RoleNo
	ac.TenantNo = normalizeTenant(holder.GetTenantNo())
	ac.OrgNo = normalizeTenant(holder.GetOrgNo())
	ac.StoreNo = holder.GetStoreNo()
	ac.OperatorType = operatorType(holder, account)
	ac.AppNo = trustedAppNo(c.GetHeader("X-App-No"))
}

func record(c *gin.Context, ac *AuditContext) {
	if ac.Skip {
		return
	}
	afterHandle(c, ac)

	logEntity := buildEntity(c, ac)
	if logEntity == nil {
		return
	}
	w := defaultWriter.Load()
	if w == nil {
		log.Error(c.Request.Context(), log.TagAppDef, nil, "审计 Writer 未初始化，日志跳过："+logEntity.No)
		return
	}
	if ac.Meta.Sync || logEntity.RiskLevel >= entityBasic.OperationLogRiskHigh {
		w.writeOneSync(c.Request.Context(), logEntity)
		return
	}
	w.enqueue(logEntity)
}

// afterHandle 计算耗时、HTTP 状态、响应体与结果。
func afterHandle(c *gin.Context, ac *AuditContext) {
	if ac.HTTPStatus == 0 {
		ac.HTTPStatus = c.Writer.Status()
	}
	// 响应体由 captureWriter 捕获
	if cw, ok := c.Writer.(*captureWriter); ok {
		ac.ResponseBody = cw.body()
	}
	parseResult(ac)
	if ac.ErrorMessage == "" && ac.PanicValue == nil && resultValue(ac) != entityBasic.OperationLogResultSuccess {
		ac.ErrorMessage = parseMessage(ac.ResponseBody)
	}
}

func buildEntity(c *gin.Context, ac *AuditContext) *entityBasic.BasicOperationLogEntity {
	meta := ac.Meta
	target := ac.Target
	if target.Type == "" {
		target.Type = meta.TargetType
	}
	if target.ID == "" {
		target.ID = resolveTargetID(c, meta, ac)
	}

	evidenceChanged := true
	if meta.SkipNoChange && ac.OperationTypeIsUpdate() {
		before, _ := snapshot(ac.Before, meta)
		after, _ := snapshot(ac.After, meta)
		if !hasMeaningfulChange(before, after, meta) {
			evidenceChanged = false
		}
		if !evidenceChanged {
			return nil
		}
	}

	effectiveRisk := meta.RiskLevel
	if ac.TypeDomain == "system" && effectiveRisk < entityBasic.OperationLogRiskHigh {
		effectiveRisk++
	}

	description := ac.Description
	if description == "" {
		description = meta.Description
	}
	if description == "" {
		description = meta.ModuleName + meta.actionName()
		if target.Name != "" {
			description += "：" + target.Name
		}
	}

	beforeMap, beforeJSON := snapshot(ac.Before, meta)
	afterMap, afterJSON := snapshot(ac.After, meta)
	var diffJSON string
	if len(ac.Diff) > 0 {
		diffJSON = marshalDiff(ac.Diff, meta.MaxBodyBytes)
	} else {
		diffJSON = marshalDiff(computeDiff(beforeMap, afterMap, meta), meta.MaxBodyBytes)
	}

	extra := maskMap(ac.Extra)
	if target.Count > 1 {
		if _, ok := extra["target_count"]; !ok {
			extra["target_count"] = target.Count
		}
	}
	if target.TenantNo != "" && target.TenantNo != ac.TenantNo {
		extra["target_tenant_no"] = target.TenantNo
	}

	requestBody := ""
	if meta.RecordRequest && !meta.SkipRequestBody {
		requestBody = maskJSON(ac.RequestBody, meta.MaxBodyBytes)
	}
	responseBody := ""
	result := resultValue(ac)
	if meta.RecordResponse || result != entityBasic.OperationLogResultSuccess {
		responseBody = maskJSON(ac.ResponseBody, meta.MaxBodyBytes)
	}

	entity := &entityBasic.BasicOperationLogEntity{
		No:              nextNo(),
		TenantNo:        truncateName(ac.TenantNo, 80),
		OrgNo:           truncateName(ac.OrgNo, 80),
		StoreNo:         truncateName(ac.StoreNo, 80),
		TypeDomain:      truncateName(ac.TypeDomain, 40),
		AppNo:           truncateName(ac.AppNo, 80),
		Module:          truncateName(meta.Module, 80),
		ModuleName:      truncateName(meta.ModuleName, 120),
		SubModule:       truncateName(meta.SubModule, 80),
		Action:          truncateName(meta.Action, 80),
		ActionName:      truncateName(meta.actionName(), 80),
		OperationType:   meta.OperationType,
		RiskLevel:       effectiveRisk,
		Description:     truncate(description, 512),
		OperatorNo:      truncateName(ac.OperatorNo, 80),
		OperatorAccount: truncateName(ac.OperatorAccount, 120),
		OperatorName:    truncateName(ac.OperatorName, 120),
		OperatorType:    ac.OperatorType,
		RoleNo:          truncateName(ac.RoleNo, 80),
		RoleName:        truncateName(ac.RoleName, 120),
		TargetType:      truncateName(target.Type, 80),
		TargetID:        truncateName(target.ID, 80),
		TargetNo:        truncateName(target.No, 80),
		TargetName:      truncate(target.Name, 255),
		TargetCount:     target.Count,
		BeforeData:      beforeJSON,
		AfterData:       afterJSON,
		DiffData:        diffJSON,
		RequestMethod:   truncateName(c.Request.Method, 16),
		RequestPath:     truncate(c.Request.URL.Path, 512),
		RequestQuery:    maskJSONQuery(c.Request.URL.RawQuery, meta.MaxBodyBytes),
		RequestBody:     requestBody,
		ResponseBody:    responseBody,
		Result:          result,
		ResultCode:      resultCodeField(c, ac),
		ErrorMessage:    truncate(ac.ErrorMessage, maxErrorMessage),
		DurationMs:      time.Since(ac.StartAt).Milliseconds(),
		IPAddress:       truncateName(c.ClientIP(), 128),
		UserAgent:       truncate(c.Request.UserAgent(), maxUserAgent),
		RequestID:       truncateName(ac.RequestID, 80),
		TraceID:         truncateName(ac.TraceID, 80),
		Source:          ac.Source,
		State:           entityBasic.OperationLogStateNormal,
		OperateAt:       &ac.OperateAt,
	}
	if len(extra) > 0 {
		bs, err := json.Marshal(extra)
		if err == nil {
			entity.ExtraData = truncate(string(bs), meta.MaxBodyBytes)
		}
	}
	return entity
}

func resultCodeField(c *gin.Context, ac *AuditContext) string {
	if ac.Result != nil && ac.Result.Code != "" {
		return truncateName(ac.Result.Code, 40)
	}
	if ac.HTTPStatus != 0 {
		return strconv.Itoa(ac.HTTPStatus)
	}
	return ""
}

// resolveTargetID 按 Meta.TargetFrom 解析对象主键。
func resolveTargetID(c *gin.Context, meta Meta, ac *AuditContext) string {
	if meta.TargetFrom == TargetFromContext {
		return ""
	}
	if meta.TargetKey == "" {
		return ""
	}
	switch meta.TargetFrom {
	case TargetFromPath:
		return c.Param(meta.TargetKey)
	case TargetFromQuery:
		return c.Query(meta.TargetKey)
	case TargetFromBody:
		if len(ac.RequestBody) == 0 {
			return ""
		}
		var m map[string]any
		if err := json.Unmarshal(ac.RequestBody, &m); err != nil {
			return ""
		}
		if list, ok := m[meta.TargetKey].([]any); ok {
			if len(list) > 0 {
				ac.Target.Count = len(list)
			}
		}
		return valueToString(m[meta.TargetKey])
	default:
		return ""
	}
}

func valueToString(v any) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return t
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case []any:
		if len(t) > 0 {
			return valueToString(t[0])
		}
		return ""
	default:
		bs, err := json.Marshal(t)
		if err != nil {
			return ""
		}
		return strings.Trim(string(bs), "\"")
	}
}

// captureRequestBody 捕获 JSON / 表单请求体；超过上限或二进制类型不读。
func captureRequestBody(c *gin.Context, ac *AuditContext) {
	if c.Request == nil || c.Request.Body == nil {
		return
	}
	if c.Request.Method == "GET" || c.Request.Method == "HEAD" {
		return
	}
	ct := c.ContentType()
	mediaType, _, _ := mime.ParseMediaType(ct)
	if mediaType != "application/json" && mediaType != "application/x-www-form-urlencoded" {
		return
	}
	if c.Request.ContentLength > maxRequestBody {
		return
	}
	body, err := io.ReadAll(io.LimitReader(c.Request.Body, maxRequestBody))
	if err != nil {
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewReader(body))
	ac.RequestBody = body
}

// parseResult 根据响应体 code 与 HTTP 状态给出失败标记。
func parseResult(ac *AuditContext) {
	if ac.Result != nil || ac.PanicValue != nil {
		return
	}
	if ac.HTTPStatus >= 400 {
		return
	}
	if len(ac.ResponseBody) == 0 {
		return
	}
	var body struct {
		Code    json.RawMessage `json:"code"`
		Message string          `json:"message"`
	}
	if err := json.Unmarshal(ac.ResponseBody, &body); err != nil {
		return
	}
	code := strings.Trim(string(body.Code), "\"")
	if code == "" || code == "200" {
		return
	}
	ac.ErrorMessage = body.Message
	if ac.ErrorMessage == "" {
		ac.ErrorMessage = "业务返回 code=" + code
	}
}

func parseMessage(body []byte) string {
	var b struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &b); err != nil {
		return ""
	}
	return b.Message
}

// OperationTypeIsUpdate 判断本请求是否为更新类操作。
func (ac *AuditContext) OperationTypeIsUpdate() bool {
	switch ac.Meta.OperationType {
	case entityBasic.OperationLogTypeUpdate, entityBasic.OperationLogTypeStateChange:
		return true
	default:
		return false
	}
}

func operatorType(holder holderPg.HolderPg, account holderPg.AccountHolder) int8 {
	if holder.IsFounder() {
		return entityBasic.OperationLogOperatorPlatform
	}
	if strings.TrimSpace(holder.GetAccountNo()) == "" {
		return entityBasic.OperationLogOperatorAnonymous
	}
	for _, role := range account.Os.Roles {
		if strings.Contains(strings.ToLower(role), "administrator") || strings.Contains(strings.ToLower(role), "admin") {
			return entityBasic.OperationLogOperatorTenant
		}
	}
	return entityBasic.OperationLogOperatorNormal
}

func domainFromPath(path string) string {
	switch {
	case strings.HasPrefix(path, "/xianfu/sys/"):
		return "system"
	case strings.HasPrefix(path, "/xianfu/manage/"):
		return "manage"
	default:
		return "customer"
	}
}

func normalizeTenant(tenantNo string) string {
	if tenantNo == "-1" {
		return ""
	}
	return tenantNo
}

func sanitizeRequestID(raw string) string {
	raw = strings.TrimSpace(raw)
	if len(raw) > 64 {
		return ""
	}
	for _, r := range raw {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			continue
		}
		return ""
	}
	return raw
}

func traceIDFromHeader(traceID, traceparent string) string {
	if id := sanitizeRequestID(traceID); id != "" {
		return id
	}
	// traceparent: 00-<trace-id>-<parent-id>-<flags>
	parts := strings.Split(traceparent, "-")
	if len(parts) == 4 && len(parts[1]) == 32 {
		return parts[1]
	}
	return ""
}

func trustedAppNo(appNo string) string {
	appNo = strings.TrimSpace(appNo)
	if len(appNo) > 80 {
		return ""
	}
	for _, r := range appNo {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			continue
		}
		return ""
	}
	return appNo
}

func truncateName(s string, max int) string {
	return truncate(strings.TrimSpace(s), max)
}

// maskJSONQuery 独立出来，便于测试。
func maskJSONQuery(raw string, max int) string {
	return maskQuery(raw, max)
}

// captureWriter 捕获响应体前 maxResponseBody 字节。
type captureWriter struct {
	gin.ResponseWriter
	buf      bytes.Buffer
	max      int
	overflow bool
}

func newCaptureWriter(w gin.ResponseWriter, max int) *captureWriter {
	return &captureWriter{ResponseWriter: w, max: max}
}

func (w *captureWriter) Write(b []byte) (int, error) {
	w.capture(b)
	return w.ResponseWriter.Write(b)
}

func (w *captureWriter) WriteString(s string) (int, error) {
	w.capture([]byte(s))
	return w.ResponseWriter.WriteString(s)
}

func (w *captureWriter) capture(b []byte) {
	if w.max <= 0 || len(b) == 0 {
		return
	}
	remain := w.max - w.buf.Len()
	if remain <= 0 {
		w.overflow = true
		return
	}
	if len(b) > remain {
		w.buf.Write(b[:remain])
		w.overflow = true
		return
	}
	w.buf.Write(b)
}

func (w *captureWriter) body() []byte {
	return w.buf.Bytes()
}
