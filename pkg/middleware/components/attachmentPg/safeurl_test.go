package attachmentPg

import (
	"context"
	"testing"
)

func TestValidateFetchURL(t *testing.T) {
	ctx := context.Background()
	cases := []struct {
		url string
		ok  bool
	}{
		{"http://127.0.0.1:8080/a", false},
		{"http://169.254.169.254/latest/meta-data/", false},
		{"file:///etc/passwd", false},
		{"ftp://example.com/a", false},
		{"http://8.8.8.8/a", true},
	}
	for _, c := range cases {
		err := ValidateFetchURL(ctx, c.url)
		if c.ok && nil != err {
			t.Fatalf("%s 应允许: %v", c.url, err)
		}
		if !c.ok && nil == err {
			t.Fatalf("%s 应拒绝", c.url)
		}
	}
}
