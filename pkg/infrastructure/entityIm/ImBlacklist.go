package entityIm

import (
	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/entityPg"
)

type ImBlacklist struct {
	entityPg.BaseCrUpAtEntity
	InitializerAno string `gorm:"column:initializer_ano;type:varchar(80);index;default:'';comment:发起人" json:"initializer_ano" comment:"发起人" `
	BlockedAno     string `gorm:"column:blocked_ano;type:varchar(80);index;default:'';comment:被加入黑名单的人" json:"blocked_ano" comment:"被加入黑名单的人" `
	Reason         string `gorm:"column:reason;type:varchar(255);default:;comment:原因" json:"reason" comment:"原因" `
	Scope          int    `gorm:"column:scope;type:smallint;default:0;comment:范围 0全部 1全站2私信 3动态 4评论" json:"scope" comment:"范围" `
}

func (*ImBlacklist) TableName() string {
	return "social_blacklist"
}

func (*ImBlacklist) TableComment() string {
	return "黑名单"
}
