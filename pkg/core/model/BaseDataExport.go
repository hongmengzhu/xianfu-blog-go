package model

// BaseDataExport
// @Description: 导出
type BaseDataExport struct {
	Count    int64       `json:"count"  label:"总数"`
	Header   []string    `json:"header" label:"标题"`
	Title    []string    `json:"title" label:"头部"`
	Foot     []string    `json:"foot" label:""`
	Export   interface{} `json:"export"`
	FileName string      `json:"filename" label:"文件名"`
}
