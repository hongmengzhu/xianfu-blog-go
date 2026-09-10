package entityRam

import (
	"time"
)

// RamAccountSessionLogEntity SSO登录日志
type RamAccountSessionLogEntity struct {
	ID         int64      `gorm:"column:id;primaryKey;autoIncrement:true" json:"id" comment:"" `
	No         string     `gorm:"column:no;type:varchar(80);index;default:;comment:编号代号" json:"no" comment:"编号代号" `
	TenantNo   string     `gorm:"column:tenant_no;type:varchar(80);index;default:;comment:租户编号" json:"tenant_no" comment:"租户编号" `
	OrgNo      string     `gorm:"column:org_no;type:varchar(80);index;default:;comment:组织编号" json:"org_no" comment:"组织编号" `
	StoreNo    string     `gorm:"column:store_no;type:varchar(80);index;default:;comment:店编号" json:"store_no" comment:"店编号" `
	CreateAt   *time.Time `gorm:"column:create_at;type:datetime;index;autoCreateTime;default:current_timestamp;comment:创建时间" json:"create_at" comment:"创建时间" `
	TypeDomain string     `gorm:"column:type_domain;type:varchar(80);index;default:'general';comment:域类型" json:"type_domain" comment:"域类型系统-商户" `
	// 关联认证源/IdP
	Idp      string `gorm:"column:idp;type:varchar(80);index;default:;comment:身份 提供商" json:"idp" comment:"身份 提供商" `
	SourceNo string `gorm:"column:source_no;index;comment:关联RamIdentitySource.No" json:"source_no"`
	Protocol string `gorm:"column:protocol;type:varchar(80);index;default:;comment:协议|openid|oauth" json:"protocol" comment:"协议|openid|oauth" `
	// 用户信息
	Ano         string `gorm:"column:ano;type:varchar(80);index;default:;comment:账号" json:"ano" comment:"账号"`
	BindAno     string `gorm:"column:bind_ano;type:varchar(80);index;default:;comment:绑定账号" json:"bind_ano" comment:"绑定账号" `
	ExternalSub string `gorm:"column:external_sub;type:varchar(255);index;default:;comment:外部唯一标识sub/NameID" json:"external_sub" comment:"外部身份唯一标识"`
	// 会话信息
	SessionId string `gorm:"column:session_id;type:varchar(255);index;default:;comment:服务端会话ID" json:"session_id" comment:"会话ID"`
	ClientNo  string `gorm:"column:client_no;type:varchar(80);index;default:;comment:客户端编号" json:"client_no" comment:"客户端编号"`
	// 事件类型与结果
	EventCategory string `gorm:"column:event_category;type:varchar(80);index;default:;comment:事件分类 login/logout/slo/mfa_sync" json:"event_category" comment:"事件分类：登录/登出/单点登出/MFA同步"`
	EventType     string `gorm:"column:event_type;type:varchar(80);index;default:;comment:事件细分 login_success/login_fail/logout_init/slo_callback" json:"event_type" comment:"事件细分类型"`
	EventResult   int8   `gorm:"column:event_result;type:int8;not null;index;default:0;comment:结果0未知 1成功 2失败" json:"event_result" comment:"执行结果：0未知,1成功,2失败"`
	FailReason    string `gorm:"column:fail_reason;type:varchar(512);comment:失败原因描述" json:"fail_reason" comment:"失败原因"`

	// 客户端信息
	IpAddress string `gorm:"column:ip_address;type:varchar(128);comment:客户端IP" json:"ip_address" comment:"客户端IP"`
	UserAgent string `gorm:"column:user_agent;type:text;comment:UA" json:"user_agent" comment:"浏览器/客户端UA"`

	// 时间维度
	OperateAt *time.Time `gorm:"column:operate_at;type:datetime;index;comment:事件发生时间" json:"operate_at" comment:"事件发生时间"`
	LoginAt   *time.Time `gorm:"column:login_at;type:datetime;index;comment:登录时间" json:"login_at" comment:"登录时间"`
	LogoutAt  *time.Time `gorm:"column:logout_at;type:datetime;index;comment:登出时间" json:"logout_at" comment:"登出时间"`
	ExpireAt  *time.Time `gorm:"column:expire_at;type:datetime;index;comment:会话过期时间" json:"expire_at" comment:"会话过期时间"`

	// 扩展原始报文（调试审计）
	RawRequest string `gorm:"column:raw_request;type:text;comment:原始请求参数JSON" json:"-"`
	ExtraData  string `gorm:"column:extra_data;type:text;comment:扩展数据JSON" json:"extra_data" comment:"扩展数据"`

	LoginSource string `gorm:"column:login_source;type:varchar(80);index;default:;comment:登陆来源" json:"login_source" comment:"登陆来源" `
	AppNo       string `gorm:"column:app_no;type:varchar(80);index;default:;comment:应用编号" json:"app_no" comment:"应用编号" `
	System      string `gorm:"column:system;type:varchar(80);index;default:;comment:系统" json:"system" comment:"系统"`
	Device      string `gorm:"column:device;type:varchar(80);index;default:'';comment:设备" json:"device" comment:"设备" `
	DeviceNo    string `gorm:"column:device_no;type:varchar(80);index;default:'';comment:设备编号" json:"device_no" comment:"设备编号" `
	Version     string `gorm:"column:version;type:varchar(80);default:'';comment:版本" json:"version" comment:"版本" `
	Ip          string `gorm:"column:ip;type:varchar(80);index;default:'';comment:ip" json:"ip" comment:"ip" `
}

// TableName RamAccountLoginLog's table name
func (*RamAccountSessionLogEntity) TableName() string {
	return "ram_account_session_log"
}

func (*RamAccountSessionLogEntity) TableComment() string {
	return "身份认证会话审计日志-SSO登录日志"
}
