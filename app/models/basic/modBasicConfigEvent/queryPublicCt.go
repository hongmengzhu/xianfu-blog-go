package modBasicConfigEvent

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/model"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
	"gorm.io/datatypes"
)

type QueryPublicCt struct {
	model.BaseQueryNodeCt
	ID    typePg.Int64String `json:"id" label:"" `
	State typePg.Int8        `json:"state" label:"状态:1启用;2禁用" `
	//
	Description     string         `json:"description" label:"描述" ` // 描述
	Name            string         `json:"name" label:"名称" `
	LangCode        string         `json:"langCode" comment:"语言" `
	TypeSys         string         `json:"typeSys" comment:"类型|普通|系统|api" `
	Module          string         `json:"module" comment:"模块" `
	ModuleSub       string         `json:"moduleSub" comment:"子模块" `
	Owner           string         `json:"owner" comment:"所属/拥有者" `
	Model           string         `json:"model" comment:"模型" `
	Show            typePg.Int8    `json:"show" comment:"1显示2隐藏" `
	ExtraData       datatypes.JSON `json:"extraData" label:"额外数据" `
	Value           string         `json:"value" comment:"值" `
	Client          string         `json:"client" comment:"端" `
	ClientProgram   string         `json:"clientProgram" comment:"端内程序|隔开" `
	ClientSync      string         `json:"clientSync" comment:"端同步" `
	LoadingLocation string         `json:"loadingLocation" comment:"加载位置" `
	Cache           string         `json:"cache" comment:"缓存" `
}
