package drive

import "testing"

func TestValidateUploadSuffix(t *testing.T) {
	for _, suffix := range []string{".html", ".htm", ".svg", ".js", ".php", ".exe"} {
		if err := validateUploadSuffix(suffix); nil == err {
			t.Fatalf("%s 应被拒绝", suffix)
		}
	}
	for _, suffix := range []string{".jpg", ".png", ".pdf", ".docx", ".zip"} {
		if err := validateUploadSuffix(suffix); nil != err {
			t.Fatalf("%s 应允许: %v", suffix, err)
		}
	}
	if err := validateUploadSuffix(""); nil == err {
		t.Fatal("无扩展名应被拒绝")
	}
}
