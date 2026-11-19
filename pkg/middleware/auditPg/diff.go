package auditPg

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityBasic"
)

// computeDiff 计算 before/after 的字段差异（仅顶层字段）。
// 忽略 meta.IgnoreFields、空值变更与未变化字段。
func computeDiff(before, after map[string]any, meta Meta) map[string]DiffItem {
	if before == nil && after == nil {
		return nil
	}
	ignore := make(map[string]struct{}, len(meta.IgnoreFields)+8)
	for _, f := range meta.IgnoreFields {
		ignore[normalizeKey(f)] = struct{}{}
	}
	// 审计字段一律不进入 diff
	for _, f := range []string{"id", "create_at", "update_at", "create_by", "update_by"} {
		ignore[normalizeKey(f)] = struct{}{}
	}

	diff := make(map[string]DiffItem)
	for key, beforeValue := range before {
		if _, ok := ignore[normalizeKey(key)]; ok {
			continue
		}
		afterValue, exists := after[key]
		if !exists {
			diff[key] = DiffItem{Before: beforeValue, After: nil}
			continue
		}
		if reflect.DeepEqual(beforeValue, afterValue) {
			continue
		}
		diff[key] = DiffItem{Before: beforeValue, After: afterValue}
	}
	for key, afterValue := range after {
		if _, ok := ignore[normalizeKey(key)]; ok {
			continue
		}
		if _, exists := before[key]; !exists {
			diff[key] = DiffItem{Before: nil, After: afterValue}
		}
	}
	if len(diff) == 0 {
		return nil
	}
	return diff
}

// hasMeaningfulChange 判断 before/after 是否存在实际业务变化（用于 SkipNoChange）。
func hasMeaningfulChange(before, after map[string]any, meta Meta) bool {
	return len(computeDiff(before, after, meta)) > 0
}

// marshalDiff 序列化 diff，空 diff 返回空串。
func marshalDiff(diff map[string]DiffItem, max int) string {
	if len(diff) == 0 {
		return ""
	}
	bs, err := json.Marshal(diff)
	if err != nil {
		return ""
	}
	return truncate(string(bs), max)
}

// resultFromPanic 把 panic 值与堆栈拼成失败原因。
func resultFromPanic(panicValue any) string {
	if panicValue == nil {
		return ""
	}
	switch v := panicValue.(type) {
	case string:
		return strings.TrimSpace(v)
	case error:
		return strings.TrimSpace(v.Error())
	default:
		bs, err := json.Marshal(v)
		if err != nil {
			return ""
		}
		return string(bs)
	}
}

// resultValue 计算最终 result 枚举。
func resultValue(ac *AuditContext) int8 {
	if ac.Result != nil {
		if ac.Result.Partial {
			return entityBasic.OperationLogResultPartial
		}
		if ac.Result.Success {
			return entityBasic.OperationLogResultSuccess
		}
		return entityBasic.OperationLogResultFail
	}
	if ac.PanicValue != nil {
		return entityBasic.OperationLogResultFail
	}
	if ac.HTTPStatus >= 400 {
		return entityBasic.OperationLogResultFail
	}
	if ac.HTTPStatus == 0 {
		return entityBasic.OperationLogResultUnknown
	}
	// JSON 响应 code 判定在 parseResponse 阶段写入
	if ac.ErrorMessage != "" {
		return entityBasic.OperationLogResultFail
	}
	return entityBasic.OperationLogResultSuccess
}
