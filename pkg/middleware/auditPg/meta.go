package auditPg

import (
	"fmt"
	"strings"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityBasic"
)

// TargetFrom 业务对象 ID 的来源。
type TargetFrom int8

const (
	// TargetFromContext 由业务通过 SetTarget 显式提供。
	TargetFromContext TargetFrom = 0
	// TargetFromPath 从路由参数 :id 读取。
	TargetFromPath TargetFrom = 1
	// TargetFromQuery 从 query 参数读取。
	TargetFromQuery TargetFrom = 2
	// TargetFromBody 从请求 JSON body 的字段读取。
	TargetFromBody TargetFrom = 3
)

// Target 审计日志中记录的业务对象。
type Target struct {
	Type     string `json:"type"`
	ID       string `json:"id"`
	No       string `json:"no"`
	Name     string `json:"name"`
	TenantNo string `json:"tenantNo"`
	Count    int    `json:"count"`
}

// DiffItem 单个字段的变更前后值。
type DiffItem struct {
	Before any `json:"before"`
	After  any `json:"after"`
}

// Result 业务显式设置的操作结果。
type Result struct {
	Code    string `json:"code"`
	Success bool   `json:"success"`
	Partial bool   `json:"partial"`
	Message string `json:"message"`
}

// Meta 一条路由的审计元数据。
type Meta struct {
	Module          string     // 模块 code，如 "blog_article"
	ModuleName      string     // 模块名称快照，如 "文章管理"
	SubModule       string     // 子功能，可空
	Action          string     // 动作 code：create/update/delete/enable/disable/export/...
	ActionName      string     // 动作名称快照
	OperationType   int8       // 见 entityBasic.OperationLogTypeXxx
	RiskLevel       int8       // 1/2/3；系统端自动 +1
	TargetType      string     // 对象类型，如 "blog_article"
	TargetFrom      TargetFrom // 对象 ID 来源
	TargetKey       string     // path/query/body 的字段名
	Description     string     // 默认描述模板，可被业务覆盖
	SnapshotFields  []string   // 快照白名单；空表示全部字段
	IgnoreFields    []string   // diff 忽略字段
	Sync            bool       // 强制同步写
	SkipNoChange    bool       // 更新无变化时跳过（默认 false；需业务提供 before/after 才能开启）
	RecordRequest   bool       // 是否保存请求体（默认 true）
	SkipRequestBody bool       // 显式关闭请求体记录
	RecordResponse  bool       // 是否保存响应体（默认只在失败时保存）
	MaxBodyBytes    int        // 单字段截断上限，默认 32KB
}

// MaxBodyBytesLimit 单字段上限，超过按上限处理。
const MaxBodyBytesLimit = 128 * 1024

func (m Meta) withDefaults() Meta {
	if m.RiskLevel < entityBasic.OperationLogRiskLow || m.RiskLevel > entityBasic.OperationLogRiskHigh {
		m.RiskLevel = entityBasic.OperationLogRiskLow
	}
	if m.MaxBodyBytes <= 0 {
		m.MaxBodyBytes = 32 * 1024
	}
	if m.MaxBodyBytes > MaxBodyBytesLimit {
		m.MaxBodyBytes = MaxBodyBytesLimit
	}
	// 请求体默认记录；业务可用 SkipRequestBody 显式关闭。
	if !m.SkipRequestBody {
		m.RecordRequest = true
	}
	if m.TargetFrom == TargetFromContext && m.TargetKey == "" {
		m.TargetFrom = TargetFromContext
	}
	return m
}

// Validate 启动期校验，缺字段直接 panic，避免上线后才发现没记。
func (m Meta) Validate() error {
	if strings.TrimSpace(m.Module) == "" {
		return fmt.Errorf("auditPg: Meta.Module 不能为空")
	}
	if strings.TrimSpace(m.Action) == "" {
		return fmt.Errorf("auditPg: Meta.Action 不能为空")
	}
	if strings.TrimSpace(m.TargetType) == "" {
		return fmt.Errorf("auditPg: Meta.TargetType 不能为空")
	}
	if m.OperationType < entityBasic.OperationLogTypeCreate || m.OperationType > entityBasic.OperationLogTypeOther {
		return fmt.Errorf("auditPg: Meta.OperationType 必须在 1-9 之间")
	}
	if m.RiskLevel < entityBasic.OperationLogRiskLow || m.RiskLevel > entityBasic.OperationLogRiskHigh {
		return fmt.Errorf("auditPg: Meta.RiskLevel 必须在 1-3 之间")
	}
	if m.TargetFrom < TargetFromContext || m.TargetFrom > TargetFromBody {
		return fmt.Errorf("auditPg: Meta.TargetFrom 非法")
	}
	if m.TargetFrom == TargetFromPath || m.TargetFrom == TargetFromQuery || m.TargetFrom == TargetFromBody {
		if strings.TrimSpace(m.TargetKey) == "" {
			return fmt.Errorf("auditPg: TargetFrom=%d 时 TargetKey 不能为空", m.TargetFrom)
		}
	}
	return nil
}

// actionNames 常用动作名称，Meta.ActionName 为空时兜底。
var actionNames = map[string]string{
	"create":         "新增",
	"update":         "修改",
	"delete":         "删除",
	"enable":         "启用",
	"disable":        "停用",
	"recovery":       "恢复",
	"grant":          "授权",
	"export":         "导出",
	"import":         "导入",
	"reset":          "重置",
	"audit":          "审核",
	"loginAs":        "代登录",
	"querySensitive": "敏感查询",
}

// ActionName 返回动作中文名。
func (m Meta) actionName() string {
	if m.ActionName != "" {
		return m.ActionName
	}
	if v, ok := actionNames[m.Action]; ok {
		return v
	}
	return m.Action
}
