package modRamAccount

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/infrastructure/entityRam"
)

type LoginLogDto struct {
	Ip          string         `json:"ip"`
	Client      string         `json:"client"`
	LoginSource string         `json:"loginSource"`
	AppNo       string         `json:"appNo"`
	Ano         string         `json:"ano"`
	ExtraData   map[string]any `json:"extraData"`
	Account     entityRam.RamAccountEntity
}
