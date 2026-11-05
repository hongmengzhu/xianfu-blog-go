package modBasicConfig

// UpdateValuesCt 按“租户 + 模型”批量保存配置值。
type UpdateValuesCt struct {
	ModelNo string            `json:"modelNo" label:"模型编号"`
	Model   string            `json:"model" label:"模型标识"`
	Data    map[string]string `json:"data" label:"字段值集合"`
}
