package modSysPub

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/core/model"
)

type AuthConfigVo struct {
	LoginEncrypt      bool                  `json:"loginEncrypt" lable:"登陆加密"`
	AuthPubPriveLogin model.AuthPubPriveDto `json:"login" label:"密钥"`
	MenuRouter        RamMenuRouterVo       `json:"menuRouter" label:"菜单路由"`
}

func NewAuthConfig(dto model.AuthPubPriveDto) AuthConfigVo {
	return AuthConfigVo{
		AuthPubPriveLogin: dto,
	}
}
func NewAuthConfigPublic(dto model.AuthPubPriveDto) AuthConfigVo {
	return AuthConfigVo{
		AuthPubPriveLogin: model.AuthPubPriveDto{
			PublicKey: dto.PublicKey,
			Type:      dto.Type,
			No:        dto.No,
		},
	}
}
