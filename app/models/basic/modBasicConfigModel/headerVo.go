package modBasicConfigModel

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/typePg"
	"gorm.io/datatypes"
)

type HeaderVo struct {
	Id            typePg.Uint64String `json:"id" form:"id" label:"id" `
	Name          string              `json:"name" binding:"required" label:"中文名称"`
	Model         string              `json:"model" binding:"required" label:"英文标识"`
	Table         string              `json:"table" label:"表名"`
	TypeSys       string              `json:"typeSys" label:"类型"`
	ModelCategory string              `json:"modelCategory" label:"模型种类"`
	ModuleSub     string              `json:"moduleSub" label:"子模块选择"`
	Description   string              `json:"description" label:"描述"`
	No            string              `json:"no" label:"编号"`
	Category      string              `json:"category" label:"配置分类"`
	CategoryTab   string              `json:"categoryTab" label:"分类选项卡"`
	Tags          []string            `json:"tags" label:"标签"`

	// Rules 模型级验证规则的简写集合（ruleMode 列表），用于快速保存。
	Rules []string `json:"rules" label:"模型验证规则"`
	// RuleItems 模型级验证规则的完整定义，保存到 basic_config_model_rules。
	RuleItems []RuleItemVo `json:"ruleItems" label:"模型验证规则明细"`
}

// RuleItemVo 字段/模型验证规则
type RuleItemVo struct {
	Id           typePg.Int64String `json:"id" label:"id"`
	Name         string             `json:"name" label:"名称"`
	RuleMode     string             `json:"ruleMode" label:"验证模式"`
	Coding       string             `json:"coding" label:"代码"`
	Condition    string             `json:"condition" label:"条件"`
	ErrorMessage string             `json:"errorMessage" label:"错误提示"`
	Structure    string             `json:"structure" label:"验证结构"`
	RuleTarget   []string           `json:"ruleTarget" label:"目标前端/后端"`
	Show         typePg.Int8        `json:"show" label:"1显示2隐藏"`
	Sort         typePg.Int64String `json:"sort" label:"排序"`
	ExtraData    datatypes.JSON     `json:"extraData" label:"额外数据"`
}
