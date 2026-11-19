package entitySocial

import "github.com/hongmengzhu/xianfu-blog-go/pkg/tools/dbHelper/entityPg"

type SocialComment struct {
	entityPg.BaseCrUpAtEntity
	PostNo      string `gorm:"column:post_no;type:varchar(80);index;default:'';comment:帖子编号" json:"post_no" comment:"帖子编号" `
	Ano         string `gorm:"column:ano;type:varchar(80);index;default:'';comment:评论人编号" json:"ano" comment:"评论人编号" `
	TargetType  string `gorm:"column:target_type;type:varchar(80);index;default:'';comment:目标类型动态|评论" json:"target_type" comment:"目标类型" `
	ParentNo    string `gorm:"column:parent_no;type:varchar(80);index;default:'';comment:父评论编号" json:"parent_no" comment:"父评论编号" `
	RootNo      string `gorm:"column:root_no;type:varchar(80);index;default:'';comment:根评论编号" json:"root_no" comment:"根评论编号" `
	ReplyToAno  string `gorm:"column:reply_to_ano;type:varchar(80);index;default:'';comment:回复的评论人编号" json:"reply_to_ano" comment:"回复的评论人编号" `
	Content     string `gorm:"column:content;type:varchar(255);default:;comment:内容" json:"content" comment:"内容" `
	ContentType string `gorm:"column:content_type;type:varchar(80);index;default:'';comment:内容类型文本|图片" json:"content_type" comment:"内容类型" `
	LikeCount   int64  `gorm:"column:like_count;type:bigint;default:0;comment:点赞数" json:"like_count" comment:"点赞数" `
	ReplyCount  int64  `gorm:"column:reply_count;type:bigint;default:0;comment:回复数" json:"reply_count" comment:"回复数" `
	State       int8   `gorm:"column:state;not null;index;default:1;comment:状态|1启用|2禁用|11审核中" json:"state" comment:"状态:1启用;2禁用" `
}

func (SocialComment) TableName() string {
	return "social_comment"
}
func (SocialComment) TableComment() string {
	return "评论"
}
