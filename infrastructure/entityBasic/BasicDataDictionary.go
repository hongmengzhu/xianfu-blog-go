package entityBasic

import "time"

// BasicDataDictionaryEntity 数据字典
type BasicDataDictionaryEntity struct {
	ID            int64      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id" comment:"" `
	No            string     `gorm:"column:no;type:varchar(80);index;default:;comment:编号" json:"no" comment:"编号" `
	Name          string     `gorm:"column:name;type:varchar(255);comment:名称" json:"name" comment:"名称" `
	NameFl        string     `gorm:"column:name_fl;type:varchar(255);comment:名称外文" json:"name_fl" comment:"名称外文" `
	Code          string     `gorm:"column:code;type:varchar(100);index;comment:标记" json:"code" comment:"标记" `
	NameFull      string     `gorm:"column:name_full;type:varchar(255);comment:全称" json:"name_full" comment:"全称" `
	State         int8       `gorm:"column:state;not null;index;default:1;comment:1有效2停用11取消(对应有效)12弃置(对应停用)13批量删除(无状态)" json:"state" comment:"1有效2停用11取消(对应有效)12弃置(对应停用)13批量删除(无状态)" ` // 1有效2停用11取消(对应有效)12弃置(对应停用)13批量删除(无状态)
	Description   string     `gorm:"column:description;type:varchar(255);comment:描述" json:"description" comment:"描述" `
	CreateAt      *time.Time `gorm:"column:create_at;type:datetime;index;autoCreateTime;default:current_timestamp;comment:创建时间" json:"create_at" comment:"创建时间" `
	UpdateAt      *time.Time `gorm:"column:update_at;type:datetime;autoUpdateTime;comment:更新时间" json:"update_at" comment:"更新时间" `
	CreateBy      string     `gorm:"column:create_by;type:varchar(80);index;default:;comment:创建人" json:"create_by" comment:"创建人" `
	UpdateBy      string     `gorm:"column:update_by;type:varchar(80);default:;comment:更新人" json:"update_by" comment:"更新人" `
	Sort          int64      `gorm:"column:sort;type:bigint;not null;index;default:0;comment:排序" json:"sort" comment:"排序" `
	TypeUniqueMd5 string     `gorm:"column:type_unique_md5;type:varchar(80);not null;index;default:;comment:当前类型md5唯一" json:"type_unique_md5" comment:"当前类型md5唯一" `
	Value         string     `gorm:"column:value;type:varchar(255);comment:值内容" json:"value" comment:"值内容" `
	Extend        string     `gorm:"column:extend;type:text;comment:扩展参数" json:"extend" comment:"扩展参数" `
	Range         string     `gorm:"column:range;type:text;comment:范围" json:"range" comment:"范围" `
	TypeCode      string     `gorm:"column:type_code;type:varchar(100);index;default:;comment:类型编号/码值" json:"type_code" comment:"类型编号/码值" `
	TenantNo      string     `gorm:"column:tenant_no;type:varchar(80);index;default:;comment:租户编号" json:"tenant_no" comment:"租户编号" `
	StoreNo       string     `gorm:"column:store_no;type:varchar(80);index;default:;comment:店铺编号" json:"store_no" comment:"店铺编号" `
	OwnerNo       string     `gorm:"column:owner_no;type:varchar(80);index;default:;comment:所属编号" json:"owner_no" comment:"所属编号" `
}

func (*BasicDataDictionaryEntity) TableName() string {
	return "basic_data_dictionary"
}

func (*BasicDataDictionaryEntity) TableComment() string {
	return "数据字典"
}
