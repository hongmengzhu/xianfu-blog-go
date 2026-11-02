package domainPg

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/consts/constsPg/terminalCodePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/consts/constsRam/typeDomainPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/state/clientPg"
)

type Option struct {
	TerminalCode terminalCodePg.TerminalCode `json:"terminalCode" label:"终端码;general;manage;merchant;shop;tenant;customer"`
	TypeDomain   typeDomainPg.TypeDomain     `json:"typeDomain" label:"域类型;普通;系统;商户"`
	Client       clientPg.Client             `json:"client" label:"客户端;general;manage;merchant;shop;tenant;customer"`
}
