package templatePg

// TemplateParameter 模版参数
type TemplateParameter struct {
	Code     int            `json:"code" label:"http 状态码"`
	HtmlObj  map[string]any `json:"htmlObj" label:"html对象"`
	Data     any            `json:"data" label:"html对象中，data数据"`
	DataIs   bool           `json:"dataIs" label:"html对象中，data数据是否存在"`
	SitePage SitePage       `json:"sitePage" label:"网站信息"`
}
type Option func(*TemplateParameter)

// WithCode 设置状态码
func WithCode(code int) Option {
	return func(tp *TemplateParameter) {
		tp.Code = code
	}
}

// WithData 设置数据
func WithData(data any) Option {
	return func(tp *TemplateParameter) {
		tp.Data = data
	}
}

// WithDataByResult 设置数据
func WithDataByResult(dataIs bool, data any) Option {
	return func(tp *TemplateParameter) {
		tp.Data = data
		tp.DataIs = dataIs
	}
}

// MakeHtmlObjDefaultMap 创建默认数据
func MakeHtmlObjDefaultMap() map[string]any {
	mapDefault := make(map[string]any)
	mapDefault["dataIs"] = false
	return mapDefault
}

// WithHtmlObjSet 这 html obj 数据
func WithHtmlObjSet(field string, value any) Option {
	return func(tp *TemplateParameter) {
		tp.HtmlObj[field] = value
	}
}

// WithSitePage 设置网站信息
func WithSitePage(obj SitePage) Option {
	return func(tp *TemplateParameter) {
		tp.SitePage = obj
	}
}

//// WithSitePageCtx 设置网站信息
//func WithSitePageCtx(obj SitePage) Option {
//	return func(tp *TemplateParameter) {
//		tp.SitePage = obj
//	}
//}
