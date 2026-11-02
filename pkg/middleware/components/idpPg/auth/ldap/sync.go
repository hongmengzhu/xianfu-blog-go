package ldap

import (
	"fmt"
	"strings"
)

// SyncResult LDAP 用户同步结果
type SyncResult struct {
	NewUsers     int        `json:"newUsers"`
	UpdatedUsers int        `json:"updatedUsers"`
	FailedUsers  int        `json:"failedUsers"`
	Errors       []string   `json:"errors,omitempty"`
}

// SyncUsers 从 LDAP 同步用户到本地（简化版）
// 返回同步结果
func SyncUsers(conn *LdapConn, cfg *LdapConfig) (*SyncResult, error) {
	users, err := conn.SearchUsers(cfg)
	if err != nil {
		return nil, fmt.Errorf("搜索 LDAP 用户失败: %w", err)
	}

	result := &SyncResult{}

	for _, user := range users {
		uuid := user.GetLdapUuid()
		if uuid == "" {
			result.FailedUsers++
			result.Errors = append(result.Errors, fmt.Sprintf("用户缺少唯一标识: cn=%s", user.Cn))
			continue
		}

		// 简化版：仅统计，实际同步逻辑由 LdapService 处理
		result.NewUsers++
	}

	return result, nil
}

// ParseBaseDn 从 BaseDN 中提取域名部分
// 例如 "dc=example,dc=com" -> "example.com"
func ParseBaseDn(baseDn string) string {
	var parts []string
	for _, component := range strings.Split(baseDn, ",") {
		component = strings.TrimSpace(component)
		if strings.HasPrefix(strings.ToLower(component), "dc=") {
			parts = append(parts, component[3:])
		}
	}
	return strings.Join(parts, ".")
}

// ParseOuFromUsername 从 bind username 中提取 OU 部分
// 例如 "cn=admin,ou=People,dc=example,dc=com" -> "People"
func ParseOuFromUsername(username string) string {
	var parts []string
	for _, component := range strings.Split(username, ",") {
		component = strings.TrimSpace(component)
		if strings.HasPrefix(strings.ToLower(component), "ou=") {
			parts = append(parts, component[3:])
		}
	}
	return strings.Join(parts, ".")
}
