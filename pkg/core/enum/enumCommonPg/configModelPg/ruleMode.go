package configModelPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// RuleMode 规则类型
type RuleMode string

const (
	RuleModeRequired     RuleMode = "required"     //不为空
	RuleModeLength       RuleMode = "length"       //长度
	RuleModePattern      RuleMode = "pattern"      //正则
	RuleModeGreaterThan0 RuleMode = "greaterThan0" //大于0
	RuleModeVerifyFunc   RuleMode = "verifyFunc"   //内置验证
	RuleModeCustom       RuleMode = "custom"       //自定义
	RuleModeCondition    RuleMode = "condition"    //条件
)

// Name 名称
func (this RuleMode) Name() string {
	switch this {
	case "required":
		return "不为空"
	case "length":
		return "长度"
	case "pattern":
		return "正则"
	case "greaterThan0":
		return "大于0"
	case "custom":
		return "自定义"
	case "verifyFunc":
		return "内置验证"
	case "condition":
		return "条件"
	default:
		return "未知"
	}
}

// 值
func (this RuleMode) String() string {
	return string(this)
}

// Index 值
func (this RuleMode) Index() string {
	return string(this)
}

// IsEqual 值是否相等
func (this RuleMode) IsEqual(id string) bool {
	return string(this) == id
}

var RuleModeMap = map[string]enumBasePg.EnumString{
	RuleModeRequired.String():     enumBasePg.EnumString{RuleModeRequired.String(), RuleModeRequired.Name()},
	RuleModeLength.String():       enumBasePg.EnumString{RuleModeLength.String(), RuleModeLength.Name()},
	RuleModePattern.String():      enumBasePg.EnumString{RuleModePattern.String(), RuleModePattern.Name()},
	RuleModeGreaterThan0.String(): enumBasePg.EnumString{RuleModeGreaterThan0.String(), RuleModeGreaterThan0.Name()},
	RuleModeVerifyFunc.String():   enumBasePg.EnumString{RuleModeVerifyFunc.String(), RuleModeVerifyFunc.Name()},
	RuleModeCustom.String():       enumBasePg.EnumString{RuleModeCustom.String(), RuleModeCustom.Name()},
	RuleModeCondition.String():    enumBasePg.EnumString{RuleModeCondition.String(), RuleModeCondition.Name()},
}

func IsExistRuleMode(id string) (RuleMode, bool) {
	_, ok := RuleModeMap[id]
	return RuleMode(id), ok
}
