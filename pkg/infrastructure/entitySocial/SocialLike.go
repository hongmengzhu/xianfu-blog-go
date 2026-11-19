package entitySocial

import "github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/entityPg"

type SocialLike struct {
	entityPg.BaseCrUpAtEntity
	Ano        string `gorm:"column:ano;type:varchar(80);index;default:'';comment:点赞人" json:"ano" comment:"点赞人" `
	TargetType string `gorm:"column:target_type;type:varchar(80);index;default:'';comment:目标类型动态|评论" json:"target_type" comment:"目标类型" `
	TargetAno  string `gorm:"column:target_ano;type:varchar(80);index;default:'';comment:目标人" json:"target_ano" comment:"目标人" `
	Scope      int    `gorm:"column:scope;type:smallint;default:0;comment:范围 0全部 1全站2私信 3动态 4评论" json:"scope" comment:"范围" `
	State      int8   `gorm:"column:state;not null;index;default:1;comment:状态|1启用|2禁用" json:"state" comment:"状态:1启用;2禁用" `
}

func (SocialLike) TableName() string {
	return "social_like"
}
func (SocialLike) TableComment() string {
	return "点赞"
}
