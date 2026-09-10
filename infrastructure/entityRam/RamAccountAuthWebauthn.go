package entityRam

import "time"

// RamAccountAuthWebauthnEntity WebAuthn Passkey通行密钥
type RamAccountAuthWebauthnEntity struct {
	ID           int64      `gorm:"column:id;primaryKey;autoIncrement:true"`
	TenantNo     string     `gorm:"column:tenant_no;index"`
	Ano          string     `gorm:"column:ano;index"`
	CredentialID string     `gorm:"column:cred_id;uniqueIndex;comment:凭证ID base64"`
	PublicKey    string     `gorm:"column:public_key;type:text;comment:公钥"`
	DeviceName   string     `gorm:"column:device_name;comment:设备名称"`
	Transport    string     `gorm:"column:transport;type:text;comment:json数组 ble/nfc/usb/internal"`
	Counter      uint64     `gorm:"column:counter;comment:签名计数器，防重放"`
	Enabled      int8       `gorm:"column:enabled;default:1"`
	LastUsedAt   *time.Time `gorm:"column:last_used_at"`
	CreateAt     *time.Time `gorm:"column:create_at;type:datetime;index;autoCreateTime;default:current_timestamp;comment:创建时间" json:"create_at" comment:"创建时间" `
}

func (*RamAccountAuthWebauthnEntity) TableName() string    { return "ram_account_auth_webauthn" }
func (*RamAccountAuthWebauthnEntity) TableComment() string { return "WebAuthn Passkey通行密钥" }
