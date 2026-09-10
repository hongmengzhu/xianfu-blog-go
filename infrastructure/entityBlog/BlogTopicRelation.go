package entityBlog

import (
	"time"
)

// BlogTopicRelationEntity 文章主题关系
type BlogTopicRelationEntity struct {
	ID          int64      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id" comment:"" `
	CreateAt    *time.Time `gorm:"column:create_at;type:datetime;index;autoCreateTime;default:current_timestamp;comment:创建时间" json:"create_at" comment:"创建时间" `
	CreateBy    string     `gorm:"column:create_by;type:varchar(80);index;default:'';comment:创建人" json:"createBy" comment:"创建人" `
	Sort        int64      `gorm:"column:sort;type:bigint;not null;default:0;;comment:排序" json:"sort" comment:"排序" `
	TenantNo    string     `gorm:"column:tenant_no;type:varchar(80);index;default:;comment:租户编号" json:"tenant_no" comment:"租户编号" `
	TopicNo     string     `gorm:"column:topic_no;type:varchar(80);index;default:;comment:话题编号" json:"topic_no" comment:"话题编号" `
	ArticleNo   string     `gorm:"column:article_no;type:varchar(80);index;default:;comment:文章编号" json:"article_no" comment:"文章编号" `
	Ano         string     `gorm:"column:ano;type:varchar(80);index;default:'';comment:操作员" json:"ano" comment:"操作员" `
	Version     string     `gorm:"column:version;type:varchar(32);comment:版本" json:"version" comment:"版本" `
	Name        string     `gorm:"column:name;type:varchar(255);comment:名称" json:"name" comment:"名称" `
	Description string     `gorm:"column:description;type:varchar(255);comment:描述" json:"description" comment:"描述" `
}

func (*BlogTopicRelationEntity) TableName() string {
	return "blog_topic_relation"
}

func (*BlogTopicRelationEntity) TableComment() string {
	return "文章主题关系"
}
