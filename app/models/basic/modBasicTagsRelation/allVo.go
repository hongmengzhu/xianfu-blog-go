package modBasicTagsRelation

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
	"time"
)

type AllVo struct {
	ID           typePg.Int64String     `json:"id" label:"id" `
	TagNo        string                 `json:"tagNo" label:"标签" `
	Name         string                 `json:"name" label:"名称" `          // 名称
	NameFl       string                 `json:"nameFl" label:"名称外文" `      // 名称外文
	Code         string                 `json:"code" label:"编号代号" `        // 编号代号
	NameFull     string                 `json:"nameFull" label:"全称" `      // 全称
	State        typePg.Int8            `json:"state" label:"状态:1启用;2禁用" ` // 状态:1启用;2禁用
	Description  string                 `json:"description" label:"描述" `   // 描述
	CreateAt     *time.Time             `json:"createAt" label:"创建时间" `    // 创建时间
	UpdateAt     *time.Time             `json:"updateAt" label:"更新时间" `    // 更新时间
	Type         string                 `json:"type" label:"类型" binding:"required,min=1,max=255" validate:"required,min=1,max=255"`
	TypeSys      string                 `json:"typeSys" label:"类型" binding:"required,min=1,max=255" validate:"required,min=1,max=255"`
	Url          string                 `json:"url" label:"建议url" binding:"required,min=1,max=255" validate:"required,min=1,max=255"`
	CategoryNo   string                 `json:"categoryNo" label:"分类" `
	Show         bool                   `json:"show" label:"显示" `
	NameShort    string                 `json:"nameShort" label:"名称简称" `
	AttributeMap map[string]interface{} `json:"attributeMap" label:"属性" `
	//CreateBy     typePg.Int64String     `json:"createBy" label:"创建人" `     // 创建人
	//UpdateBy     typePg.Int64String     `json:"updateBy" label:"更新人" `     // 更新人
	//OrgId        string                 `json:"orgId" label:"组织id" `       // 组织id
}
