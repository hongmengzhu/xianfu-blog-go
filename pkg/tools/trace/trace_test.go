package trace

import (
	"regexp"
	"testing"
)

// traceID 格式：IP 十六进制 + 进程号 + 毫秒时间戳 + 递增计数。
var traceIDPattern = regexp.MustCompile(`^[0-9a-f]+-\d+-\d+-\d+$`)

func TestGenerateTraceIDFormat(t *testing.T) {
	id := GenerateTraceID()
	if !traceIDPattern.MatchString(id) {
		t.Fatalf("traceID 格式不符合预期: %q", id)
	}
}

func TestGenerateTraceIDUnique(t *testing.T) {
	const n = 1000
	seen := make(map[string]struct{}, n)
	for i := 0; i < n; i++ {
		id := GenerateTraceID()
		if _, ok := seen[id]; ok {
			t.Fatalf("traceID 重复: %q", id)
		}
		seen[id] = struct{}{}
	}
}
