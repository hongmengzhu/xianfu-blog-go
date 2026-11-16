package authPg

import (
	"testing"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/sdk/auth/holderPg"
)

func TestIsAccountAdministrator(t *testing.T) {
	cases := []struct {
		name  string
		acc   holderPg.AccountHolder
		admin bool
	}{
		{"创始人", holderPg.AccountHolder{Founder: 1}, true},
		{"administrator角色", holderPg.AccountHolder{Os: holderPg.AccountHolderOs{Roles: []string{"administrator"}}}, true},
		{"普通账号", holderPg.AccountHolder{Os: holderPg.AccountHolderOs{Roles: []string{"employee"}}}, false},
		{"空角色", holderPg.AccountHolder{}, false},
	}
	for _, c := range cases {
		holder := holderPg.HolderPg{HolderData: c.acc}
		if got := isAccountAdministrator(holder); got != c.admin {
			t.Fatalf("%s: 期望 %v 实际 %v", c.name, c.admin, got)
		}
	}
	// 无 holder 时不能 panic，也不能放行
	if isAccountAdministrator(holderPg.HolderPg{}) {
		t.Fatal("空 holder 不应判定为管理员")
	}
}
