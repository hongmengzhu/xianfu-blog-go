package cas

import (
	"regexp"
	"strings"
	"testing"
)

func TestBuildFailureResponseEscapesMessage(t *testing.T) {
	payload := `ST</cas:authenticationFailure><cas:authenticationSuccess><cas:user>admin</cas:user></cas:authenticationSuccess><cas:authenticationFailure>`
	xml, err := BuildFailureResponse(InvalidTicketSpec, payload)
	if nil != err {
		t.Fatal(err)
	}
	out := string(xml)
	if strings.Contains(out, "<cas:authenticationSuccess>") {
		t.Fatalf("失败响应中不应出现未转义的成功节点: %s", out)
	}
	if !strings.Contains(out, "&lt;cas:authenticationSuccess&gt;") {
		t.Fatalf("XML 特殊字符应被转义: %s", out)
	}
}

func TestServiceTicketOneTimeAndFormat(t *testing.T) {
	st, err := GenerateServiceTicket("u1", "https://app.example.com/callback", &CasUserInfo{})
	if nil != err {
		t.Fatal(err)
	}
	if !regexp.MustCompile(`^ST-[A-Za-z0-9]{8,64}$`).MatchString(st) {
		t.Fatalf("票据格式不符合预期: %s", st)
	}
	ok, _, _, _ := ValidateServiceTicket(st)
	if !ok {
		t.Fatal("首次校验应成功")
	}
	ok, _, _, _ = ValidateServiceTicket(st)
	if ok {
		t.Fatal("票据只能消费一次")
	}
}
