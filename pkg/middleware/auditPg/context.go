package auditPg

import (
	"time"

	"github.com/gin-gonic/gin"
)

const (
	// contextKey gin.Context 中保存 *AuditContext 的键。
	contextKey = "auditPg:context"
	// sourceApi 来源为接口。
	sourceApi int8 = 1
)

// AuditContext 单次请求的审计上下文。
type AuditContext struct {
	Meta      Meta
	RequestID string
	TraceID   string
	StartAt   time.Time
	OperateAt time.Time

	TenantNo   string
	OrgNo      string
	StoreNo    string
	TypeDomain string
	AppNo      string

	OperatorNo      string
	OperatorAccount string
	OperatorName    string
	OperatorType    int8
	RoleNo          string
	RoleName        string

	Target      Target
	Before      any
	After       any
	Diff        map[string]DiffItem
	Description string
	Extra       map[string]any
	Result      *Result
	HighRisk    bool
	NoChange    bool
	Skip        bool
	Source      int8

	RequestBody  []byte
	ResponseBody []byte
	HTTPStatus   int
	PanicValue   any
	ErrorMessage string
}

// Get 返回当前请求的审计上下文；没有挂审计中间件时返回 nil。
func Get(ctx *gin.Context) *AuditContext {
	if ctx == nil {
		return nil
	}
	v, ok := ctx.Get(contextKey)
	if !ok {
		return nil
	}
	ac, _ := v.(*AuditContext)
	return ac
}

func set(ctx *gin.Context, fn func(ac *AuditContext)) {
	if ac := Get(ctx); ac != nil {
		fn(ac)
	}
}

// SetTarget 设置业务对象；批量操作时 Count 为总数，Extra 里可存 target_ids。
func SetTarget(ctx *gin.Context, target Target) {
	set(ctx, func(ac *AuditContext) { ac.Target = target })
}

// SetBefore 设置变更前快照（struct / map 均可，defer 时统一序列化）。
func SetBefore(ctx *gin.Context, before any) {
	set(ctx, func(ac *AuditContext) { ac.Before = before })
}

// SetAfter 设置变更后快照。
func SetAfter(ctx *gin.Context, after any) {
	set(ctx, func(ac *AuditContext) { ac.After = after })
}

// SetDiff 直接提供 diff，跳过框架计算。
func SetDiff(ctx *gin.Context, diff map[string]DiffItem) {
	set(ctx, func(ac *AuditContext) { ac.Diff = diff })
}

// SetDescription 覆盖默认描述。
func SetDescription(ctx *gin.Context, description string) {
	set(ctx, func(ac *AuditContext) { ac.Description = description })
}

// SetExtra 追加扩展信息；同一 key 以最后一次为准。
func SetExtra(ctx *gin.Context, key string, value any) {
	set(ctx, func(ac *AuditContext) {
		if ac.Extra == nil {
			ac.Extra = make(map[string]any)
		}
		ac.Extra[key] = value
	})
}

// MarkHighRisk 标记高风险，强制同步写。
func MarkHighRisk(ctx *gin.Context) {
	set(ctx, func(ac *AuditContext) { ac.HighRisk = true })
}

// SetNoChange 标记本次更新没有实际变化；SkipNoChange 生效时跳过落库。
func SetNoChange(ctx *gin.Context, noChange bool) {
	set(ctx, func(ac *AuditContext) { ac.NoChange = noChange })
}

// SetSkip 跳过本次落库（例如业务确认该请求无需记录）。
func SetSkip(ctx *gin.Context, skip bool) {
	set(ctx, func(ac *AuditContext) { ac.Skip = skip })
}

// SetResult 显式指定结果，一般不需要。
func SetResult(ctx *gin.Context, result Result) {
	set(ctx, func(ac *AuditContext) { ac.Result = &result })
}

// SetErrorMessage 覆盖失败原因。
func SetErrorMessage(ctx *gin.Context, message string) {
	set(ctx, func(ac *AuditContext) { ac.ErrorMessage = message })
}
