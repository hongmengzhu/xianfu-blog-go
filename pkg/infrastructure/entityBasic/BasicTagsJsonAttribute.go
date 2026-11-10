package entityBasic

// BasicTagsJsonAttribute
// @Description: 属性
type BasicTagsJsonAttribute struct {
	Type     string `json:"type" label:"快速颜色类型"`
	Bordered bool   `json:"bordered" label:"边框"`
	Color    string `json:"color" label:"自定义颜色"`
	Size     string `json:"size" label:"尺寸"`
	Strong   bool   `json:"strong" label:"粗体"`
	Round    bool   `json:"round" label:"圆角"`
}
