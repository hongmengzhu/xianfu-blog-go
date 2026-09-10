package entityRam

import "time"

// RamAccountAuthMfaEntity TOTP MFA、FaceID、指纹等二次验证
type RamAccountAuthMfaEntity struct {
	ID               int64      `gorm:"column:id;primaryKey;autoIncrement:true"`
	TenantNo         string     `gorm:"column:tenant_no;index"`
	Ano              string     `gorm:"column:ano;index"`
	MfaType          string     `gorm:"column:mfa_type;comment:MFA类型 app/sms/email"`
	Secret           string     `gorm:"column:secret;type:text;comment:TOTP密钥"`
	CredentialID     string     `gorm:"column:cred_id;uniqueIndex;comment:凭证ID base64"`
	PublicKey        string     `gorm:"column:public_key;type:text;comment:公钥"`
	DeviceName       string     `gorm:"column:device_name;comment:设备名称"`
	Transport        string     `gorm:"column:transport;type:text;comment:json数组 ble/nfc/usb(internal"`
	Counter          uint64     `gorm:"column:counter;comment:签名计数器，防重放"`
	RecoveryCodes    string     `gorm:"column:recovery_codes;type:text;comment:恢复码JSON数组"`
	MfaToken         string     `gorm:"column:mfa_token;index;comment:MFA临时令牌"`
	MfaTokenExpireAt *time.Time `gorm:"column:mfa_token_expire_at;comment:MFA令牌过期时间"`
	State            int8       `gorm:"column:state;not null;index;default:1;comment:1有效2停用" json:"state" comment:"1有效2停用" `
	LastUsedAt       *time.Time `gorm:"column:last_used_at"`
	CreateAt         *time.Time `gorm:"column:create_at;type:datetime;index;autoCreateTime;default:current_timestamp;comment:创建时间" json:"create_at" comment:"创建时间" `
}

func (*RamAccountAuthMfaEntity) TableName() string { return "ram_account_auth_mfa" }
func (*RamAccountAuthMfaEntity) TableComment() string {
	return "用户多因素认证表 TOTP / 生物识别"
}
