package main

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestGinEngineTemplateFuncs 确保 gin 渲染器注册的模板函数与 echo 侧一致。
//
// nav/tool/footer-mobile.html 使用了 hasPrefix；若 gin 的 FuncMap 缺少该函数，
// 整份模板解析失败，所有 gin 前台路由都会返回状态码 200 的空响应。
func TestGinEngineTemplateFuncs(t *testing.T) {
	gin.SetMode(gin.TestMode)
	e := gin.New()
	NewGinEngine(nil)(e)
	if nil == e.HTMLRender {
		t.Fatal("gin HTMLRender 未初始化，data/templates 目录可能不存在")
	}

	w := httptest.NewRecorder()
	render := e.HTMLRender.Instance("nav/tool/footer-mobile.html", map[string]any{"pageUrl": "/category"})
	if err := render.Render(w); nil != err {
		t.Fatalf("模板渲染失败: %v", err)
	}
	if !strings.Contains(w.Body.String(), "mobile-bottom-nav") {
		t.Fatalf("模板渲染结果异常: %q", w.Body.String())
	}
	if !strings.Contains(w.Body.String(), `class="item active"`) {
		t.Fatal("hasPrefix 未生效，分类导航项未高亮")
	}
}
