package entitySocial

import "github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/entityPg"

type SocialFllower struct {
	entityPg.BaseCrUpAtEntity
	ToAno       string `gorm:"column:to_ano;type:varchar(80);index;default:'';comment:被关注者" json:"to_ano" comment:"被关注者" `
	FollowerAno string `gorm:"column:follower_ano;type:varchar(80);index;default:'';comment:关注者" json:"follower_ano" comment:"关注者" `
	Status      int    `gorm:"column:status;type:smallint;default:0;comment:关系状态1关注2特别关注-1取关-2拉黑" json:"status" comment:"关系状态" `
	HotFollower int    `gorm:"column:hot_follower;type:smallint;default:0;comment:关注热度;1特别活跃-1冷" json:"hot_follower" comment:"关注热度" `
	ExtraInfo   string `gorm:"column:extra_info;type:varchar(255);default:;comment:扩展信息" json:"extra_info" comment:"扩展信息" `
}

func (*SocialFllower) TableName() string {
	return "social_follower"
}

func (*SocialFllower) TableComment() string {
	return "粉丝"
}
