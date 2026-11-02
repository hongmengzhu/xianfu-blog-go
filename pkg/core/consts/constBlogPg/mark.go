package constBlogPg

import "strings"

func AttachmentMark(module, typ, value string) string {
	var b strings.Builder
	b.WriteString(module)
	b.WriteString("|")
	b.WriteString(typ)
	b.WriteString("|")
	b.WriteString(value)
	return b.String()
}
