package service

import (
	"testing"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/enumCommonPg/configModelPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/state/enumStatePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/state/yesNoPg/yesNoIntPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityBasic"
)

func TestParseLengthCondition(t *testing.T) {
	cases := []struct {
		in       string
		min, max int
		exact    int
	}{
		{"", 0, 0, 0},
		{"10", 0, 0, 10},
		{"1,100", 1, 100, 0},
		{" 2 , 5 ", 2, 5, 0},
	}
	for _, c := range cases {
		min, max, exact := parseLengthCondition(c.in)
		if min != c.min || max != c.max || exact != c.exact {
			t.Fatalf("parseLengthCondition(%q) = (%d,%d,%d), want (%d,%d,%d)", c.in, min, max, exact, c.min, c.max, c.exact)
		}
	}
}

func TestNormalizeYesNo(t *testing.T) {
	if got := normalizeYesNo(0, yesNoIntPg.Yes.IndexInt8()); got != yesNoIntPg.Yes.IndexInt8() {
		t.Fatalf("normalizeYesNo(0, yes) = %d", got)
	}
	if got := normalizeYesNo(yesNoIntPg.No.IndexInt8(), yesNoIntPg.Yes.IndexInt8()); got != yesNoIntPg.No.IndexInt8() {
		t.Fatalf("normalizeYesNo(no, yes) = %d", got)
	}
	if got := normalizeYesNo(0, 0); got != yesNoIntPg.No.IndexInt8() {
		t.Fatalf("normalizeYesNo(0, 0) = %d", got)
	}
}

func TestValidateRulesRequired(t *testing.T) {
	field := &entityBasic.BasicConfigModelFieldsEntity{Name: "站点名称"}
	rule := &entityBasic.BasicModelRulesEntity{
		State:    enumStatePg.ENABLE.IndexInt8(),
		RuleMode: configModelPg.RuleModeRequired.Index(),
	}
	if msg := validateRules("", []*entityBasic.BasicModelRulesEntity{rule}, field); msg != "站点名称 不能为空" {
		t.Fatalf("required message = %q", msg)
	}
	if msg := validateRules("abc", []*entityBasic.BasicModelRulesEntity{rule}, field); msg != "" {
		t.Fatalf("required with value = %q", msg)
	}
}

func TestValidateRulesLengthAndPattern(t *testing.T) {
	field := &entityBasic.BasicConfigModelFieldsEntity{Name: "编码"}
	lengthRule := &entityBasic.BasicModelRulesEntity{
		State:        enumStatePg.ENABLE.IndexInt8(),
		RuleMode:     configModelPg.RuleModeLength.Index(),
		Condition:    "2,5",
		ErrorMessage: "编码长度必须在 2-5",
	}
	if msg := validateRules("a", []*entityBasic.BasicModelRulesEntity{lengthRule}, field); msg != "编码长度必须在 2-5" {
		t.Fatalf("length message = %q", msg)
	}
	patternRule := &entityBasic.BasicModelRulesEntity{
		State:     enumStatePg.ENABLE.IndexInt8(),
		RuleMode:  configModelPg.RuleModePattern.Index(),
		Condition: "^[0-9]+$",
	}
	if msg := validateRules("abc", []*entityBasic.BasicModelRulesEntity{patternRule}, field); msg == "" {
		t.Fatal("pattern should reject abc")
	}
	if msg := validateRules("123", []*entityBasic.BasicModelRulesEntity{patternRule}, field); msg != "" {
		t.Fatalf("pattern accepted 123: %q", msg)
	}
}

func TestValidateRulesSkipsDisabled(t *testing.T) {
	rule := &entityBasic.BasicModelRulesEntity{
		State:    enumStatePg.DISABLE.IndexInt8(),
		RuleMode: configModelPg.RuleModeRequired.Index(),
	}
	if msg := validateRules("", []*entityBasic.BasicModelRulesEntity{rule}, nil); msg != "" {
		t.Fatalf("disabled rule should be skipped: %q", msg)
	}
}
