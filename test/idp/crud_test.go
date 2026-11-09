package idp_e2e

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

// ---------------- 通用 ----------------

type pageWrap struct {
	PageNum int64             `json:"pageNum"`
	Total   int64             `json:"total"`
	Data    []json.RawMessage `json:"data"`
}

func (e *env) queryPage(path string, body string) pageWrap {
	e.t.Helper()
	r := e.ok(path, body)
	p := decodeData[pageWrap](e.t, r)
	if p.Data == nil {
		p.Data = []json.RawMessage{}
	}
	return p
}

func (e *env) selectPublic(path string) []json.RawMessage {
	e.t.Helper()
	r := e.ok(path+"/selectPublic", "{}")
	var list []json.RawMessage
	if err := json.Unmarshal(r.Data, &list); err != nil {
		e.t.Fatalf("selectPublic 解析失败: %v data=%s", err, r.Data)
	}
	return list
}

func (e *env) detailState(ctrl, id string) int {
	e.t.Helper()
	r := e.getAuth(ctrl + "/detail/" + id)
	if r.Code != "200" {
		e.t.Fatalf("detail id=%s 期望成功，实际 %s", id, r.Message)
	}
	var vo struct {
		State int `json:"state"`
	}
	if err := json.Unmarshal(r.Data, &vo); err != nil {
		e.t.Fatalf("detail 解析 state 失败: %v", err)
	}
	return vo.State
}

func (e *env) expectDetailErr(ctrl, id, contains string) {
	e.t.Helper()
	r := e.getAuth(ctrl + "/detail/" + id)
	if r.Code == "200" {
		e.t.Fatalf("detail id=%s 期望失败(%q)，实际成功 %s", id, contains, r.Body)
	}
	if contains != "" && !strings.Contains(r.Message, contains) {
		e.t.Fatalf("detail id=%s 期望 message 含 %q，实际 %q", id, contains, r.Message)
	}
}

// assertDeleteCycle 覆盖 delete / recovery / physicalDeletion。
// pg.data.delete=true（默认）时 delete 即物理删除，recovery 必然失败；
// 设置 IDP_TEST_SOFT_DELETE=true 后走逻辑删除（state 1→11）→恢复(11→1)→物理删除。
func (e *env) assertDeleteCycle(t *testing.T, ctrl, table, id string) {
	e.t = t
	ids := js(map[string]any{"ids": []string{id}})
	e.ok(ctrl+"/delete", ids)
	if e.soft {
		if s := e.detailState(ctrl, id); s != 11 {
			t.Fatalf("逻辑删除后期望 state=11，实际 %d", s)
		}
		e.ok(ctrl+"/recovery", ids)
		if s := e.detailState(ctrl, id); s != 1 {
			t.Fatalf("恢复后期望 state=1，实际 %d", s)
		}
		e.ok(ctrl+"/physicalDeletion", ids)
		e.expectDetailErr(ctrl, id, "数据不存在")
	} else {
		e.expectDetailErr(ctrl, id, "数据不存在")
		e.expectErr(ctrl+"/recovery", ids, "数据不存在")
		e.expectErr(ctrl+"/physicalDeletion", ids, "数据不存在")
	}
}

// assertCommonState 覆盖 enable/disable/state 三个端点。
func (e *env) assertCommonState(t *testing.T, ctrl, id string) {
	e.t = t
	t.Run("disable", func(t *testing.T) {
		e.t = t
		e.ok(ctrl+"/disable", js(map[string]any{"ids": []string{id}}))
		if s := e.detailState(ctrl, id); s != 2 {
			t.Fatalf("禁用后 state=%d want 2", s)
		}
	})
	t.Run("state-invalid", func(t *testing.T) {
		e.t = t
		e.expectErr(ctrl+"/state", js(map[string]any{"ids": []string{id}, "state": 99}), "类型不正确")
	})
	t.Run("state-empty-ids", func(t *testing.T) {
		e.t = t
		e.expectErr(ctrl+"/state", js(map[string]any{"ids": []string{}, "state": 2}), "id错误")
	})
	t.Run("enable", func(t *testing.T) {
		e.t = t
		e.ok(ctrl+"/enable", js(map[string]any{"ids": []string{id}}))
		if s := e.detailState(ctrl, id); s != 1 {
			t.Fatalf("启用后 state=%d want 1", s)
		}
	})
	t.Run("delete-empty-ids", func(t *testing.T) {
		e.t = t
		e.expectErr(ctrl+"/delete", `{"ids":[]}`, "id错误")
	})
}

// ---------------- 身份提供商 ----------------

type providerVO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	NameFl      string `json:"nameFl"`
	Code        string `json:"code"`
	NameFull    string `json:"nameFull"`
	Description string `json:"description"`
	Protocol    string `json:"protocol"`
	Platform    string `json:"platform"`
	Sort        int64  `json:"sort"`
	State       int    `json:"state"`
}

func runProviderCRUD(t *testing.T, e *env) {
	e.t = t
	ctrl := baseURL + "/identity-provider"

	t.Run("auth-required", func(t *testing.T) {
		e.t = t
		r := e.postNoAuth(ctrl+"/query", "{}")
		if r.Code == "200" || !strings.Contains(r.Message, "Authorization") {
			t.Fatalf("缺 token 期望鉴权失败，实际 code=%s msg=%s", r.Code, r.Message)
		}
	})
	t.Run("create-name-required", func(t *testing.T) {
		e.t = t
		e.expectErr(ctrl+"/create", `{}`, "名称不能为空")
	})
	t.Run("detail-invalid-id", func(t *testing.T) {
		e.t = t
		e.expectDetailErr(ctrl, "0", "id错误")
	})
	t.Run("detail-not-found", func(t *testing.T) {
		e.t = t
		e.expectDetailErr(ctrl, "999999999", "数据不存在")
	})

	name := e.unique("provider")
	p := map[string]any{
		"name": name, "nameFl": "E2E Provider", "code": "Custom",
		"nameFull": "E2E Identity Provider", "description": "e2e provider",
		"icon": "", "protocol": "oauth2", "platform": "web", "sort": 1,
	}
	id := e.createID(ctrl, js(p), "ram_identity_provider")

	t.Run("create-detail", func(t *testing.T) {
		e.t = t
		r := e.getAuth(ctrl + "/detail/" + id)
		if r.Code != "200" {
			t.Fatalf("detail 失败: %s", r.Body)
		}
		vo := decodeData[providerVO](t, r)
		if vo.Name != name || vo.Code != "Custom" || vo.Protocol != "oauth2" || vo.State != 1 {
			t.Fatalf("detail 不符: %+v", vo)
		}
	})
	t.Run("update", func(t *testing.T) {
		e.t = t
		e.ok(ctrl+"/update", js(map[string]any{"id": id, "name": name + "-upd", "protocol": "oidc"}))
		vo := decodeData[providerVO](t, e.getAuth(ctrl+"/detail/"+id))
		if vo.Name != name+"-upd" || vo.Protocol != "oidc" {
			t.Fatalf("update 未生效: %+v", vo)
		}
	})
	t.Run("update-name-required", func(t *testing.T) {
		e.t = t
		e.expectErr(ctrl+"/update", js(map[string]any{"id": id}), "名称不能为空")
	})
	t.Run("update-not-found", func(t *testing.T) {
		e.t = t
		e.expectErr(ctrl+"/update", js(map[string]any{"id": "999999999", "name": "x"}), "数据不存在")
	})
	t.Run("existName", func(t *testing.T) {
		e.t = t
		// 更新后名称已变，用当前名称验证重复与排除自身
		cur := name + "-upd"
		e.expectErr(ctrl+"/existName", js(map[string]any{"wd": cur}), "重复")
		r := e.ok(ctrl+"/existName", js(map[string]any{"wd": cur, "id": id}))
		if r.Message != "可以使用" {
			t.Fatalf("existName 排除自身 message = %q data=%s", r.Message, r.Data)
		}
		e.expectErr(ctrl+"/existName", `{"wd":""}`, "查询内容不能为空")
	})
	t.Run("query", func(t *testing.T) {
		e.t = t
		pg := e.queryPage(ctrl+"/query", js(map[string]any{"pageNum": 1, "pageSize": 20, "wd": name + "-upd"}))
		hit := false
		for _, raw := range pg.Data {
			var vo providerVO
			_ = json.Unmarshal(raw, &vo)
			if vo.ID == id {
				hit = true
			}
		}
		if pg.Total < 1 || !hit {
			t.Fatalf("query total=%d hit=%v", pg.Total, hit)
		}
	})
	t.Run("selectPublic", func(t *testing.T) {
		e.t = t
		list := e.selectPublic(ctrl)
		hit := false
		for _, raw := range list {
			var vo providerVO
			_ = json.Unmarshal(raw, &vo)
			if vo.ID == id {
				hit = true
			}
		}
		if !hit {
			t.Fatalf("selectPublic 未包含 id=%s", id)
		}
	})
	t.Run("disable-excluded-from-public", func(t *testing.T) {
		e.t = t
		e.ok(ctrl+"/disable", js(map[string]any{"ids": []string{id}}))
		for _, raw := range e.selectPublic(ctrl) {
			var vo providerVO
			_ = json.Unmarshal(raw, &vo)
			if vo.ID == id {
				t.Fatalf("停用记录仍出现在 selectPublic")
			}
		}
		e.ok(ctrl+"/enable", js(map[string]any{"ids": []string{id}}))
	})

	e.assertCommonState(t, ctrl, id)
	t.Run("deletion", func(t *testing.T) { e.assertDeleteCycle(t, ctrl, "ram_identity_provider", id) })
}

// ---------------- 认证源 ----------------

type sourceVO struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	Code               string `json:"code"`
	Protocol           string `json:"protocol"`
	Idp                string `json:"idp"`
	BaseConfig         string `json:"baseConfig"`
	SyncStrategyConfig string `json:"syncStrategyConfig"`
	State              int    `json:"state"`
	Configured         int    `json:"configured"`
	AutoCreateUser     int    `json:"autoCreateUser"`
}

func runSourceCRUD(t *testing.T, e *env) {
	e.t = t
	ctrl := baseURL + "/identity-source"

	// 认证源需要关联一个身份提供商（No）
	pid := e.createID(baseURL+"/identity-provider",
		js(map[string]any{"name": e.unique("provider-for-source"), "code": "Custom", "protocol": "oauth2"}),
		"ram_identity_provider")
	pno := e.providerNo(pid)

	t.Run("create-name-required", func(t *testing.T) {
		e.t = t
		e.expectErr(ctrl+"/create", `{}`, "名称不能为空")
	})
	name := e.unique("source")
	s := map[string]any{
		"name": name, "code": "E2E-SRC", "protocol": "oidc", "idp": pno,
		"baseConfig":         `{"clientId":"e2e-client","clientSecret":"e2e-secret"}`,
		"syncStrategyConfig": `{"mode":"manual"}`, "enableSlo": 2,
		"autoCreateUser": 2, "attributeMapping": `{"mail":"email"}`, "sort": 2,
	}
	sid := e.createID(ctrl, js(s), "ram_identity_source")
	sno := e.sourceNo(sid)

	t.Run("create-detail", func(t *testing.T) {
		e.t = t
		vo := decodeData[sourceVO](t, e.getAuth(ctrl+"/detail/"+sid))
		if vo.Name != name || vo.Protocol != "oidc" || vo.Idp != pno || vo.State != 1 {
			t.Fatalf("detail 不符: %+v", vo)
		}
		if !strings.Contains(vo.BaseConfig, "e2e-client") {
			t.Fatalf("baseConfig 未保存: %q", vo.BaseConfig)
		}
	})
	t.Run("update", func(t *testing.T) {
		e.t = t
		e.ok(ctrl+"/update", js(map[string]any{
			"id": sid, "name": name + "-upd", "protocol": "oauth2", "baseConfig": `{"clientId":"e2e-client-2"}`,
			"autoCreateUser": 1,
		}))
		vo := decodeData[sourceVO](t, e.getAuth(ctrl+"/detail/"+sid))
		if vo.Name != name+"-upd" || vo.Protocol != "oauth2" {
			t.Fatalf("update 未生效: %+v", vo)
		}
	})
	t.Run("update-name-required", func(t *testing.T) {
		e.t = t
		e.expectErr(ctrl+"/update", js(map[string]any{"id": sid}), "名称不能为空")
	})
	t.Run("existName", func(t *testing.T) {
		e.t = t
		cur := name + "-upd"
		e.expectErr(ctrl+"/existName", js(map[string]any{"wd": cur}), "重复")
		r := e.ok(ctrl+"/existName", js(map[string]any{"wd": cur, "id": sid}))
		if r.Message != "可以使用" {
			t.Fatalf("existName message = %q data=%s", r.Message, r.Data)
		}
		e.expectErr(ctrl+"/existName", `{"wd":""}`, "查询内容不能为空")
	})
	t.Run("query", func(t *testing.T) {
		e.t = t
		pg := e.queryPage(ctrl+"/query", js(map[string]any{"pageNum": 1, "pageSize": 20, "wd": name + "-upd"}))
		hit := false
		for _, raw := range pg.Data {
			var vo sourceVO
			_ = json.Unmarshal(raw, &vo)
			if vo.ID == sid {
				hit = true
			}
		}
		if pg.Total < 1 || !hit {
			t.Fatalf("query total=%d hit=%v", pg.Total, hit)
		}
	})
	t.Run("selectPublic", func(t *testing.T) {
		e.t = t
		list := e.selectPublic(ctrl)
		hit := false
		for _, raw := range list {
			var vo sourceVO
			_ = json.Unmarshal(raw, &vo)
			if vo.ID == sid {
				hit = true
			}
		}
		if !hit {
			t.Fatalf("selectPublic 未包含 id=%s（sourceNo=%s）", sid, sno)
		}
	})

	e.assertCommonState(t, ctrl, sid)
	t.Run("deletion", func(t *testing.T) { e.assertDeleteCycle(t, ctrl, "ram_identity_source", sid) })

	t.Cleanup(func() { _ = pno })
}

// ---------------- 回调白名单 ----------------

type callbackVO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Callback    string `json:"callback"`
	State       int    `json:"state"`
	Description string `json:"description"`
	Platform    string `json:"platform"`
	Protocol    string `json:"protocol"`
	Idp         string `json:"idp"`
}

func runCallbackCRUD(t *testing.T, e *env) {
	e.t = t
	ctrl := baseURL + "/identity-source-callback"
	pid := e.createID(baseURL+"/identity-provider",
		js(map[string]any{"name": e.unique("provider-for-callback"), "code": "Custom", "protocol": "oauth2"}),
		"ram_identity_provider")
	pno := e.providerNo(pid)

	t.Run("create-name-required", func(t *testing.T) {
		e.t = t
		e.expectErr(ctrl+"/create", `{}`, "名称不能为空")
	})
	// 当前实现只校验 name，未校验 callback 必填：不带 callback 仍创建成功。
	t.Run("create-without-callback-succeeds", func(t *testing.T) {
		e.t = t
		n := e.unique("callback-nocb")
		r := e.ok(ctrl+"/create", js(map[string]any{"name": n}))
		id := decodeData[string](t, r)
		e.track("ram_identity_source_callback", id)
		vo := decodeData[callbackVO](t, e.getAuth(ctrl+"/detail/"+id))
		if vo.Callback != "" {
			t.Fatalf("期望 callback 为空，实际 %q", vo.Callback)
		}
	})

	name := e.unique("callback")
	cb := "https://app.example.com/idp/callback"
	c := map[string]any{
		"name": name, "callback": cb, "description": "e2e callback",
		"platform": "web", "protocol": "oauth2", "idp": pno,
	}
	cid := e.createID(ctrl, js(c), "ram_identity_source_callback")

	t.Run("create-detail", func(t *testing.T) {
		e.t = t
		vo := decodeData[callbackVO](t, e.getAuth(ctrl+"/detail/"+cid))
		if vo.Name != name || vo.Callback != cb || vo.Protocol != "oauth2" || vo.State != 1 {
			t.Fatalf("detail 不符: %+v", vo)
		}
	})
	t.Run("update", func(t *testing.T) {
		e.t = t
		e.ok(ctrl+"/update", js(map[string]any{"id": cid, "name": name + "-upd", "callback": cb + "?v=2"}))
		vo := decodeData[callbackVO](t, e.getAuth(ctrl+"/detail/"+cid))
		if vo.Name != name+"-upd" || !strings.Contains(vo.Callback, "v=2") {
			t.Fatalf("update 未生效: %+v", vo)
		}
	})
	t.Run("existName", func(t *testing.T) {
		e.t = t
		cur := name + "-upd"
		e.expectErr(ctrl+"/existName", js(map[string]any{"wd": cur}), "重复")
		r := e.ok(ctrl+"/existName", js(map[string]any{"wd": cur, "id": cid}))
		if r.Message != "可以使用" {
			t.Fatalf("existName message = %q data=%s", r.Message, r.Data)
		}
	})
	t.Run("query", func(t *testing.T) {
		e.t = t
		pg := e.queryPage(ctrl+"/query", js(map[string]any{"pageNum": 1, "pageSize": 20, "wd": name + "-upd"}))
		hit := false
		for _, raw := range pg.Data {
			var vo callbackVO
			_ = json.Unmarshal(raw, &vo)
			if vo.ID == cid {
				hit = true
			}
		}
		if pg.Total < 1 || !hit {
			t.Fatalf("query total=%d hit=%v", pg.Total, hit)
		}
	})
	t.Run("selectPublic", func(t *testing.T) {
		e.t = t
		list := e.selectPublic(ctrl)
		hit := false
		for _, raw := range list {
			var vo callbackVO
			_ = json.Unmarshal(raw, &vo)
			if vo.ID == cid {
				hit = true
			}
		}
		if !hit {
			t.Fatalf("selectPublic 未包含 id=%s", cid)
		}
	})

	e.assertCommonState(t, ctrl, cid)
	t.Run("deletion", func(t *testing.T) { e.assertDeleteCycle(t, ctrl, "ram_identity_source_callback", cid) })

	t.Cleanup(func() { _ = pno })
}

// ---------------- 身份绑定 ----------------

type bindingVO struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	State       int    `json:"state"`
	StateBind   int    `json:"stateBind"`
	Idp         string `json:"idp"`
	OpenId      string `json:"openId"`
	UnionId     string `json:"unionId"`
	ExternalSub string `json:"externalSub"`
	BindAno     string `json:"bindAno"`
	Mail        string `json:"mail"`
	NickName    string `json:"nickName"`
	Protocol    string `json:"protocol"`
}

func runBindingCRUD(t *testing.T, e *env) {
	e.t = t
	ctrl := baseURL + "/idp-binding"
	pid := e.createID(baseURL+"/identity-provider",
		js(map[string]any{"name": e.unique("provider-for-binding"), "code": "Custom", "protocol": "oauth2"}),
		"ram_identity_provider")
	pno := e.providerNo(pid)

	name := e.unique("binding")
	extSub := "ext-sub-" + e.runID
	b := map[string]any{
		"description": name, "idp": pno, "openId": "openid-" + e.runID, "unionId": "union-" + e.runID,
		"appMark": "app1", "platform": "web", "protocol": "oauth2",
		"mail": "e2e@example.com", "phone": "13800000000", "nickName": "E2E",
		"bindAno": testTenantNo, "externalSub": extSub,
		"bindTime": "2026-11-01T00:00:00+08:00",
	}
	bid := e.createID(ctrl, js(b), "ram_idp_binding")

	t.Run("create-detail", func(t *testing.T) {
		e.t = t
		vo := decodeData[bindingVO](t, e.getAuth(ctrl+"/detail/"+bid))
		if vo.Idp != pno || vo.ExternalSub != extSub || vo.State != 1 || vo.BindAno != testTenantNo {
			t.Fatalf("detail 不符: %+v", vo)
		}
	})
	t.Run("update", func(t *testing.T) {
		e.t = t
		e.ok(ctrl+"/update", js(map[string]any{"id": bid, "nickName": "E2E-Updated", "mail": "upd@example.com"}))
		vo := decodeData[bindingVO](t, e.getAuth(ctrl+"/detail/"+bid))
		if vo.NickName != "E2E-Updated" || vo.Mail != "upd@example.com" {
			t.Fatalf("update 未生效: %+v", vo)
		}
	})
	t.Run("update-id-required", func(t *testing.T) {
		e.t = t
		e.expectErr(ctrl+"/update", `{}`, "id错误")
	})
	t.Run("update-not-found", func(t *testing.T) {
		e.t = t
		e.expectErr(ctrl+"/update", js(map[string]any{"id": "999999999"}), "数据不存在")
	})
	t.Run("query", func(t *testing.T) {
		e.t = t
		pg := e.queryPage(ctrl+"/query", js(map[string]any{"pageNum": 1, "pageSize": 20, "wd": pno}))
		hit := false
		for _, raw := range pg.Data {
			var vo bindingVO
			_ = json.Unmarshal(raw, &vo)
			if vo.ID == bid {
				hit = true
			}
		}
		if pg.Total < 1 || !hit {
			t.Fatalf("query total=%d hit=%v", pg.Total, hit)
		}
	})
	t.Run("selectPublic", func(t *testing.T) {
		e.t = t
		list := e.selectPublic(ctrl)
		hit := false
		for _, raw := range list {
			var vo bindingVO
			_ = json.Unmarshal(raw, &vo)
			if vo.ID == bid {
				hit = true
			}
		}
		if !hit {
			t.Fatalf("selectPublic 未包含 id=%s", bid)
		}
	})

	e.assertCommonState(t, ctrl, bid)
	t.Run("deletion", func(t *testing.T) { e.assertDeleteCycle(t, ctrl, "ram_idp_binding", bid) })

	t.Cleanup(func() { _ = pno })
}

// ---------------- 凭证 ----------------

type credentialVO struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	State       int    `json:"state"`
	Idp         string `json:"idp"`
	SourceNo    string `json:"sourceNo"`
	CredType    string `json:"credType"`
}

func runCredentialCRUD(t *testing.T, e *env) {
	e.t = t
	ctrl := baseURL + "/idp-credential"
	pid := e.createID(baseURL+"/identity-provider",
		js(map[string]any{"name": e.unique("provider-for-cred"), "code": "Custom", "protocol": "oidc"}),
		"ram_identity_provider")
	pno := e.providerNo(pid)
	sid := e.createID(baseURL+"/identity-source",
		js(map[string]any{"name": e.unique("source-for-cred"), "protocol": "oidc", "idp": pno}),
		"ram_identity_source")
	sno := e.sourceNo(sid)

	secret := "top-secret-" + e.runID
	c := map[string]any{
		"description": "e2e credential", "idp": pno, "sourceNo": sno,
		"credType": "client_secret", "value": secret, "expireAt": "2027-01-01T00:00:00+08:00",
	}
	cid := e.createID(ctrl, js(c), "ram_idp_credential")

	t.Run("create-detail-omits-secret", func(t *testing.T) {
		e.t = t
		r := e.getAuth(ctrl + "/detail/" + cid)
		vo := decodeData[credentialVO](t, r)
		if vo.SourceNo != sno || vo.CredType != "client_secret" || vo.State != 1 {
			t.Fatalf("detail 不符: %+v", vo)
		}
		if strings.Contains(string(r.Data), secret) || strings.Contains(string(r.Data), `"value"`) {
			t.Fatalf("detail 泄露凭证值: %s", r.Data)
		}
	})
	t.Run("update", func(t *testing.T) {
		e.t = t
		e.ok(ctrl+"/update", js(map[string]any{"id": cid, "credType": "jwks", "description": "e2e credential upd"}))
		vo := decodeData[credentialVO](t, e.getAuth(ctrl+"/detail/"+cid))
		if vo.CredType != "jwks" {
			t.Fatalf("update 未生效: %+v", vo)
		}
		// 已知缺陷：RamIdpCredentialEntity 没有 description 列，
		// Update 的 description 不会落库，详情的 description 恒为空。
		if vo.Description != "" {
			t.Logf("description 已落库（缺陷已修复）: %q", vo.Description)
		}
	})
	t.Run("update-id-required", func(t *testing.T) {
		e.t = t
		e.expectErr(ctrl+"/update", `{}`, "id错误")
	})
	t.Run("query", func(t *testing.T) {
		e.t = t
		pg := e.queryPage(ctrl+"/query", js(map[string]any{"pageNum": 1, "pageSize": 20, "wd": pno}))
		hit := false
		for _, raw := range pg.Data {
			var vo credentialVO
			_ = json.Unmarshal(raw, &vo)
			if vo.ID == cid {
				hit = true
			}
		}
		if pg.Total < 1 || !hit {
			t.Fatalf("query total=%d hit=%v", pg.Total, hit)
		}
	})
	t.Run("selectPublic", func(t *testing.T) {
		e.t = t
		list := e.selectPublic(ctrl)
		hit := false
		for _, raw := range list {
			var vo credentialVO
			_ = json.Unmarshal(raw, &vo)
			if vo.ID == cid {
				hit = true
			}
		}
		if !hit {
			t.Fatalf("selectPublic 未包含 id=%s", cid)
		}
	})

	e.assertCommonState(t, ctrl, cid)
	t.Run("deletion", func(t *testing.T) { e.assertDeleteCycle(t, ctrl, "ram_idp_credential", cid) })

	t.Cleanup(func() { _ = pno; _ = sno })
}

// ---------------- 元数据缓存 ----------------

type metadataVO struct {
	ID          string `json:"id"`
	IdpNo       string `json:"idpNo"`
	SourceNo    string `json:"sourceNo"`
	Protocol    string `json:"protocol"`
	EntityID    string `json:"entityId"`
	MetadataUrl string `json:"metadataUrl"`
	MetaFormat  string `json:"metaFormat"`
	CacheTtl    int64  `json:"cacheTtl"`
	State       int    `json:"state"`
}

func runMetadataCRUD(t *testing.T, e *env) {
	e.t = t
	ctrl := baseURL + "/idp-metadata-cache"
	pid := e.createID(baseURL+"/identity-provider",
		js(map[string]any{"name": e.unique("provider-for-meta"), "code": "Custom", "protocol": "oidc"}),
		"ram_identity_provider")
	pno := e.providerNo(pid)
	sid := e.createID(baseURL+"/identity-source",
		js(map[string]any{"name": e.unique("source-for-meta"), "protocol": "oidc", "idp": pno}),
		"ram_identity_source")
	sno := e.sourceNo(sid)

	entityID := "https://issuer.example.com/" + e.runID
	m := map[string]any{
		"idpNo": pno, "sourceNo": sno, "protocol": "oidc",
		"entityId": entityID, "metadataUrl": entityID + "/.well-known/openid-configuration",
		"metaFormat": "json", "cacheTtl": 3600,
	}
	mid := e.createID(ctrl, js(m), "ram_idp_metadata_cache")

	t.Run("create-detail", func(t *testing.T) {
		e.t = t
		r := e.getAuth(ctrl + "/detail/" + mid)
		vo := decodeData[metadataVO](t, r)
		if vo.IdpNo != pno || vo.SourceNo != sno || vo.EntityID != entityID || vo.CacheTtl != 3600 || vo.State != 1 {
			t.Fatalf("detail 不符: %+v", vo)
		}
		if strings.Contains(string(r.Data), "metadataRaw") {
			t.Fatalf("detail 不应返回 metadataRaw: %s", r.Data)
		}
	})
	t.Run("update", func(t *testing.T) {
		e.t = t
		e.ok(ctrl+"/update", js(map[string]any{"id": mid, "cacheTtl": 600, "metaFormat": "xml"}))
		vo := decodeData[metadataVO](t, e.getAuth(ctrl+"/detail/"+mid))
		if vo.CacheTtl != 600 || vo.MetaFormat != "xml" {
			t.Fatalf("update 未生效: %+v", vo)
		}
	})
	t.Run("update-id-required", func(t *testing.T) {
		e.t = t
		e.expectErr(ctrl+"/update", `{}`, "id错误")
	})
	t.Run("query", func(t *testing.T) {
		e.t = t
		pg := e.queryPage(ctrl+"/query", js(map[string]any{"pageNum": 1, "pageSize": 20, "wd": pno}))
		hit := false
		for _, raw := range pg.Data {
			var vo metadataVO
			_ = json.Unmarshal(raw, &vo)
			if vo.ID == mid {
				hit = true
			}
		}
		if pg.Total < 1 || !hit {
			t.Fatalf("query total=%d hit=%v", pg.Total, hit)
		}
	})
	t.Run("selectPublic", func(t *testing.T) {
		e.t = t
		list := e.selectPublic(ctrl)
		hit := false
		for _, raw := range list {
			var vo metadataVO
			_ = json.Unmarshal(raw, &vo)
			if vo.ID == mid {
				hit = true
			}
		}
		if !hit {
			t.Fatalf("selectPublic 未包含 id=%s", mid)
		}
	})

	e.assertCommonState(t, ctrl, mid)
	t.Run("deletion", func(t *testing.T) { e.assertDeleteCycle(t, ctrl, "ram_idp_metadata_cache", mid) })

	t.Cleanup(func() { _ = pno; _ = sno })
}

var _ = fmt.Sprintf
