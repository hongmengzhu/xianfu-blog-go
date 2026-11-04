package templatePg

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin/render"
)

// liveCheckInterval 模板文件变更检查的最小间隔，避免每个请求都去 stat 全部模板。
const liveCheckInterval = time.Second

// LiveSet 支持热加载的模板集合。
//
// 依据模板文件的 mtime + size 组成指纹，指纹变化才重新解析；
// 重新解析失败时保留上一份可用模板并打印错误，不会让整个站点失效。
type LiveSet struct {
	pattern string
	funcs   template.FuncMap

	mu      sync.Mutex
	cur     atomic.Pointer[template.Template]
	sig     string
	checked time.Time
}

// NewLiveSet 创建热加载模板集合，并立即解析一次。
//
//	@Description: 创建热加载模板集合
//	@param pattern 模板 glob，例如 data/templates/**/**/*
//	@param funcs 模板函数
//	@return *LiveSet
func NewLiveSet(pattern string, funcs template.FuncMap) *LiveSet {
	l := &LiveSet{pattern: pattern, funcs: funcs}
	t, err := l.parse()
	if nil != err {
		fmt.Printf("[template] 首次解析模板失败（%s）: %v\n", pattern, err)
		return l
	}
	l.cur.Store(t)
	l.sig = l.signature()
	l.checked = time.Now()
	return l
}

// parse 解析全部模板文件。
func (l *LiveSet) parse() (*template.Template, error) {
	return template.New("").Funcs(l.funcs).ParseGlob(l.pattern)
}

// signature 计算模板文件指纹：路径 + mtime + size。
//
//	@Description: 计算模板文件指纹
//	@return string
func (l *LiveSet) signature() string {
	files, err := filepath.Glob(l.pattern)
	if nil != err {
		return ""
	}
	sort.Strings(files)
	var b strings.Builder
	for _, file := range files {
		info, err := os.Stat(file)
		if nil != err {
			continue
		}
		b.WriteString(file)
		b.WriteByte(':')
		b.WriteString(strconv.FormatInt(info.ModTime().UnixNano(), 10))
		b.WriteByte(':')
		b.WriteString(strconv.FormatInt(info.Size(), 10))
		b.WriteByte(';')
	}
	return b.String()
}

// reload 检查指纹并按需重新解析（调用方需持有锁）。
func (l *LiveSet) reload() {
	if time.Since(l.checked) < liveCheckInterval {
		return
	}
	l.checked = time.Now()
	sig := l.signature()
	if sig == l.sig || "" == sig {
		return
	}
	t, err := l.parse()
	if nil != err {
		// 保留旧模板，避免改坏一个模板导致整站白屏
		fmt.Printf("[template] 重新解析失败，继续使用上一版模板（%s）: %v\n", l.pattern, err)
		return
	}
	l.cur.Store(t)
	l.sig = sig
	fmt.Printf("[template] 模板已热加载（%s）\n", l.pattern)
}

// Template 返回当前可用的模板集合，必要时先热加载。
//
//	@Description: 获取当前模板集合
//	@return *template.Template
func (l *LiveSet) Template() *template.Template {
	l.mu.Lock()
	l.reload()
	l.mu.Unlock()
	if t := l.cur.Load(); nil != t {
		return t
	}
	return template.New("")
}

// Refresh 立即检查并热加载（忽略最小检查间隔），便于测试与手动触发。
func (l *LiveSet) Refresh() {
	l.mu.Lock()
	l.checked = time.Time{}
	l.reload()
	l.mu.Unlock()
}

// GinLiveHTML 基于 LiveSet 的 gin HTML 渲染器，实现 render.HTMLRender。
type GinLiveHTML struct {
	live *LiveSet
}

// NewGinLiveHTML 创建 gin 热加载渲染器。
//
//	@Description: 创建 gin 热加载渲染器
//	@param pattern 模板 glob
//	@param funcs 模板函数
//	@return render.HTMLRender
func NewGinLiveHTML(pattern string, funcs template.FuncMap) render.HTMLRender {
	return &GinLiveHTML{live: NewLiveSet(pattern, funcs)}
}

// Instance 实现 render.HTMLRender。
func (g *GinLiveHTML) Instance(name string, data any) render.Render {
	return &ginLiveHTML{Template: g.live.Template(), Name: name, Data: data}
}

// ginLiveHTML 单次渲染实例。
type ginLiveHTML struct {
	Template *template.Template
	Name     string
	Data     any
}

// Render 实现 render.Render。
func (h *ginLiveHTML) Render(w http.ResponseWriter) error {
	h.WriteContentType(w)
	if "" == h.Name {
		return h.Template.Execute(w, h.Data)
	}
	return h.Template.ExecuteTemplate(w, h.Name, h.Data)
}

// WriteContentType 实现 render.Render。
func (h *ginLiveHTML) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	if len(header["Content-Type"]) == 0 {
		header["Content-Type"] = []string{"text/html; charset=utf-8"}
	}
}
