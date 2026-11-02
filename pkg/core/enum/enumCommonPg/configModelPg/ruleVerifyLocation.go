package configModelPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// RuleVerifyLocation 规则验证位置
type RuleVerifyLocation string

const (
	RuleVerifyLocationCreate RuleVerifyLocation = "create" //创建
	RuleVerifyLocationUpdate RuleVerifyLocation = "update" //更新
)

// Name 名称
func (this RuleVerifyLocation) Name() string {
	switch this {
	case "create":
		return "创建"
	case "update":
		return "更新"
	default:
		return "未知"
	}
}

// 值
func (this RuleVerifyLocation) String() string {
	return string(this)
}

// Index 值
func (this RuleVerifyLocation) Index() string {
	return string(this)
}

// IsEqual 值是否相等
func (this RuleVerifyLocation) IsEqual(id string) bool {
	return string(this) == id
}

var RuleVerifyLocationMap = map[string]enumBasePg.EnumString{
	RuleVerifyLocationCreate.String(): enumBasePg.EnumString{RuleVerifyLocationCreate.String(), RuleVerifyLocationCreate.Name()},
	RuleVerifyLocationUpdate.String(): enumBasePg.EnumString{RuleVerifyLocationUpdate.String(), RuleVerifyLocationUpdate.Name()},
}

func IsExistRuleVerifyLocation(id string) (RuleVerifyLocation, bool) {
	_, ok := RuleVerifyLocationMap[id]
	return RuleVerifyLocation(id), ok
}
