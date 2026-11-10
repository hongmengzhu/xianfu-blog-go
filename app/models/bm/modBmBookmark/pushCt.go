package modBmBookmark

// PushCt 书签批量推送到工具模块
// @Description:
type PushCt struct {
	// 要推送的书签 id
	IDs []string `json:"ids" label:"书签id" validate:"required,min=1"`
	// 目的模块编号，未传时取分类所属模块
	Module string `json:"module" label:"模块"`
	// 目的分类（工具分类编号）
	CategoryNo string `json:"categoryNo" label:"目的分类" validate:"required,min=1,max=80"`
	// 目的标签（标签编码），为空时继承书签原标签
	Tags []string `json:"tags" label:"目的标签"`
}

// PushVo 推送结果
// @Description:
type PushVo struct {
	// 总数
	Total int `json:"total" label:"总数"`
	// 成功数
	Success int `json:"success" label:"成功数"`
	// 已存在而跳过数
	Skipped int `json:"skipped" label:"跳过数"`
	// 失败数
	Failed int `json:"failed" label:"失败数"`
	// 新生成的工具 id
	IDs []string `json:"ids" label:"工具id"`
	// 失败/跳过原因
	Errors []string `json:"errors" label:"失败原因"`
}
