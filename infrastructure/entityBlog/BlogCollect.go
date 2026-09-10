package entityBlog

import (
	"time"

	"gorm.io/datatypes"
)

// BlogCollectEntity 收集数据
type BlogCollectEntity struct {
	ID               int64                        `gorm:"column:id;primaryKey;autoIncrement:true" json:"id" comment:"" `
	No               string                       `gorm:"column:no;type:varchar(80);default:;comment:编号代号" json:"no" comment:"编号代号" `
	Code             string                       `gorm:"column:code;type:varchar(80);default:;comment:标志" json:"code" comment:"标志" `
	Name             string                       `gorm:"column:name;type:varchar(955);comment:名称" json:"name" comment:"名称" `                                                                                // 名称
	NameFl           string                       `gorm:"column:name_fl;type:varchar(255);comment:名称外文" json:"name_fl" comment:"名称外文" `                                                                      // 名称外文
	NameFull         string                       `gorm:"column:name_full;type:varchar(255);comment:全称" json:"name_full" comment:"全称" `                                                                      // 全称
	State            int8                         `gorm:"column:state;not null;index;default:1;comment:1有效2停用11取消(对应有效)12弃置(对应停用)13批量删除(无状态)" json:"state" comment:"1有效2停用11取消(对应有效)12弃置(对应停用)13批量删除(无状态)" ` // 1有效2停用11取消(对应有效)12弃置(对应停用)13批量删除(无状态)
	Description      string                       `gorm:"column:description;type:text;comment:描述" json:"description" comment:"描述" `                                                                          // 描述
	CreateAt         *time.Time                   `gorm:"column:create_at;type:datetime;index;autoCreateTime;default:current_timestamp;comment:创建时间" json:"create_at" comment:"创建时间" `                       // 创建时间
	UpdateAt         *time.Time                   `gorm:"column:update_at;type:datetime;autoUpdateTime;comment:更新时间" json:"update_at" comment:"更新时间" `                                                       // 更新时间
	CreateBy         string                       `gorm:"column:create_by;type:varchar(80);index;default:'';comment:创建人" json:"createBy" comment:"创建人" `
	UpdateBy         string                       `gorm:"column:update_by;type:varchar(80);default:'';comment:更新人" json:"updateBy" comment:"更新人" `
	Sort             int64                        `gorm:"column:sort;type:bigint;not null;default:0;;comment:排序" json:"sort" comment:"排序" ` // 排序
	TenantNo         string                       `gorm:"column:tenant_no;type:varchar(80);index;default:;comment:租户编号" json:"tenant_no" comment:"租户编号" `
	OrgNo            string                       `gorm:"column:org_no;type:varchar(80);index;default:;comment:组织编号" json:"org_no" comment:"组织编号" `
	StoreNo          string                       `gorm:"column:store_no;type:varchar(80);index;default:;comment:店铺编号" json:"store_no" comment:"店铺编号" `
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
	Attachments      string                       `gorm:"column:attachments;type:text;comment:附件" json:"attachments" comment:"附件" `
	UrlSource        string                       `gorm:"column:url_source;type:varchar(2000);comment:来源地址" json:"url_source" comment:"来源地址"`
	UrlSourceMd5     string                       `gorm:"column:url_source_md5;type:varchar(80);index;default:;comment:来源地址md5" json:"url_source_md5" comment:"来源地址md5"`
	UrlRewriting     string                       `gorm:"column:url_rewriting;type:varchar(2000);comment:重写地址" json:"url_rewriting" comment:"重写地址"`
	Jump             int8                         `gorm:"column:jump;type:int8;not null;default:2;comment:跳转" json:"jump" comment:"跳转"`
	Source           string                       `gorm:"column:source;type:varchar(255);comment:来源" json:"source" comment:"来源"`
	TypeComment      string                       `gorm:"column:type_comment;type:varchar(80);index;default:no;comment:评论类型|允许评论|不允许评论;" json:"type_comment" comment:"评论类型|允许评论|不允许评论" `
	OperationTime    *time.Time                   `gorm:"column:operation_time;type:datetime;comment:操作时间" json:"operation_time" comment:"操作时间" `
	PublishTime      *time.Time                   `gorm:"column:publish_time;type:datetime;index;comment:发布时间|源文章发布时间" json:"publish_time" comment:"发布时间|源文章发布时间" `
	TypeReading      string                       `gorm:"column:type_reading;type:varchar(80);index;default:no;comment:阅读类型|未看|在看|已看;" json:"type_reading" comment:"阅读类型||未看|在看|已看" `
	TypeDomain       string                       `gorm:"column:type_domain;type:varchar(80);index;default:no;comment:领域|;" json:"type_domain" comment:"领域" `
	Remark           string                       `gorm:"column:remark;type:varchar(955);comment:备注" json:"remark" comment:"备注"`
	Comment          int64                        `gorm:"column:comment;type:bigint;not null;default:0;comment:评论" json:"comment" comment:"评论" `
	Read             int64                        `gorm:"column:read;type:bigint;not null;default:0;comment:阅读" json:"read" comment:"阅读" `
	SeoKeywords      string                       `gorm:"column:seo_keywords;type:varchar(255);default:;comment:seo关键词" json:"seo_keywords" comment:"seo关键词" `
	SeoDescription   string                       `gorm:"column:seo_description;type:varchar(255);default:;comment:seo描述" json:"seo_description" comment:"seo描述" `
	PageTitle        string                       `gorm:"column:page_title;type:varchar(255);comment:网页标题" json:"page_title" comment:"网页标题" `
}

func (*BlogCollectEntity) TableName() string {
	return "blog_collect"
}

func (*BlogCollectEntity) TableComment() string {
	return "收集数据"
}
