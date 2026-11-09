package idp_e2e

import (
	"strings"
	"testing"
)

type ldapTestVO struct {
	Connected bool   `json:"connected"`
	IsAD      bool   `json:"isAd"`
	Message   string `json:"message"`
	UserCount int    `json:"userCount"`
}

func runLDAP(t *testing.T, e *env) {
	e.t = t
	ctrl := baseURL + "/ldap"

	pid, pno := e.makeIdpProvider(e.unique("provider-ldap"), "ldap")
	t.Cleanup(func() { _ = pid })
	sid := e.createID(baseURL+"/identity-source", js(map[string]any{
		"name": e.unique("source-ldap"), "protocol": "ldap", "idp": pno,
		"baseConfig": js(map[string]any{
			"hostUrl": "127.0.0.1", "clientId": "cn=admin,dc=e2e,dc=com", "clientSecret": "bind-pwd",
		}),
		"syncStrategyConfig": js(map[string]any{
			"port": 1, "baseDn": "dc=e2e,dc=com", "filter": "(uid=*)", "enableSsl": false,
		}),
	}), "ram_identity_source")
	sno := e.sourceNo(sid)

	// 使用未开放的 127.0.0.1:1，确保连接快速失败且不依赖外部 LDAP。
	t.Run("test-connection-fails", func(t *testing.T) {
		e.t = t
		r := e.postNoAuth(ctrl+"/test", js(map[string]any{"sourceNo": sno}))
		if r.Code != "200" {
			t.Fatalf("连接测试接口期望返回业务结果，实际 code=%s msg=%s", r.Code, r.Message)
		}
		vo := decodeData[ldapTestVO](t, r)
		if vo.Connected {
			t.Fatalf("127.0.0.1:1 不应连接成功: %+v", vo)
		}
		if !strings.Contains(vo.Message, "连接失败") && !strings.Contains(vo.Message, "搜索用户失败") {
			t.Fatalf("连接失败提示异常: %+v", vo)
		}
	})
	t.Run("test-unknown-source", func(t *testing.T) {
		e.t = t
		r := e.postNoAuth(ctrl+"/test", js(map[string]any{"sourceNo": "not-exist-" + e.runID}))
		if r.Code == "200" || !strings.Contains(r.Message, "认证源不存在") {
			t.Fatalf("未知认证源期望失败，实际 code=%s msg=%s", r.Code, r.Message)
		}
	})
	t.Run("test-missing-source", func(t *testing.T) {
		e.t = t
		r := e.postNoAuth(ctrl+"/test", `{}`)
		if r.Code == "200" {
			t.Fatalf("缺 sourceNo 期望失败: %s", r.Body)
		}
	})
	t.Run("search-connection-fails", func(t *testing.T) {
		e.t = t
		r := e.postNoAuth(ctrl+"/search", js(map[string]any{"sourceNo": sno, "filter": "(cn=e2e)"}))
		if r.Code == "200" || !strings.Contains(r.Message, "LDAP 连接失败") {
			t.Fatalf("搜索期望连接失败，实际 code=%s msg=%s", r.Code, r.Message)
		}
	})
	t.Run("sync-connection-fails", func(t *testing.T) {
		e.t = t
		r := e.postNoAuth(ctrl+"/sync", js(map[string]any{"sourceNo": sno}))
		if r.Code == "200" || !strings.Contains(r.Message, "LDAP 连接失败") {
			t.Fatalf("同步期望连接失败，实际 code=%s msg=%s", r.Code, r.Message)
		}
	})
	t.Run("login-connection-fails", func(t *testing.T) {
		e.t = t
		r := e.postNoAuth(ctrl+"/login", js(map[string]any{"sourceNo": sno, "username": "e2e", "password": "x"}))
		if r.Code == "200" || !strings.Contains(r.Message, "LDAP 连接失败") {
			t.Fatalf("登录期望连接失败，实际 code=%s msg=%s", r.Code, r.Message)
		}
	})
	t.Run("login-source-disabled", func(t *testing.T) {
		e.t = t
		e.ok(baseURL+"/identity-source/disable", js(map[string]any{"ids": []string{sid}}))
		r := e.postNoAuth(ctrl+"/login", js(map[string]any{"sourceNo": sno, "username": "e2e", "password": "x"}))
		if r.Code == "200" || !strings.Contains(r.Message, "认证源已禁用") {
			t.Fatalf("停用认证源期望失败，实际 code=%s msg=%s", r.Code, r.Message)
		}
		e.ok(baseURL+"/identity-source/enable", js(map[string]any{"ids": []string{sid}}))
	})
	t.Run("login-missing-password", func(t *testing.T) {
		e.t = t
		r := e.postNoAuth(ctrl+"/login", js(map[string]any{"sourceNo": sno, "username": "e2e"}))
		if r.Code == "200" {
			t.Fatalf("缺 password 期望失败: %s", r.Body)
		}
	})
	t.Run("login-does-not-check-protocol", func(t *testing.T) {
		e.t = t
		// 当前实现未校验认证源 protocol，非 LDAP 源也会走 LDAP 连接逻辑。
		_, pno2 := e.makeIdpProvider(e.unique("provider-oauth-for-ldap"), "oauth2")
		oauthSid := e.createID(baseURL+"/identity-source", js(map[string]any{
			"name": e.unique("source-oauth-for-ldap"), "protocol": "oauth2", "idp": pno2,
			"baseConfig":         js(map[string]any{"hostUrl": "127.0.0.1"}),
			"syncStrategyConfig": js(map[string]any{"port": 1}),
		}), "ram_identity_source")
		r := e.postNoAuth(ctrl+"/login", js(map[string]any{
			"sourceNo": e.sourceNo(oauthSid), "username": "e2e", "password": "x",
		}))
		if r.Code == "200" || !strings.Contains(r.Message, "LDAP 连接失败") {
			t.Fatalf("非 LDAP 协议期望连接失败（说明未校验 protocol），实际 code=%s msg=%s", r.Code, r.Message)
		}
	})
}
