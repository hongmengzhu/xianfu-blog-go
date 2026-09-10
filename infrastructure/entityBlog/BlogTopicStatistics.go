package entityBlog

import (
	"time"
)

// BlogTopicStatisticsEntity 统计
type BlogTopicStatisticsEntity struct {
	ID             int64      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id" comment:"" `
	TopicNo        string     `gorm:"column:topic_no;type:varchar(80);index;default:;comment:编号" json:"topic_no" comment:"编号" `
	CreateAt       *time.Time `gorm:"column:create_at;type:datetime;index;autoCreateTime;default:current_timestamp;comment:创建时间" json:"create_at" comment:"创建时间" `
	UpdateAt       *time.Time `gorm:"column:update_at;type:datetime;autoUpdateTime;comment:更新时间" json:"update_at" comment:"更新时间" `
	Comment        int64      `gorm:"column:comment;type:bigint;not null;default:0;comment:评论" json:"comment" comment:"评论" `
	Read           int64      `gorm:"column:read;type:bigint;not null;default:0;comment:阅读" json:"read" comment:"阅读" `
	SeoKeywords    string     `gorm:"column:seo_keywords;type:varchar(255);default:;comment:seo关键词" json:"seo_keywords" comment:"seo关键词" `
	SeoDescription string     `gorm:"column:seo_description;type:varchar(255);default:;comment:seo描述" json:"seo_description" comment:"seo描述" `
	PageTitle      string     `gorm:"column:page_title;type:varchar(255);comment:网页标题" json:"page_title" comment:"网页标题" `
}

func (*BlogTopicStatisticsEntity) TableName() string {
	return "blog_topic_statistics"
}

func (*BlogTopicStatisticsEntity) TableComment() string {
	return "统计"
}
