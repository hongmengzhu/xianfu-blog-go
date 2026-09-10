package entityRam

import (
	"time"
)

// RamIdpMetadataCacheEntity IdP OIDC/SAML元数据缓存
type RamIdpMetadataCacheEntity struct {
	ID         int64  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id" comment:"主键ID"`
	No         string `gorm:"column:no;type:varchar(80);index;default:;comment:缓存编号" json:"no" comment:"缓存编号"`
	TenantNo   string `gorm:"column:tenant_no;type:varchar(80);index;default:;comment:租户编号" json:"tenant_no" comment:"租户编号"`
	TypeDomain string `gorm:"column:type_domain;type:varchar(80);index;default:'general';comment:域类型" json:"type_domain" comment:"域类型"`

	IdpNo    string `gorm:"column:idp_no;type:varchar(80);index;default:;comment:身份提供商编号" json:"idp_no" comment:"身份提供商编号"`
	SourceNo string `gorm:"column:source_no;type:varchar(80);index;default:;comment:认证源编号RamIdentitySource.No" json:"source_no" comment:"认证源编号"`
	Protocol string `gorm:"column:protocol;type:varchar(80);index;default:;comment:协议 oidc/saml2" json:"protocol" comment:"协议类型"`

	// 元数据标识
	EntityID    string `gorm:"column:entity_id;type:varchar(512);index;default:;comment:SAML EntityID / OIDC Issuer地址" json:"entity_id" comment:"Issuer/EntityID"`
	MetadataUrl string `gorm:"column:metadata_url;type:varchar(1024);comment:远程元数据拉取地址" json:"metadata_url" comment:"远程元数据地址"`

	// 缓存内容
	MetaFormat  string `gorm:"column:meta_format;type:varchar(40);default:json;comment:格式 json/xml" json:"meta_format" comment:"元数据格式 json/xml"`
	MetadataRaw string `gorm:"column:metadata_raw;type:text;comment:原始元数据内容" json:"-"`

	// 缓存生命周期
	LastFetchAt *time.Time `gorm:"column:last_fetch_at;type:datetime;comment:上次拉取时间" json:"last_fetch_at" comment:"上次拉取时间"`
	ExpireAt    *time.Time `gorm:"column:expire_at;type:datetime;index;comment:缓存过期时间" json:"expire_at" comment:"缓存过期时间"`
	CacheTtl    int64      `gorm:"column:cache_ttl;type:bigint;default:3600;comment:缓存有效期(秒)" json:"cache_ttl" comment:"缓存TTL，单位秒"`

	State int8 `gorm:"column:state;not null;index;default:1;comment:状态1有效2失效" json:"state" comment:"缓存状态：1有效 2失效"`

	CreateAt *time.Time `gorm:"column:create_at;type:datetime;index;autoCreateTime;default:current_timestamp;comment:创建时间" json:"create_at" comment:"创建时间"`
	UpdateAt *time.Time `gorm:"column:update_at;type:datetime;autoUpdateTime;comment:更新时间" json:"update_at" comment:"更新时间"`
}

// TableName RamIdpMetadataCacheEntity's table name
func (*RamIdpMetadataCacheEntity) TableName() string {
	return "ram_idp_metadata_cache"
}

func (*RamIdpMetadataCacheEntity) TableComment() string {
	return "IdP OIDC/SAML元数据缓存表"
}
