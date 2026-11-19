package auditPg

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestMaskJSON(t *testing.T) {
	raw := []byte(`{"username":"fox","password":"secret","nested":{"token":"abc","mobile":"13800001111"}}`)
	out := maskJSON(raw, 0)
	if strings.Contains(out, "secret") || strings.Contains(out, "abc") || strings.Contains(out, "13800001111") {
		t.Fatalf("sensitive value not masked: %s", out)
	}
	if !strings.Contains(out, "fox") {
		t.Fatalf("normal value should be kept: %s", out)
	}
}

func TestComputeDiff(t *testing.T) {
	meta := Meta{IgnoreFields: []string{"createBy"}}
	before := map[string]any{"id": float64(1), "name": "old", "state": float64(1), "createBy": "a"}
	after := map[string]any{"id": float64(1), "name": "new", "state": float64(1), "createBy": "b"}
	diff := computeDiff(before, after, meta)
	if len(diff) != 1 {
		t.Fatalf("want 1 diff, got %d: %+v", len(diff), diff)
	}
	if _, ok := diff["name"]; !ok {
		t.Fatalf("expected name diff: %+v", diff)
	}
	if hasMeaningfulChange(before, map[string]any{"id": float64(1), "name": "old", "state": float64(1), "createBy": "b"}, meta) {
		t.Fatalf("ignored field change should not be meaningful")
	}
}

func TestTruncateUTF8(t *testing.T) {
	s := "审计日志"
	got := truncate(s, 3)
	if !utf8.ValidString(got) {
		t.Fatalf("invalid utf8: %q", got)
	}
}

func TestMetaValidate(t *testing.T) {
	meta := Meta{Module: "basic.area", Action: "delete", OperationType: 3, RiskLevel: 3, TargetType: "area", TargetFrom: TargetFromBody, TargetKey: "ids"}
	if err := meta.Validate(); err != nil {
		t.Fatalf("valid meta rejected: %v", err)
	}
	if err := (Meta{Module: "x"}).Validate(); err == nil {
		t.Fatal("invalid meta accepted")
	}
}
