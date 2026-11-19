package entitySocial

import "github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/entityPg"

type SocialFollow struct {
	entityPg.BaseCrUpAtEntity
	ToAno        string `gorm:"column:to_ano;type:varchar(80);index;default:'';comment:被关注者" json:"to_ano" comment:"被关注者" `
	FollowerAno  string `gorm:"column:follower_ano;type:varchar(80);index;default:'';comment:关注者" json:"follower_ano" comment:"关注者" `
	State        int8   `gorm:"column:state;not null;index;default:1;comment:状态|1关注|2取关" json:"state" comment:"状态:1关注|2取关" `
	RelationType int    `gorm:"column:relation_type;type:smallint;default:0;comment:关系类型1关注2特别关注" json:"relation_type" comment:"关系类型" `
	Source       string `gorm:"column:source;type:varchar(80);default:;comment:来源profile/post/search" json:"source" comment:"来源" `
	HotFollower  int    `gorm:"column:hot_follower;type:smallint;default:0;comment:关注热度;1特别活跃-1冷" json:"hot_follower" comment:"关注热度" `
	ExtraInfo    string `gorm:"column:extra_info;type:varchar(255);default:;comment:扩展信息" json:"extra_info" comment:"扩展信息" `
}

func (*SocialFollow) TableName() string {
	return "social_follower"
}

func (*SocialFollow) TableComment() string {
	return "关注/粉丝"
}
