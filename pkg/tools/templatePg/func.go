package templatePg

import "html/template"

func Unescaped(str string) interface{} {
	return template.HTML(str)
}
