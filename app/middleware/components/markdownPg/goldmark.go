package markdownPg

import (
	"bytes"

	"github.com/yuin/goldmark"
	emoji "github.com/yuin/goldmark-emoji"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	"go.abhg.dev/goldmark/anchor"
	"go.abhg.dev/goldmark/frontmatter"
	"go.abhg.dev/goldmark/toc"
	"go.abhg.dev/goldmark/wikilink"
)

var markdownRender = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,             // 语法扩展
		meta.Meta,                 // 获取markdown中的meta数据
		highlighting.Highlighting, // 语法高亮
		emoji.Emoji,               // emoji
		&wikilink.Extender{},
		&anchor.Extender{},      // 添加锚点
		&frontmatter.Extender{}, // 添加frontmatter
		&toc.Extender{},         // 添加目录
	),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
	),
	goldmark.WithRendererOptions(
		html.WithHardWraps(),
		html.WithXHTML(),
	),
)

func Markdown(markdown []byte) (raw bytes.Buffer) {
	if err := markdownRender.Convert(markdown, &raw); err != nil {
		panic(err)
	}
	return
}
