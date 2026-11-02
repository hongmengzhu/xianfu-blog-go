package markdownPg

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// MarkdownContentBlock 定义Markdown内容块的结构
type MarkdownContentBlock struct {
	Type    string // 类型：text(普通文本)、image(图片)、code(代码段)
	Content string // 内容
}

// ParseMarkdown 解析Markdown文本，拆分为内容块
func ParseMarkdown(text string) []MarkdownContentBlock {
	var blocks []MarkdownContentBlock
	lines := strings.Split(text, "\n")
	inCodeBlock := false
	codeContent := []string{}

	// 匹配Markdown图片的正则（![alt](url) 格式）
	imgRegex := regexp.MustCompile(`!\[.*?\]\(.*?\)`)

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// 处理代码段开始/结束标记（``` 或 ```go/python等）
		if strings.HasPrefix(trimmedLine, "```") {
			if !inCodeBlock {
				// 代码段开始
				inCodeBlock = true
				codeContent = []string{line}
			} else {
				// 代码段结束
				inCodeBlock = false
				codeContent = append(codeContent, line)
				blocks = append(blocks, MarkdownContentBlock{
					Type:    "code",
					Content: strings.Join(codeContent, "\n"),
				})
				codeContent = []string{}
				continue
			}
		}

		if inCodeBlock {
			// 代码段内的行，暂存
			codeContent = append(codeContent, line)
			continue
		}

		// 处理普通行：先移除图片，区分普通文本和图片行
		// 移除行内的Markdown图片
		lineWithoutImg := imgRegex.ReplaceAllString(line, "")
		if lineWithoutImg == line && trimmedLine != "" {
			// 无图片的普通文本行
			blocks = append(blocks, MarkdownContentBlock{
				Type:    "text",
				Content: line,
			})
		} else if lineWithoutImg != "" {
			// 行内有图片，移除图片后剩余文本作为普通文本
			blocks = append(blocks, MarkdownContentBlock{
				Type:    "text",
				Content: lineWithoutImg,
			})
		}
		// 纯图片行（移除后为空）则跳过，不加入blocks
	}

	return blocks
}

// TruncateMarkdown 截取Markdown文本，满足：
// 1. 前1000个有效文字（图片不计）；
// 2. 代码段不截断；
// 3. 移除图片内容
func TruncateMarkdown(markdownText string, maxChars int) string {
	blocks := ParseMarkdown(markdownText)
	var result []string
	currentCount := 0

	for _, block := range blocks {
		switch block.Type {
		case "image":
			// 图片块，跳过，不计入字数
			continue
		case "text":
			// 普通文本块：计算字符数，判断是否超限制
			charCount := utf8.RuneCountInString(block.Content)
			if currentCount+charCount <= maxChars {
				// 加入后不超限制，直接添加
				result = append(result, block.Content)
				currentCount += charCount
			} else {
				// 加入后超限制，截取剩余字符，终止遍历
				remaining := maxChars - currentCount
				if remaining > 0 {
					// 截取剩余字数的文本
					truncated := SubstrCN(block.Content, remaining)
					result = append(result, truncated)
					currentCount = maxChars
				}
				// 已达上限，终止
				break
			}
		case "code":
			// 代码段：判断加入后是否超限制，不超则完整加入，超则跳过
			// 代码段本身不计入「有效文字数」，仅判断是否截断
			if currentCount >= maxChars {
				// 已达上限，跳过代码段
				break
			}
			// 代码段完整加入，且不计入字数统计
			result = append(result, block.Content)
		}

		// 达到字数上限，终止遍历
		if currentCount >= maxChars {
			break
		}
	}

	return strings.Join(result, "\n")
}

// SubstrCN 安全截取包含中文的字符串，避免乱码（复用之前的函数）
func SubstrCN(str string, length int) string {
	if length <= 0 {
		return ""
	}
	if utf8.RuneCountInString(str) <= length {
		return str
	}
	runes := []rune(str)
	return string(runes[:length])
}
