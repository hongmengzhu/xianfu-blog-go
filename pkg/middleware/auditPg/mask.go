package auditPg

import (
	"encoding/json"
	"net/url"
	"sort"
	"strings"
	"unicode/utf8"
)

// maskedValue 敏感字段的占位值。
const maskedValue = "******"

// sensitiveKeyNames 敏感字段（大小写不敏感，含常见变体）。
var sensitiveKeyNames = []string{
	"password", "passwd", "pwd", "old_password", "new_password",
	"confirm_password", "token", "access_token", "refresh_token",
	"id_token", "authorization", "cookie", "set-cookie",
	"secret", "client_secret", "private_key", "api_key",
	"session_id", "sessionid", "captcha", "sms_code", "otp",
	"mfa_code", "id_card", "idcard", "bank_card", "bankcard",
	"cvv", "expiry", "phone", "mobile", "email", "address",
	// 边界场景补充
	"passport", "credit_card", "ssn", "signature", "webhook_secret",
}

var sensitiveKeys = func() map[string]struct{} {
	m := make(map[string]struct{}, len(sensitiveKeyNames))
	for _, k := range sensitiveKeyNames {
		m[normalizeKey(k)] = struct{}{}
	}
	return m
}()

func normalizeKey(key string) string {
	normalized := strings.ToLower(strings.TrimSpace(key))
	normalized = strings.ReplaceAll(normalized, "-", "")
	normalized = strings.ReplaceAll(normalized, "_", "")
	return normalized
}

func isSensitiveKey(key string) bool {
	_, ok := sensitiveKeys[normalizeKey(key)]
	return ok
}

// maskValue 递归脱敏：敏感 key 的值整体替换为占位串。
func maskValue(key string, v any) any {
	if isSensitiveKey(key) && v != nil {
		return maskedValue
	}
	switch t := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(t))
		for k, item := range t {
			out[k] = maskValue(k, item)
		}
		return out
	case []any:
		out := make([]any, 0, len(t))
		for _, item := range t {
			out = append(out, maskValue("", item))
		}
		return out
	default:
		return v
	}
}

// maskJSON 脱敏并序列化；非法 JSON 原样截断返回。
func maskJSON(raw []byte, max int) string {
	if len(raw) == 0 {
		return ""
	}
	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		return truncate(bytesToString(raw), max)
	}
	masked := maskValue("", v)
	bs, err := json.Marshal(masked)
	if err != nil {
		return truncate(string(raw), max)
	}
	return truncate(bytesToString(bs), max)
}

// maskQuery 对 query string 做脱敏；保留结构，敏感值替换。
func maskQuery(raw string, max int) string {
	if raw == "" {
		return ""
	}
	values, err := url.ParseQuery(raw)
	if err != nil {
		return truncate(raw, max)
	}
	for key := range values {
		if isSensitiveKey(key) {
			values[key] = []string{maskedValue}
		}
	}
	return truncate(values.Encode(), max)
}

// maskMap 对 map 做脱敏，用于 extra_data。
func maskMap(in map[string]any) map[string]any {
	if in == nil {
		return nil
	}
	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = maskValue(k, v)
	}
	return out
}

// truncate 按字节安全截断，末尾加省略号。
func truncate(s string, max int) string {
	if max <= 0 || len(s) <= max {
		return s
	}
	cut := max
	for cut > 0 && !utf8.RuneStart(s[cut]) {
		cut--
	}
	return strings.TrimSuffix(s[:cut], "\x00") + "…"
}

func bytesToString(b []byte) string {
	return string(b)
}

// snapshot 将业务对象转为可读快照：白名单、忽略字段、脱敏，最后序列化。
// 返回 nil map 表示没有可用快照。
func snapshot(value any, meta Meta) (map[string]any, string) {
	if value == nil {
		return nil, ""
	}
	raw, err := json.Marshal(value)
	if err != nil {
		return nil, ""
	}
	var m map[string]any
	if err := json.Unmarshal(raw, &m); err != nil {
		// 非对象（数组/标量）直接脱敏返回
		var anyValue any
		if err := json.Unmarshal(raw, &anyValue); err != nil {
			return nil, ""
		}
		bs, _ := json.Marshal(maskValue("", anyValue))
		return nil, truncate(bytesToString(bs), meta.MaxBodyBytes)
	}

	whitelist := make(map[string]struct{}, len(meta.SnapshotFields))
	for _, f := range meta.SnapshotFields {
		whitelist[strings.TrimSpace(f)] = struct{}{}
	}
	ignore := make(map[string]struct{}, len(meta.IgnoreFields))
	for _, f := range meta.IgnoreFields {
		ignore[strings.TrimSpace(f)] = struct{}{}
	}

	out := make(map[string]any, len(m))
	for k, v := range m {
		if len(whitelist) > 0 {
			if _, ok := whitelist[k]; !ok {
				continue
			}
		}
		if _, ok := ignore[k]; ok {
			continue
		}
		out[k] = maskValue(k, v)
	}
	bs, err := json.Marshal(out)
	if err != nil {
		return out, ""
	}
	return out, truncate(bytesToString(bs), meta.MaxBodyBytes)
}

// sortedKeys 仅用于测试/稳定输出。
func sortedKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
