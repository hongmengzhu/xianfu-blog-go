package entityRam

import "time"

// RamIdpCredentialEntity 身份 第三方认证凭证、证书密钥
type RamIdpCredentialEntity struct {
	ID         int64      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id" comment:""`
	No         string     `gorm:"column:no;type:varchar(80);index;default:;comment:编号" json:"no" comment:"编号" `
	TenantNo   string     `gorm:"column:tenant_no;type:varchar(80);index;default:;comment:租户编号" json:"tenant_no" comment:"租户编号" ` // 租户
	OrgNo      string     `gorm:"column:org_no;type:varchar(80);index;default:;comment:组织编号" json:"org_no" comment:"组织编号" `
	StoreNo    string     `gorm:"column:store_no;type:varchar(80);index;default:;comment:店编号" json:"store_no" comment:"店编号" `
	TypeDomain string     `gorm:"column:type_domain;type:varchar(80);index;default:'general';comment:域类型" json:"type_domain" comment:"域类型系统-商户" `
	CreateAt   *time.Time `gorm:"column:create_at;type:datetime;index;autoCreateTime;default:current_timestamp;comment:创建时间" json:"create_at" comment:"创建时间" `
	CreateBy   string     `gorm:"column:create_by;type:varchar(80);index;default:;comment:创建人" json:"create_by" comment:"创建人" `
	Idp        string     `gorm:"column:idp;type:varchar(80);index;default:;comment:身份 提供商" json:"idp" comment:"身份 提供商" `
	SourceNo   string     `gorm:"column:source_no;index;comment:关联RamIdentitySource.No" json:"source_no"`
	CredType   string     `gorm:"column:cred_type;comment:类型 client_secret/sp_cert/sp_private_key/ldap_bind_pwd/jwks" json:"cred_type"`
	Value      string     `gorm:"column:value;type:text;comment:加密存储密钥内容" json:"-"` // 不返回前端
	ExpireAt   *time.Time `gorm:"column:expire_at;comment:过期时间，证书使用" json:"expire_at"`
	State      int8       `gorm:"column:state;not null;index;default:1;comment:状态|1启用|2禁用" json:"state" comment:"状态:1启用;2禁用" ` // 状态:1启用;2禁用
}

// TableName RamIdpCredentialEntity's table name
func (*RamIdpCredentialEntity) TableName() string {
	return "ram_idp_credential"
}

func (*RamIdpCredentialEntity) TableComment() string {
	return "认证源密钥证书凭证"
}
