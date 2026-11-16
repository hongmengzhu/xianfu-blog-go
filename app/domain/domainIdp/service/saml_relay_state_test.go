package service

import (
	"testing"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/configPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/configPg/pg"
)

func TestSamlRelayStateSignedAndOneTime(t *testing.T) {
	s := &SamlLoginService{}
	s.pg = configPg.Pg{Jwt: pg.Jwt{System: pg.JwtConfig{Secret: "unit-test-secret"}}}

	state, err := s.makeRelayState("source-1")
	if nil != err {
		t.Fatal(err)
	}
	source, err := s.verifyRelayState(state)
	if nil != err || source != "source-1" {
		t.Fatalf("首次校验失败: source=%s err=%v", source, err)
	}
	if _, err := s.verifyRelayState(state); nil == err {
		t.Fatal("RelayState 必须只能使用一次")
	}
	if _, err := s.verifyRelayState(state + "x"); nil == err {
		t.Fatal("被篡改的 RelayState 必须拒绝")
	}
}
