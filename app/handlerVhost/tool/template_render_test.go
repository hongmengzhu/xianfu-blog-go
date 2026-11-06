package tool

// 工具站模板冒烟测试：把 handler 传给模板的数据结构固化成用例，
// 保证「模版改名/字段调整」不会在测试里静默退化。
//
// 运行：go test ./app/handlerVhost/tool/

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
	"time"

	blogWeb "github.com/hongmengzhu/xianfu-blog-go/app/models/blog/modBlogArticleWeb"
	"github.com/hongmengzhu/xianfu-blog-go/app/models/bm/modBmToolCategory"
	toolWeb "github.com/hongmengzhu/xianfu-blog-go/app/models/bm/modBmToolWeb"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/configPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/templatePg"
	"github.com/pangu-2/go-tools/tools/dbPg/pagePg"
)

// toolPages 工具站的页面模板名（主题 + 模块 + 页面）
var toolPages = []string{"index", "category", "category-code", "detail", "news", "tutorial", "tutorial-detail", "search"}

// loadTemplate 复刻 staterEchoV5 的模板加载方式
func loadTemplate(t *testing.T) *template.Template {
	t.Helper()
	pattern := filepath.Join("..", "..", "..", "data", "templates", "**", "**", "*")
	tpl, err := template.New("").Funcs(templatePg.DefaultFuncMap()).ParseGlob(pattern)
	if nil != err {
		t.Fatalf("模板解析失败: %v", err)
	}
	return tpl
}

// baseData 复刻 templatePg.buildTemplateParam 写入的基础数据
func baseData(t *testing.T, path string) map[string]any {
	t.Helper()
	req, err := http.NewRequest("GET", "http://tool.test"+path, nil)
	if nil != err {
		t.Fatalf("构造请求失败: %v", err)
	}
	m := templatePg.MakeHtmlObjDefaultMap()
	m["ctxRequest"] = templatePg.NewHttpPgRequest(req)
	m["sitePage"] = templatePg.SitePage{
		Title:       "测试页面 - 分宝呀",
		Description: "测试描述",
		Keywords:    "AI工具,AI导航",
		SiteName:    "分宝呀",
		SiteUrl:     "http://tool.test",
	}
	m["pg"] = configPg.Pg{}
	m["server"] = configPg.Server{}
	m["pageUrl"] = path
	return m
}

func toolCategory(no, code, name, desc string) *modBmToolCategory.Cache {
	c := &modBmToolCategory.Cache{}
	c.No = no
	c.Code = code
	c.Name = name
	c.Description = desc
	return c
}

// categoryTrees 工具分类树（pageCategory）
func categoryTrees() []model.BaseNodeTree {
	return []model.BaseNodeTree{
		{Id: "c1", Name: "AI聊天助手", Extend: toolCategory("c1", "ai-chat", "AI聊天助手", "聊天类工具")},
		{Id: "c2", Name: "AI绘画", Extend: toolCategory("c2", "ai-paint", "AI绘画", "绘画类工具")},
	}
}

// toolGroups QueryListBYCategory 的返回结构
func toolGroups() []*model.BaseNodeDataTree {
	return []*model.BaseNodeDataTree{
		{No: "c1", Code: "ai-chat", Name: "AI聊天助手", Extend: toolCategory("c1", "ai-chat", "AI聊天助手", "聊天类工具"), Data: toolVos(2)},
		{No: "c2", Code: "ai-paint", Name: "AI绘画", Extend: toolCategory("c2", "ai-paint", "AI绘画", "绘画类工具"), Data: toolVos(1)},
	}
}

func toolVos(n int) []toolWeb.Vo {
	now := time.Now()
	out := make([]toolWeb.Vo, 0, n)
	for i := 1; i <= n; i++ {
		vo := toolWeb.Vo{}
		vo.ID = 100
		vo.No = fmt.Sprintf("tool-%d", i)
		vo.Code = fmt.Sprintf("chat-%d", i)
		vo.Name = fmt.Sprintf("工具%d", i)
		vo.Description = "工具描述"
		vo.CategoryName = "AI聊天助手"
		vo.CategoryNo = "c1"
		vo.UrlSource = "https://example.com"
		vo.Tags = []string{"对话", "写作"}
		vo.CreateAt = &now
		vo.AttachmentsMap = map[string]string{
			"list": "/assets/static/nav-001/img/tool-icons/chatwiz.png",
			"main": "/assets/static/nav-001/img/tool-icons/chatwiz.png",
		}
		out = append(out, vo)
	}
	return out
}

// toolPaginator QueryByContext / 搜索结果
func toolPaginator(total int64) pagePg.Paginator[toolWeb.Vo] {
	return pagePg.Paginator[toolWeb.Vo]{
		PageNum:  2,
		PageSize: 20,
		Total:    total,
		Data:     toolVos(3),
		Pageable: pagePg.MakePageable(total, 2, 20),
	}
}

func toolDetailVo() toolWeb.DetailVo {
	now := time.Now()
	vo := toolWeb.DetailVo{}
	vo.ID = 100
	vo.No = "tool-1"
	vo.Code = "chat-1"
	vo.Name = "工具1"
	vo.Description = "工具描述"
	vo.CategoryName = "AI聊天助手"
	vo.CategoryNo = "c1"
	vo.Version = "1.0"
	vo.UpdateAt = &now
	vo.ContentConvHtml = template.HTML("<p>详细介绍</p>")
	vo.UrlSource = "https://example.com"
	vo.CategoryObj = []*modBmToolCategory.Cache{toolCategory("c1", "ai-chat", "AI聊天助手", "聊天类工具")}
	vo.Tags = []string{"对话", "写作"}
	vo.AttachmentsMap = map[string]string{"main": "/assets/static/nav-001/img/tool-icons/chatwiz.png"}
	return vo
}

// articlePaginator 快讯/教程列表
func articlePaginator(total int64) pagePg.Paginator[blogWeb.Vo] {
	list := make([]blogWeb.Vo, 0, 2)
	for i := 1; i <= 2; i++ {
		vo := blogWeb.Vo{}
		vo.ID = 1
		vo.No = fmt.Sprintf("a-%d", i)
		vo.Code = fmt.Sprintf("art-%d", i)
		vo.Name = fmt.Sprintf("文章%d", i)
		vo.Description = "文章摘要"
		vo.Author = "作者"
		vo.CategoryName = "专栏"
		vo.UrlSource = "https://example.com/post"
		vo.CreateAt = time.Now()
		list = append(list, vo)
	}
	return pagePg.Paginator[blogWeb.Vo]{
		PageNum:  1,
		PageSize: 20,
		Total:    total,
		Data:     list,
		Pageable: pagePg.MakePageable(total, 1, 20),
	}
}

func articleDetailVo() blogWeb.DetailVo {
	vo := blogWeb.DetailVo{}
	vo.ID = 1
	vo.No = "a-1"
	vo.Name = "教程标题"
	vo.Description = "教程摘要"
	vo.Author = "作者"
	vo.CategoryName = "教程"
	vo.CreateAt = time.Now()
	vo.UpdateAt = time.Now()
	vo.ContentConvHtml = template.HTML("<p>正文</p>")
	vo.UrlSource = "https://example.com/post"
	vo.Tags = []string{"tag1", "tag2"}
	return vo
}

type renderCase struct {
	name    string
	page    string
	build   func(t *testing.T) map[string]any
	expects []string
}

func renderCases() []renderCase {
	return []renderCase{
		{"首页-有数据", "index", func(t *testing.T) map[string]any {
			m := baseData(t, "/")
			m["dataIs"] = true
			m["data"] = toolGroups()
			m["pageCategory"] = categoryTrees()
			return m
		}, []string{"工具集", "/category/ai-chat", "/detail/tool-1", "个分类"}},
		{"首页-无分类", "index", func(t *testing.T) map[string]any {
			m := baseData(t, "/")
			m["pageCategory"] = make([]model.BaseNodeTree, 0)
			return m
		}, []string{"暂无工具"}},
		{"分类总览", "category", func(t *testing.T) map[string]any {
			m := baseData(t, "/category")
			m["dataIs"] = true
			m["data"] = toolGroups()
			m["pageCategory"] = categoryTrees()
			return m
		}, []string{"工具库", "快速定位", "/category/ai-chat"}},
		{"分类详情-分页", "category-code", func(t *testing.T) map[string]any {
			m := baseData(t, "/category/ai-chat")
			m["dataIs"] = true
			m["data"] = toolPaginator(45)
			m["pageCategory"] = categoryTrees()
			m["pageCategoryItem"] = toolCategory("c1", "ai-chat", "AI聊天助手", "聊天类工具")
			return m
		}, []string{"AI聊天助手", "共找到", "45", "pageNum=1", "pageNum=3"}},
		{"分类详情-分类不存在", "category-code", func(t *testing.T) map[string]any {
			m := baseData(t, "/category/unknown")
			m["pageCategory"] = categoryTrees()
			return m
		}, []string{"分类不存在"}},
		{"分类详情-空列表", "category-code", func(t *testing.T) map[string]any {
			m := baseData(t, "/category/ai-chat")
			m["dataIs"] = true
			m["data"] = pagePg.Paginator[toolWeb.Vo]{Pageable: pagePg.MakePageable(0, 1, 20)}
			m["pageCategory"] = categoryTrees()
			m["pageCategoryItem"] = toolCategory("c1", "ai-chat", "AI聊天助手", "聊天类工具")
			return m
		}, []string{"该分类下暂无工具"}},
		{"工具详情", "detail", func(t *testing.T) map[string]any {
			m := baseData(t, "/detail/tool-1")
			m["dataIs"] = true
			m["data"] = toolDetailVo()
			m["pageCategory"] = categoryTrees()
			return m
		}, []string{"工具1", "工具描述", "访问官网", "/category/ai-chat", "详细介绍", "<p>详细介绍</p>", "版本 1.0"}},
		{"工具详情-不存在", "detail", func(t *testing.T) map[string]any {
			m := baseData(t, "/detail/unknown")
			m["pageCategory"] = categoryTrees()
			return m
		}, []string{"工具不存在"}},
		{"快讯列表", "news", func(t *testing.T) map[string]any {
			m := baseData(t, "/news")
			m["dataIs"] = true
			m["data"] = articlePaginator(30)
			m["categorys"] = categoryTrees()
			return m
		}, []string{"AI 快讯", "文章1", "原文", "条快讯"}},
		{"快讯列表-空", "news", func(t *testing.T) map[string]any {
			m := baseData(t, "/news")
			m["dataIs"] = true
			m["data"] = articlePaginator(0)
			return m
		}, []string{"暂无快讯"}},
		{"教程列表", "tutorial", func(t *testing.T) map[string]any {
			m := baseData(t, "/tutorial")
			m["dataIs"] = true
			m["data"] = articlePaginator(30)
			m["categorys"] = categoryTrees()
			return m
		}, []string{"AI 教程资源", "/tutorial/art-1", "份学习资源"}},
		{"教程列表-空", "tutorial", func(t *testing.T) map[string]any {
			m := baseData(t, "/tutorial")
			m["dataIs"] = true
			m["data"] = articlePaginator(0)
			return m
		}, []string{"暂无教程"}},
		{"教程详情", "tutorial-detail", func(t *testing.T) map[string]any {
			m := baseData(t, "/tutorial/art-1")
			m["dataIs"] = true
			m["data"] = articleDetailVo()
			m["categorys"] = categoryTrees()
			return m
		}, []string{"教程标题", "<p>正文</p>", "更新于"}},
		{"教程详情-不存在", "tutorial-detail", func(t *testing.T) map[string]any {
			m := baseData(t, "/tutorial/unknown")
			m["categorys"] = categoryTrees()
			return m
		}, []string{"教程不存在"}},
		{"搜索-有结果", "search", func(t *testing.T) map[string]any {
			m := baseData(t, "/search")
			m["dataIs"] = true
			m["data"] = toolPaginator(45)
			m["pageCategory"] = categoryTrees()
			m["searchQuery"] = "聊天"
			// template.URL 避免 html/template 二次转义查询串
			m["pageQuery"] = template.URL("q=%E7%AE%94%E9%A2%8A&")
			return m
		}, []string{"搜索结果", "聊天", "pageNum=3", "q=%E7%AE%94"}},
		{"搜索-无结果", "search", func(t *testing.T) map[string]any {
			m := baseData(t, "/search")
			m["dataIs"] = true
			m["data"] = toolPaginator(0)
			m["pageCategory"] = categoryTrees()
			m["searchQuery"] = "不存在的词"
			m["pageQuery"] = template.URL("q=%E4%B8%8D%E5%AD%98%E5%9C%A8&")
			return m
		}, []string{"没有找到相关工具"}},
		{"搜索-无关键字", "search", func(t *testing.T) map[string]any {
			m := baseData(t, "/search")
			return m
		}, []string{"搜索工具", "没有找到相关工具"}},
	}
}

// TestToolTemplatesRender 渲染工具站所有页面，校验数据契约与空数据处理
func TestToolTemplatesRender(t *testing.T) {
	tpl := loadTemplate(t)
	for _, c := range renderCases() {
		t.Run(c.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := tpl.ExecuteTemplate(&buf, "nav/tool/"+c.page+".html", c.build(t)); nil != err {
				t.Fatalf("渲染 %s 失败: %v", c.page, err)
			}
			out := buf.String()
			if strings.Contains(out, "<no value>") {
				t.Errorf("输出出现 <no value>，说明模板引用了未设置的数据")
			}
			for _, w := range c.expects {
				if !strings.Contains(out, w) {
					t.Errorf("缺少期望内容 %q", w)
				}
			}
		})
	}
}

// TestToolTemplatesBareRender 只给基础数据（模拟缓存未命中/查询失败）也不能报错
func TestToolTemplatesBareRender(t *testing.T) {
	tpl := loadTemplate(t)
	for _, page := range toolPages {
		t.Run(page, func(t *testing.T) {
			var buf bytes.Buffer
			if err := tpl.ExecuteTemplate(&buf, "nav/tool/"+page+".html", baseData(t, "/"+page)); nil != err {
				t.Fatalf("空数据渲染 %s 失败: %v", page, err)
			}
			if strings.Contains(buf.String(), "<no value>") {
				t.Errorf("空数据输出出现 <no value>")
			}
		})
	}
}

// TestToolTemplatesNisargStub 默认主题 nisarg 的转接桩必须与 nav 实现一致
func TestToolTemplatesNisargStub(t *testing.T) {
	tpl := loadTemplate(t)
	for _, c := range renderCases() {
		t.Run(c.name, func(t *testing.T) {
			d := c.build(t)
			var nav, nisarg bytes.Buffer
			if err := tpl.ExecuteTemplate(&nav, "nav/tool/"+c.page+".html", d); nil != err {
				t.Fatalf("渲染 nav 模板失败: %v", err)
			}
			if nil == tpl.Lookup("nisarg/tool/"+c.page+".html") {
				t.Fatalf("缺少 nisarg 主题转接桩: nisarg/tool/%s.html", c.page)
			}
			if err := tpl.ExecuteTemplate(&nisarg, "nisarg/tool/"+c.page+".html", d); nil != err {
				t.Fatalf("渲染 nisarg 转接桩失败: %v", err)
			}
			if nav.String() != nisarg.String() {
				t.Errorf("nisarg 转接桩渲染结果与 nav 实现不一致")
			}
		})
	}
}
