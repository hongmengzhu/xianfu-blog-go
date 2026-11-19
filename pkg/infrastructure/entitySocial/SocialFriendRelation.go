package entitySocial

import (
	"time"

	"github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/entityPg"
)

type SocialFriendRelation struct {
	entityPg.BaseCrUpAtEntity
	RequesterAno string     `gorm:"column:requester_ano;type:varchar(80);index;default:'';comment:请求者" json:"requester_ano" comment:"请求者" `
	ToAno        string     `gorm:"column:to_ano;type:varchar(80);index;default:'';comment:被请求者" json:"to_ano" comment:"被请求者" `
	Source       string     `gorm:"column:source;type:varchar(80);default:;comment:来源profile/post/search" json:"source" comment:"来源" `
	AcceptedAt   *time.Time `gorm:"column:accepted_at;default:;comment:接受时间" json:"accepted_at" comment:"接受时间" `
	State        int8       `gorm:"column:state;not null;index;default:1;comment:关系状态|1同意2拒绝" json:"state" comment:"关系状态:1同意2拒绝" `
}

func (*SocialFriendRelation) TableName() string {
	return "social_friend_relation"
}

func (*SocialFriendRelation) TableComment() string {
	return "好友关系"
}
