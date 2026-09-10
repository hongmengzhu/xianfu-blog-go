package entityBasic

import (
	"time"
)

// BasicAttachmentEntity 附件
type BasicAttachmentEntity struct {
	ID            int64      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id" comment:"" `
	Name          string     `gorm:"column:name;type:varchar(255);comment:名称" json:"name" comment:"名称" `
	OriginalName  string     `gorm:"column:original_name;type:varchar(255);comment:原始名称" json:"original_name" comment:"原始名称" `
	State         int8       `gorm:"column:state;not null;index;default:1;comment:1有效2停用11取消(对应有效)12弃置(对应停用)13批量删除(无状态)" json:"state" comment:"1有效2停用11取消(对应有效)12弃置(对应停用)13批量删除(无状态)" ` // 1有效2停用11取消(对应有效)12弃置(对应停用)13批量删除(无状态)
	Description   string     `gorm:"column:description;type:varchar(255);comment:描述" json:"description" comment:"描述" `
	CreateAt      *time.Time `gorm:"column:create_at;type:datetime;index;autoCreateTime;default:current_timestamp;comment:创建时间" json:"create_at" comment:"创建时间" `
	UpdateAt      *time.Time `gorm:"column:update_at;type:datetime;autoUpdateTime;comment:更新时间" json:"update_at" comment:"更新时间" `
	CreateBy      string     `gorm:"column:create_by;type:varchar(80);index;default:;comment:创建人" json:"create_by" comment:"创建人" `
	UpdateBy      string     `gorm:"column:update_by;type:varchar(80);default:;comment:更新人" json:"update_by" comment:"更新人" `
	Sort          int64      `gorm:"column:sort;type:bigint;not null;default:0;;comment:排序" json:"sort" comment:"排序" `
	TenantNo      string     `gorm:"column:tenant_no;type:varchar(80);index;default:;comment:租户编号" json:"tenant_no" comment:"租户编号" `
	StoreNo       string     `gorm:"column:store_no;type:varchar(80);index;default:;comment:店编号" json:"store_no" comment:"店编号" `
	OrgNo         string     `gorm:"column:org_no;type:varchar(80);index;default:;comment:组织编号" json:"org_no" comment:"组织编号" `
	OwnerNo       string     `gorm:"column:owner_no;type:varchar(80);index;default:;comment:所有者编号" json:"owner_no" comment:"所有者编号"`
	Size          int64      `gorm:"column:size;type:bigint;not null;default:0;;comment:大小" json:"size" comment:"大小" `
	Module        string     `gorm:"column:module;type:varchar(180);index;comment:模块" json:"module" comment:"模块" `
	Tag           string     `gorm:"column:tag;type:varchar(80);comment:标签" json:"tag" comment:"标签" `
	Label         string     `gorm:"column:label;type:varchar(80);comment:标记" json:"label" comment:"标记" `
	File          string     `gorm:"column:file;type:varchar(980);comment:相对路径" json:"file" comment:"相对路径" `
	Domain        string     `gorm:"column:domain;type:varchar(380);comment:域名" json:"domain" comment:"域名" `
	Url           string     `gorm:"column:url;type:varchar(980);comment:路径" json:"url" comment:"路径" `
	Mark          string     `gorm:"column:mark;type:varchar(80);index;default:;comment:mark标记" json:"mark" comment:"mark标记" `
	Type          string     `gorm:"column:type;type:varchar(80);index;comment:类型" json:"type" comment:"类型" `
	No            string     `gorm:"column:no;type:varchar(80);index;comment:流水号" json:"no" comment:"流水号" `
	Method        string     `gorm:"column:method;type:varchar(80);index;comment:方式" json:"method" comment:"方式" `
	Ext           string     `gorm:"column:ext;type:varchar(80);index;comment:文件扩展名" json:"ext" comment:"文件扩展名" `
	Category      string     `gorm:"column:category;type:varchar(80);index;comment:上传分类|自己上传|远程下载|系统图" json:"category" comment:"上传分类" `
	Client        string     `gorm:"column:client;type:varchar(80);index;comment:客户端" json:"client" comment:"客户端" `
	ProtocolSpace string     `gorm:"column:protocol_space;type:varchar(80);index;comment:协议空间" json:"protocol_space" comment:"协议空间" `
	FileOwner     string     `gorm:"column:file_owner;type:varchar(80);index;comment:拥有者" json:"file_owner" comment:"拥有者/文件属于谁" `
	FileOwnerSub  string     `gorm:"column:file_owner_sub;type:varchar(80);index;comment:子类型" json:"file_owner_sub" comment:"子类型" `
	FileOwnerCopy string     `gorm:"column:file_owner_copy;type:varchar(80);index;comment:拥有者源" json:"file_owner_copy" comment:"拥有者源" `
	TypeData      string     `gorm:"column:type_data;type:varchar(80);comment:类型数据" json:"type_data" comment:"类型数据/源/复制" `
	RefNo         string     `gorm:"column:ref_no;type:varchar(80);comment:关联业务号" json:"ref_no" comment:"关联业务号" `
}

func (*BasicAttachmentEntity) TableName() string {
	return "basic_attachment"
}

func (*BasicAttachmentEntity) TableComment() string {
	return "附件"
}
