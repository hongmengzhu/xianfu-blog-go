package constsCacheRam

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/consts/constsCachePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/consts/constsRam/typeDomainPg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/consts/constsRam/typePubPrivePg"
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/enum/state/clientPg"
)

const (
	pubprive = "pubprive:"
)

func SessionPubPrive_system() string {
	return constsCachePg.Prefix_Auth + pubprive + typeDomainPg.System.Code()
}
func SessionPubPrive_system_desktop() string {
	return constsCachePg.Prefix_Auth + pubprive + typeDomainPg.System.Code()
}

func SessionPubPrive_key(tp typeDomainPg.TypeDomain, client clientPg.Client) string {
	return constsCachePg.Prefix_Auth + pubprive + tp.Code() + "-" + client.Index()
}

// SessionPubPrive_login 登录
func SessionPubPrive_login() string {
	return constsCachePg.Prefix_Auth + pubprive + typePubPrivePg.Login.Index()
}
