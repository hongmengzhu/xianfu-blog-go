package modEventBasicTags

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/sdk/auth/holderPg"
)

type TagsRelation struct {
	Category string            `json:"category"`
	Tags     []string          `json:"tags"`
	Holder   holderPg.HolderPg `json:"holder"`
}
