package entityBlog

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
)

// BlogTopicEntity 文章主题
type BlogTopicEntity struct {
	ID               int64                        `gorm:"column:id;primaryKey;autoIncrement:true" json:"id" comment:"" `
	No               string                       `gorm:"column:no;type:varchar(80);default:;comment:编号代号" json:"no" comment:"编号代号" `
	Code             string                       `gorm:"column:code;type:varchar(80);default:;comment:标志" json:"code" comment:"标志" `
	Name             string                       `gorm:"column:name;type:varchar(255);comment:名称" json:"name" comment:"名称" `                                                                                // 名称
	NameFl           string                       `gorm:"column:name_fl;type:varchar(255);comment:名称外文" json:"name_fl" comment:"名称外文" `                                                                      // 名称外文
	NameFull         string                       `gorm:"column:name_full;type:varchar(255);comment:全称" json:"name_full" comment:"全称" `                                                                      // 全称
	State            int8                         `gorm:"column:state;not null;index;default:1;comment:1有效2停用11取消(对应有效)12弃置(对应停用)13批量删除(无状态)" json:"state" comment:"1有效2停用11取消(对应有效)12弃置(对应停用)13批量删除(无状态)" ` // 1有效2停用11取消(对应有效)12弃置(对应停用)13批量删除(无状态)
	Description      string                       `gorm:"column:description;type:varchar(255);comment:描述" json:"description" comment:"描述" `                                                                  // 描述
	CreateAt         *time.Time                   `gorm:"column:create_at;type:datetime;index;autoCreateTime;default:current_timestamp;comment:创建时间" json:"create_at" comment:"创建时间" `                       // 创建时间
	UpdateAt         *time.Time                   `gorm:"column:update_at;type:datetime;autoUpdateTime;comment:更新时间" json:"update_at" comment:"更新时间" `                                                       // 更新时间
	CreateBy         string                       `gorm:"column:create_by;type:varchar(80);index;default:'';comment:创建人" json:"createBy" comment:"创建人" `
	UpdateBy         string                       `gorm:"column:update_by;type:varchar(80);default:'';comment:更新人" json:"updateBy" comment:"更新人" `
	Sort             int64                        `gorm:"column:sort;type:bigint;not null;default:0;;comment:排序" json:"sort" comment:"排序" `
	TenantNo         string                       `gorm:"column:tenant_no;type:varchar(80);index;default:;comment:租户编号" json:"tenant_no" comment:"租户编号" `
	Version          string                       `gorm:"column:version;type:varchar(32);comment:版本" json:"version" comment:"版本" `
	Ano              string                       `gorm:"column:ano;type:varchar(80);index;default:'';comment:操作员" json:"ano" comment:"操作员" `
	PlatformApproved string                       `gorm:"column:platform_approved;type:varchar(80);index;default:draft;comment:平台审批状态" json:"platform_approved" comment:"平台审批状态|通过|未审批|驳回|进行中|提交" `
	Content          string                       `gorm:"column:content;type:text;comment:内容" json:"content" comment:"内容" `
	Editor           string                       `gorm:"column:editor;type:varchar(80);index;default:markdown;comment:编辑器类型" json:"editor" comment:"编辑器类型" `
	CategoryNo       string                       `gorm:"column:category_no;type:varchar(80);index;comment:上级" json:"category_no" comment:"上级" `
	Tags             datatypes.JSONType[[]string] `gorm:"column:tags;type:jsonb;index;default:'[]';comment:标签" json:"tags" comment:"标签" `
	Author           string                       `gorm:"column:author;type:varchar(80);index;default:'';comment:作者" json:"author" comment:"作者" `
	TypeContent      string                       `gorm:"column:type_content;type:varchar(80);index;default:original;comment:内容类型|原创|翻译|转载;" json:"type_content" comment:"类型|原创|翻译|转载" `
	TypeSource       string                       `gorm:"column:type_source;type:varchar(80);index;default:handwritten;comment:类型源|采集|手写;" json:"type_source" comment:"类型源|采集|手写" `
	TypeDataSource   string                       `gorm:"column:type_data_source;type:varchar(80);index;default:'general';comment:数据源类型|本平台|外部;" json:"type_data_source" comment:"数据源类型|本平台|外部" `
	Where            datatypes.JSONType[[]string] `gorm:"column:where;type:jsonb;index;default:'[]';comment:可用范围" json:"where" comment:"可用范围" `
	Points           decimal.Decimal              `gorm:"column:points;type:bigint;not null;default:0;comment:积分" json:"points" comment:"积分" `
	PriceShop        decimal.Decimal              `gorm:"column:price_shop;type:numeric(20,6);not null;default:0;comment:商城价" json:"price_shop" comment:"商城价" `
	Attachments      string                       `gorm:"column:attachments;type:text;comment:附件" json:"attachments" comment:"附件" `
	UrlSource        string                       `gorm:"column:url_source;type:varchar(2000);comment:来源地址" json:"url_source" comment:"来源地址"`
	UrlRewriting     string                       `gorm:"column:url_rewriting;type:varchar(2000);comment:重写地址" json:"url_rewriting" comment:"重写地址"`
	Jump             int8                         `gorm:"column:jump;type:int8;not null;default:2;comment:跳转" json:"jump" comment:"跳转"`
	Source           string                       `gorm:"column:source;type:varchar(255);comment:来源" json:"source" comment:"来源"`
	TypeComment      string                       `gorm:"column:type_comment;type:varchar(80);index;default:no;comment:评论类型|允许评论|不允许评论;" json:"type_comment" comment:"评论类型|允许评论|不允许评论" `
	OperationTime    *time.Time                   `gorm:"column:operation_time;type:datetime;comment:操作时间" json:"operation_time" comment:"操作时间" `
	TypeReading      string                       `gorm:"column:type_reading;type:varchar(80);index;default:no;comment:阅读类型|未看|在看|已看;" json:"type_reading" comment:"阅读类型||未看|在看|已看" `
}

func (*BlogTopicEntity) TableName() string {
	return "blog_topic"
}

func (*BlogTopicEntity) TableComment() string {
	return "文章主题"
}
