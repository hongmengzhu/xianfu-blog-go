package constModulePg

// Item 模块（库表模型）映射项。
type Item struct {
	Name string `json:"name" label:"名称"`
	// Module 点分模块编码，如 basic.area；审计日志的 module 字段取自该值。
	Module string `json:"module" label:"模块"`
	// ModuleRoot 模块根目录；子模块填写上级模块编码，无上级时与 Module 相同。
	ModuleRoot string `json:"moduleRoot" label:"模块根目录"`
	// PathMark 额外的路由路径标识（<domain>/<resource>），用于路由与模块名不一致的场景，
	// 如 idp 模块的路由挂在 ram-idp 下。为空时默认使用 ModuleMap 的 key。
	PathMark []string `json:"pathMark" label:"路径标识"`
	// TargetType 对象类型，取库表表名，如 blog_article。
	TargetType string `json:"targetType" label:"目标类型|对象类型 blog_article"`
}
