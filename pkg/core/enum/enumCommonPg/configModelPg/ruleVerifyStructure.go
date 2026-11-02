package configModelPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumBasePg"
)

// RuleVerifyStructure 规则验证结构
type RuleVerifyStructure string

const (
	RuleVerifyStructureFront   RuleVerifyStructure = "front"   //前端
	RuleVerifyStructureBackend RuleVerifyStructure = "backend" //后端
)

// Name 名称
func (this RuleVerifyStructure) Name() string {
	switch this {
	case "front":
		return "前端"
	case "backend":
		return "后端"
	default:
		return "未知"
	}
}

// 值
func (this RuleVerifyStructure) String() string {
	return string(this)
}

// Index 值
func (this RuleVerifyStructure) Index() string {
	return string(this)
}

// IsEqual 值是否相等
func (this RuleVerifyStructure) IsEqual(id string) bool {
	return string(this) == id
}

var RuleVerifyStructureMap = map[string]enumBasePg.EnumString{
	RuleVerifyStructureFront.String():   enumBasePg.EnumString{RuleVerifyStructureFront.String(), RuleVerifyStructureFront.Name()},
	RuleVerifyStructureBackend.String(): enumBasePg.EnumString{RuleVerifyStructureBackend.String(), RuleVerifyStructureBackend.Name()},
}

func IsExistRuleVerifyStructure(id string) (RuleVerifyStructure, bool) {
	_, ok := RuleVerifyStructureMap[id]
	return RuleVerifyStructure(id), ok
}
