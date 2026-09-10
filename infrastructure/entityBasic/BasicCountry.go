package entityBasic

import (
	"time"
)

// BasicCountryEntity 国家
type BasicCountryEntity struct {
	ID           int64      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id" comment:"" `
	No           string     `gorm:"column:no;type:varchar(80);index;default:;comment:编号" json:"no" comment:"编号" `
	Name         string     `gorm:"column:name;type:varchar(255);comment:名称" json:"name" comment:"名称" `                                                                                // 名称
	NameFl       string     `gorm:"column:name_fl;type:varchar(255);comment:名称外文" json:"name_fl" comment:"名称外文" `                                                                      // 名称外文
	Code         string     `gorm:"column:code;type:varchar(80);comment:标志" json:"code" comment:"标志" `                                                                                 // 编号代号
	NameFull     string     `gorm:"column:name_full;type:varchar(255);comment:全称" json:"name_full" comment:"全称" `                                                                      // 全称
	State        int8       `gorm:"column:state;not null;index;default:1;comment:1有效2停用11取消(对应有效)12弃置(对应停用)13批量删除(无状态)" json:"state" comment:"1有效2停用11取消(对应有效)12弃置(对应停用)13批量删除(无状态)" ` // 1有效2停用11取消(对应有效)12弃置(对应停用)13批量删除(无状态)
	Description  string     `gorm:"column:description;type:varchar(255);comment:描述" json:"description" comment:"描述" `                                                                  // 描述
	CreateAt     *time.Time `gorm:"column:create_at;type:datetime;index;autoCreateTime;default:current_timestamp;comment:创建时间" json:"create_at" comment:"创建时间" `                       // 创建时间
	UpdateAt     *time.Time `gorm:"column:update_at;type:datetime;autoUpdateTime;comment:更新时间" json:"update_at" comment:"更新时间" `                                                       // 更新时间
	CreateBy     string     `gorm:"column:create_by;type:varchar(80);index;default:;comment:创建人" json:"create_by" comment:"创建人" `
	UpdateBy     string     `gorm:"column:update_by;type:varchar(80);default:;comment:更新人" json:"update_by" comment:"更新人" `
	Sort         int64      `gorm:"column:sort;type:bigint;not null;default:0;;comment:排序" json:"sort" comment:"排序" `  // 排序
	ParentId     string     `gorm:"column:parent_id;type:varchar(80);index;comment:上级" json:"parent_id" comment:"上级" ` // 上级
	IdLink       string     `gorm:"column:id_link;type:varchar(5000);comment:上级" json:"id_link" comment:"上级链" `
	ParentNo     string     `gorm:"column:parent_no;type:varchar(80);index;default:;comment:上级编号" json:"parent_no" comment:"上级编号" `
	NoLink       string     `gorm:"column:no_link;type:text;comment:上级编号" json:"no_link" comment:"上级编号" `
	Continent    string     `gorm:"column:continent;type:varchar(80);comment:所属洲" json:"continent" comment:"所属洲" `
	Iso3         string     `gorm:"column:iso3;type:varchar(80);comment:ISO三字代码" json:"iso3" comment:"ISO三字代码" `
	CountryCode  string     `gorm:"column:country_code;type:varchar(20);comment:国际区号" json:"country_code" comment:"国际区号" `
	DomainSuffix string     `gorm:"column:domain_suffix;type:varchar(20);comment:域名后缀" json:"domain_suffix" comment:"域名后缀" `
	PhoneUse     int8       `gorm:"column:phone_use;type:int8;not null;index;default:2;comment:电话使用1是2否" json:"phone_use" comment:"电话使用1是2否" `
}

//in international number 国际号码

// TableName IamGroupEntity's table name
func (*BasicCountryEntity) TableName() string {
	return "basic_country"
}

func (*BasicCountryEntity) TableComment() string {
	return "国家"
}
