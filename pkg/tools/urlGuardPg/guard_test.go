package urlGuardPg

import (
	"context"
	"testing"
)

func TestValidateFetchURL(t *testing.T) {
	ctx := context.Background()
	cases := []struct {
		url  string
		pass bool
	}{
		{"http://127.0.0.1/x", false},
		{"http://localhost/x", false},
		{"http://169.254.169.254/latest", false},
		{"http://10.0.0.2/x", false},
		{"file:///etc/passwd", false},
		{"ftp://example.com/x", false},
		{"http://8.8.8.8/x", true},
		{"https://example.com/x", true},
	}
	for _, c := range cases {
		err := ValidateFetchURL(ctx, c.url)
		if c.pass && nil != err {
			t.Fatalf("%s should pass, got %v", c.url, err)
		}
		if !c.pass && nil == err {
			t.Fatalf("%s should be blocked", c.url)
		}
	}
}
