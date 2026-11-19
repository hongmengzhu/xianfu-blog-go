package constModulePg

type Item struct {
	Name       string   `json:"name" label:"名称"`
	Module     string   `json:"module" label:"模块"`
	ModuleRoot string   `json:"moduleRoot" label:"模块根目录"`
	PathMark   []string `json:"pathMark" label:"路径标识"`
	TargetType string   `json:"targetType" label:"目标类型|对象类型 blog_article"`
}
