package authTokenPg

import "testing"

func TestParseUnverifiedRejectsShortPayload(t *testing.T) {
	cases := []string{
		"",
		"v4.public.e30",          // {}
		"v4.public.eyJhIjoiYiJ9", // {"a":"b"}
		"v4.public.e30.extra.sig",
	}
	for _, token := range cases {
		if _, ok := ParseUnverified(token); ok {
			t.Fatalf("短载荷令牌不应解析成功: %s", token)
		}
	}
}
