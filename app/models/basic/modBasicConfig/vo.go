package modBasicConfig

import (
	"time"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
)

type Vo struct {
	ID       typePg.Int64String `json:"id" label:"id"`
	CreateAt *time.Time         `json:"createAt" label:"创建时间"`
	UpdateAt *time.Time         `json:"updateAt" label:"更新时间"`
	CreateBy string             `json:"createBy" label:"创建人"`
	UpdateBy string             `json:"updateBy" label:"更新人"`
	State    typePg.Int8        `json:"state" label:"状态"`
	Sort     typePg.Int64String `json:"sort" label:"排序"`

	TenantNo string `json:"tenantNo" label:"租户编号"`
	No       string `json:"no" label:"编号"`
	Name     string `json:"name" label:"名称"`

	ModelNo   string `json:"modelNo" label:"模型编号"`
	Model     string `json:"model" label:"模型标识"`
	Module    string `json:"module" label:"模块"`
	ModuleSub string `json:"moduleSub" label:"子模块"`

	Field       string `json:"field" label:"字段标识"`
	FieldPath   string `json:"fieldPath" label:"路径字段"`
	Description string `json:"description" label:"描述"`

	Show      typePg.Int8 `json:"show" label:"1显示2隐藏"`
	Binary    typePg.Int8 `json:"binary" label:"值二进制1是2否"`
	ValueType string      `json:"valueType" label:"字段值类型"`
	FormCode  string      `json:"formCode" label:"表单控件"`

	DefaultValue string `json:"defaultValue" label:"默认值"`
	Value        string `json:"value" label:"值"`
	ValueBinary  string `json:"valueBinary" label:"值二进制"`

	CategoryNo string `json:"categoryNo" label:"配置分类编号"`
}

// ItemVo 用于按租户 + 模型返回的配置键值项
type ItemVo struct {
	Field       string      `json:"field" label:"字段标识"`
	Name        string      `json:"name" label:"名称"`
	Description string      `json:"description" label:"描述"`
	Value       string      `json:"value" label:"值"`
	ValueType   string      `json:"valueType" label:"字段值类型"`
	Show        typePg.Int8 `json:"show" label:"1显示2隐藏"`
	CategoryNo  string      `json:"categoryNo" label:"配置分类编号"`
}
